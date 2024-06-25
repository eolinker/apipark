/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 19:11:43
 * @FilePath: \frontend\packages\core\src\pages\team\TeamInsideAccess.tsx
 */
import {App, Collapse, Empty} from "antd";
import {AccessMemberList, AccessOptionItem, PermissionListItem} from "../access/AccessList.tsx";
import styles from "./Team.module.css";
import {DownOutlined, UpOutlined} from "@ant-design/icons";
import {useEffect, useRef, useState} from "react";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {Link, useParams} from "react-router-dom";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {DefaultOptionType} from "antd/es/cascader";
import { PermissionOptionType } from "../access/AccessPage.tsx";
import { useBreadcrumb } from "@common/contexts/BreadcrumbContext.tsx";

export default function TeamInsideAccess(){
    const [accessList, setAccessList] = useState<Array<unknown>>([])
    const { message } = App.useApp()
    const {teamId} = useParams<RouterParams>();
    const {fetchData} = useFetch()
    const { setBreadcrumb} = useBreadcrumb()
    const requestState = useRef<{ keyword?: string, promise?: Promise<DefaultOptionType[]>}>({ });
    const [activeKey, setActiveKey] = useState<string[]>([])
      
    const getMemberOptionList = ( keyword?: string): Promise<DefaultOptionType[]> => {
      // 如果type或keyword与缓存中的相同，返回缓存的promise
      if (requestState.current.keyword === keyword && requestState.current.promise) {
        return requestState.current.promise!;
      }
  
      // 创建新的promise并更新引用
      const newPromise = new Promise((resolve, reject)=>{
        fetchData<BasicResponse<{ options: PermissionOptionType[] }>>(
            'team/setting/permission/options',
            { method: 'GET', eoParams:{team:teamId,keyword:keyword||""} }
          ).then((response) => {
            const { code, data, msg } = response;
            if (code === STATUS_CODE.SUCCESS) {
                resolve(data);
            } else {
              message.error(msg || '操作失败');
              reject(msg || '操作失败');
            }
          }).catch((error) => {
            // 处理错误
            reject(error)
            throw error;
          });
      }) 
  
      // 更新引用
      requestState.current = {keyword, promise: newPromise };
  
        // 返回新的promise
      return newPromise;
    };
  

    const addMember = (value:{[k:string]:unknown})=>{
        //console.log(value)
        return fetchData<BasicResponse<null>>('team/setting/permission',{method:'POST',eoBody:({...value}),eoParams:{team:teamId}}).then(response=>{
            const {code,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功')
                    getAccessList()
                    return Promise.resolve(true)
            }else{
                message.error(msg || '操作失败')
                return Promise.reject(msg)
            }
        }).catch(errInfo=>Promise.reject(errInfo))
    }

    const removeMember :(entity:AccessOptionItem & {access:string})=>Promise<boolean> = (entity:AccessOptionItem& {access:string})=>{
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>('team/setting/permission',{method:'DELETE',eoParams:{access:entity.access,key:entity.key,team:teamId}}).then(response=>{
                const {code,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功')
                    resolve(true)
                }else{
                    message.error(msg || '操作失败')
                    reject(false)
                }
            }).catch(()=>reject(false))
        })
    }

    const getAccessList = ()=>{
        fetchData<BasicResponse<{ permissions: PermissionListItem[] }>>('team/setting/permissions',{method:'GET',eoParams:{team:teamId}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setAccessList(data.permissions?.map((x:PermissionListItem)=>({
                    key:x.access, label:<><span className="font-bold my-btnybase mr-btnbase">{x.name}</span><span className="text-SECOND_TEXT mb-btnybase"> {x.description}</span></>,
                    children: <AccessMemberList permissionStr="team.myTeam.access.edit" access={x.access} items = {x.grant} requestOption={getMemberOptionList} addNewItem={addMember} removeItem={removeMember}/>
                })) || [])
                setActiveKey(data.permissions?.map((x:PermissionListItem)=>(x.access)))
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    useEffect(() => {
        setBreadcrumb([
            {title:<Link to="/team/list">团队</Link>},
            {title:'权限'}
        ])
        getAccessList()
    }, [teamId]);

    return (
        <div className="m-btnbase h-full">
            <p className="text-[16px] font-bold leading-[26px] mb-btnbase">权限设置</p>
            <div className = "h-[calc(100%-62px)] overflow-y-auto">
                { accessList && accessList.length > 0 ?
                    <>
                <Collapse  
                    className={`${styles['collapse-without-padding']} p-[0px] mb-btnybase`} 
                    items={accessList}  
                    activeKey={activeKey}
                    onChange={(val)=>{setActiveKey(val as string[])}}
                    expandIcon={({isActive})=>(isActive? <UpOutlined className="w-[23px]"/> : <DownOutlined className="w-[23px]"/>)}/></>
                :<div className="block h-full align-middle"><Empty className="mt-[20%]" image={Empty.PRESENTED_IMAGE_SIMPLE}/></div>
                }
            </div>
        </div>
    )
}