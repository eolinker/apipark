/*
 * @Date: 2024-03-11 15:12:43
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-07-29 16:27:56
 * @FilePath: \frontend\packages\common\src\utils\permission.ts
 */
import { PERMISSION_DEFINITION } from "@common/const/permissions"
import { AccessDataType } from "@common/const/type"


export const checkAccess:(access:AccessDataType, accessData:Map<string,string[]>)=>boolean = (access, accessData)=>{
  // TODO 暂时隐藏权限
    return true
    if(!access){
      return true
    }
    const accLevel = access.split('.')[0]
    if(['system','team','project'].indexOf(accLevel) === -1){
        console.warn('权限字段有误：',access)
        return false
    }
    const neededBackendAccessArr = PERMISSION_DEFINITION[0]?.[access]?.granted.anyOf[0].backend || []
    return accessData?.has(accLevel)&& accessData.get(accLevel)!.length > 0 ? hasIntersection(neededBackendAccessArr, accessData.get(accLevel)!) : false
}

const hasIntersection = (arr1:string[], arr2:string[])=> {
    const set = new Set(arr1.length > arr2.length ? arr2:arr1) 
    const arr = arr1.length > arr2.length ? arr1:arr2
    for (const item of arr) {
      if (set.has(item)) {
        return true; // 发现交集
      }
    }
    return false; // 没有交集
  }