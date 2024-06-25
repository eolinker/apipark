import { Descriptions, Collapse, Space, Input, App, Empty } from "antd";
import ApiPreview from "@common/components/postcat/ApiPreview";
import ApiMatch from "@common/components/postcat/api/ApiPreview/components/ApiMatch";
import ApiProxy from "@common/components/postcat/api/ApiPreview/components/ApiProxy";
import { useState, useEffect, useMemo } from "react";
import { BasicResponse, STATUS_CODE } from "@common/const/const";
import { ServiceBasicInfoType, ServiceDetailType } from "../../../const/serviceHub/type";
import { useFetch } from "@common/hooks/http";
import { SystemApiDetail } from "../../../const/system/type";

type SubServiceDetail = {
    serviceId?:string
}

export default function SubServiceDetail(props:SubServiceDetail){
    const { serviceId} = props
    const { message } = App.useApp()
    const [apiDocs,setApiDocs ] = useState<SystemApiDetail[]>()
    const [serviceBasicInfo, setServiceBasicInfo] = useState<ServiceBasicInfoType>()
    const {fetchData} = useFetch()
    const [activeKey, setActiveKey] = useState<string[]>([])

    const getBasicInfo = useMemo(() => [
        {
            key: 'organzation',
            label: '所属组织',
            children: serviceBasicInfo?.organization.name,
            style: {paddingBottom: '10px'},
        },
        {
            key: 'project',
            label: '所属服务',
            children: serviceBasicInfo?.project.name,
            style: {paddingBottom: '10px'},
        },
        {
            key: 'team',
            label: '所属团队',
            children: serviceBasicInfo?.team.name,
            style: {paddingBottom: '10px'},
        },
        // {
        //     key: 'master',
        //     label: '负责人',
        //     children: serviceBasicInfo?.master?.name,
        //     style: {paddingBottom: '10px'},
        // }
    ], [serviceBasicInfo]);


    const getServiceBasicInfo = ()=>{
        fetchData<BasicResponse<{service:ServiceDetailType}>>('catalogue/service',{method:'GET',eoParams:{service:serviceId}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setServiceBasicInfo(data.service.basic)
                setApiDocs(data.service.apis)
                setActiveKey(data.service.apis.map((x)=>x.id))
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    
    useEffect(() => {
        if(!serviceId){
            console.warn('缺少serviceId')
            return
        }
        serviceId && getServiceBasicInfo()
    }, [serviceId]);

    return (
        <div className="flex flex-col h-full flex-1">
                    <div className="bg-[#fff] rounded  mb-[16px]">
                        <p className="font-bold text-[16px] leading-[22px] mb-[8px] h-[22px]">基本信息</p>
                        <Descriptions className="bg-bar-theme p-[16px] rounded service-hub-description" title="" items={getBasicInfo} column={4} labelStyle={{width:'80px',justifyContent:'flex-end',fontWeight:'bold'}}  contentStyle={{color:'#333'}}/>
                    </div>
                    <div  className='bg-[#fff] rounded '>
                        <div className="" >
                        <p className="font-bold  text-[16px] leading-[22px] mb-[12px] h-[22px]" id="apiDocument-list">API 列表</p>
                            <div className="">
                                {apiDocs?.map((apiDetail)=>(
                                    <div  className="mb-btnbase "  key={apiDetail.id} id={`apiDocument-${apiDetail.id}`}>
                                    <Collapse key={`apiDocument-${apiDetail.id}`} 
                                        expandIcon={({isActive})=>(isActive?  <iconpark-icon name="shouqi-2"></iconpark-icon>:<iconpark-icon name="zhankai"></iconpark-icon> )}
                                        items={[{
                                            key: apiDetail.id,
                                            label: <span><span className="text-status_update font-bold mr-[8px]">{apiDetail.method}</span><span>{apiDetail.name}</span></span>,
                                            children:<div className="scroll-area h-[calc(100%-84px)] overflow-auto">
                                                    <Space direction="vertical" className="mb-btnybase w-full mt-btnybase">
                                                    <Input
                                                        readOnly
                                                        addonBefore={apiDetail?.method}
                                                        value={apiDetail?.path}
                                                        // enterButton={<SearchBtn  entity={apiDetail}/>}
                                                        // onSearch={handleTest}
                                                    />
                                                </Space>
                                            {
                                                apiDetail?.match && apiDetail.match?.length > 0 &&
                                                <ApiMatch title='高级匹配' rows={apiDetail?.match}  />
                                            }
                            
                                            {
                                                apiDetail?.proxy && Object.keys(apiDetail?.proxy).length > 0 &&
                                                <ApiProxy title='转发规则' proxyInfo={apiDetail?.proxy}  />
                                            }
                            
                                            {apiDetail && <ApiPreview entity={{...apiDetail.doc,name:apiDetail.name, method:apiDetail.method,uri:apiDetail.path, protocol:apiDetail.protocol||'HTTP'}}  />}
                                        </div>
                                            // <ApiPreview testClick={()=>testClick(apiDocs.id)} entity={doc}  /> 
                                        }]} 
                                        activeKey={activeKey}
                                        onChange={(val)=>{setActiveKey(val as string[])}}
                                            />
                                    </div>
                                ))}

                        </div>
                            {/* <div className="h-[16px] "></div>
                            <div className='bg-[#fff] rounded  mb-[20px]'>
                                <p className="font-bold text-[16px] leading-[22px] mb-[12px] h-[22px]" id="apiDocument-statusCode">状态码</p>
                                <CodePage />
                            </div> */}
                            </div>
                        </div>
                    </div>
    )
}