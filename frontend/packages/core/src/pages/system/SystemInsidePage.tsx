/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-07-12 20:58:39
 * @FilePath: \frontend\packages\core\src\pages\system\SystemInsidePage.tsx
 */
import  {FC, useEffect, useMemo, useState} from "react";
import {Outlet, useLocation, useNavigate, useParams} from "react-router-dom";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {App, Menu, MenuProps} from "antd";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import { useSystemContext} from "../../contexts/SystemContext.tsx";
import { SYSTEM_PAGE_MENU_ITEMS } from "../../const/system/const.tsx";
import { SystemConfigFieldType } from "../../const/system/type.ts";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";
import { PERMISSION_DEFINITION } from "@common/const/permissions.ts";
import InsidePage from "@common/components/aoplatform/InsidePage.tsx";
import Paragraph from "antd/es/typography/Paragraph";
import { ItemType, MenuItemGroupType, MenuItemType } from "antd/es/menu/hooks/useItems";
import { cloneDeep } from "lodash-es";

const SystemInsidePage:FC = ()=> {
    const { message } = App.useApp()
    const {orgId, teamId,systemId,partitionId,apiId} = useParams<RouterParams>();
    const location = useLocation()
    const currentUrl = location.pathname
    const {fetchData} = useFetch()
    const { setPartitionList,setPrefixForce,setApiPrefix ,systemInfo,setSystemInfo} = useSystemContext()
    const { accessData,checkPermission,projectDataFlushed} = useGlobalContext()
    const [activeMenu, setActiveMenu] = useState<string>()
    const navigateTo = useNavigate()

    const getSystemInfo = ()=>{
        fetchData<BasicResponse<{ project:SystemConfigFieldType }>>('project/info',{method:'GET',eoParams:{project:systemId}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setSystemInfo(data.project)
                setPartitionList(data.project.partition)
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const getApiDefine = ()=>{
        setApiPrefix('')
        setPrefixForce(false)
        fetchData<BasicResponse<{ prefix:string, force:boolean }>>('project/api/define',{method:'GET',eoParams:{project:systemId}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setApiPrefix(data.prefix)
                setPrefixForce(data.force)
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const menuData = useMemo(()=>{
        const filterMenu = (menu:MenuItemGroupType<MenuItemType>[])=>{
            const newMenu = cloneDeep(menu)
            return newMenu!.filter((m:MenuItemGroupType )=>{
                if(m.children && m.children.length > 0){
                     m.children = m.children.filter(
                        (c)=>(c&&(c as MenuItemType&{access:string} ).access ? 
                            checkPermission((c as MenuItemType&{access:string} ).access as keyof typeof PERMISSION_DEFINITION[0]): 
                            true))
                }
                return m.children && m.children.length > 0
            })
        }
        const filteredMenu = filterMenu(SYSTEM_PAGE_MENU_ITEMS as MenuItemGroupType<MenuItemType>[])
        setActiveMenu((pre)=>{
            if(!pre && projectDataFlushed){
                const activeMenu = filteredMenu?.[0]?.children?.[0]?.key as string
                return activeMenu
            }
            return pre
        })
        return  filteredMenu || []
    },[accessData,projectDataFlushed])
    
    const onMenuClick: MenuProps['onClick'] = ({key}) => {
        setActiveMenu(key)
    };
    
    useEffect(() => {
        if(partitionId !== undefined){
            setActiveMenu('upstream')
        }else if(apiId !== undefined){
            setActiveMenu('api')
        }else if(systemId !== currentUrl.split('/')[currentUrl.split('/').length - 1]){ 
            setActiveMenu(currentUrl.split('/')[currentUrl.split('/').length - 1])
        }
    }, [currentUrl]);

    useEffect(()=>{
        if(accessData && accessData.get('system') && accessData.get('system')?.indexOf('project.mySystem.api.view') !== -1){
            getApiDefine()
        }
    },[accessData])

    useEffect(()=>{
        if( activeMenu && systemId === currentUrl.split('/')[currentUrl.split('/').length - 1]){
            navigateTo(`/system/${orgId}/${teamId}/inside/${systemId}/${activeMenu}`)
        }
    },[activeMenu])

    useEffect(() => {
        systemId && getSystemInfo()
    }, [systemId]);

    return (
        <>
        <InsidePage pageTitle={systemInfo?.name || '-'} 
                tagList={[{label:
                    <Paragraph className="mb-0" copyable={systemId ? { text: systemId } : false}>服务 ID：{systemId || '-'}</Paragraph>
                }]}
                backUrl="/system/list">
                <div className="flex flex-1 h-full">
                    <Menu
                        onClick={onMenuClick}
                        className="h-full overflow-y-auto"
                        style={{ width: 182 }}
                        selectedKeys={[activeMenu!]}
                        mode="inline"
                        items={menuData as unknown as ItemType<MenuItemType>[] } 
                    />
                    <div  className={` ${activeMenu?.indexOf('setting')  !== -1   ? 'pt-[20px] pl-[10px] pr-btnrbase' :''} w-full h-full flex flex-1 flex-col overflow-auto bg-MAIN_BG`}>
                            <Outlet/>
                    </div>
                </div>
            </InsidePage>

        </>
    )
}
export default SystemInsidePage