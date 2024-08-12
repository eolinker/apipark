import PageList from "@common/components/aoplatform/PageList.tsx"
import {ActionType, ProColumns} from "@ant-design/pro-components";
import  {FC, useEffect, useMemo, useRef, useState} from "react";
import {Link, useParams} from "react-router-dom";
import {App, Divider, Drawer, Tabs, Select, Switch} from "antd";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {useSystemContext} from "../../../contexts/SystemContext.tsx";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import { SYSTEM_MYSERVICE_TABLE_COLUMNS } from "../../../const/system/const.tsx";
import { MyServiceInsideConfigHandle, MyServiceTableListItem } from "../../../const/system/type.ts";
import { EntityItem } from "@common/const/type.ts";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import { DrawerWithFooter } from "@common/components/aoplatform/DrawerWithFooter.tsx";
import MyServiceInsideConfig from "./MyServiceInsideConfig.tsx";
import MyServiceInsideApi from "./MyServiceInsideApi.tsx";
import MyServiceInsideDocument from "./MyServiceInsideDocument.tsx";
import type { TabsProps } from 'antd';


const SystemInsideMyService:FC = ()=>{
    const { message } = App.useApp()
    const {orgId, teamId,systemId} = useParams<RouterParams>()
    const [searchWord, setSearchWord] = useState<string>('')
    const { setBreadcrumb } = useBreadcrumb()
    const [init, setInit] = useState<boolean>(true)
    const [tableListDataSource, setTableListDataSource] = useState<MyServiceTableListItem[]>([]);
    const [originTableListDataSource, setOriginTableListDataSource] = useState<MyServiceTableListItem[]>([]);
    const [tableHttpReload, setTableHttpReload] = useState(true);
    const pageListRef = useRef<ActionType>(null);
    const {fetchData} = useFetch()
    const {partitionList } = useSystemContext()
    const [selectedPartition, setSelectedPartition] = useState<string[]>([])
    const drawerFormRef = useRef<MyServiceInsideConfigHandle>(null)
    const [open, setOpen] = useState(false);
    const [editDrawerOpen, setEditDrawerOpen] =useState<boolean>(false)
    const [curService, setCurService] = useState<MyServiceTableListItem|undefined>()
    const [switchLoading, setSwitchLoading] = useState<Set<string>>(new Set())
    const operation:ProColumns<MyServiceTableListItem>[] =[
        {
            title: '操作',
            key: 'option',
            width: 60,
            fixed:'right',
            valueType: 'option',
            render: (_: React.ReactNode, entity: MyServiceTableListItem) => [
                <TableBtnWithPermission  access="project.mySystem.service.view" key="api" onClick={()=>openDrawer('api',entity)} btnTitle="编辑"/>,
                // <Divider type="vertical" className="mx-0"  key="div1" />,
                // <TableBtnWithPermission  access="project.mySystem.service.view" key="detail" onClick={()=>openDrawer('detail',entity)} btnTitle="服务详情"/>,
                // <Divider type="vertical" className="mx-0"  key="div2"/>,
                // <TableBtnWithPermission  access="project.mySystem.service.view" key="setting" onClick={()=>openDrawer('setting',entity)} btnTitle="服务设置"/>,
            ],
        }
    ]
    
    const SYSTEM_MYSERVICE_EDIT_DRAWER_ITEM:TabsProps['items'] = [
        {label:'API 管理',key:'api',children:<MyServiceInsideApi systemId={systemId!} serviceId={curService?.id || ''}  />},
        {label:'服务详情',key:'detail',children:<MyServiceInsideDocument systemId={systemId!} serviceId={curService?.id || ''}  />},
        {label:'服务设置',key:'setting',children:<MyServiceInsideConfig teamId={teamId!} systemId={systemId!} serviceId={curService?.id || ''} closeDrawer={()=>{setCurService(undefined); setEditDrawerOpen(false);manualReloadTable()}}/>}
    ]

    const getMyServiceList =(): Promise<{ data: MyServiceTableListItem[], success: boolean }>=> {
        if(!tableHttpReload){
            setTableHttpReload(true)
            //console.log(selectedPartition,originTableListDataSource)
            const newTableListData = selectedPartition.length > 0 ? originTableListDataSource.filter((x:MyServiceTableListItem)=>x.partition.filter((p:EntityItem)=>selectedPartition.includes(p.id)).length >0 ):originTableListDataSource
            setTableListDataSource(newTableListData)
            return Promise.resolve({
                data: newTableListData,
                success: true,
            });
        }
        return fetchData<BasicResponse<{services:MyServiceTableListItem}>>('project/services',{method:'GET',eoParams:{project:systemId, keyword:searchWord},eoTransformKeys:['partition_id','service_type','api_num','create_time','update_time']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setOriginTableListDataSource(data.services)
                setTableListDataSource(selectedPartition.length > 0 ? data.services.filter((x:MyServiceTableListItem)=>x.partition.filter((x:EntityItem)=>selectedPartition.includes(x.id)).length > 0):data.services)
                setInit((prev)=>prev ? false : prev)
                setTableHttpReload(false)
                return  {data:data.services, success: true}
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

    const handlePartitionSelectedChange = (e:string[])=>{
        setSelectedPartition(e)
        setTableHttpReload(false)
        pageListRef.current?.reload()
    }

    useEffect(() => {
        setBreadcrumb([
            {
                title:<Link to={`/system/list`}>内部数据服务</Link>
            },
            {
                title:'提供的服务列表'
            }
        ])
        manualReloadTable()
    }, [systemId]);
    

    const openDrawer = (type:'add'|'api',entity?:MyServiceTableListItem)=>{
        if(type !== 'add'){
            setCurService(entity)
            setEditDrawerOpen(true)
        }else{
            setOpen(true)
        }
    }

    const newColumns = useMemo(() => {
        const columns = [...SYSTEM_MYSERVICE_TABLE_COLUMNS, ...operation];
      
        return columns.map(column => {
          if (column.dataIndex === 'status') {
            return {
              ...column,
              render: (_:unknown,entity:MyServiceTableListItem) =>
                <Switch size="small" checked={entity.status === 'on'}  loading={switchLoading.has(entity.id)} onChange={(checked)=>handleChangeServiceStatus(entity,checked)} onClick={(checked, e)=>e?.stopPropagation()} />
            };
          }
      
          return column;
        });
      }, []);

      const handleChangeServiceStatus = (entity:MyServiceTableListItem,checked:boolean)=>{
        setSwitchLoading(prev => {prev.add(entity.id);return prev})
        fetchData<BasicResponse<null>>(`project/service/${checked?'enable':'disable'}`,{method:'PUT',eoParams:{service:entity.id,project:systemId}}).then(response=>{
            const {code,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                message.success(msg || '操作成功')
                manualReloadTable()
            }else{
                message.error(msg || '操作失败')
            }
        }).finally(()=>{setSwitchLoading(prev => {prev.delete(entity.id); return prev})})
      }
    
    return (
        <>
            <PageList
                id="global_system_myService"
                ref={pageListRef}
                request={()=>getMyServiceList()}
                dataSource={tableListDataSource}
                columns = {newColumns}
                addNewBtnTitle="添加服务"
                addNewBtnAccess="project.mySystem.service.add"
                tableClickAccess="project.mySystem.service.view"
                beforeSearchNode={[ <Select key="zoneSelect"
                    mode="multiple"
                    allowClear
                    style={{ width: '100%' }}
                    placeholder="全部网络区域"
                    value={selectedPartition}
                    onChange={handlePartitionSelectedChange}
                    options={partitionList?.map((x)=>({label:x.name, value:x.id})) || []}
                />]}
                onChange={() => {
                    setTableHttpReload(false)
                }}
                searchPlaceholder="输入名称、ID、负责人查找服务"
                onAddNewBtnClick={()=>{
                    openDrawer('add')
                }}
                onSearchWordChange={(e)=>{setSearchWord(e.target.value)}}
                onRowClick={(row:MyServiceTableListItem)=>openDrawer('api',row)}
            />
            <DrawerWithFooter title="添加服务" open={open} onClose={()=>{setOpen(false);}} onSubmit={()=>drawerFormRef.current?.save()?.then((res)=>{res && manualReloadTable();return res})} >
                <MyServiceInsideConfig ref={drawerFormRef} systemId={systemId!} teamId={teamId!} />
            </DrawerWithFooter>
            <Drawer 
                className="full-tabs"
                destroyOnClose={true} 
                maskClosable={false}
                title={curService?.name || ''} 
                width={'60%'}
                onClose={()=>{setEditDrawerOpen(false);setTimeout(()=>setCurService(undefined),100); manualReloadTable()}}
                open={editDrawerOpen}>
                <>
                    <Tabs 
                    tabBarStyle={{paddingLeft:'20px'}}
                    className="h-full full-tabs" rootClassName="h-full" defaultActiveKey="1" items={SYSTEM_MYSERVICE_EDIT_DRAWER_ITEM} />
                </>
              </Drawer>
        </>
    )

}
export default SystemInsideMyService