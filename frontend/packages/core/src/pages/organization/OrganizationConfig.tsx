import  { forwardRef, useEffect, useImperativeHandle, useState} from "react";
import {App, Checkbox, CheckboxOptionType, Form, Input, Select} from "antd";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {DefaultOptionType} from "antd/es/cascader";
import { v4 as uuidv4 } from 'uuid'
import {PartitionItem, MemberItem} from "@common/const/type.ts";
import { OrganizationFieldType } from "../../const/organization/type.ts";
import { validateUrlSlash } from "@common/utils/validate.ts";
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";

export type OrganizationConfigHandle = {
    save:()=>Promise<string|boolean>
}

type OrganizationConfigProps = {
    entity?:OrganizationFieldType
}

const OrganizationConfig = forwardRef<OrganizationConfigHandle,OrganizationConfigProps>((props,ref)=>{
    const { message } = App.useApp()
    const { entity } = props
    const [type, setType] = useState<'add'|'edit'>(entity === undefined? 'add' : 'edit')
    const [form] = Form.useForm();
    const {fetchData} = useFetch()
    const [managerOption, setManagerOption] = useState<DefaultOptionType[]>([])
    const [partitionOption,setPartitionOption] = useState<CheckboxOptionType[]>([])
    // 获取表单默认值
    // const getOrgInfo = () => {
    //     setLoading(true)
    //     fetchData<BasicResponse<{ organization: OrganizationFieldType }>>('manager/organization',{method:'GET',eoParams:{id:orgId},eoTransformKeys:['create_time','master_id','update_time']}).then(response=>{
    //         const {code,data,msg} = response
    //         if(code === STATUS_CODE.SUCCESS){
    //            setTimeout(()=>{form.setFieldsValue({...data.organization,partitions:data.organization.partitions?.map((x:PartitionItem)=>(x.id)),master:data.organization.master.id})},0)
                
    //         }else{
    //             message.error(msg || '操作失败')
    //         }
    //     }).finally(()=>{setLoading(false)})
    // };
    
    useImperativeHandle(ref, () => ({
        save:onFinish
    }));

    const onFinish = () => {
        return form.validateFields().then((value)=>{
            return fetchData<BasicResponse<null>>('manager/organization',{method:type === 'add'? 'POST' : 'PUT',eoBody:({...value,prefix:value.prefix?.trim()}),...type === 'edit' ?{eoParams:{id:entity!.id}}:{}}).then(response=>{
                const {code,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功！')
                    return Promise.resolve(true)
                }else{
                    message.error(msg || '操作失败')
                    return Promise.reject(msg || '操作失败')
                }
            }).catch(errInfo=>Promise.reject(errInfo))
        }).catch(errInfo=>Promise.reject(errInfo))
    };

    const getManagerList = ()=>{
        setManagerOption([])
        fetchData<BasicResponse<{ members: MemberItem[] }>>('simple/member',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setManagerOption(data.members?.map((x:MemberItem)=>{return {
                    label:x.name, value:x.id
                }}))
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const getPartitionList = ()=>{
        setPartitionOption([])
        fetchData<BasicResponse<{partitions:PartitionItem[]}>>('simple/partitions',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setPartitionOption(data.partitions?.map((x:PartitionItem)=>{return {
                    label:x.name,value:x.id
                }}))
                if(type === 'add' && data.partitions?.length === 1){
                    form.setFieldValue('partitions',[data.partitions[0].id])
                }
            }else{
                message.error(msg || '操作失败')
            }
        })
    }
    
    useEffect(() => {
        getManagerList()
        getPartitionList()

        if (entity !== undefined) {
            setType('edit')
            form.setFieldsValue(entity)
        } else {
            setType('add')
            form.setFieldValue('id',uuidv4()); // 清空 initialValues
        }
        return (form.setFieldsValue({}))
    }, [entity]);

    return (
        <>
                <WithPermission access={type === 'add' ? 'system.organization.self.add':'system.organization.self.edit'}>
                <Form
                    layout='vertical'
                    labelAlign='left'
                    scrollToFirstError
                    form={form}
                    className="mx-auto "
                    name="OrganizationConfig"
                    // labelCol={{ offset:1,span: 5 }}
                    // wrapperCol={{ span: 18}}
                    onFinish={onFinish}
                    autoComplete="off"
                >
                    <Form.Item<OrganizationFieldType>
                        label="组织名称"
                        name="name"
                        rules={[{ required: true, message: '必填项' ,whitespace:true }]}
                    >
                        <Input className="w-INPUT_NORMAL" placeholder="请输入组织名称"/>
                    </Form.Item>

                    <Form.Item<OrganizationFieldType>
                        label="组织ID"
                        name="id"
                        extra="组织ID（org_id）可用于检索组织，一旦保存无法修改。"
                        rules={[{ required: true, message: '必填项' ,whitespace:true }]}
                    >
                        <Input className="w-INPUT_NORMAL" disabled={type === 'edit'} placeholder="请输入组织ID"/>
                    </Form.Item>

                    <Form.Item<OrganizationFieldType>
                        label="描述"
                        name="description"
                    >
                        <Input.TextArea className="w-INPUT_NORMAL" placeholder="请输入描述"/>
                    </Form.Item>

                    <Form.Item<OrganizationFieldType>
                        label="负责人"
                        name="master"
                        extra="负责人对组织内的团队、服务、成员有管理权限"
                        rules={[{ required: true, message: '必填项' }]}
                    >
                        <Select className="w-INPUT_NORMAL" placeholder="请选择负责人" options={managerOption}>
                        </Select>
                    </Form.Item>

                    <Form.Item<OrganizationFieldType>
                        label="组织请求前缀"
                        name="prefix"
                        extra="该请求前缀将会拼接到API请求路径中，格式为：{协议}{主机地址}{组织前缀}{分区前缀}{系统前缀}{API请求路径}"
                        rules={[
                        {
                          validator: validateUrlSlash,
                        }]}
                    >
                        <Input prefix={type === 'edit' ? '' : '/'} className="w-INPUT_NORMAL" disabled={type === 'edit'} placeholder="请输入组织请求前缀" />
                    </Form.Item>

                    <Form.Item<OrganizationFieldType>
                        label="分区权限"
                        name="partitions"
                        rules={[{ required: true, message: '必填项' }]}
                    >
                        <Checkbox.Group className="flex flex-col gap-[8px] mt-[5px]" options={partitionOption} />
                    </Form.Item>

                </Form>
                </WithPermission>
        </>
    )
})
export default OrganizationConfig