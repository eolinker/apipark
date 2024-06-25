import PageList from "@common/components/aoplatform/PageList.tsx"
import {ActionType, ProColumns} from "@ant-design/pro-components";
import  {FC, ReactNode, useEffect, useMemo, useRef, useState} from "react";
import {Link, useNavigate, useParams} from "react-router-dom";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {App, Divider, Drawer} from "antd";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import { RouterParams } from "@core/components/aoplatform/RenderRoutes.tsx";
import { SimpleMemberItem, SubSubscribeApprovalModalHandle, SystemSubServiceTableListItem } from "../../../const/system/type.ts";
import { SYSTEM_SUBSERVICE_TABLE_COLUMNS } from "../../../const/system/const.tsx";
import { SubSubscribeApprovalModalContent } from "./SubSubscribeApprovalDetailModalContent.tsx";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";
import { PERMISSION_DEFINITION } from "@common/const/permissions.ts";
import { checkAccess } from "@common/utils/permission.ts";
import { SubscribeApprovalInfoType } from "@common/const/approval/type.tsx";
import SubServiceDetail from "./SubServiceDetail.tsx";

const SystemInsideSubService:FC = ()=>{
    const [searchWord, setSearchWord] = useState<string>('')
    const navigate = useNavigate();
    const { setBreadcrumb } = useBreadcrumb()
    const { modal,message } = App.useApp()
    // const [confirmLoading, setConfirmLoading] = useState(false);
    const pageListRef = useRef<ActionType>(null);
    const {fetchData} = useFetch()
    const [init, setInit] = useState<boolean>(true)
    const [tableListDataSource, setTableListDataSource] = useState<SystemSubServiceTableListItem[]>([]);
    const [tableHttpReload, setTableHttpReload] = useState(true);
    const {systemId,orgId, teamId} = useParams<RouterParams>();
    const subSubscribeRef = useRef<SubSubscribeApprovalModalHandle>(null)
    const [memberValueEnum, setMemberValueEnum] = useState<{[k:string]:{text:string}}>({})
    const {accessData} = useGlobalContext()
    const [drawerOpen, setDrawerOpen] = useState<boolean>(false)
    const [curSubService, setCurSubService] = useState<SystemSubServiceTableListItem>()
    const getSubServiceList = ()=>{

        if(!tableHttpReload){
            setTableHttpReload(true)
            return Promise.resolve({
                data: tableListDataSource,
                success: true,
            });
        }

        return fetchData<BasicResponse<{subscriptions:SystemSubServiceTableListItem[]}>>('project/subscriptions',{method:'GET',eoParams:{project:systemId,keyword:searchWord},eoTransformKeys:['apply_status','create_time']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setTableListDataSource(data.subscriptions)
                setInit((prev)=>prev ? false : prev)
                setTableHttpReload(false)
                return  {data:data.subscriptions, success: true}
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        }).catch(() => {
            return {data:[], success:false}
        })
    }

    const handlerService = (type:'reApply'|'delete'|'cancelSub'|'cancelApply',entity:SystemSubServiceTableListItem)=>{
        let url:string
        let method:string
        let eoParams:{[k:string]:string|number}
        switch (type){
            case 'reApply':
                url = 'project/subscription/cancel_apply'
                method = 'POST'
                break
            case 'delete':
                url = 'project/subscription'
                method = 'DELETE'
                eoParams = {subscription:entity.id!,project:systemId!}
                break
            case 'cancelSub':
                url = 'project/subscription/cancel'
                eoParams = {subscription:entity.id!,project:systemId!}
                method = 'POST'
                break
            case 'cancelApply':
                url = 'project/subscription/cancel/application'
                eoParams = {application:entity.id!,project:systemId!}
                method = 'POST'
                break
        }
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>(url,{method,eoParams}).then(response=>{
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

    
    const isActionAllowed = (type:'view'|'cancelSub'|'cancelApply'|'delete'|'reApply') => {
        const actionToPermissionMap = {
            'view': 'viewApproval',
            'cancelSub': 'cancelSubscribe',
            'cancelApply': 'cancelApply',
            'reApply': 'subscribe',
            'delete': 'delete'
        };
        
        const action = actionToPermissionMap[type];
        const permission = `project.mySystem.subservice.${action}` as keyof typeof PERMISSION_DEFINITION[0];
        
        return !checkAccess(permission, accessData);
        };

    
    const openDrawer = (entity:SystemSubServiceTableListItem)=>{
        setCurSubService(entity)
        setDrawerOpen(true)
    }

    const openModal = async (type:'view'|'cancelSub'|'cancelApply'|'delete'|'reApply',entity:SystemSubServiceTableListItem)=>{
        //console.log(type,entity)
        let title:string = ''
        let content:string|ReactNode = ''
        switch (type){
            case 'view':
            case 'reApply':{ 
                message.loading('正在加载数据')
                const {code,data,msg} = await fetchData<BasicResponse<{approval:SubscribeApprovalInfoType}>>('project/subscription/approval',{method:'GET',eoParams:{subscription:entity!.id, project:systemId},eoTransformKeys:['apply_project','apply_team','apply_time','approval_time']})
                message.destroy()
                if(code === STATUS_CODE.SUCCESS){
                    title=type === 'view' ? '审批详情':'重新申请'
                        content = <SubSubscribeApprovalModalContent ref={subSubscribeRef} data={data.approval} type={type} systemId={systemId}/>;
                }else{
                    message.error(msg || '操作失败')
                    return
                }
                break;
            }
            case 'cancelSub':
                title='取消订阅'
                content='请确认是否取消订阅？'
                break;
            case 'cancelApply':
                title='取消申请'
                content='请确认是否取消申请？'
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
                if(type === 'reApply'){
                    return subSubscribeRef.current?.reApply().then((res)=>{
                        if(res === true) manualReloadTable()
                    })
                }
                if(type !== 'view'){
                    return handlerService(type,entity).then((res)=>{if(res === true) manualReloadTable()})
                }
            },
            okText:type === 'reApply' ? '重新申请' :'确认',
            okButtonProps:{
                disabled : isActionAllowed(type)
            },
            cancelText:'取消',
            closable:true,
            icon:<></>,
            width:600
        })
    }
    const operation:ProColumns<SystemSubServiceTableListItem>[] =[
        {
            title: '操作',
            key: 'option',
            width: 194,
            fixed:'right',
            valueType: 'option',
            render: (_: React.ReactNode, entity: SystemSubServiceTableListItem) => [
                <TableBtnWithPermission  access="project.mySystem.subservice.view" key="view"  onClick={()=>{openDrawer(entity)}} btnTitle="查看"/>,
                ...(entity.applyStatus === 2 ? [
                    <Divider type="vertical" className="mx-0"  key="div1" />,
                    <TableBtnWithPermission  access="project.mySystem.subservice.viewApproval" key="detail"  onClick={()=>{openModal('view',entity)}} btnTitle="审批详情"/>,
                    <Divider type="vertical" className="mx-0"  key="div2"/>,
                    <TableBtnWithPermission  access="project.mySystem.subservice.cancelSubscribe" key="cancelSub" onClick={()=>{openModal('cancelSub',entity)}} btnTitle="取消订阅"/>
                ]:[]),
                ...( entity.applyStatus === 1? [
                    <Divider type="vertical" className="mx-0"  key="div3"/>,
                    <TableBtnWithPermission  access="project.mySystem.subservice.cancelApply" key="cancelApply" onClick={()=>{openModal('cancelApply',entity)}} btnTitle="取消申请"/>
                ]:[]),
                ...( entity.applyStatus !== 1 &&  entity.applyStatus !== 2? [
                    <TableBtnWithPermission  access="project.mySystem.subservice.subscribe"  key="reApply" onClick={()=>{openModal('reApply',entity)}} btnTitle="重新申请"/>,
                    <Divider type="vertical" className="mx-0"  key="div4"/>,
                    <TableBtnWithPermission  access="project.mySystem.subservice.delete" key="delete"  onClick={()=>{openModal('delete',entity)}} btnTitle="删除"/>]:[]),
            ],
        }
    ]

    useEffect(() => {
        setBreadcrumb([
            {
                title:<Link to={`/system/list`}>内部数据服务</Link>
            },
            {
                title:'使用的服务列表'
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
        return SYSTEM_SUBSERVICE_TABLE_COLUMNS.map(x=>{if(x.filters &&(x.dataIndex as string[])?.indexOf('applier') !== -1){x.valueEnum = memberValueEnum} return x})
    },[memberValueEnum])

    return (
        <>
        <PageList
            id="global_system_subService"
            ref={pageListRef}
            request={()=>getSubServiceList()}
            dataSource={tableListDataSource}
            columns = {[...columns,...operation]}
            addNewBtnTitle="订阅第三方服务"
            addNewBtnAccess="project.mySystem.subservice.subscribe"
            searchPlaceholder="输入名称、ID、所属服务、负责人查找服务"
            onAddNewBtnClick={()=>{navigate(`/serviceHub/list?callbackUrl=${location.pathname}${location.search}`)}}
            tableClickAccess="project.mySystem.subservice.view"
            onSearchWordChange={(e)=>{setSearchWord(e.target.value)}}
            onChange={() => {
                setTableHttpReload(false)
            }}
            onRowClick={(row:SystemSubServiceTableListItem)=>openDrawer(row)}
            />
            <Drawer 
              destroyOnClose={true} 
              maskClosable={false}
              title={curSubService?.service.name || '-'}
              width={'60%'}
              onClose={()=>{setDrawerOpen(false); setCurSubService(undefined)}}
              open={drawerOpen}
              footer={null}>
                <SubServiceDetail serviceId={curSubService?.service?.id}/>
            </Drawer>
            </>
    )

}
export default SystemInsideSubService