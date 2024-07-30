import PageList from "@common/components/aoplatform/PageList.tsx"
import {ActionType} from "@ant-design/pro-components";
import  {FC, useEffect, useMemo, useRef, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";
import { App} from "antd";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import { SimpleTeamItem ,SimpleMemberItem} from "@common/const/type.ts";
import { SystemConfigHandle, SystemTableListItem } from "../../const/system/type.ts";
import { SYSTEM_TABLE_COLUMNS } from "../../const/system/const.tsx";
import { DrawerWithFooter } from "@common/components/aoplatform/DrawerWithFooter.tsx";
import SystemConfig from "./SystemConfig.tsx";

const SystemList:FC = ()=>{
    const {orgId }  = useParams<RouterParams>();
    const navigate = useNavigate();
    const [tableSearchWord, setTableSearchWord] = useState<string>('')
    const { setBreadcrumb } = useBreadcrumb()
    const [teamList, setTeamList] = useState<{ [k: string]: { text: string; }; }>()
    const {fetchData} = useFetch()
    const [init, setInit] = useState<boolean>(true)
    const [tableListDataSource, setTableListDataSource] = useState<SystemTableListItem[]>([]);
    const [tableHttpReload, setTableHttpReload] = useState(true);
    const { message } = App.useApp()
    const pageListRef = useRef<ActionType>(null);
    const [loading, setLoading] = useState<boolean>(true)
    const [memberValueEnum, setMemberValueEnum] = useState<{[k:string]:{text:string}}>({})
    const [open, setOpen] = useState(false);
    const drawerFormRef = useRef<SystemConfigHandle>(null)
    // const operation:ProColumns<SystemTableListItem>[] =[
    //     {
    //         title: '操作',
    //         key: 'option',
    //         width: 56,
    //         fixed:'right',
    //         valueType: 'option',
    //         render: (_: React.ReactNode, entity: SystemTableListItem) => [
    //             <TableBtnWithPermission  access="" key="view" navigateTo={`/system/${entity.organization.id}/${entity.team.id}/inside/${entity.id}`} btnTitle="查看"/>,
    //             // <TableBtnWithPermission  access="team.mySystem.self.view" key="view" navigateTo={`/system/${entity.organization.id}/${entity.team.id}/inside/${entity.id}`} btnTitle="查看"/>,
    //             // <TableBtnWithPermission  access="" key="publish" navigateTo={`/system/${entity.organization.id}/${teamId}/inside/${entity.id}/publish`} btnTitle="发布"/>
    //         ],
    //     }
    // ]

    const getSystemList = ()=>{
        if(!tableHttpReload){
            setTableHttpReload(true)
            return Promise.resolve({
                data: tableListDataSource,
                success: true,
            });
        }
        return fetchData<BasicResponse<{projects:SystemTableListItem[]}>>('my_projects',{method:'GET',eoParams:{keyword:tableSearchWord},eoTransformKeys:['api_num','service_num','create_time']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setTableListDataSource(data.projects)
                setInit((prev)=>prev ? false : prev)
                setTableHttpReload(false)
                return  {data:data.projects, success: true}
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        }).catch(() => {
            return {data:[], success:false}
        })
    }

    const getTeamsList = ()=>{
        fetchData<BasicResponse<{teams:SimpleTeamItem[]}>>('simple/teams/mine',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            setTeamList(data.teams)
            if(code === STATUS_CODE.SUCCESS){
                    setLoading(false)
                    const tmpValueEnum:{[k:string]:{text:string}} = {}
                    data.teams?.forEach((x:SimpleMemberItem)=>{
                        tmpValueEnum[x.name] = {text:x.name}
                    })
                    setTeamList(tmpValueEnum)
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        })
    }

    const manualReloadTable = () => {
        setTableHttpReload(true); // 表格数据需要从后端接口获取
        pageListRef.current?.reload()
    };

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
            getTeamsList();
            getMemberList()
            setBreadcrumb([
                {
                    title: '内部数据服务'
                }])
    }, [orgId]);

    // useEffect(() => {
    //     if(!init){
    //         manualReloadTable()
    //     }
    // }, []);

    const onClose = () => {
        setOpen(false);
      };
    
    const columns = useMemo(()=>{
        const res =  SYSTEM_TABLE_COLUMNS.map(x=>{
            if(x.filters &&((x.dataIndex as string[])?.indexOf('master') !== -1 ) ){
                x.valueEnum = memberValueEnum
            } 
            if(x.filters &&((x.dataIndex as string[])?.indexOf('team') !== -1 ) ){
                x.valueEnum = teamList
            } 
            
            return x})
            return res
    },[memberValueEnum,teamList])

    return (
        // <Skeleton  className='m-btnbase w-[calc(100%-20px)]' loading={loading} active>
          <div className="h-full w-full">
            <PageList
                id="global_system"
                ref={pageListRef}
                columns={[...columns]}
                request={()=>getSystemList()}
                addNewBtnTitle="添加服务"
                searchPlaceholder="输入名称、ID、所属团队、负责人查找服务"
                onAddNewBtnClick={() => {
                    setOpen(true) 
                }}
                // addNewBtnAccess="team.mySystem.self.add"
                // tableClickAccess="team.mySystem.self.view"
                onChange={() => {
                    setTableHttpReload(false)
                }}
                onSearchWordChange={(e) => {
                    setTableSearchWord(e.target.value)
                }}
                onRowClick={(row:SystemTableListItem)=>navigate(`/system/${row.organization.id}/${row.team.id}/inside/${row.id}`)}
                />
                    <DrawerWithFooter title="添加服务" open={open} onClose={onClose} onSubmit={()=>drawerFormRef.current?.save()?.then((res)=>{res && manualReloadTable();return res})} >
                        <SystemConfig ref={drawerFormRef} />
                    </DrawerWithFooter>
                </div>
            // </Skeleton>
    )

}
export default SystemList