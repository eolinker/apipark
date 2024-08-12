/*
 * @Date: 2024-06-05 16:00:58
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-07 09:18:38
 * @FilePath: \frontend\packages\core\src\pages\dashboard\DashboardList.tsx
 */
import { useParams } from "react-router-dom"
import { RouterParams } from "@core/components/aoplatform/RenderRoutes"
import DashboardApiList from "./DashboardApiList"
import DashboardProjectList from "./DashboardProjectList"
import DashboardApplicationList from "./DashboardApplicationList"

export default function DashboardList(){
    const {dashboardType} = useParams<RouterParams>()
    
    return (
        <>
        {
            dashboardType === 'api' && <DashboardApiList />
        }
        {
            dashboardType === 'subscriber' && <DashboardProjectList />
        }
        {
            dashboardType === 'provider' && <DashboardApplicationList />
        }
        </>
    )
}