/*
 * @Date: 2024-02-04 10:25:41
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-03-06 17:00:40
 * @FilePath: \frontend\packages\core\src\const\role\const.tsx
 */
export const ROLE_TABLE_COLUMNS = [
    {
        title: '角色名称',
        dataIndex: 'name',
        copyable: true,
        ellipsis:true,
        width:160,
        fixed:'left',
        sorter: (a,b)=> {
            return a.name.localeCompare(b.name)
        },
    }

]