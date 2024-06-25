/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-02-07 18:47:17
 * @FilePath: \frontend\packages\core\src\pages\system\myService\MyServiceInsideMenu.tsx
 */
import  {FC, useEffect, useState} from "react";
import {Link, Outlet, useLocation} from "react-router-dom";
import {Menu, MenuProps} from "antd";
import { getItem } from "@common/utils/navigation.tsx";

const MyServiceInsideMenu:FC = ()=> {
    const [selectedKeys, setSelectedKey] = useState('')
    const location = useLocation()
    const currentLocation = location.pathname

    const items: MenuProps['items'] = [
        getItem('管理', 'grp', null,
            [getItem(<Link to={`./api`}>API</Link>, 'api'),
                getItem(<Link to={`./document`}>服务详情</Link>, 'document'),
                getItem(<Link to={`./setting`}>服务设置</Link>, 'setting')],
            'group'),
    ];

    const onMenuClick: MenuProps['onClick'] = ({key}) => {
        setSelectedKey(key)
    };

    useEffect(()=>{
        const currentKey = currentLocation.split('/')[currentLocation.split('/').length -1]
        setSelectedKey(currentKey)
    },[currentLocation])

    return (
        <>
                <div className="flex h-full">
                    <Menu
                        selectedKeys={[selectedKeys]}
                        onClick={onMenuClick}
                        style={{ width: 176 }}
                        mode="inline"
                        items={items}
                    />
                    <div className="w-[calc(100%-175px)]">
                        <Outlet />
                    </div>
                </div>
        </>
    )
}
export default MyServiceInsideMenu