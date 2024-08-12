import { App, Form, Row, Col, Checkbox, Select, Input } from "antd";
import { forwardRef, useEffect, useImperativeHandle, useMemo } from "react";
import WithPermission from "@common/components/aoplatform/WithPermission";
import { BasicResponse, STATUS_CODE } from "@common/const/const";
import { ApplyServiceHandle, ApplyServiceProps, ServiceHubApplyModalFieldType } from "../../const/serviceHub/type";
import { EntityItem } from "@common/const/type";
import { useFetch } from "@common/hooks/http";

export const ApplyServiceModal = forwardRef<ApplyServiceHandle,ApplyServiceProps>((props,ref)=>{
    const { message } = App.useApp()
    const {entity,mySystemOptionList,reApply} = props
    const [form] = Form.useForm();
    const {fetchData} = useFetch()
    // const [avaliablePartitions, setAvaPartitions] = useState<Array<string>>([])

    useEffect(() => {
        form.setFieldsValue(reApply ? {applications:entity?.project.id}:{})
    }, []);

    const apply: ()=>Promise<boolean | string> =  ()=>{
        return new Promise((resolve, reject)=>{
            form.validateFields().then((value)=>{
                fetchData<BasicResponse<null>>('catalogue/service/subscribe',{method:'POST',eoBody:({...value,service:entity.id,})}).then(response=>{
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

    // const onProjectsChange = (_: ParamsType & {
    //     pageSize?: number | undefined;
    //     current?: number | undefined;
    //     keyword?: string | undefined;
    // },projectInfo: DefaultOptionType&{partition:string[]}[])=>{
    //     if (Array.isArray(projectInfo) && projectInfo.every((x) => Array.isArray(x.partition))) {
    //         const tmpPartitinList = projectInfo?.map((x) => x.partition)
    //         setAvaPartitions(tmpPartitinList&&tmpPartitinList.length > 0 ? tmpPartitinList.reduce((acc, curr) => acc.filter((p) => curr.includes(p))):[]);
    //       } else {
    //         setAvaPartitions([]);
    //       }
    // }

    useImperativeHandle(ref, ()=>({
            apply
        })
    )

    const partitionsList = useMemo(()=>{
        const newList = entity.partition?.map((x:EntityItem)=>({label:x.name, value:x.id})) || []
        if(newList?.length === 1) {
            form.setFieldValue('partitions',[newList[0].value])
        }
        return newList
    },[entity])

    return (<WithPermission access="">
        <Form
            layout='vertical'
            scrollToFirstError
            form={form}
            className=" w-full mt-[20px]"
            name="applyServiceModal"
            // labelCol={{ span: 6 }}
            // wrapperCol={{ span: 18}}
            autoComplete="off"
        >
            <Row className="mb-btnybase h-[32px]" >
                <Col span={6} className="pb-[8px] text-left">服务名称：</Col>
                <Col span={18}>{entity.name}</Col>
            </Row>
            <Row className="h-[32px] mb-btnybase">
                <Col span={6} className="pb-[8px]  text-left">服务 ID：</Col>
                <Col span={18}>{entity.id}</Col>
            </Row>
            <Row className="h-[32px] mb-btnybase">
                <Col span={6} className=" pb-[8px]  text-left">所属系统：</Col>
                <Col span={18}>{entity.project.name}</Col>
            </Row>
            <Form.Item<ServiceHubApplyModalFieldType>
                label="申请的环境"
                name="partitions"
                rules={[{ required: true, message: '必填项' }]}
            >
                <Checkbox.Group options={partitionsList}/>
            </Form.Item>
            <Form.Item
                label="应用"
                name="applications"
                rules={[{ required: true, message: '必填项' }]}
            >
                <Select className="w-INPUT_NORMAL" disabled={reApply} placeholder="搜索或选择应用" mode="multiple" options={mySystemOptionList.filter((x)=>x.value !== entity.project.id)}/>
            </Form.Item>

            <Form.Item
                label="申请理由"
                name="reason"
            >
                <Input.TextArea className="w-INPUT_NORMAL" placeholder=""/>
            </Form.Item>

        </Form>
    </WithPermission>)
})