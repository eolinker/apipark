/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 19:09:10
 * @FilePath: \frontend\packages\core\src\pages\partitions\PartitionInsidePage.tsx
 */
import  {useEffect, useState} from "react";
import {   Outlet, useLocation, useParams} from "react-router-dom";
import {App, Menu, MenuProps} from "antd";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {useFetch} from "@common/hooks/http.ts";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import { PartitionConfigFieldType } from "../../const/partitions/types.ts";
import { PARTITIONS_INNER_MENU } from "../../const/partitions/const.tsx";
import { usePartitionContext } from "../../contexts/PartitionContext.tsx";
import InsidePage from "@common/components/aoplatform/InsidePage.tsx";
import Paragraph from "antd/es/typography/Paragraph";

const PartitionInsidePage = ()=> {
    const { message } = App.useApp()
    const {partitionId,moduleId} = useParams<RouterParams>();
    const {fetchData} = useFetch()
    const location = useLocation()
    const { partitionInfo,setPartitionInfo} = usePartitionContext()
    const [activeMenu, setActiveMenu] = useState<string>()

    const onMenuClick: MenuProps['onClick'] = ({key}) => {
        setActiveMenu(key)
    };

    const getPartitionInfo = ()=>{
        fetchData<BasicResponse<{ partition:PartitionConfigFieldType }>>(`partition`,{method:'GET',eoParams:{id:partitionId}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setPartitionInfo(data.partition)
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    useEffect(() => {
        const menuStr = location.pathname.split('/')[location.pathname.split('/').length -1]
        let currentMenu:string = menuStr
        if(menuStr === 'node'){currentMenu = 'cluster'}
        if(moduleId &&location.pathname.split('/')[location.pathname.split('/').length -2] === 'template'){
            currentMenu = `template/${moduleId}`
        }
        setActiveMenu(currentMenu)
    }, [location, moduleId]);

    useEffect(() => {
        getPartitionInfo()
    }, [partitionId]);

    return (
        <>
            <InsidePage 
                pageTitle={partitionInfo?.name || '-'} 
                tagList={[{label:
                    <Paragraph className="mb-0" copyable={partitionId ? { text: partitionId } : false}>环境 ID：{partitionId || '-'}</Paragraph>
                }]}
                backUrl="/partition/list">
                <div className="flex h-full">
                    <Menu
                        className="h-full overflow-y-auto"
                        selectedKeys={[activeMenu || '']}
                        onClick={onMenuClick}
                        style={{ width: 182 }}
                        mode="inline"
                        items={PARTITIONS_INNER_MENU}
                    />
                    <div className={` ${activeMenu?.indexOf('setting') !== -1 || activeMenu?.indexOf('dashboard_setting') !== -1   ? 'pt-[20px] pl-[10px] pr-btnrbase' :''} w-full flex flex-1 flex-col h-full overflow-auto bg-MAIN_BG`}>
                        <Outlet />
                    </div>
                </div>
            </InsidePage>
        </>
    )
}
export default PartitionInsidePage