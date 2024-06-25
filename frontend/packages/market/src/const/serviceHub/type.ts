/*
 * @Date: 2024-02-27 11:03:59
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-06 15:40:56
 * @FilePath: \frontend\packages\market\src\const\serviceHub\type.ts
 */
import { DefaultOptionType } from "antd/es/select"
import { ApiDetail } from "@common/const/api-detail"
import { EntityItem } from "@common/const/type"
import { SubscribeEnum, SubscribeFromEnum } from "@core/const/system/const"
import WithPermission from "@common/components/aoplatform/WithPermission"

export type ServiceBasicInfoType = {
    organization:EntityItem
    project:EntityItem
    team:EntityItem
    master:EntityItem
    apiNum:number
}

export type ServiceDetailType = {
    name:string
    description:string
    basic:ServiceBasicInfoType
    apis:ApiDetail[]
    applied:boolean
    partition:Array<{id:string, name:string}>
}

export type ServiceHubCategoryConfigFieldType = {
    id?:string
    name:string
    parent?:string
};

export type ServiceHubCategoryConfigProps = {
    type:'addCate'|'addChildCate'|'renameCate'
    entity?:{[k:string]:unknown}
    WithPermission: typeof WithPermission
}

export type ServiceHubCategoryConfigHandle = {
    save:()=>Promise<boolean|string>
}


export type CategorizesType = {
    id:string
    name:string
    children:CategorizesType[]
}

export type TagType = {
    id:string
    name:string
}


export type ServiceHubTableListItem = {
    id:string;
    name: string;
    tags?:EntityItem[];
    catalogue:EntityItem
    apiNum:number
    subscribeNum:number
    description:string
    logo:string
};


export type ApplyServiceProps = {
    entity:ServiceHubTableListItem & {partition:EntityItem[]} & {project:EntityItem}
    mySystemOptionList:DefaultOptionType[]
    reApply?:boolean
}

export type ApplyServiceHandle = {
    apply:()=>Promise<boolean|string>
}


export type ServiceHubApplyModalFieldType = {
    partitions?:string[],
    projects?:string;
    reason?:string;
};

export type ServiceHubAppListItem = {
    id:string,
    name:string,
    team:EntityItem,
    organization:EntityItem,
    subscribeNum:number,
    subscribeVerifyNum:number,
    description:string,
    master:EntityItem,
    createTime:string,
}

export type TenantManagementServiceListItem = {
    id:string
    service:EntityItem
    partition:EntityItem
    applyStatus:SubscribeEnum
    project:EntityItem
    team:EntityItem
    from:SubscribeFromEnum
    createTime:string
}