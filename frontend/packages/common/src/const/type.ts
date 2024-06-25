import { PERMISSION_DEFINITION } from "./permissions"
import { MatchPositionEnum, MatchTypeEnum } from "@core/const/system/const"

/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 11:23:49
 * @FilePath: \frontend\packages\common\src\const\type.ts
 */
export type UserInfoType = {
    username: string
    nickname: string
    email: string
    phone: string
    avatar: string
}

export type UserProfileProps = {
    entity?:UserInfoType
}

export type UserProfileHandle = {
    save:()=>Promise<boolean|string>
}

export type ClusterSimpleOption = {
    id:string
    name:string
    description:string
}


export type ClusterEnumData = {
    name:string,
    uuid:string,
    title:string
}

export interface ClusterEnum{
    clusters:Array<ClusterEnumData>
    name:string
}

export type TeamSimpleMemberItem = {
    user:EntityItem
    mail:string
    department:EntityItem
}

export type MemberItem = {
    id:string;
    name:string;
    email:string;
    department:Array<{id:string,name:string}>
}

export type DashboardPartitionItem = {
    id:string;
    name:string
    enableMonitor:boolean
}

export type PartitionItem = {
    id:string;
    name:string
}

export type OrganizationItem = {
    id:string
    name:string
    description:string
}

export type SimpleTeamItem = {
    id:string
    name:string
    description:string
    organization:EntityItem
    appNum:number
}

export type MatchItem = {
    position:MatchPositionEnum
    matchType:MatchTypeEnum
    key:string
    pattern:string
    id?:string
}

export type EntityItem = {
    id:string
    name:string
}

export type DynamicMenuItem = {
    name:string
    title:string
    path:string
}

export type AccessDataType = keyof typeof PERMISSION_DEFINITION[0]
