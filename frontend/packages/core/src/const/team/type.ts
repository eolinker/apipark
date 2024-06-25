/*
 * @Date: 2024-02-04 11:09:15
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-04-24 18:17:56
 * @FilePath: \frontend\packages\core\src\const\team\type.ts
 */
import { EntityItem } from "@common/const/type";

export type TeamTableListItem = {
    id:string;
    name: string;
    description:string;
    systemNum:number;
    creator:EntityItem;
    createTime:string;
    canDelete:boolean
    organization:EntityItem
};

export type TeamConfigProps = {
    entity?:TeamConfigFieldType
}
export type TeamConfigHandle = {
    save:()=>Promise<boolean|string>
}

export type TeamConfigType = {
    name: string;
    id?: string;
    description: string;
    organization?:EntityItem;
    master:EntityItem;
    orgId?:string;
    canDelete:boolean
};

export type TeamConfigFieldType = {
    name: string;
    id?: string;
    description: string;
    organization?:string;
    master:string;
    orgId?:string
};

export type TeamMemberTableListItem = {
    userId:string;
    name: EntityItem;
    role:string;
    userGroup:EntityItem;
    attachTime:string;
    canDelete:boolean
};