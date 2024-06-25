/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 19:12:16
 * @FilePath: \frontend\packages\core\src\pages\team\TeamInsidePage.tsx
 */
import  {FC, useEffect, useMemo, useState} from "react";
import { Outlet, useLocation, useNavigate, useParams} from "react-router-dom";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {App, Menu, MenuProps} from "antd";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {OrganizationItem} from "@common/const/type.ts";
import { TEAM_INSIDE_MENU_ITEMS } from "../../const/team/const.tsx";
import { useTeamContext } from "../../contexts/TeamContext.tsx";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";
import InsidePage from "@common/components/aoplatform/InsidePage.tsx";
import Paragraph from "antd/es/typography/Paragraph";
import { MenuItemGroupType, MenuItemType } from "antd/es/menu/hooks/useItems";
import { cloneDeep } from "lodash-es";
import { PERMISSION_DEFINITION } from "@common/const/permissions.ts";

const TeamInsidePage:FC = ()=> {
    const { message } = App.useApp()
    const {orgId,teamId} = useParams<RouterParams>();
    const {fetchData} = useFetch()
    const location = useLocation()
    const { teamInfo ,setTeamInfo } = useTeamContext()
    const {getTeamAccessData,cleanTeamAccessData,accessData,checkPermission,teamDataFlushed} = useGlobalContext()
    const navigateTo = useNavigate()
    const [activeMenu, setActiveMenu] = useState<string>()

    const onMenuClick: MenuProps['onClick'] = ({key}) => {
        setActiveMenu(key)
    };

    const menuData = useMemo(()=>{
        const filterMenu = (menu:MenuItemGroupType<MenuItemType>[])=>{
            const newMenu = cloneDeep(menu)
            return newMenu!.filter((m:MenuItemGroupType )=>{
                if(m.children && m.children.length > 0){
                     m.children = m.children.filter(
                        (c)=>(c as MenuItemType&{access:string | string[]} ).access ? 
                            checkPermission((c as MenuItemType&{access:string | string[]} ).access as keyof typeof PERMISSION_DEFINITION[0]): 
                            true)
                }
                return m.children && m.children.length > 0
            })
        }
        const filteredMenu = filterMenu(TEAM_INSIDE_MENU_ITEMS as MenuItemGroupType<MenuItemType>[])
        setActiveMenu((pre)=>{
            if(!pre && teamDataFlushed){
                const activeMenu = filteredMenu?.[0]?.children?.[0]?.key as string
                return activeMenu
            }
            return pre
        })
        return  filteredMenu || []
    },[accessData])

    const getTeamInfo = ()=>{
        setTeamInfo?.(undefined)
        fetchData<BasicResponse<{ team:OrganizationItem }>>('team',{method:'GET',eoParams:{team:teamId}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setTeamInfo?.(data.team)
            }else{
                message.error(msg || '操作失败')
            }
        })
    }
    

    useEffect(() => {
        if(location.pathname.split('/')[location.pathname.split('/').length -1] !== teamId){
            setActiveMenu(location.pathname.split('/')[location.pathname.split('/').length -1])
        }
    }, [location]);

    useEffect(()=>{
        if( activeMenu && teamId === location.pathname.split('/')[location.pathname.split('/').length - 1]){
            navigateTo(`/team/inside/${orgId}/${teamId}/${activeMenu}`)
        }
    },[activeMenu])

    useEffect(()=>{
        getTeamInfo()
        teamId && getTeamAccessData(teamId)
        return ()=>{
            cleanTeamAccessData()
        }
    },[teamId])

    return (
        <>
            <InsidePage 
                pageTitle={teamInfo?.name || '-'} 
                tagList={[{label:
                    <Paragraph className="mb-0" copyable={teamId ? { text: teamId } : false}>团队 ID：{teamId || '-'}</Paragraph>
            }]}
                backUrl="/team/list">
                <div className="flex h-full">
                    <Menu
                        style={{ width: 176 }}
                        mode="inline"
                        items={menuData}
                        onClick={onMenuClick}
                        selectedKeys={[activeMenu || '']}
                    />
                    <div className={`flex flex-1 flex-col h-full overflow-auto bg-MAIN_BG ${activeMenu === 'setting' ? 'pt-[20px] pl-[10px] pr-btnrbase ':''}`}>
                        <Outlet  />
                    </div>
                </div>
            </InsidePage>
            </>
    )
}
export default TeamInsidePage