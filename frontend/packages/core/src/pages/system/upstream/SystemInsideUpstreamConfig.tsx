import  { forwardRef, useEffect, useImperativeHandle, useState} from "react";
import {App, Button, Divider, Form, FormItemProps, Input, InputNumber, Radio,Select} from "antd";
import {useParams} from "react-router-dom";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import EditableTableWithModal from "@common/components/aoplatform/EditableTableWithModal.tsx";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import { typeOptions, schemeOptions,balanceOptions,passHostOptions, PROXY_HEADER_CONFIG, SYSTEM_UPSTREAM_GLOBAL_CONFIG_TABLE_COLUMNS } from "../../../const/system/const.tsx";
import { ServiceUpstreamFieldType, NodeItem, ProxyHeaderItem, GlobalNodeItem, SystemInsideUpstreamConfigHandle, SystemInsideUpstreamConfigProps } from "../../../const/system/type.ts";
import EditableTable from "@common/components/aoplatform/EditableTable.tsx";
import { v4 as uuidv4} from 'uuid'
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";


const SystemInsideUpstreamConfig = forwardRef<SystemInsideUpstreamConfigHandle,SystemInsideUpstreamConfigProps>((props,ref) => {
    const {upstreamNameForm,partitionId,setLoading} = props
    const { message } = App.useApp()
    const { serviceId,systemId } = useParams<RouterParams>();
    const [form] = Form.useForm();
    const {fetchData} = useFetch()
    const [, forceUpdate] = useState<unknown>(null);
    const [showHost, setShowHost] = useState(false)

    useImperativeHandle(ref, () => ({
        save:onFinish
    }));

    const saveUpstream = (value:ServiceUpstreamFieldType)=>{
        value.proxyHeaders ||= []
        if(value.globalConfig?.nodes){
            value.globalConfig.nodes = value.globalConfig.nodes.filter((x:GlobalNodeItem)=>x.address)?.map((x:GlobalNodeItem)=>({address:x.address, weight:x.weight ?? 100}))
        }
        if(value.tmpConfig?.nodes){
            const clusterConfig:Map<string,{nodes:GlobalNodeItem[]}> = new Map()
            value.tmpConfig.nodes.forEach((x:NodeItem)=>{
                if(!x.address) return
                if (!clusterConfig.has(x.cluster)) {
                    clusterConfig.set(x.cluster, { nodes: [] });
                  }
                  clusterConfig.get(x.cluster)!.nodes.push({ address: x.address, weight: x.weight ?? 100 });
            })
            value.config = Array.from(clusterConfig.entries())?.map(([cluster, nodes]) => ({ cluster, nodes: nodes.nodes }));
            delete value.tmpConfig;
        }

        return fetchData<BasicResponse<null>>(
            'project/upstream',
            {
                method:'PUT',
                eoBody:({
                    ...value,
                    name:upstreamNameForm.getFieldValue('name'),
                    limitPeerSecond:Number(value.limitPeerSecond)||0,
                    retry:Number(value.retry)||0,
                    timeout:Number(value.timeout)||0}),
                    eoParams:{project:systemId},
                    eoTransformKeys:['limitPeerSecond','proxyHeaders','optType','globalConfig','passHost','upstreamHost']
                }).then(response=>{
                    const {code,msg} = response
                    if(code === STATUS_CODE.SUCCESS){
                        message.success(msg || '操作成功！')
                        return Promise.resolve(true)
                    }else{
                        message.error(msg || '操作失败')
                        return Promise.reject(msg || '操作失败')
                    }
                }).catch((errorInfo)=> {return Promise.reject(errorInfo)})
    }

    const onFinish:()=>Promise<boolean|string> = async () => {
    try {
        form.setFieldValue('partition',partitionId)
        await upstreamNameForm.validateFields();
        const value: ServiceUpstreamFieldType = await form.validateFields();
        return saveUpstream(value);
    } catch (errorInfo) {
        return Promise.reject(errorInfo);
    }
};
    // 获取表单默认值
    const getUpstreamInfo = () => {
            setLoading(true)
            fetchData<BasicResponse<{ upstream: ServiceUpstreamFieldType }>>('project/upstream',{method:'GET',eoParams:{project:systemId},eoTransformKeys:['limit_peer_second','proxy_headers','opt_type','global_config','pass_host','upstream_host']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setShowHost(data.upstream.passHost == 'rewrite')
                if(data.upstream.globalConfig?.nodes){
                    data.upstream.globalConfig.nodes.push({_id:uuidv4()})
                }
                const addEmptyMergedData:NodeItem[] = []
                //   const mergeGlobalData =  data.upstream?.config?.default?.length > 0 ? generateResultList(mergedData, data.upstream.config.default):mergedData
                setTimeout(()=>{form.setFieldsValue({...data.upstream,partition:partitionId,balance: data.upstream.balance || 'round-robin',limitPeerSecond:data.upstream.limitPeerSecond || 5,timeout:data.upstream.timeout || 10000,proxyHeaders: data.upstream.proxyHeaders || [],retry:data.upstream.retry || 3,config:data.config, tmpConfig:{nodes:addEmptyMergedData}})},0)
                upstreamNameForm.setFieldValue('name',data.upstream.name)
            }else{
                message.error(msg || '操作失败')
            }
        }).finally(()=>{
            setLoading(false)
            forceUpdate({})
        })
    };

   // 自定义校验规则
const globalConfigNodesRule: FormItemProps['rules'] = [
    {
      validator: (_, value) => {
        if (!value || !Array.isArray(value)) {
            return Promise.resolve();
        }
        const filteredValue = value.filter((item) => item.address && item.weight!== '' && item.weight!== null);
        if (filteredValue.length > 0) {
          return Promise.resolve();
        } else {
          return Promise.reject(new Error('必填项'));
        }
      },
    },
  ];

    useEffect(() => {
        form.setFieldsValue({partition:partitionId,driver:'static', scheme:'HTTP', balance:'round-robin',passHost:'pass', limitPeerSecond:5,timeout:10000,retry:3,globalConfig:{nodes:undefined}}); // 清空 initialValues
            getUpstreamInfo();
    }, [partitionId]);

    return (
            <WithPermission access={'project.mySystem.upstream.edit'}>
                    <Form
                        layout='vertical'
                        labelAlign='left'
                        name="systemInsideUpstreamConfig"
                        scrollToFirstError
                        form={form}
                        className="mx-auto   h-[calc(100%-98px)] py-[20px] overflow-y-auto"
                        // labelCol={{ span: 4 }}
                        // wrapperCol={{ span: 20}}
                        onFinish={onFinish}
                        autoComplete="off"
                        >

                        <Form.Item<ServiceUpstreamFieldType>
                            label="环境 ID"
                            name="partition"
                            hidden
                            rules={[{ required: true, message: '必填项',whitespace:true  }]}
                        >
                            <Input className="w-INPUT_NORMAL" placeholder="请输入上游名称"/>
                        </Form.Item>

                        <Form.Item<ServiceUpstreamFieldType>
                            label="上游类型"
                            name="driver"
                            rules={[{ required: true, message: '必填项' }]}
                        >
                            <Radio.Group options={typeOptions} />
                        </Form.Item>


                        <Form.Item<ServiceUpstreamFieldType>
                            label="服务地址"
                            name={["globalConfig","nodes"]}
                            tooltip="后端默认使用的IP地址"
                            rules={[{ required: true, message: '必填项' },
                            ...globalConfigNodesRule]}
                        >
                            <EditableTable<GlobalNodeItem & {_id:string}>
                                className="w-[528px]"
                                configFields={SYSTEM_UPSTREAM_GLOBAL_CONFIG_TABLE_COLUMNS}
                            />
                        </Form.Item>

                        <Form.Item<ServiceUpstreamFieldType>
                            label="请求协议"
                            name="scheme"
                            rules={[{ required: true, message: '必填项' }]}
                        >
                        <Select className="w-INPUT_NORMAL" placeholder="请选择" options={schemeOptions}>
                        </Select>
                        </Form.Item>

                        <Form.Item<ServiceUpstreamFieldType>
                            label="负载均衡"
                            name="balance"
                            rules={[{ required: true, message: '必填项' }]}
                        >
                            <Radio.Group className="flex flex-col gap-[8px] mt-[5px]" options={balanceOptions} />
                        </Form.Item>

                        <Form.Item<ServiceUpstreamFieldType>
                            label="转发 Host"
                            name="passHost"
                            rules={[{ required: true, message: '必填项' }]}
                        >
                            <Select className="w-INPUT_NORMAL" placeholder="请选择" options={passHostOptions} onChange={(val)=>setShowHost(val === 'rewrite')}>
                            </Select>
                        </Form.Item>

                        {showHost && <Form.Item<ServiceUpstreamFieldType>
                            label="重写域名"
                            name="upstreamHost"
                            rules={[{ required: true, message: '必填项',whitespace:true  }]}
                        >
                            <Input className="w-INPUT_NORMAL" placeholder="请输入上游名称"/>
                        </Form.Item>
                    }

                        <Divider />

                        

                        <Form.Item<ServiceUpstreamFieldType>
                            label="超时时间"
                            name="timeout"
                            rules={[{ required: true, message: '必填项' }]}
                        >
                            <InputNumber className="w-INPUT_NORMAL" min={1} addonAfter={<span className="whitespace-nowrap">ms</span> }/> 
                        </Form.Item>

                        <Form.Item<ServiceUpstreamFieldType>
                            label="超时重试次数"
                            name="retry"
                            rules={[{ required: true, message: '必填项' }]}
                        >
                            <InputNumber className="w-INPUT_NORMAL" min={1} addonAfter={<span>次</span>} /> 
                        </Form.Item>

                        <Form.Item<ServiceUpstreamFieldType>
                            label="调用频率限制"
                            name="limitPeerSecond"
                            rules={[{ required: true, message: '必填项' }]}
                        >
                            <InputNumber className="w-INPUT_NORMAL"  min={1} addonAfter={<span className="whitespace-nowrap">次/秒</span> } />
                        </Form.Item>

                        <Form.Item<ServiceUpstreamFieldType>
                            label="转发上游请求头"
                            name="proxyHeaders"
                            className="mb-0"
                        >
                            <EditableTableWithModal<ProxyHeaderItem & {_id:string}>
                                className="w-[528px]"
                                configFields={PROXY_HEADER_CONFIG}
                            />
                        </Form.Item>

                        <Form.Item 
                        // wrapperCol={{ offset: 7, span: 16 }}
                        >
                            <WithPermission access='project.mySystem.upstream.edit'><Button type="primary" htmlType="submit">
                                保存
                            </Button></WithPermission>
                        </Form.Item>
                    </Form>
            </WithPermission>
    )
})

export default SystemInsideUpstreamConfig;