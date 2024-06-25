import {ActionType, ProColumns} from "@ant-design/pro-components";
import  {FC, useEffect, useMemo, useRef, useState} from "react";
import {Link, useParams} from "react-router-dom";
import PageList from "@common/components/aoplatform/PageList.tsx";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {App, Divider} from "antd";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import { SimpleMemberItem, SystemInsideUpstreamContentHandle, SystemUpstreamTableListItem } from "../../../const/system/type.ts";
import { SYSTEM_UPSTREAM_TABLE_COLUMNS } from "../../../const/system/const.tsx";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";
import { checkAccess } from "@common/utils/permission.ts";
import { DrawerWithFooter } from "@common/components/aoplatform/DrawerWithFooter.tsx";
import SystemInsideUpstreamContent from "./SystemInsideUpstreamContent.tsx";

const SystemInsideUpstreamList:FC = ()=>{
    const { modal,message } = App.useApp()
    const { setBreadcrumb } = useBreadcrumb()
    const pageListRef = useRef<ActionType>(null);
    const [tableListDataSource, setTableListDataSource] = useState<SystemUpstreamTableListItem[]>([]);
    const [tableHttpReload, setTableHttpReload] = useState(true);
    const {fetchData} = useFetch()
    const {orgId, teamId, systemId } = useParams<RouterParams>()
    const [memberValueEnum, setMemberValueEnum] = useState<{[k:string]:{text:string}}>({})
    const {accessData} = useGlobalContext()
    const [open, setOpen] = useState(false);
    const [curUpstreamId, setCurUpstreamId] = useState<string>()
    const drawerFormRef = useRef<SystemInsideUpstreamContentHandle>(null)

    const deleteUpstream = (entity:SystemUpstreamTableListItem)=>{
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>('project/upstream',{method:'DELETE',eoParams:{upstream:entity!.id,project:systemId}}).then(response=>{
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

    const manualReloadTable = () => {
        setTableHttpReload(true); // 表格数据需要从后端接口获取
        pageListRef.current?.reload()
    };

    const openModal = async (entity:SystemUpstreamTableListItem)=>{
        modal.confirm({
            title:'删除',
            content:'该数据删除后将无法找回，请确认是否删除？',
            onOk:()=> {
                return deleteUpstream(entity).then((res)=>{if(res === true) manualReloadTable()})
            },
            width:600,
            okText:'确认',
            okButtonProps:{
                disabled:!checkAccess('project.mySystem.upstream.delete',accessData)
            },
            cancelText:'取消',
            closable:true,
            icon:<></>
        })
    }

    const operation:ProColumns<SystemUpstreamTableListItem>[] =[
        {
            title: '操作',
            key: 'option',
            width: 93,
            fixed:'right',
            valueType: 'option',
            render: (_: React.ReactNode, entity: SystemUpstreamTableListItem) => [
                <TableBtnWithPermission  access="project.mySystem.upstream.edit" key="edit" onClick={()=>{openDrawer(entity.id)}} btnTitle="编辑"/>,
                entity.canDelete &&  <Divider type="vertical" className="mx-0"  key="div1"/>,
                entity.canDelete &&  <TableBtnWithPermission  access="project.mySystem.upstream.delete" key="delete"  onClick={()=>{openModal(entity)}} btnTitle="删除"/>,
            ],
        }
    ]

    const getUpstreamList =(): Promise<{ data: SystemUpstreamTableListItem[], success: boolean }>=> {
        if(!tableHttpReload){
            setTableHttpReload(true)
            return Promise.resolve({
                data: tableListDataSource,
                success: true,
            });
        }
        return fetchData<BasicResponse<{upstreams:SystemUpstreamTableListItem}>>('project/upstreams',{method:'GET',eoParams:{project:systemId},eoTransformKeys:['update_time','create_time','can_delete']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setTableListDataSource(data.upstreams)
                setTableHttpReload(false)
                return  {data:data.upstreams, success: true}
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        }).catch(() => {
            return {data:[], success:false}
        })
    }

    useEffect(() => {
        setBreadcrumb([
            {
                title:<Link to={`/system/list`}>内部数据服务</Link>
            },
            {
                title:'上游'
            }
        ])
        getMemberList()
        manualReloadTable()
    }, [systemId]);

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
    
    const columns = useMemo(()=>{
        return SYSTEM_UPSTREAM_TABLE_COLUMNS.map(x=>{if(x.filters &&((x.dataIndex as string[])?.indexOf('creator') !== -1 || (x.dataIndex as string[])?.indexOf('updater') !== -1) ){x.valueEnum = memberValueEnum} return x})
    },[memberValueEnum])

    const openDrawer = (id?:string)=>{
        id === undefined? setOpen(true) : setCurUpstreamId(id)
    }

    useEffect(()=>{
        curUpstreamId !== undefined &&  setOpen(true)
    },[curUpstreamId])

    const onClose = () => {
        setOpen(false);
        setCurUpstreamId(undefined)
      };

    return (
        <>
            <PageList
                id="global_system_upstream"
                ref={pageListRef}
                columns = {[...columns,...operation]}
                request={()=>getUpstreamList()}
                addNewBtnTitle="添加上游"
                addNewBtnAccess="project.mySystem.upstream.add"
                tableClickAccess="project.mySystem.upstream.edit"
                onChange={() => {
                    setTableHttpReload(false)
                }}
                onAddNewBtnClick={()=>{
                    openDrawer()
                }}
                onRowClick={(row:SystemUpstreamTableListItem)=>openDrawer(row.id)}
            />
            <DrawerWithFooter title={`${curUpstreamId === undefined ? '添加' : '编辑'}上游`} open={open} onClose={onClose} onSubmit={()=>drawerFormRef.current?.save()?.then((res)=>{res && manualReloadTable(); return res})} >
                <SystemInsideUpstreamContent ref={drawerFormRef}  />
            </DrawerWithFooter>
        </>
    )

}
export default SystemInsideUpstreamList