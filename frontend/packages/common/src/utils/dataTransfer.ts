/*
 * @Date: 2024-03-06 16:58:39
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 10:30:55
 * @FilePath: \frontend\packages\common\src\utils\dataTransfer.ts
 */
import { ColumnFilterItem } from 'antd/es/table/interface'
import {DepartmentListItem} from '@core/const/member/type'

export const handleDepartmentListToFilter:(departmentList:DepartmentListItem[])=>ColumnFilterItem[]   = (departmentList:DepartmentListItem[])=>{
    return departmentList?.map((x:DepartmentListItem)=>(
        {
            text:x.name,
            value:x.id,
            children:x.children ? handleDepartmentListToFilter(x.children):null
        }
    ))
}

export const getImgBase64 = (img: RcFile, callback: (url: string) => void) => {
    const reader = new FileReader();
    reader.addEventListener('load', () => callback(reader.result as string));
    reader.readAsDataURL(img);
};