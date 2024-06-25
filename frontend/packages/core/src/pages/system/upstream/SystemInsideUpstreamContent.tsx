/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-07 17:36:28
 * @FilePath: \frontend\packages\core\src\pages\system\upstream\SystemInsideUpstreamContent.tsx
 */
import { App, Badge, Button, Col, Divider, Form, Input, InputNumber, Radio, Row, Select, Spin, Switch, Tabs} from "antd";
import  {forwardRef, useEffect, useImperativeHandle, useRef, useState} from "react";
import styles from './SystemInsideUpstream.module.css'
import { LoadingOutlined } from "@ant-design/icons";
import { GlobalNodeItem, ProxyHeaderItem, ServiceUpstreamFieldType, SystemInsideUpstreamConfigHandle, SystemInsideUpstreamContentHandle } from "../../../const/system/type.ts";
import TabPane from "antd/es/tabs/TabPane";
import { FormItemProps } from "antd/es/form/index";
import EditableTable from "@common/components/aoplatform/EditableTable.tsx";
import EditableTableWithModal from "@common/components/aoplatform/EditableTableWithModal.tsx";
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";
import { typeOptions, SYSTEM_UPSTREAM_GLOBAL_CONFIG_TABLE_COLUMNS, schemeOptions, balanceOptions, passHostOptions, PROXY_HEADER_CONFIG } from "../../../const/system/const.tsx";
import { Link, useParams } from "react-router-dom";
import { RouterParams } from "@core/components/aoplatform/RenderRoutes.tsx";
import { BasicResponse, STATUS_CODE } from "@common/const/const.ts";
import { useFetch } from "@common/hooks/http.ts";
import { v4 as uuidv4} from 'uuid'
import { cloneDeep } from "lodash-es";
import { EntityItem } from "@common/const/type.ts";
import { useBreadcrumb } from "@common/contexts/BreadcrumbContext.tsx";

