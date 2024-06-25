import PageList from "@common/components/aoplatform/PageList.tsx"
import {ActionType, ProColumns} from "@ant-design/pro-components";
import  {FC, useEffect, useMemo, useRef, useState} from "react";
import {Link, useParams} from "react-router-dom";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import { App, Divider } from "antd";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {SystemAuthorityView} from "./SystemAuthorityView.tsx";
import {
    SystemAuthorityConfig,
} from "./SystemAuthorityConfig.tsx";
import { SYSTEM_AUTHORITY_TABLE_COLUMNS } from "../../../const/system/const.tsx";
import { SystemAuthorityTableListItem, SystemAuthorityConfigHandle, EditAuthFieldType, SimpleMemberItem } from "../../../const/system/type.ts";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import { checkAccess } from "@common/utils/permission.ts";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";



const SystemInsideAuthority:FC = ()=>{
    const { setBreadcrumb } = useBreadcrumb()
    const { modal,message } = App.useApp()
    const {fetchData} = useFetch()
    const [init, setInit] = useState<boolean>(true)
    const [tableListDataSource, setTableListDataSource] = useState<SystemAuthorityTableListItem[]>([]);
    const [tableHttpReload, setTableHttpReload] = useState(true);
    const {systemId} = useParams<RouterParams>();
    const addRef = useRef<SystemAuthorityConfigHandle>(null)
    const editRef = useRef<SystemAuthorityConfigHandle>(null)
    const pageListRef = useRef<ActionType>(null);
    const [memberValueEnum, setMemberValueEnum] = useState<{[k:string]:{text:string}}>({})
    const {accessData} = useGlobalContext()
    const getSystemAuthority = ()=>{
        if(!tableHttpReload){
            setTableHttpReload(true)
            return Promise.resolve({
                data: tableListDataSource,
                success: true,
            });
        }
        return fetchData<BasicResponse<{authorizations:SystemAuthorityTableListItem[]}>>('project/authorizations',{method:'GET',eoParams:{project:systemId},eoTransformKeys:['hide_credential','create_time','update_time','expire_time']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setTableListDataSource(data.authorizations)
                setInit((prev)=>prev ? false : prev)
                setTableHttpReload(false)
                return  {data:data.authorizations, success: true}
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        }).catch(() => {
            return {data:[], success:false}
        })
    }

    const manualReloadTable = () => {
        setTableHttpReload(true); // 表格数据需要从后端接口获取
        pageListRef.current?.reload()
    };

    const deleteAuthority = (entity:SystemAuthorityTableListItem)=>{
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>('project/authorization',{method:'DELETE',eoParams:{authorization:entity!.id,project:systemId}}).then(response=>{
                const {code,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功！')
                    resolve(true)
                }else{
                    message.error(msg || '操作失败')
                    reject(msg || '操作失败')
                }
            }).catch((errorInfo)=> reject(errorInfo))
        })
    }

    const openModal =async (type:'view'|'delete'|'add'|'edit',entity?:unknown)=>{
        //console.log(type,entity)
        let title:string = ''
        let content:string|React.ReactNode = ''
        switch (type){
            case 'view':{
                title='鉴权详情'
                message.loading('正在加载数据')
                const {code,data,msg} = await fetchData<BasicResponse<{details:{[k:string]:string}}>>('project/authorization/details',{method:'GET',eoParams:{authorization:entity!.id,project:systemId}})
                message.destroy()
                if(code === STATUS_CODE.SUCCESS){
                    content=<SystemAuthorityView entity={data.details}/>
                }else{
                    message.error(msg || '操作失败')
                    return
                }}
                break;
            case 'add':
                title='添加鉴权'
                content=<SystemAuthorityConfig ref={addRef} type={type} systemId={systemId!}/>
                break;
            case 'edit':{
                title='编辑鉴权'
                message.loading('正在加载数据')
                const {code,data,msg} = await fetchData<BasicResponse<{authorization:EditAuthFieldType}>>('project/authorization',{method:'GET',eoParams:{authorization:entity!.id,project:systemId},eoTransformKeys:['hide_credential','token_name','expire_time','user_name','public_key','user_path','claims_to_verify','signature_is_base64']})
                message.destroy()
                if(code === STATUS_CODE.SUCCESS){
                    content=<SystemAuthorityConfig ref={editRef} type={type} data={data.authorization} systemId={systemId!} />
                }else{
                    message.error(msg || '操作失败')
                    return
                }}
                break;
            case 'delete':
                title='删除'
                content='该数据删除后将无法找回，请确认是否删除？'
                break;
        }

        modal.confirm({
            title,
            content,
            onOk:()=>{
                switch (type){
                    case 'add':
                        return addRef.current?.save().then((res)=>{if(res === true) manualReloadTable()})
                    case 'edit':
                        return editRef.current?.save().then((res)=>{if(res === true) manualReloadTable()})
                    case 'delete':
                        return deleteAuthority(entity).then((res)=>{if(res === true) manualReloadTable()})
                    case 'view':
                        return true
                }
            },
            width:600,
            okText:'确认',
            okButtonProps:{
                disabled : !checkAccess( `project.mySystem.auth.${type}`, accessData)
            },
            cancelText:'取消',
            closable:true,
            icon:<></>,
        })
    }

    const operation:ProColumns<SystemAuthorityTableListItem>[] =[
        {
            title: '操作',
            key: 'option',
            width: 138,
            fixed:'right',
            valueType: 'option',
            render: (_: React.ReactNode, entity: SystemAuthorityTableListItem) => [
           <TableBtnWithPermission  access="project.mySystem.auth.view" key="view" onClick={()=>{openModal('view',entity)}} btnTitle="查看"/>,
            <Divider type="vertical" className="mx-0"  key="div1" />,
            <TableBtnWithPermission  access="project.mySystem.auth.edit" key="edit" onClick={()=>{openModal('edit',entity)}} btnTitle="编辑"/>,
            <Divider type="vertical" className="mx-0" key="div2" />,
            <TableBtnWithPermission  access="project.mySystem.auth.delete" key="delete" onClick={()=>{openModal('delete',entity)}} btnTitle="删除"/>,
            ],
        }
    ]

    
    const getMemberList = async ()=>{
        setMemberValueEnum({})
        const {code,data,msg}  = await fetchData<BasicResponse<{ members: SimpleMemberItem[] }>>('simple/member',{method:'GET'})
        if(code === STATUS_CODE.SUCCESS){
            const tmpValueEnum:{[k:string]:{text:string}} = {}
            data.members?.forEach((x:SimpleMemberItem)=>{
                tmpValueEnum[x.name] = {text:x.name}
            })
            setMemberValueEnum(tmpValueEnum)
        }else{
            message.error(msg || '操作失败')
        }
    }

    useEffect(() => {
        setBreadcrumb([
            {
                title:<Link to={`/system/list`}>内部数据服务</Link>
            },
            {
                title:'身份认证'
            }
        ])
        getMemberList()
        manualReloadTable()
    }, [systemId]);

    const columns = useMemo(()=>{
        return SYSTEM_AUTHORITY_TABLE_COLUMNS.map(x=>{if(x.filters &&((x.dataIndex as string[])?.indexOf('creator') !== -1 ) ){x.valueEnum = memberValueEnum} return x})
    },[memberValueEnum])

    return (
        <PageList
            id="global_system_auth"
            ref={pageListRef}
            columns = {[...columns,...operation]}
            request={()=>getSystemAuthority()}
            dataSource={tableListDataSource}
            showPagination={false}
            addNewBtnTitle="添加鉴权"
            onAddNewBtnClick={()=>{openModal('add')}}
            addNewBtnAccess="project.mySystem.auth.add"
            onRowClick={(row:SystemAuthorityTableListItem)=>openModal('view',row)}
            tableClickAccess="project.mySystem.auth.view"
        />
    )
}

export default SystemInsideAuthority