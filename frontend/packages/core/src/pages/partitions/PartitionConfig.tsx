import  {forwardRef, useEffect, useImperativeHandle, useReducer, useState} from "react";
import {App, Button, Divider, Form, Input, Row, Table} from "antd";
import {Link, useNavigate, useParams} from "react-router-dom";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import { PartitionClusterFieldType, PartitionClusterNodeModalTableListItem, PartitionClusterNodeTableListItem, PartitionConfigFieldType } from "../../const/partitions/types.ts";
import { v4 as uuidv4} from 'uuid'
import { validateUrlSlash } from "@common/utils/validate.ts";
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";
import { useBreadcrumb } from "@common/contexts/BreadcrumbContext.tsx";
import { usePartitionContext } from "../../contexts/PartitionContext.tsx";
import { NODE_MODAL_COLUMNS } from "../../const/partitions/const.tsx";

export type PartitionConfigProps = {
    mode:'config' | 'retry' | 'result' | 'edit',
}
export interface PartitionConfigHandle{
    save:()=>Promise<boolean|string>
    check:()=>Promise<boolean>
}

const formReducer = (state: PartitionClusterFieldType, action: { type: 'UPDATE_FIELD'; value: PartitionClusterFieldType }) => {
    switch (action.type) {
        case 'UPDATE_FIELD':
            return { ...action.value };
        default:
            return state;
    }
};


