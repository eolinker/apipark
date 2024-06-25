/*
 * @Date: 2024-02-04 10:01:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 11:18:49
 * @FilePath: \frontend\packages\core\src\const\member\type.ts
 */
import { EntityItem } from "@common/const/type"

export type DepartmentListItem = {
    id:string
    name:string
    number?:string
    children:DepartmentListItem[]
    departmentIds?:string[]
    key?:string
}

export type MemberTableListItem = {
    id:string;
    name: string;
    email:string;
    department:Array<EntityItem>;
    userGroup:Array<EntityItem>;
    enable:boolean
    departmentId:string
};

export type AddToDepartmentProps = {
    selectedUserIds:string[]
}

export type AddToDepartmentHandle = {
    save:()=>Promise<boolean|string>
}

export type MemberDropdownModalFieldType = {
    id?:string
    name:string
    parent?:string
    email?:string
    departmentIds?:string[]
};

export type MemberDropdownModalProps = {
    type:'addDep'|'addChild'|'addMember'|'editMember'|'rename'
    entity?:(MemberTableListItem & {departmentIds:string[]}) | ({id?:string, departmentIds?:string[],name?:string})
    selectedMemberGroupId?:string
}

export type MemberDropdownModalHandle = {
    save:()=>Promise<boolean|string>
}
