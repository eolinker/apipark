/*
 * @Date: 2024-02-27 11:03:54
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-06 17:27:37
 * @FilePath: \frontend\packages\market\src\const\serviceHub\const.tsx
 */

import { ProColumns } from "@ant-design/pro-components";
import { MenuProps } from "antd";
import { getItem } from "@common/utils/navigation";
import { ServiceHubTableListItem } from "./type";
import { ApiOutlined, KeyOutlined } from "@ant-design/icons";

export const SERVICE_HUB_TABLE_COLUMNS: ProColumns<ServiceHubTableListItem>[] = [
    {
        title: '服务名称',
        dataIndex: 'name',
        copyable: true,
        ellipsis:true,
        width:160,
        fixed:'left',
        sorter: (a:ServiceHubTableListItem,b:ServiceHubTableListItem)=> {
            return a.name.localeCompare(b.name)
        },
    },
    {
        title: '服务ID',
        dataIndex: 'id',
        width: 140,
        copyable: true,
        ellipsis:true
    },
    {
        title: '服务标签',
        dataIndex: 'tags',
        ellipsis:true,
        renderText:(_,entity:ServiceHubTableListItem)=>entity.tags?.map(x=>x.name).join(',') || '-'
    },
    {
        title: '所属组织',
        dataIndex: ['organization','name'],
        copyable: true,
        ellipsis:true
    },
    {
        title: '所属系统',
        dataIndex: ['project','name'],
        copyable: true,
        ellipsis:true
    },
    {
        title: '所属团队',
        dataIndex: ['team','name'],
        copyable: true,
        ellipsis:true
    },
    {
        title: '服务分类',
        dataIndex: ['catalogue','name'],
        copyable: true,
        ellipsis:true
    },
    {
        title: '可用环境',
        dataIndex: 'partition',
        ellipsis:true,
        renderText:(_,entity:ServiceHubTableListItem)=>entity.partition?.map((x)=>x.name).join(',')
    }
];


export const TENANT_MANAGEMENT_APP_MENU: MenuProps['items'] = [
   
    getItem('访问授权', 'authorization',<KeyOutlined rotate={225}/>),
    getItem('应用管理', 'setting',<span><iconpark-icon  className="" name="auto-generate-api"></iconpark-icon></span>),
];
