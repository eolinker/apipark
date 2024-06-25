/*
 * @Date: 2024-04-19 15:22:46
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 19:11:51
 * @FilePath: \frontend\packages\core\src\pages\access\AccessPage.tsx
 */
import  {useEffect, useState} from "react";
import {App, Spin, Tabs, TabsProps} from "antd";
import {Outlet, useNavigate, useParams} from "react-router-dom";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import { PermissionType, PermissionListItem } from "./AccessList.tsx";
import { useBreadcrumb } from "@common/contexts/BreadcrumbContext.tsx";
import { RouterParams } from "@core/components/aoplatform/RenderRoutes.tsx";
import { LoadingOutlined } from "@ant-design/icons";

export type PermissionOptionType = {
    key:string
    label:string 
    name:string 
    tag:string 
    type:string
}

const ACCESS_TAB_ITEMS:TabsProps['items'] = [
    {
        label: '平台权限',
        key: 'system',
    },{
        label: '团队权限模板',
        key: 'team'
    },{
        label: '服务权限模板',
        key: 'project'
    }
]

export default function AccessPage(){
    const {message } = App.useApp()
    const navigate = useNavigate()
    const {fetchData} = useFetch()
    const [selectedPermissionGrp,SetSelectedPermissionGrp] = useState<string>('system')
    const [accessMap, setAccessMap] = useState<{system:PermissionListItem[],team:PermissionListItem[],project:PermissionListItem[]}>()
    const {setBreadcrumb} = useBreadcrumb()
    const {accessType} = useParams<RouterParams>();
    const [loading, setLoading] = useState<boolean>(true)

    const getAccessList = ()=>{
        setLoading(true)
        fetchData<BasicResponse< PermissionType>>('system/permissions',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setAccessMap(data)
                if((selectedPermissionGrp.length === 0 || !accessType) && data.system.length > 0 ){
                    SetSelectedPermissionGrp('system')
                    navigate('/access/system')
                }
            }else{
                message.error(msg || '操作失败')
            }
        }).finally(()=>{
            setLoading(false)
        })
    }

    const onChange = (key: string) => {
        SetSelectedPermissionGrp(key)
        navigate(`/access/${key}`)
    };

    useEffect(() => {
        setBreadcrumb([
            {title:'权限配置'}
        ])
        getAccessList()
    }, []);

    return (
        <div className="flex flex-1 flex-col h-full ">
                <Tabs 
                    defaultActiveKey={'system'} 
                    size="small" 
                    className="h-auto bg-MAIN_BG" 
                    tabBarStyle={{paddingLeft:'10px',marginTop:'0px',marginBottom:'0px'}} 
                    tabBarGutter={20} 
                    items={ACCESS_TAB_ITEMS}
                    onChange={onChange} 
                    destroyInactiveTabPane={true}/>
                    <div className="h-full overflow-hidden">
                        <Spin size="large" wrapperClassName="h-[calc(100vh-94px)] overflow-auto" indicator={<LoadingOutlined style={{ fontSize: 24 }} spin/>} spinning={loading}>
                            {!loading &&<Outlet context={{accessMap,getAccessList}}/>}
                            </Spin>
                            </div>
        </div>)
}