const PartitionConfig = forwardRef<PartitionConfigHandle,PartitionConfigProps>((props,ref)=> {
    const {mode} = props
    const { message,modal } = App.useApp()
    const { partitionId } = useParams<RouterParams>();
    const [ onEdit, setOnEdit] = useState<boolean>(!!partitionId)
    const [ form ] = Form.useForm();
    const { fetchData} = useFetch()
    const [formData, dispatch] = useReducer(formReducer, {});
    const navigate = partitionId === undefined ? ()=>{} : useNavigate();
    const { setBreadcrumb } = partitionId === undefined ?{setBreadcrumb:()=>{}} : useBreadcrumb()
    const { setPartitionInfo} = partitionId === undefined ?{setPartitionInfo:()=>{}}: usePartitionContext()
    const [dataSource,setDataSource] = useState<PartitionClusterNodeModalTableListItem[]>([])
    const [canDelete, setCanDelete] = useState<boolean>(false)
    useImperativeHandle(ref, ()=>({
        save:onFinish,
        check
    }))

    const onFinish =async () => {
        // eslint-disable-next-line no-async-promise-executor
            let body
            if(mode === 'edit') {
                body = await form.validateFields()
            }else{
                body = {...formData
                }
            }

            return fetchData<BasicResponse<{partition:PartitionConfigFieldType}>>('partition',{method:partitionId === undefined? 'POST' : 'PUT',eoTransformKeys:['managerAddress'],eoBody:({...body}),...(partitionId !== undefined?{eoParams:{id:partitionId}}:{})}).then(response=>{
                const {code,data,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功！')
                    setPartitionInfo(data.partition)
                    return Promise.resolve(true)
                }else{
                    message.error(msg || '操作失败')
                    return Promise.reject(msg || '操作失败')
                }
            }).catch((errInfo)=>Promise.reject(errInfo))
    };

    
    const check:()=>Promise<boolean> =  ()=>{
        setDataSource([])
        return form.validateFields().then((value)=>{
                dispatch({type:'UPDATE_FIELD',value})
                return fetchData<BasicResponse<{ nodes: PartitionClusterNodeTableListItem[] }>>('partition/cluster/check',{method:'POST',eoBody:({address: value.managerAddress}),eoTransformKeys:['manager_address','service_address','peer_address']}).then(response=>{
                    const {code,data,msg} = response
                    if(code === STATUS_CODE.SUCCESS){
                        setDataSource(data.nodes)
                        return Promise.resolve(true)
                    }else{
                        form.setFields([{
                            name:'address',errors:[msg]
                        }])
                        message.error(msg || '操作失败')
                        return Promise.reject(msg || '操作失败')
                    }
                }).catch((errorInfo)=> Promise.reject(errorInfo))
            }).catch((err)=> {form.scrollToField(err.errorFields[0].name[0]); return Promise.reject(err)})
    }

    // 获取表单默认值
    const getPartitionInfo = () => {
        fetchData<BasicResponse<{ partition: PartitionConfigFieldType }>>('partition',{method:'GET',eoParams:{id:partitionId}, eoTransformKeys:['can_delete']},).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setTimeout(()=>{form.setFieldsValue(data.partition)},0)
                setCanDelete(data.partition.canDelete)
            }else{
                message.error(msg || '操作失败')
            }
        })
    };

    const deletePartitionModal = async ()=>{
        modal.confirm({
            title:'删除',
            content:'该数据删除后将无法找回，请确认是否删除？',
            onOk:()=> {
                return deletePartition()
            },
            width:600,
            okText:'确认',
            okButtonProps:{
                danger:true
            },
            cancelText:'取消',
            closable:true,
            icon:<></>
        })
    }

    const deletePartition = ()=>{
            fetchData<BasicResponse<null>>('partition',{method:'DELETE',eoParams:{id:partitionId}}).then(response=>{
                const {code,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功')
                    navigate('/partition/list')
                }else{
                    message.error(msg || '操作失败')
                }
            })
    }

    useEffect(() => {
        if (partitionId !== undefined) {
            setBreadcrumb([
                {title:<Link to="/partition/list">部署管理</Link>},
                {title:'环境设置'}
            ])
            setOnEdit(true);
            getPartitionInfo();
        } else {
            setOnEdit(false);
            form.setFieldsValue({id:uuidv4()}); // 清空 initialValues
        }

        return (form.setFieldsValue({}))
    }, [partitionId]);

    return (
        <>
            <div className="h-full min-w-[560px]">
                { mode !== 'result' ? 
                <WithPermission access={onEdit ? 'system.partition.self.edit':'system.partition.self.add'} >
                    <Form
                        layout='vertical'
                        labelAlign='left'
                        scrollToFirstError
                        form={form}
                        className="mx-auto   flex flex-col justify-between h-full"
                        name="partitionConfig"
                        // labelCol={{ offset:1,span: 4 }}
                        // wrapperCol={{ span: 19}}
                        onFinish={onFinish}
                        autoComplete="off"
                    >
                        <div>
                            <Form.Item<PartitionConfigFieldType>
                                label="环境名称"
                                name="name"
                                rules={[{ required: true, message: '必填项' }]}
                            >
                                <Input className="w-INPUT_NORMAL" placeholder="请输入环境名称"/>
                            </Form.Item>

                            <Form.Item<PartitionConfigFieldType>
                                label="环境 ID"
                                name="id"
                                extra="环境 ID（partition_id）可用于检索环境，一旦保存无法修改。"
                                rules={[{ required: true, message: '必填项' }]}
                            >
                                <Input className="w-INPUT_NORMAL" disabled={onEdit} placeholder="请输入环境标识"/>
                            </Form.Item>

                            <Form.Item<PartitionConfigFieldType>
                                label="环境请求前缀"
                                name="prefix"
                                extra="该请求前缀将会拼接到API请求路径中，格式为：{协议}{主机地址}{组织前缀}{环境前缀}{服务前缀}{API请求路径}"
                                rules={[
                                {
                                validator: validateUrlSlash,
                                }]}
                            >
                                <Input  prefix={onEdit ? '' : "/"} className="w-INPUT_NORMAL" disabled={onEdit} placeholder="请输入环境请求前缀"/>
                            </Form.Item>

                    {!onEdit && 
                            <Form.Item<PartitionConfigFieldType>
                                label="集群地址"
                                name="managerAddress"
                                rules={[{ required: true, message: '必填项' }]}
                            >
                            {/* <Space> */}
                                <Input className="w-INPUT_NORMAL" placeholder="请输入"/>
                                {/* <Button type='primary' htmlType='submit'>测试</Button> */}
                            {/* </Space> */}
                            </Form.Item>}

                            <Form.Item<PartitionConfigFieldType>
                                label="描述"
                                name="description"
                            >
                                <Input.TextArea className="w-INPUT_NORMAL" placeholder="请输入描述"/>
                            </Form.Item>

                            { onEdit && <Row className="mb-[10px]" 
                            // wrapperCol={{ offset: 5, span: 16 }}
                            >
                                <WithPermission access={onEdit ? 'system.partition.self.edit':'system.partition.self.add'}>
                                    <Button type="primary" htmlType="submit">
                                        保存
                                    </Button>
                                </WithPermission>
                            </Row>
                            }
                        </div>
                        {onEdit && 
                        <div>
                            <Divider />
                            
                            <p className="text-center">删除环境：删除操作不可恢复，请谨慎操作！</p>
                            <div className="text-center">
                                <WithPermission access="system.partition.self.delete" disabled={!canDelete}  tooltip={canDelete ? '':'环境已被使用，不可删除'}><Button className="m-auto mt-[16px] mb-[20px]" type="default" onClick={deletePartitionModal}>删除</Button></WithPermission>
                            </div>
                        </div>
                        }
                    </Form>
                </WithPermission>
                :
                <div className="">
                    <p className="mb-btnybase">检查通过。该集群有一个节点</p>
                    <Table
                        bordered={true}
                        columns={[...NODE_MODAL_COLUMNS]}
                        size="small"
                        rowKey="id"
                        dataSource={dataSource}
                        pagination={false}
                    />
                </div>
                }
            </div>
        </>
    )
})
export default PartitionConfig