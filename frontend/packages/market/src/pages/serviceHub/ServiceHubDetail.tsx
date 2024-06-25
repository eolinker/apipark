/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-05 18:00:51
 * @FilePath: \frontend\packages\market\src\pages\serviceHub\ServiceHubDetail.tsx
 */
import {Link, useLocation, useParams} from "react-router-dom";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {Anchor, App, Button, Collapse, Descriptions, Drawer, FloatButton, Input, Space} from "antd";
import  { useEffect, useMemo, useRef, useState} from "react";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import ApiPreview from "@common/components/postcat/ApiPreview.tsx";
import ApiTestGroup from "./ApiTestGroup.tsx";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {DefaultOptionType} from "antd/es/cascader";
import {ApiDetail} from "@common/const/api-detail";
import { ApplyServiceHandle, ServiceBasicInfoType, ServiceDetailType } from "../../const/serviceHub/type.ts";
import { SimpleSystemItem } from "@core/const/system/type.ts";
import { EntityItem } from "@common/const/type.ts";
import ApiMatch from "@common/components/postcat/api/ApiPreview/components/ApiMatch/index.tsx";
import ApiProxy from "@common/components/postcat/api/ApiPreview/components/ApiProxy/index.tsx";
import InsidePageForHub from "@common/components/aoplatform/InsidePageForHub.tsx";
import { ApplyServiceModal } from "./ApplyServiceModal.tsx";

