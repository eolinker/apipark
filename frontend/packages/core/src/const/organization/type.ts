/*
 * @Date: 2024-02-04 11:05:42
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-03-04 18:34:13
 * @FilePath: \frontend\packages\core\src\const\organization\type.ts
 */
import { EntityItem } from "@common/const/type";

export type OrganizationFieldType = {
    name?: string;
    id?: string;
    description?: string;
    master?:string;
    partitions?:Array<string>
    prefix?:string
};

export type OrganizationTableListItem = {
    id:string;
    name: string;
    description:string;
    master:EntityItem;
    partition:EntityItem[]
    createTime:string;
    updateTime:string;
    canDelete:boolean
};
