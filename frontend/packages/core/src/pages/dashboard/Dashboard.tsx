/*
 * @Date: 2024-02-26 19:03:06
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-07 10:36:14
 * @FilePath: \frontend\packages\core\src\pages\dashboard\Dashboard.tsx
 */

import { useEffect, useState } from "react";
import { useBreadcrumb } from "@common/contexts/BreadcrumbContext";
import { App, Spin, Tabs, TabsProps } from "antd";
import { BasicResponse, STATUS_CODE } from "@common/const/const";
import { DashboardPartitionItem } from "@common/const/type";
import { useFetch } from "@common/hooks/http";
import DashboardPage from "./DashboardTabPage";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import { RouterParams } from "@core/components/aoplatform/RenderRoutes";
import DashboardInstruction from "./DashboardInstruction";
import { LoadingOutlined } from "@ant-design/icons";

export default function Dashboard(){
    const { message } = App.useApp()
    const { setBreadcrumb } = useBreadcrumb()
    const [partitionOption,setPartitionOption] = useState<TabsProps['items']>([])
    const {fetchData} = useFetch()
    const {partitionId} = useParams<RouterParams>()
    const navigateTo = useNavigate()
    const [loading, setLoading] = useState<boolean>(true)
    const location = useLocation().pathname
 
    const getPartitionList = ()=>{
        setLoading(true)
        setPartitionOption([])
        fetchData<BasicResponse<{partitions:DashboardPartitionItem[]}>>('simple/monitor/partitions',{method:'GET',eoTransformKeys:['enable_monitor']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                const filteredPartitions = data.partitions?.filter((x:DashboardPartitionItem)=>x.enableMonitor)
                setPartitionOption(filteredPartitions.map((x:DashboardPartitionItem)=>{return {
                    label:x.name,key:x.id,children:<DashboardPage />,destroyInactiveTabPane:true
                }}))
               partitionId === undefined&& filteredPartitions.length > 0 && navigateTo(`/dashboard/${filteredPartitions[0].id}/total`)
                
            }else{
                message.error(msg || '操作失败')
            }
        }).finally(()=>setLoading(false))
    }

    useEffect(() => {
        setBreadcrumb([
            {
                title:'运行视图'
            },
        ])
        getPartitionList()

    }, [partitionId]);

    return (
        <>
            <Spin wrapperClassName="h-[calc(100vh-50px)]" indicator={<LoadingOutlined style={{ fontSize: 24 }} spin/>} spinning={loading}>
                {!loading && <>
                {partitionOption && partitionOption?.length > 0 && partitionId !== undefined? 
                <Tabs   
                    activeKey={partitionId ?? partitionOption?.[0]?.key} 
                    items={partitionOption}  
                    className="h-auto"  
                    type="card" 
                    tabBarGutter={0} 
                    destroyInactiveTabPane={true}
                    tabBarStyle={{marginTop:'0px',marginBottom:'0px',border:'none'}} 
                    onChange={(val)=>{
                        partitionId !== val && navigateTo(location.replace(partitionId,val))
                        }} />
                :<DashboardInstruction />} </> }
                
            </Spin>
            
        </>
    )
}