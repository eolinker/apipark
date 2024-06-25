/*
 * @Date: 2024-03-11 16:34:40
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 19:12:29
 * @FilePath: \frontend\packages\core\src\pages\system\SystemOutlet.tsx
 */
import { Outlet, useParams } from "react-router-dom"
import { RouterParams } from "@core/components/aoplatform/RenderRoutes"
import { useEffect } from "react"
import { useGlobalContext } from "@common/contexts/GlobalStateContext"

export default function SystemOutlet(){
    const {teamId, systemId} = useParams<RouterParams>()
    const {getTeamAccessData,cleanTeamAccessData, getProjectAccessData,cleanProjectAccessData} = useGlobalContext()

    useEffect(()=>{
        teamId ? getTeamAccessData(teamId) : cleanTeamAccessData()
        return ()=>{
            cleanTeamAccessData()
        }
    },[teamId])

    useEffect(()=>{
        systemId ? getProjectAccessData(systemId) : cleanProjectAccessData()
        return ()=>{
            cleanProjectAccessData()
        }
    },[systemId])

    return (<Outlet />)
}