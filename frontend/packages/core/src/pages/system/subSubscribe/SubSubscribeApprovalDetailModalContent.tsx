/*
 * @Date: 2024-04-19 15:22:46
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-05 15:15:29
 * @FilePath: \frontend\packages\core\src\pages\system\subSubscribe\SubSubscribeApprovalDetailModalContent.tsx
 */
import { App, Form, Row, Col, Checkbox, Input } from "antd"
import { forwardRef, useImperativeHandle, useEffect } from "react"
import WithPermission from "@common/components/aoplatform/WithPermission"
import { BasicResponse, STATUS_CODE } from "@common/const/const"
import { useFetch } from "@common/hooks/http"
import { SYSTEM_SUBSCRIBE_APPROVAL_DETAIL_LIST } from "../../../const/system/const"
import { SubSubscribeApprovalModalHandle, SubSubscribeApprovalModalProps } from "../../../const/system/type"
import { FieldType } from "../../../const/user/types"

export const SubSubscribeApprovalModalContent = forwardRef<SubSubscribeApprovalModalHandle,SubSubscribeApprovalModalProps>((props, ref) => {
    const { message } = App.useApp()
    const {data, type, systemId} = props
    const [form] = Form.useForm();
    const {fetchData} = useFetch()

    const reApply:()=>Promise<boolean | string> =  ()=>{
        return new Promise((resolve, reject)=>{
            if(type === 'view'){
                resolve(true)
                return
            }
            form.validateFields().then((value)=>{
                fetchData<BasicResponse<null>>('catalogue/service/subscribe',{method: 'POST',eoBody:({partitions:value.partitions,service:data!.service.id, applications:[systemId], reason:value.reason})}).then(response=>{
                    const {code,msg} = response
                    if(code === STATUS_CODE.SUCCESS){
                        message.success(msg || '操作成功！')
                        resolve(true)
                    }else{
                        message.error(msg || '操作失败')
                        reject(msg || '操作失败')
                    }
                }).catch((errorInfo)=> reject(errorInfo))
            }).catch((errorInfo)=> reject(errorInfo))
        })
    }

    useImperativeHandle(ref, ()=>({
        reApply
        })
    )

    useEffect(()=>{
        form.setFieldsValue({...data,partitions:data?.partition?.map((x:{id:string})=>x.id)})
    },[])


    return (
        <div className="my-btnybase">
        <WithPermission access="">
            <Form
                layout='vertical'
                labelAlign='left'
                scrollToFirstError
                form={form}
                className="mx-auto "
                name="subSubscribeApprovalDetailModalContent"
                // labelCol={{ span: 8 }}
                // wrapperCol={{ span: 12}}
                autoComplete="off"
                disabled={type === 'view'}
            >

            {SYSTEM_SUBSCRIBE_APPROVAL_DETAIL_LIST?.map((x)=>{
                if(x.dataType === 'checkbox'){
                    return(
                        <Form.Item<FieldType>
                        className="mb-btnbase"
                        label={x.title}
                        name="partitions"
                        rules={[{ required: true, message: '必填项' }]}
                    >
                        <Checkbox.Group disabled={type === 'view'}
                            options={data?.partition?.map((x:{id:string,name:string})=>({label:x.name,value:x.id}))}
                        />
                    </Form.Item>
                    )
                }
                return (
                    <Row key={x.key} className="leading-[32px] mb-btnbase">
                        <Col className="text-left" span={8}>{x.title}：</Col>
                        {/* <Col >{showData(x)}</Col> */}
                        <Col >{x.nested ? data?.[x.key]?.[x.nested] : ( (data as {[k:string]:unknown})?.[x.key] || '-')}</Col>
                    </Row>)
                })}

                <Form.Item<FieldType>
                    label="申请原因"
                    name="reason"
                >
                    <Input.TextArea className="w-INPUT_NORMAL" disabled={type === 'view'} placeholder="请输入"  />
                </Form.Item>
                <Form.Item<FieldType>
                    label="审核意见"
                    name="opinion"
                >
                    <Input.TextArea className="w-INPUT_NORMAL" placeholder="请输入" disabled={true} />
                </Form.Item>
            </Form>
            </WithPermission>
        </div>
    )
})