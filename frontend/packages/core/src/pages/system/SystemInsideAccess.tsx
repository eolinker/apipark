/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 19:09:27
 * @FilePath: \frontend\packages\core\src\pages\system\SystemInsideAccess.tsx
 */
import {App, Collapse, Empty} from "antd";
import styles from "../team/Team.module.css";
import {DownOutlined, UpOutlined} from "@ant-design/icons";
import {AccessMemberList, AccessOptionItem, PermissionListItem} from "../access/AccessList.tsx";
import {useEffect, useRef, useState} from "react";
import {Link, useParams} from "react-router-dom";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {useFetch} from "@common/hooks/http.ts";
import {DefaultOptionType} from "antd/es/cascader";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import { PermissionOptionType } from "../access/AccessPage.tsx";
import { useBreadcrumb } from "@common/contexts/BreadcrumbContext.tsx";

export default function SystemInsideAccess(){
    const [accessList, setAccessList] = useState<Array<unknown>>([])
    const { message } = App.useApp()
    const {systemId} = useParams<RouterParams>();
    const {fetchData} = useFetch()
    const {setBreadcrumb} = useBreadcrumb()
    const [activeKey, setActiveKey] = useState<string[]>([])
    
    const requestState = useRef<{ keyword?: string, promise?: Promise<DefaultOptionType[]>}>({ });
      
    const getMemberOptionList = ( keyword?: string): Promise<DefaultOptionType[]> => {
      // 如果type或keyword与缓存中的相同，返回缓存的promise
      if (requestState.current.keyword === keyword && requestState.current.promise) {
        return requestState.current.promise!;
      }
  
      // 创建新的promise并更新引用
      const newPromise = new Promise((resolve, reject)=>{
        fetchData<BasicResponse<{ options: PermissionOptionType[] }>>('project/setting/permission/options',{method:'GET',eoParams:{project:systemId, keyword:keyword||""}}
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
        return fetchData<BasicResponse<null>>('project/setting/permission',{method:'POST',eoBody:({...value}), eoParams:{project:systemId}}).then(response=>{
            const {code,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功')
                    getAccessList()
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const removeMember :(entity:AccessOptionItem & {access:string})=>Promise<boolean> = (entity:AccessOptionItem& {access:string})=>{
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>('project/setting/permission',{method:'DELETE',eoParams:{access:entity.access,key:entity.key,project:systemId}}).then(response=>{
                const {code,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功')
                    getAccessList()
                    resolve(true)
                }else{
                    message.error(msg || '操作失败')
                    reject(false)
                }
            }).catch(()=>reject(false))
        })
    }

    const getAccessList = ()=>{
        fetchData<BasicResponse<{ permissions: PermissionListItem[] }>>('project/setting/permissions',{method:'GET',eoParams:{project:systemId}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setAccessList(data.permissions?.map((x:PermissionListItem)=>({
                    key:x.access, label: <><span className="font-bold my-btnybase mr-btnbase">{x.name}</span><span className="text-SECOND_TEXT mb-btnybase"> {x.description}</span></>,
                    children: <AccessMemberList  permissionStr="project.mySystem.access.edit" access={x.access} items = {x.grant} requestOption={getMemberOptionList} addNewItem={addMember} removeItem={removeMember}/>
                })))
                setActiveKey(data.permissions?.map((x:PermissionListItem)=>x.access))
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    useEffect(()=>{
        setBreadcrumb([
            {title: <Link to={`/system/list`}>内部数据服务</Link>},
            {title:'权限'}
        ])
    },[])

    useEffect(() => {
        getAccessList()
    }, [systemId]);

    return (
        <div className="p-btnbase h-full bg-MAIN_BG">
            <p className="text-[16px] font-bold leading-[26px] mb-btnbase">权限设置</p>
                <div className = "h-[calc(100%-50px)] overflow-y-auto">
                    { accessList?.length > 0 ? <>
                    <Collapse  
                        className={`${styles['collapse-without-padding']} p-[0px]`} 
                        items={accessList}
                        activeKey={activeKey}
                        onChange={(val)=>{setActiveKey(val as string[])}}
                        expandIcon={({isActive})=>(isActive? <UpOutlined className="w-[23px] text-MAIN_TEXT hover:text-MAIN_HOVER_TEXT"/> : <DownOutlined className="w-[23px] text-MAIN_TEXT hover:text-MAIN_HOVER_TEXT"/>)}/>
                    </>:<div className="block h-full align-middle"><Empty className="mt-[20%]" image={Empty.PRESENTED_IMAGE_SIMPLE}/></div>
                    }
                </div>
        </div>
    )
}

