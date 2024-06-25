/*
 * @Date: 2024-04-19 10:21:03
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-05-29 18:48:57
 * @FilePath: \frontend\packages\common\src\components\aoplatform\TagWithPermission.tsx
 */
import { Tag, TagProps } from "antd";
import { useState, useMemo, useEffect } from "react";
import { PERMISSION_DEFINITION } from "@common/const/permissions";
import { useGlobalContext } from "@common/contexts/GlobalStateContext";

export interface TagWithPermission extends TagProps{
    access?:string
}
export default function TagWithPermission(props:TagWithPermission){
    const {access,onClose} = props
    const [editAccess, setEditAccess] = useState<boolean>(access ? false:true)
    const {accessData,checkPermission} = useGlobalContext()
    const lastAccess = useMemo(()=>{
      if(!access) return true
      return checkPermission(access as keyof typeof PERMISSION_DEFINITION[0])
  },[access, accessData,checkPermission])

    useEffect(()=>{
        access ? setEditAccess(lastAccess) :  setEditAccess(true)
    },[lastAccess])
    
    const handleTagClose = (e: React.MouseEvent<HTMLElement>)=>{
        e.preventDefault();
        if(!editAccess) return
        onClose?.(e)
    }

    return  <Tag  
        closeIcon
        {...props}
        className={` rounded-SEARCH_RADIUS h-[32px] text-[14px] leading-[22px] py-[5px] px-btnbase bg-transparent mb-[8px] ${props.className}`}
        onClose={handleTagClose}>
    {props.children}
  </Tag>

}