/*
 * @Date: 2024-04-19 15:22:46
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-05-14 11:01:34
 * @FilePath: \frontend\packages\core\src\const\organization\const.tsx
 */


import { ProColumns } from "@ant-design/pro-components";
import { OrganizationTableListItem } from "./type";

export const ORGANIZATION_TABLE_COLUMNS: ProColumns<OrganizationTableListItem>[] = [
    {
        title: '名称',
        dataIndex: 'name',
        copyable: true,
        ellipsis:true,
        width:160,
        fixed:'left',
        sorter: (a,b)=> {
            return a.name.localeCompare(b.name)
        },
    },
    {
        title: 'ID',
        dataIndex: 'id',
        width: 140,
        copyable: true,
        ellipsis:true
    },
    {
        title: '描述',
        dataIndex: 'description',
        copyable: true,
        ellipsis:true
    },
    {
        title: '分区权限',
        dataIndex: 'partition',
        ellipsis:true,
        renderText:(_,entity:OrganizationTableListItem)=>(entity.partition?.map(x=>x.name).join('，') || '-')
    },
    {
        title: '负责人',
        dataIndex: ['master','name'],
        ellipsis: true,
        width:108,
        filters: true,
        onFilter: true,
        valueType: 'select',
        filterSearch: true
    },
    {
        title: '创建时间',
        key: 'createTime',
        dataIndex: 'createTime',
        ellipsis: true,
        width:176,
        sorter: (a,b)=> {
            return a.createTime.localeCompare(b.createTime)
        },
    },
];