/*
 * @Date: 2024-05-30 18:18:40
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-07 17:37:01
 * @FilePath: \frontend\packages\market\src\pages\serviceHub\management\ManagementInsidePage.tsx
 */
import { ApiOutlined, ArrowLeftOutlined, LoadingOutlined } from "@ant-design/icons";
import { App, Button, Menu, MenuProps, Spin } from "antd";
import { useState, useEffect, useMemo } from "react";
import { Link, Outlet, useLocation, useNavigate, useParams } from "react-router-dom";
import { BasicResponse, STATUS_CODE } from "@common/const/const";
import { useBreadcrumb } from "@common/contexts/BreadcrumbContext";
import { useFetch } from "@common/hooks/http";
import { ItemType } from "antd/es/breadcrumb/Breadcrumb";
import { TENANT_MANAGEMENT_APP_MENU } from "../../../const/serviceHub/const";
import { EntityItem, SimpleTeamItem } from "@common/const/type";
import { RouterParams } from "@core/components/aoplatform/RenderRoutes";
import { getItem } from "@common/utils/navigation";
import { useTenantManagementContext } from "@market/contexts/TenantManagementContext";
import { ManagementConfigFieldType } from "./ManagementConfig";

export default function ManagementInsidePage(){
    const { message } = App.useApp()
    const {fetchData} = useFetch()
    const { setBreadcrumb} = useBreadcrumb()
    const [activeMenu, setActiveMenu] = useState<string>()
    const [partitionList, setPartitionList] = useState<EntityItem[]>([])
    const {partitionId,appId,teamId} = useParams<RouterParams>()
    const navigateTo = useNavigate()
    const currentUrl = useLocation().pathname
    const [openKeys, setOpenKeys] = useState<string[]>([])
    const [loading, setLoading] = useState<boolean>(false)
    const {appName,setAppName} = useTenantManagementContext()
    
  const getPartitionList = ()=>{
    setLoading(true)
    fetchData<BasicResponse<{partitions:(EntityItem &{ serviceNum:number})[]}>>('simple/application/partitions',{method:'GET',eoParams:{application:appId},eoTransformKeys:['service_num']}).then(response=>{
        const {code,data,msg} = response
        if(code === STATUS_CODE.SUCCESS){
            setPartitionList(data.partitions?.map((x:SimpleTeamItem)=>({label:<div className="flex items-center justify-between "><span className="w-[calc(100%-42px)] truncate" title={x.name}>{x.name}</span><span className="bg-[#fff] rounded-[5px] h-[20px] w-[30px] flex items-center justify-center">{x.serviceNum || 0}</span></div>, key:x.id})))
            if(!partitionId && !activeMenu){
                data.partitions&&data.partitions.length > 0 ?  navigateTo(`service/${data.partitions[0].id}`) :navigateTo('authorization')
                return 
            }
        }else{
            message.error(msg || '操作失败')
        }
    }).finally(()=>{
        setLoading(false)
    })
}

    const menuData = useMemo(()=>{
        if(!partitionList || partitionList.length === 0) return TENANT_MANAGEMENT_APP_MENU
        setOpenKeys(['service'])
        const serviceMenu = getItem('订阅的服务', 'service', <ApiOutlined />, partitionList as unknown as ItemType[])
        return  [serviceMenu,...TENANT_MANAGEMENT_APP_MENU as unknown[]] 
    },[partitionList])

    const onMenuClick: MenuProps['onClick'] = (node) => {
            setActiveMenu(node.key)
        if(['authorization','setting'].includes(node.key)){
            navigateTo(`/tenantManagement/${teamId}/inside/${appId}/${node.key}`)
        }else{
            navigateTo(`/tenantManagement/${teamId}/inside/${appId}/service/${node.key}`)
        }
    };


    useEffect(()=>{
        if(!partitionId && currentUrl.includes('authorization')){
            setActiveMenu('authorization')
        }else if(partitionId){
            setActiveMenu(partitionId)
        }
    },[currentUrl])

    useEffect(()=>{
        const fetchDataAsync = async () => {
            let _appName = appName
            if(appId && !appName  && !currentUrl.includes('setting')){
                const {code,data} = await fetchData<BasicResponse<{ project: ManagementConfigFieldType }>>('app/info',{method:'GET',eoParams:{app:appId},eoTransformKeys:['as_app']})
                if(code === STATUS_CODE.SUCCESS){
                    _appName = data.project.name
                    setAppName(_appName)
                }
            }
            setBreadcrumb(
                [
                    {title:<Link to={`/tenantManagement/list/${teamId}`}>应用</Link>},
                   ...(_appName ? [{title:_appName}] : [])
                ]
            )
        };
        fetchDataAsync();
    },
    [appId,appName])

    useEffect(() => {
        getPartitionList()
    }, [appId]);

    return (<>
        <Spin className="h-full" wrapperClassName="h-full"  indicator={<LoadingOutlined style={{ fontSize: 24 }} spin />} spinning={loading}>
        <div className="flex flex-1 h-full">
            <div className="w-[224px] border-0 border-solid border-r-[1px] border-r-BORDER">
            <div className="text-[18px] leading-[25px] pl-[12px] py-[12px]"><Button type="text" onClick={()=>navigateTo(`/tenantManagement/list/${teamId}`)}><ArrowLeftOutlined />返回</Button></div>
            <Menu
                onClick={onMenuClick}
                openKeys={openKeys}
                onOpenChange={(e)=>{setOpenKeys(e)}}
                className="h-[calc(100%-59px)] overflow-auto"
                style={{ width: 224, paddingLeft:'8px', paddingRight:'8px' }}
                selectedKeys={[activeMenu!]}
                mode="inline"
                items={menuData as unknown as ItemType<MenuItemType>[] } 
                />
        </div>
        <div className="w-[calc(100%-224px)] pb-[10px] overflow-auto">
            <Outlet context={{refreshGroup:()=>getPartitionList()}}></Outlet>
        </div>
    </div>
    </Spin></>)
}