/*
 * @Date: 2024-02-04 10:01:02
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-05-14 09:52:10
 * @FilePath: \frontend\packages\core\src\const\member\const.tsx
 */
import { ProColumns } from "@ant-design/pro-components";
import { MemberTableListItem } from "./type";


export const MEMBER_TABLE_COLUMNS: ProColumns<MemberTableListItem>[] = [
    {
        title: '用户名',
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
        title: '邮箱',
        dataIndex: 'email',
        copyable: true,
        ellipsis:true,
    },
    {
        title: '部门',
        dataIndex: 'department',
        copyable: true,
        ellipsis:true,
        filterMode:'tree',
        renderText:(_,entity:MemberTableListItem)=>(entity.department?.map(x=>x.name).join('，') || '-'),
        filters: [],
        onFilter: true,
        // valueType: 'select',
        filterSearch: true,
    },
    {
        title: '用户组',
        dataIndex: 'userGroup',
        copyable: true,
        ellipsis:true,
        renderText:(_,entity:MemberTableListItem)=>entity.userGroup?.map((x)=>x.name).join('，') || '-'
    },
    {
        title:'状态',
        dataIndex:'enable',
        valueType: 'select',
        filters: true,
        onFilter: true,
        valueEnum:new Map([
            [true,<span className="text-status_success">启用</span>],
            [false,<span className="text-status_fail">禁用</span>],
        ])
    }
];