const ServiceHubDetail = ()=>{
    const {tagId, categoryId, teamId,serviceId} = useParams<RouterParams>();
    const  cluster:string  = (new URLSearchParams(useLocation().search)).get('name') || '-'
    const {setBreadcrumb} = useBreadcrumb()
    const [apiTestDrawOpen, setApiTestDrawOpen] = useState(false);
    const [serviceBasicInfo, setServiceBasicInfo] = useState<ServiceBasicInfoType>()
    const [serviceName, setServiceName] = useState<string>()
    const [serviceDesc, setServiceDesc] = useState<string>()
    const [apiDocs,setApiDocs ] = useState<ApiDetail[]>()
    const {fetchData} = useFetch()
    const applyRef = useRef<ApplyServiceHandle>(null)
    const { modal,message } = App.useApp()
    const [partitionsList, setPartitionsList ] = useState<Array<{id:string, name:string}>>()
    const [mySystemOptionList, setMySystemOptionList] = useState<DefaultOptionType[]>()
    const [applied,setApplied] = useState<boolean>(false)
    const [selectedTestApi,setSelectedTestApi] = useState<string>()
    // const callbackUrl = new URLSearchParams(window.location.search).get('callbackUrl');
    const [activeKey, setActiveKey] = useState<string[]>([])

    const getServiceBasicInfo = ()=>{
        fetchData<BasicResponse<{service:ServiceDetailType}>>('catalogue/service',{method:'GET',eoParams:{service:serviceId}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                //console.log(data)
                setServiceBasicInfo(data.service.basic)
                setServiceName(data.service.name)
                setServiceDesc(data.service.description)
                setApiDocs(data.service.apis)
                setApplied(data.service.applied)
                setPartitionsList(data.service.partition)
                setActiveKey(data.service.apis.map((x)=>x.id))
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const getBasicInfo = useMemo(() => [
        {
            key: 'organzation',
            label: '所属组织',
            children: serviceBasicInfo?.organization.name,
            style: {padding: 0},
        },
        {
            key: 'project',
            label: '所属系统',
            children: serviceBasicInfo?.project.name,
            style: {padding: 0},
        },
        {
            key: 'team',
            label: '所属团队',
            children: serviceBasicInfo?.team.name,
            style: {padding: 0},
        },
        // {
        //     key: 'master',
        //     label: '负责人',
        //     children: serviceBasicInfo?.master?.name,
        //     style: {paddingBottom: '10px'},
        // }
    ], [serviceBasicInfo]);


    const category = useMemo(() => [
        {
            key: 'apiDocument-list',
            href: '#apiDocument-list',
            title: 'API 列表',
            children:apiDocs?.map((x)=>({
                key:x.id,
                href:`#apiDocument-${x.id}`,
                title:x.name
            })) || []
        },
        // {
        //     key: 'apiDocument-statusCode',
        //     href: '#apiDocument-statusCode',
        //     title: '状态码',
        // },
    ], [apiDocs]);

    const floatButtonStyle = { top:'10px',position:'sticky', width:'180px',height:'200px'}
    const onClick = async (e: unknown, link: {href: string}) => {
        // const arr = [...(collapseDefaultKeyNew as string[]), link.href];
        // setCollapseDefaultKeyNew(Array.from(new Set(arr)))

    }

    useEffect(() => {
        if(!serviceId){
            console.warn('缺少serviceId')
            return
        }
        serviceId && getServiceBasicInfo()
    }, [serviceId]);

    useEffect(() => {
        getMySelectList()
        setBreadcrumb(
            [
                {title:<Link to={`/serviceHub/list`}>服务市场</Link>},
                // {title:<Link to={`/serviceHub/list`}>服务市场</Link>},
                {title:'服务详情'}
            ]
        )

        // setTimeout(()=>{
        //     const element = document.querySelectorAll('.MuiDataGrid-main');
        //     if(element?.length > 0){
        //         for(const x of element){
        //             x.childNodes[x.childNodes.length - 1 ].textContent === 'MUI X Missing license key' ?  x.childNodes[x.childNodes.length - 1 ].textContent = '' :null
        //         }
        //     }

        // },500)
    }, []);

    const testClick = (id:string)=>{//console.log('test');
        setApiTestDrawOpen(true)
        setSelectedTestApi(id)}

    const onClose = () => {
        setApiTestDrawOpen(false);
    };


    const getMySelectList = ()=>{
        setMySystemOptionList([])
        fetchData<BasicResponse<{ projects: SimpleSystemItem[] }>>('simple/apps/mine',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setMySystemOptionList(data.projects?.map((x:SimpleSystemItem)=>{return {
                    label:x.name, value:x.id, partition:x.partition?.map((x:EntityItem)=>x.id)
                }}))
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const openModal = (type:'apply')=>{
        modal.confirm({
            title:'申请服务',
            content:<ApplyServiceModal ref={applyRef} entity={{...serviceBasicInfo!,partition:partitionsList!, name:serviceName!, id:serviceId!}}  mySystemOptionList={mySystemOptionList!}/>,
            onOk:()=>{
                return applyRef.current?.apply().then((res)=>{
                    if(res === true) setApplied(true)
                })
            },
            okText:'确认',
            cancelText:'取消',
            closable:true,
            icon:<></>,
            width:600
        })
    }

    return (
        <InsidePageForHub pageTitle={serviceName || '-' } tagList={[]} description={serviceDesc || '暂无服务描述'} showBtn={true} btnTitle="申请订阅" onBtnClick={()=>openModal('apply')} backUrl={tagId === undefined && categoryId === undefined ? `/serviceHub/list` :`/serviceHub/list/${tagId === undefined ? 'category' : 'tag' }/${tagId ?? categoryId}`}>
                <div className="flex flex-col p-btnbase pt-[4px] h-full flex-1 overflow-auto" id='layout-ref'>
                    <div className="bg-[#fff] rounded p-btnbase pl-0   mb-[16px]">
                        <Descriptions className="bg-bar-theme p-[16px] rounded service-hub-description" title="" items={getBasicInfo} column={4} labelStyle={{width:'80px',justifyContent:'flex-end',fontWeight:'bold'}}  contentStyle={{color:'#333'}}/>
                    </div>
                    <div  className='bg-[#fff] rounded p-btnbase  pl-0  flex justify-between'>
                        <div className="w-[calc(100%-224px)]" >
                        <p className="font-bold text-[20px] leading-[32px] mb-[12px] h-[32px]" id="apiDocument-list">API 列表</p>
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
                            {/* <div className="h-[16px] bg-[#f7f8fa] mx-[-16px]"></div>
                            <div className='bg-[#fff] rounded  pt-btnbase'>
                                <p className="font-bold text-[20px] leading-[32px] mb-[12px] h-[32px]" id="apiDocument-statusCode">状态码</p>
                                <CodePage />
                            </div> */}
                            </div>

                            <FloatButton.Group shape="circle" style={floatButtonStyle}>
                                <Anchor
                                    // className='absolute py-5 px-btnbase left-0 z-[13]'
                                    // affix={false}
                                    // showInkInFixed={true}
                                    targetOffset={60}
                                    getContainer = {()=> document.getElementById('layout-ref')!}
                                    items={category}
                                />
                            </FloatButton.Group>
                        </div>
                    </div>

            <Drawer 
            title={serviceName} 
              maskClosable={false}
              width="100%" placement="right" onClose={onClose} open={apiTestDrawOpen}
                    extra={
                            <Button onClick={onClose}>退出测试</Button>
                    }
                    closeIcon={false}
            >
                <ApiTestGroup apiInfoList={apiDocs} selectedApiId={selectedTestApi}/>
            </Drawer>
        </InsidePageForHub>
    )
}

export default ServiceHubDetail