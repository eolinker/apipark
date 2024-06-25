import { App, Collapse, List, Select} from "antd";
import styles from './AccessList.module.css'
import { DownOutlined, UpOutlined} from "@ant-design/icons";
import {DefaultOptionType} from "antd/es/cascader";
import {useEffect, useRef, useState} from "react";
import {debounce} from "lodash-es";
import {useOutletContext, useParams} from "react-router-dom";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";
import { BasicResponse, STATUS_CODE } from "@common/const/const.ts";
import { useFetch } from "@common/hooks/http.ts";
import { PermissionOptionType } from "./AccessPage.tsx";
import TagWithPermission from "@common/components/aoplatform/TagWithPermission.tsx";

export type PermissionType = {
    system : PermissionListItem[]
    team : PermissionListItem[]
    project : PermissionListItem[]
}

export type PermissionListItem = {
    access:string
    name:string
    description:string
    grant:AccessOptionItem[]
    key?:string
    children:Array<{key:string,label:string,children:JSX.Element}>
}

export type AccessOptionItem = {
    key:string, label:string, name:string, tag:string, type:string
}

const AccessRoleOrder: {[key: string]: number} = {
    '特殊角色': 1,
    '用户组': 2,
    '用户': 3,
    '成员': 4,
  };

const AddNewMember = ({requestOption,addNewItem,selectedMember,permissionStr}:{permissionStr:string,requestOption: (keyword?:string)=>Promise<{options:AccessOptionItem[]}>,addNewItem:(member:DefaultOptionType)=>Promise<BasicResponse<unknown>>,selectedMember:string[]})=>{
    const [keyword,setKeyword] = useState<string>()
    const [oriOptions,setOriOptions] = useState<AccessOptionItem[]>()
    const [options,setOptions] = useState<DefaultOptionType[]>([])
    const [value, setValue] = useState<string>()
    const onSelect = async (member:DefaultOptionType)=>{
        setValue(member.value as string)
        if(addNewItem){
            const res = await addNewItem(member)
            res && setValue(undefined)
        }
    }

    
    const getGroup = (objArr: Array<AccessOptionItem>)=> {
        const map = new Map();
        for (const obj of objArr) {
          const tag = obj.tag;
          let options = map.get(tag);
          if (!options) {
            options = [];
            map.set(tag, options);
          }
          options.push({
            label: obj.label,
            value: obj.key,
            disabled:selectedMember.indexOf(obj.key) !== -1
          });
        }
      
        const result = Array.from(map.entries())?.map(([tag, options]) => ({
          label: tag,
          options,
        }));
        return result;
      }
    
    useEffect(()=>{
        const newOpt = oriOptions ? getGroup(oriOptions) : []
        setTimeout(()=>{setOptions(newOpt)})
    },[oriOptions,selectedMember])

    const getOptionList = (curKeyword?:string)=>{
            requestOption(curKeyword ?? keyword).then((res)=>{
                setOriOptions(res.options)
            })
        }

    useEffect(() => {
        getOptionList()
    }, []);
    return (<WithPermission access={permissionStr} ><Select showSearch  allowClear value={value}  filterOption={false} className=" ml-btnbase w-[270px]" options={options} placeholder="输入搜索更多成员" onSelect={(value:string,option:DefaultOptionType)=>onSelect(option)} onSearch={(value:string)=>{debounce((value)=>{setKeyword(value);getOptionList(value)},600)(value)}} onClear={()=>setValue('')}/></WithPermission>)
}

export const AccessMemberList = ({access,permissionStr, items,requestOption,addNewItem,removeItem}:{access:string,permissionStr:string, items:AccessOptionItem[],requestOption:(keyword?:string)=>Promise<DefaultOptionType[]>,addNewItem:(obj:{[k:string]:unknown})=>Promise<BasicResponse<unknown>>,removeItem:(entity:AccessOptionItem & { access: string; })=>Promise<boolean>})=>{
    const [dataSource,setDataSource] = useState<Array<AccessOptionItem[]>>()
    const transferItemsByGroup = (items:AccessOptionItem[])=>{
        const map = new Map();
        if(!items || items.length === 0) return []
        for (const obj of items) {
          const tag = obj.tag;
          let options = map.get(tag);
          if (!options) {
            options = [];
            map.set(tag, options);
          }
          options.push(obj);
        }
        const result = Array.from(map.entries()).sort((a,b)=>AccessRoleOrder[a[0]]-AccessRoleOrder[b[0]]).map(([_,options])=>options)
        return result
    }

    const handleAddNewItem = (member:DefaultOptionType)=>{
        return addNewItem({access,key:member.value})
    }

    const handleRemoveItem = (entity:AccessOptionItem)=>{
        removeItem({...entity,access})
    }

    useEffect(()=>{setDataSource(transferItemsByGroup(items))},[items])

    return ( <List
        header={<AddNewMember permissionStr={permissionStr} addNewItem={handleAddNewItem} requestOption={requestOption} selectedMember={items?.map(x=>x.key)}/>}
        dataSource={dataSource}
        size="small"
        renderItem={(item) => (
            <List.Item className="pb-0">
                <div>
                {item.map(x=>(
                        <TagWithPermission key={x.key} access={permissionStr} onClose={()=>handleRemoveItem(x)} className="mr-btnrbase">
                        <>
                        <span className="max-w-[150px] truncate">{x.label}</span>
                                <span className="ml-[6px] text-[#b3b3b3] px-[6px]">
                                    [{x.tag || '-'}]
                                </span>
                                </>
                        </TagWithPermission>
                ))}
                </div>
                {/* <div className="flex justify-between w-full">
                    <div className="flex items-center ml-btnbase">
                        <span className="max-w-[150px] truncate">{item.label}</span>
                        {item.key &&<span className="text-[#b3b3b3] max-w-[150px] truncate">({item.key})</span> }
                        <span className="ml-[6px] border-[1px] border-solid border-BORDER text-[#b3b3b3] rounded-SEARCH_RADIUS px-[6px]">
                            {item.tag || '-'}
                        </span>
                    </div>
                    <>
                        { <WithPermission access={permissionStr}><Button type="text"  className="h-[22px] border-none p-0 flex items-center bg-transparent " onClick={()=>handleRemoveItem(item)}> <CloseOutlined className="mr-btnbase text-[#8888800] hover:text-[#888] cursor-pointer" /></Button></WithPermission>}
                    </>
                </div> */}
            </List.Item>
        )}
    />)
}

