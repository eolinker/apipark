import {forwardRef, useEffect, useImperativeHandle, useReducer, useState} from "react";
import {App, Form, Input, Table} from "antd";
import {useFetch} from "@common/hooks/http.ts";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import { PartitionClusterFieldType, ClusterConfigHandle, ClusterConfigProps, PartitionClusterNodeTableListItem, PartitionClusterNodeModalTableListItem } from "../../const/partitions/types.ts";
import { useParams } from "react-router-dom";
import { RouterParams } from "@core/components/aoplatform/RenderRoutes.tsx";
import { NODE_MODAL_COLUMNS } from "../../const/partitions/const.tsx";
import { v4 as uuidv4} from 'uuid'
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";

const formReducer = (state: PartitionClusterFieldType, action: { type: 'UPDATE_FIELD'; value: PartitionClusterFieldType }) => {
    switch (action.type) {
        case 'UPDATE_FIELD':
            return { ...action.value };
        default:
            return state;
    }
};

export const PartitionClusterConfig = forwardRef<ClusterConfigHandle, ClusterConfigProps>((props, ref)=>{
    const {clusterId,mode,initFormValue} = props
    const { message } = App.useApp()
    const [form] = Form.useForm();
    const [formData, dispatch] = useReducer(formReducer, {});
    const [dataSource,setDataSource] = useState<PartitionClusterNodeModalTableListItem[]>([])
    const {fetchData} = useFetch()
    const {partitionId} = useParams<RouterParams>()

    const save:()=>Promise<boolean | string> =  ()=>{
        // eslint-disable-next-line no-async-promise-executor
        return new Promise( async (resolve, reject)=>{
            let body
            if(mode === 'edit') {
                body = await form.validateFields()
            }else{
                body = {
                    name:formData.name,
                    description:formData.description,
                    managerAddress: formData.address
                }
            }

            fetchData<BasicResponse<null>>('partition/cluster',{
                    method:mode === 'edit' ? 'PUT':'POST' ,
                    eoBody:(body), 
                    eoTransformKeys:['managerAddress'], 
                    eoParams:(mode === 'edit' ? {id:clusterId}:{partition:partitionId})
                    }
                ).then(response=>{
                const {code,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功！')
                    resolve(true)
                }else{
                    message.error(msg || '操作失败')
                    reject(msg || '操作失败')
                }
            }).catch((errorInfo)=> reject(errorInfo))
        })
    }

    const check:()=>Promise<boolean> =  ()=>{
        setDataSource([])
        return new Promise((resolve, reject)=>{
            form.validateFields().then((value)=>{
                dispatch({type:'UPDATE_FIELD',value})
                fetchData<BasicResponse<{ nodes: PartitionClusterNodeTableListItem[] }>>('partition/cluster/check',{method:'POST',eoBody:({address: value.address}),eoTransformKeys:['manager_address','service_address','peer_address']}).then(response=>{
                    const {code,data,msg} = response
                    if(code === STATUS_CODE.SUCCESS){
                        setDataSource(data.nodes)
                        resolve(true)
                    }else{
                        form.setFields([{
                            name:'address',errors:[msg]
                        }])
                        message.error(msg || '操作失败')
                        reject(msg || '操作失败')
                    }
                }).catch((errorInfo)=> reject(errorInfo))
            }).catch((errorInfo)=> reject(errorInfo))
        })
    }

    useImperativeHandle(ref, ()=>({
            save, check
        })
    )

    useEffect(() => {
        if(mode === 'edit' && initFormValue && Object.keys(initFormValue).length > 0 ){
            form.setFieldsValue(initFormValue)
        }else{
            form.setFieldsValue({id:uuidv4()})

        }
    }, [partitionId]);

    return  (<>
        {  mode !== 'result' ?
            <WithPermission access={mode === 'edit' ? 'system.partition.cluster.edit': 'system.partition.cluster.add'}>
                <Form
                    layout='vertical'
                    labelAlign='left'
                    scrollToFirstError
                    form={form}
                    className="mx-auto "
                    name="partitionClusterConfig"
                    // labelCol={{ span: 8 }}
                    // wrapperCol={{ span: 17}}
                    autoComplete="off"
                >
                    <Form.Item<PartitionClusterFieldType>
                        label="集群名称"
                        name="name"
                        rules={[{ required: true, message: '必填项',whitespace:true  }]}
                    >
                        <Input className="w-INPUT_NORMAL" placeholder="请输入集群名称"/>
                    </Form.Item>

                    <Form.Item
                        label="集群 ID"
                        name="id"
                        rules={[{ required: true, message: '必填项',whitespace:true  }]}
                    >
                        <Input className="w-INPUT_NORMAL" disabled={mode === 'edit'} placeholder="请输入集群 ID"/>
                    </Form.Item>

                    <Form.Item
                        label="描述"
                        name="description"
                    >
                        <Input.TextArea className="w-INPUT_NORMAL" placeholder="请输入描述"/>
                    </Form.Item>
                    { mode !== 'edit' &&
                        <Form.Item
                            label="集群地址（网关节点）"
                            name="address"
                            rules={[{ required: true, message: '必填项' ,whitespace:true }]}
                        >
                            <Input className="w-INPUT_NORMAL"  style={{width: '100%'}}  placeholder="请输入集群地址（网关节点）"/>
                        </Form.Item> }
                </Form>
            </WithPermission>
            :
            <div className="mt-mbase">
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
    </>)
})
