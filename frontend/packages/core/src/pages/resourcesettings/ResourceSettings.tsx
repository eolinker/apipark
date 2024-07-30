/*
 * @Date: 2024-05-15 14:59:07
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-07 15:57:31
 * @FilePath: \frontend\packages\core\src\pages\resourcesettings\ResourceSettings.tsx
 */
import { Menu, MenuProps, Skeleton, message } from "antd";
import { Link, Outlet, useNavigate, useParams } from "react-router-dom";
import InsidePage from "@common/components/aoplatform/InsidePage";
import { useEffect, useState } from "react";
import { BasicResponse, STATUS_CODE } from "@common/const/const";
import { DynamicMenuItem, PartitionItem } from "@common/const/type";
import { useFetch } from "@common/hooks/http";
import { getItem } from "@common/utils/navigation";
import { RouterParams } from "@core/components/aoplatform/RenderRoutes";
import { DefaultOptionType } from "antd/es/select";
import ResourceSettingsInstruction from "./ResourceSettingsInstruction";

const LogSettings = ()=>{
    const {moduleId} = useParams<RouterParams>();
    const [menuItems, setMenuItems ] = useState<MenuProps['items']>([])
    const [activeMenu, setActiveMenu] = useState<string>()
    const {fetchData} = useFetch()
    const [loading, setLoading] = useState<boolean>(true)
    const navigateTo = useNavigate()
    const [partitionOptions, setPartitionOption] = useState<DefaultOptionType[]>([])

    const getPartitionList = ()=>{
        setPartitionOption([])
        return fetchData<BasicResponse<{partitions:PartitionItem[]}>>('simple/partitions',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setPartitionOption(data.partitions?.map((x:PartitionItem)=>{return {
                    label:x.name,value:x.id
                }}))
                return Promise.resolve(data.partitions)
            }else{
                message.error(msg || '操作失败')
                return Promise.reject(msg || '操作失败')
            }
        })
    }

    
    const getDynamicMenuList = ()=>{
        setLoading(true)
        fetchData<BasicResponse<{ dynamics:DynamicMenuItem[] }>>(`simple/dynamics/resource`,{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                const newMenu:MenuProps['items'] =  data.dynamics.map((x:DynamicMenuItem)=>
                    getItem(
                        <Link to={`template/${x.name}`}>{x.title}</Link>, 
                        x.name,
                        undefined,
                        undefined,
                        undefined,
                        'system.partition.self.view')) 
                
                    setMenuItems(newMenu)
                    if(!activeMenu || activeMenu.length === 0){
                        navigateTo(`/resourcesettings/template/${data.dynamics[0].name}`)
                    }
            }else{
                message.error(msg || '操作失败')
            }
        }).finally(()=>setLoading(false))
    }

    const onMenuClick: MenuProps['onClick'] = ({key}) => {
        setActiveMenu(key)
    };
    
    useEffect(() => {
        setActiveMenu(moduleId)
    }, [ moduleId]);
    
    useEffect(()=>{
        setLoading(true)
        Promise.all([getPartitionList(),getDynamicMenuList()]).finally(()=>setLoading(false))
    },[])
    
    
    return (
        <> 
          <Skeleton className='m-btnbase w-[calc(100%-20px)]' active loading={loading}>
            {loading ? null : partitionOptions && partitionOptions.length > 0 ? 
            <InsidePage 
                pageTitle='资源配置'
                >
                <div className="flex h-full">
                    <Menu
                        className="h-full overflow-y-auto"
                        selectedKeys={[activeMenu || '']}
                        onClick={onMenuClick}
                        style={{ width: 182 }}
                        mode="inline"
                        items={menuItems}
                    />
                    <div className={`w-full flex flex-1 flex-col h-full overflow-auto bg-MAIN_BG`}>
                        <Outlet />
                    </div>
                </div>
            </InsidePage>: 
                <ResourceSettingsInstruction/>
            }
            </Skeleton>
        </>
    )
}

export default LogSettings;