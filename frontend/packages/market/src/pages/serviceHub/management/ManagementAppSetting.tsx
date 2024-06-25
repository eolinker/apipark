/*
 * @Date: 2024-06-06 11:47:47
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-06 14:13:32
 * @FilePath: \frontend\packages\market\src\pages\serviceHub\management\ManagementAppSetting.tsx
 */
import { RouterParams } from "@core/components/aoplatform/RenderRoutes";
import { useParams } from "react-router-dom";
import ManagementConfig from "./ManagementConfig";

export default function ManagementAppSetting(){
    const {teamId,appId} = useParams<RouterParams>()
    
    return (
        <div className="w-[70%] mx-auto h-full pt-[32px]">
        <div className="flex items-center justify-between w-full ml-[10px] text-[18px] leading-[25px] pb-[16px]" ><span className="font-bold">应用管理</span></div>
        <div className="h-[calc(100%-41px)] flex flex-col ">
            <ManagementConfig type='edit' teamId={teamId!} appId={appId!}/>
        </div>
       </div>
    )
}