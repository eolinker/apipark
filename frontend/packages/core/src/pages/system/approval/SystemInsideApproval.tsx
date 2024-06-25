/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-04-19 16:56:54
 * @FilePath: \frontend\packages\core\src\pages\system\approval\SystemInsideApproval.tsx
 */
import {Tabs} from "antd";
import {Outlet, useLocation, useNavigate} from "react-router-dom";
import './SystemInsideApproval.module.css'
import  {FC, useEffect, useState} from "react";
import { SYSTEM_INSIDE_APPROVAL_TAB_ITEMS } from "../../../const/system/const";


const SystemInsideApproval:FC = ()=>{
    const navigateTo = useNavigate()
    const location = useLocation()
    const query =new URLSearchParams(useLocation().search)
    const currentUrl = location.pathname
    const [pageStatus,setPageStatus] = useState<0|1>(Number(query.get('status') ||0) as 0|1)
    const onChange = (key: string) => {
        setPageStatus(Number(key) as 0|1)
        navigateTo(`${currentUrl}?status=${key}`);
    };

    useEffect(() => {
        setPageStatus(Number(query.get('status') ||0) as 0|1)
    }, [currentUrl]);

    return (
        <>
        <Tabs defaultActiveKey={pageStatus.toString()} size="small" className="h-auto  bg-MAIN_BG" tabBarStyle={{paddingLeft:'10px',marginTop:'0px',marginBottom:'0px'}} tabBarGutter={20} items={SYSTEM_INSIDE_APPROVAL_TAB_ITEMS} onChange={onChange} destroyInactiveTabPane={true}/>
        <Outlet />
        </>
    )
}

export default SystemInsideApproval