const SystemInsideUpstreamContent= forwardRef<SystemInsideUpstreamContentHandle>((props,ref) => {
    const formRef = useRef<SystemInsideUpstreamConfigHandle>(null)
    const [loading, setLoading] = useState<boolean>(false)
    const [formStatus, setFormStatus] = useState<{ [key: string]: boolean }>({});
    const [formEmptyRequired, setFormEmptyRequired] = useState<{ [key: string]: number }>({});
    const formRefs = useRef<{ [key: string]: any }>({});
    const { message } = App.useApp()
    const { serviceId,systemId } = useParams<RouterParams>();
    const [form] = Form.useForm();
    const {fetchData} = useFetch()
    const [, forceUpdate] = useState<unknown>(null);
    const [formShowHost, setFormShowHost] =  useState<{ [key: string]: boolean }>({});
    const [formDataFromApi, setFormDataFromApi] = useState<{ [key: string]: unknown }>({});
    const [partitionList, setPartitionList] = useState<EntityItem[]>([]);
    const [afterSubmit, setAfterSubmit] = useState<boolean>(false)
    const { setBreadcrumb } = useBreadcrumb()

    useImperativeHandle(ref, () => ({
        save:()=>formRef.current?.save()
    }));

    const handleSwitchChange = (id: string, checked: boolean) => {
        setFormStatus({ ...formStatus, [id]: checked });
    };

    const saveUpstream = (value:{[k:string]:ServiceUpstreamFieldType})=>{
        const upstreamValueMap = cloneDeep(value)
        for (const tab in upstreamValueMap) {
            if (formStatus[tab]) {
                // upstreamValueMap[tab].proxyHeaders  []
                if(upstreamValueMap[tab]?.nodes){
                    upstreamValueMap[tab].nodes = upstreamValueMap[tab].nodes.filter((x:GlobalNodeItem)=>x.address)?.map((x:GlobalNodeItem)=>({address:x.address, weight:x.weight ?? 100}))
                }
                upstreamValueMap[tab].limitPeerSecond = Number(upstreamValueMap[tab].limitPeerSecond)||0,
                upstreamValueMap[tab].retry = Number(upstreamValueMap[tab].retry)||0,
                upstreamValueMap[tab].timeout = Number(upstreamValueMap[tab].timeout)||0
            } else {
              delete upstreamValueMap[tab]; // 删除整个 tab
            }
          }

        return fetchData<BasicResponse<null>>(
            'project/upstream',
            {
                method:'PUT',
                eoBody:({...upstreamValueMap}),
                eoParams:{project:systemId},
                eoTransformKeys:['limitPeerSecond','proxyHeaders','optType','passHost','upstreamHost']
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
            // form.setFieldValue('partition',partitionId)
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
                setPartitionList(data.partitions)
                if(!data.partitions || data.partitions.length === 0 ) return
                for ( const p of data.partitions){
                    setFormShowHost(pre => ({...pre, [p.id]: p.passHost == 'rewrite'}))
                    if(data.upstream?.[p.id]){
                        setFormStatus(pre=>({...pre, [p.id]:true}))
                    }
                    data.upstream = {...data.upstream, [p.id] : (data.upstream?.[p.id] ? {...data.upstream[p.id], nodes:data.upstream[p.id].nodes ?? [{_id:uuidv4()}]} : {driver:'static', scheme:'HTTP', balance: 'round-robin',limitPeerSecond:5,timeout: 10000,proxyHeaders: [],retry:3})}
                }
                setFormDataFromApi(data.upstream)
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
        
        setBreadcrumb([
            {
                title: <Link to={`/system/list`}>内部数据服务</Link>
            },
            {
                title: '上游'
            }])

            getUpstreamInfo();
    }, [systemId]);


    const handleSubmit = async () => {
        setAfterSubmit(true)
        const finalData:{[k:string]:ServiceUpstreamFieldType} = {}
        for (const partition of partitionList) {
          const { id } = partition;
          if (formStatus[id]) {
            try {
              const res = await formRefs.current[id].validateFields();
              finalData[id] = res;
            } catch (error) {
                setFormEmptyRequired(prev=> ({...prev, [id]: error?.errorFields?.length }));
                return
            }
          }
        }
        if(Object.keys(finalData).length > 0){
            saveUpstream(finalData)
        }
      };
    
      const renderTabTitle = (name: string, id: string) => {
        if (formStatus[id]) {
          return (
            <Badge count={formEmptyRequired[id] || 0} overflowCount={99} size='small' offset={[8,0]}>
              {name}
            </Badge>
          );
        }
        return name;
      };

    const onFormChange = async (changedValues: any, allValues: any, id: string) => {
        // let count = 0;
        // if (formStatus[id]) {
        // for (const key in allValues) {
        //     if (!allValues[key]) {
        //     count += 1;
        //     }
        // }
        // }
        setFormDataFromApi(prev => ({...prev, [id]: allValues}))
        if(!afterSubmit) return
        try {
            await formRefs.current[id].validateFields();
          } catch (error) {
              setFormEmptyRequired(prev=> ({...prev, [id]: error?.errorFields?.length }));
          }

    };

    

    return (
        <Spin indicator={<LoadingOutlined style={{ fontSize: 24 }} spin />} spinning={loading}>
            <div className={`flex-1 h-full ${styles['upstream-tabs']} min-w-[800px] overflow-auto`} >
                <Tabs 
                    className="h-auto" 
                    size="small" 
                    tabBarStyle={{paddingLeft:'10px',marginTop:'0px',marginBottom:'0px'}} 
                    tabBarGutter={20} 
                    >
                    {partitionList.map(partition => (
                      <TabPane className="pl-btnbase pr-btnrbase" tab={renderTabTitle(partition.name, partition.id)} key={partition.id}>
                            {/* <Row  className="mt-[20px] h-[32px] pb-[8px]">启用：</Row> */}
                            <Row className="mt-[20px] mb-[20px] h-[32px]"><Switch
                            checkedChildren="启用"
                            unCheckedChildren="禁用"
                            value={formStatus[partition.id]}
                            onChange={(checked) => handleSwitchChange(partition.id, checked)}
                            /></Row>
                                    
                        {formStatus[partition.id] && (
                        <WithPermission access={'project.mySystem.upstream.edit'}>
                            <Form
            layout='vertical'
            labelAlign='left'
            name="systemInsideUpstreamContent"
            scrollToFirstError
                                className="mx-auto mb-[20px]  overflow-hidden"
                                // labelCol={{ offset:1, span: 4 }}
                                // wrapperCol={{ span: 19}}
                                ref={(form) => { if(form){
                                     formRefs.current[partition.id] = form;form.setFieldsValue(formDataFromApi[partition.id])}}  }
                                onValuesChange={(changedValues, allValues) => onFormChange(changedValues, allValues, partition.id)}
                                autoComplete="off"
                                >

                                <Form.Item<ServiceUpstreamFieldType>
                                    label="上游类型"
                                    name="driver"
                                    rules={[{ required: true, message: '必填项' }]}
                                >
                                    <Radio.Group options={typeOptions} />
                                </Form.Item>


                                <Form.Item<ServiceUpstreamFieldType>
                                    label="服务地址"
                                    name="nodes"
                                    tooltip="后端默认使用的IP地址"
                                    rules={[{ required: true, message: '必填项' },
                                    ...globalConfigNodesRule]}
                                >
                                    <EditableTable<GlobalNodeItem & {_id:string}>
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
                                    <Select className="w-INPUT_NORMAL" placeholder="请选择" options={passHostOptions} onChange={(val)=>setFormShowHost(prev => ({...prev, [partition.id]:val === 'rewrite'}))}>
                                    </Select>
                                </Form.Item>

                                {formShowHost[partition.id] && <Form.Item<ServiceUpstreamFieldType>
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
                                        configFields={PROXY_HEADER_CONFIG}
                                    />
                                </Form.Item>

                            </Form>
                        </WithPermission>
                        )}
                      </TabPane>
                    ))}
                </Tabs>
                
                    <Row className="px-btnbase mb-[20px]"><Col offset={0} span={24}><WithPermission access='project.mySystem.upstream.edit'><Button type="primary" onClick={handleSubmit} >
                                        保存
                                    </Button></WithPermission></Col></Row>
                                    
            </div>
        </Spin>
    )
})

export default SystemInsideUpstreamContent