export default function AccessList(){
    const { accessMap,getAccessList} = useOutletContext<{accessMap:{system:PermissionListItem[],team:PermissionListItem[],project:PermissionListItem[]},getAccessList:()=>void}>()
    const { accessType } = useParams<RouterParams>()
    const [activeKey, setActiveKey] = useState<string[]>([])
    const {message } = App.useApp()
    const {fetchData} = useFetch()

    useEffect(()=>{
        setActiveKey(accessMap[accessType!]?.map(x=>x.access))
    },[accessMap, accessType])

    const requestState = useRef<{type?: 'system' | 'project' | 'team', keyword?: string, promise?: Promise<DefaultOptionType[]>}>({ });
      
    const getMemberOptionList = (type: 'system' | 'project' | 'team', keyword?: string): Promise<DefaultOptionType[]> => {
      // 如果type或keyword与缓存中的相同，返回缓存的promise
      if (requestState.current.type === type && requestState.current.keyword === keyword) {
        return requestState.current.promise!;
      }
  
      // 创建新的promise并更新引用
      const newPromise = new Promise((resolve, reject)=>{
        fetchData<BasicResponse<{ options: PermissionOptionType[] }>>(
            `system/permission/options${type === 'system' ? '' : `/${type}`}`,
            { method: 'GET', eoParams: { keyword: keyword || "" } }
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
      requestState.current = { type, keyword, promise: newPromise };
  
        // 返回新的promise
      return newPromise;
    };
  

    const addMember = (type:'system'|'project'|'team',value:{[k:string]:unknown})=>{
        return fetchData<BasicResponse<null>>(`system/permission${type === 'system' ? '': `/${type}`}`,{method:'POST',eoBody:({...value})}).then(response=>{
            const {code,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                message.success(msg || '操作成功')
                getAccessList()
            }else{
                message.error(msg || '操作失败')
            }
            return response
        })
}

const removeMember :(type:'system'|'project'|'team', entity:AccessOptionItem & {access:string})=>Promise<boolean> = (type:'system'|'project'|'team', entity:AccessOptionItem & {access:string})=>{
    return new Promise((resolve, reject)=>{
        fetchData<BasicResponse<null>>(`system/permission${type === 'system' ? '': `/${type}`}`,{method:'DELETE',eoParams:{access:entity.access,key:entity.key}}).then(response=>{
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


    return (
        <div className="h-full overflow-auto p-btnbase">
            {/* <p className="text-[16px] font-bold leading-[26px] mb-btnbase">权限设置</p> */}
             <Collapse className={`${styles['collapse-without-padding']} p-[0px] mb-btnybase`} 
                                expandIcon={({isActive})=>(isActive?  <UpOutlined className="w-[23px] text-MAIN_TEXT hover:text-MAIN_HOVER_TEXT"/>:<DownOutlined className="w-[23px] text-MAIN_TEXT hover:text-MAIN_HOVER_TEXT"/> )}
                                items={accessMap[accessType as 'project'|'team'|'system']?.map(x=>{
                                    return {
                                        label:<><span className="font-bold my-btnybase mr-btnbase" id={`${x.access}`}>{x.name}</span><span className="text-SECOND_TEXT mb-btnybase"> {x.description}</span></>,
                                        key:x.access,
                                        children:<AccessMemberList  access={x.access} permissionStr="system.access.self.edit" items = {x.grant} requestOption={(keyword?: string | undefined)=>getMemberOptionList(accessType!,keyword)} addNewItem={(value:{[k:string]:unknown})=>addMember(accessType!,value)} removeItem={(entity:AccessOptionItem & { access: string; })=>removeMember(accessType!, entity)}/>
                                    }
                                })}
                                activeKey={activeKey}
                                onChange={(val)=>{setActiveKey(val as string[])}}
                    />
        </div>
    )
}