/*
 * @Date: 2024-04-22 10:45:49
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 13:45:47
 * @FilePath: \frontend\packages\core\src\pages\system\publish\SystemInsidePublishOnline.tsx
 */
import { App, Table } from "antd";
import { SYSTEM_PUBLISH_ONLINE_COLUMNS } from "../../../const/system/const";
import { useEffect, useState } from "react";
import { useFetch } from "@common/hooks/http";
import { BasicResponse, STATUS_CODE } from "@common/const/const";
import { EntityItem } from "@common/const/type";

type SystemInsidePublishOnlineProps = {
    systemId:string
    id:string
}

export type SystemInsidePublishOnlineItems = {
    partition:EntityItem
    cluster:EntityItem
    status:'done' | 'error' | 'publishing'
    error:string
}
export default function SystemInsidePublishOnline(props:SystemInsidePublishOnlineProps ){
    const {systemId, id} = props
    const {message} = App.useApp()
    const [dataSource, setDataSource] = useState<[]>()
    const {fetchData} = useFetch()
    const [isStopped, setIsStopped] = useState(false);

    const getOnlineStatus = ()=>{
        fetchData<BasicResponse<{publishStatusList:SystemInsidePublishOnlineItems[]}>>('project/publish/status',{method:'GET',eoParams:{project:systemId, id}, eoTransformKeys:['publish_status_list']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setDataSource(data.publishStatusList)
                if(data.publishStatusList.filter((x:SystemInsidePublishOnlineItems)=>x.status === 'publishing').length === 0){
                    setIsStopped(true)
                }
            }else{
                message.error(msg || '操作失败')
            }
        }).catch((errorInfo)=> message.error(errorInfo))
    }

    useEffect(()=>{
        getOnlineStatus();
    },[])

    useEffect(() => {
        let intervalId: NodeJS.Timeout;
        if (!isStopped) {
            intervalId = setInterval(() => {
                !isStopped && getOnlineStatus();
            }, 5000);
        }

        return () => {
            clearInterval(intervalId);
        };
    }, [isStopped]);
    
    return (
        <Table
            className="min-h-[100px] h-full"
            bordered={true}
            columns={[...SYSTEM_PUBLISH_ONLINE_COLUMNS]}
            size="small"
            rowKey="id"
            dataSource={dataSource}
            pagination={false}
        />
    )
}