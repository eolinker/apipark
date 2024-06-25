/*
 * @Date: 2024-06-05 16:00:58
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-07 15:59:46
 * @FilePath: \frontend\packages\core\src\pages\dashboard\DashboardInstruction.tsx
 */
import { Link } from "react-router-dom";

export default function DashboardInstruction() {
    return (
        <div className="h-full w-full overflow-auto">
            <div className=" m-auto mt-[10%] flex flex-col items-center  p-[20px]">
                <p className="text-[20px] font-medium leading-[32px] text-MAIN_TEXT">分区配置并开启监控</p>
                <p className="text-[12px] font-normal leading-[20px] text-DESC_TEXT mt-[12px]" >监控功能用于辅助管理分区内信息，请选择或创建分区，设置监控信息后查看当前分区监控情况；</p>
                {/* <p className="text-[12px] font-normal leading-[20px] text-DESC_TEXT mt-[8px]">更多配置问题，请点击帮助中心
                    {/* <a>查看更多</a> *
                    </p> */}
                <div className="flex mt-[28px]">
                    <div className="h-[208px] w-[384px] flex flex-col items-center py-[32px] px-[24px] gap-[16px] rounded-DEFAULT bg-MENU_BG mr-[24px]">
                        <p className="text-[20px] font-medium leading-[32px] text-MAIN_TEXT">分区配置</p>
                        <p className="text-[12px] font-normal leading-[20px] text-DESC_TEXT">新增分区的名称、描述和其他相关属性，并设置该分区内的集群地址，以确保监控系统能够正确识别和连接到这些集群</p>
                        <p><Link to="/partition/list">添加分区信息</Link></p>
                    </div>
                    <div className="h-[208px] w-[384px] flex flex-col items-center py-[32px] px-[24px] gap-[16px] rounded-DEFAULT bg-MENU_BG ">
                        <p className="text-[20px] font-medium leading-[32px] text-MAIN_TEXT">监控设置</p>
                        <p className="text-[12px] font-normal leading-[20px] text-DESC_TEXT">配置监控参数，以确保监控系统可以正确监控和收集集群的性能数据。这包括设置监控指标、阈值、报警规则等。监控系统将定期收集数据并生成报告，帮助用户了解集群运行状态和性能表现</p>
                        <p><Link to="/partition/list">配置监控信息</Link></p>
                    </div>
                </div>
            </div></div>
    )
}