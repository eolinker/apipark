import PageList from "@common/components/aoplatform/PageList.tsx";
import {App, Button, Divider, Select, Spin} from "antd";
import  {useEffect, useRef, useState} from "react";
import {Link, useLocation, useParams} from "react-router-dom";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {ActionType, ParamsType, ProColumns} from "@ant-design/pro-components";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {DefaultOptionType} from "antd/es/cascader";
import {IntelligentPluginConfig, IntelligentPluginConfigHandle} from "./IntelligentPluginConfig.tsx";
import {IntelligentPluginPublish, IntelligentPluginPublishHandle} from "./IntelligentPluginPublish.tsx";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {EntityItem, PartitionItem} from "@common/const/type.ts";
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import { DrawerWithFooter } from "@common/components/aoplatform/DrawerWithFooter.tsx";
import { LoadingOutlined } from "@ant-design/icons";

 type DynamicTableField = {
    name: string,
    title: string,
    attr: string,
    enum: Array<string>
}

 type DynamicDriverData = {
    name:string, title:string
}

export type DynamicTableConfig = {
    basic:{
        id:string,
        name: string,
        title: string,
        drivers: Array<DynamicDriverData>,
        fields: Array<DynamicTableField>,
    }
    list: Array<DynamicTableItem>,
    total:number
}

export type DynamicRender = {
    render:unknown,
    basic:{
        id:string,
        name:string,
        title:string
    }
}

export type DynamicPublishCluster = {
    name:string,
    title:string,
    status:string,
    updater:EntityItem,
    update_time:string,
    checked?:boolean
}

export type DynamicPublishData = {
    id:string,
    name:string,
    title:string,
    description:string
    clusters:DynamicPublishCluster[]
}

export type DynamicTableItem = {[k:string]:unknown}

export const StatusColorClass = {
    "已发布":'text-[#03a9f4]',
    "待发布":'text-[#46BE11]',
    "未发布":'text-[#03a9f4]'
}
export default function IntelligentPluginList(){
    const { modal,message } = App.useApp()
    const [searchWord, setSearchWord] = useState<string>('')
    const { moduleId }  = useParams<RouterParams>();
    const [pluginName,setPluginName] = useState<string>('-')
    const [partitionOptions, setPartitionOption] = useState<DefaultOptionType[]>([])
    const { setBreadcrumb } = useBreadcrumb()
    // const [confirmLoading, setConfirmLoading] = useState(false);
    const [renderSchema ,setRenderSchema] = useState<{[k:string]:unknown}>({})
    const [tableStatusLoading, setTableStatusLoading] = useState<boolean>(true)
    // const [currentModalContentRef, setCurrentModalContentRef] = useState<unknown>()
    const drawerFormRef = useRef<IntelligentPluginConfigHandle>(null);
    const publishRef = useRef<IntelligentPluginPublishHandle>(null);
    const [driverOptions, setDriverOptions] = useState<DefaultOptionType[]>([])
    const [tableListDataSource, setTableListDataSource] = useState<DynamicTableItem[]>([]);
    const [tableListStatus, setTableListStatus] = useState<{[k:string]:{[key:string]:string}}|null>(null);

    const [tableHttpReload, setTableHttpReload] = useState(true);
    const [columns,setColumns] = useState<ProColumns<DynamicTableItem>[] >([])
    const {fetchData} = useFetch()
    const [partition, setCluster] = useState<string[]>([])
    const pageListRef = useRef<ActionType>(null);
    const [publishBtnLoading, setPublishBtnLoading] = useState<boolean>(false)
    const [curDetail,setCurDetail] = useState<{[k: string]: unknown;}|undefined>()
    const [drawerType, setDrawerType]  = useState<'add'|'edit'>('add')
    const [drawerOpen, setDrawerOpen] = useState<boolean>(false)
    const [drawerLoading, setDrawerLoading] = useState<boolean>(false)
    const [tableStatusError, setTableStatusError] = useState<boolean>(false)
    const location = useLocation().pathname

    const getIntelligentPluginTableList=(params:ParamsType & {
        pageSize?: number | undefined;
        current?: number | undefined;
        keyword?: string | undefined;
    },): Promise<{ data: DynamicTableItem[], success: boolean }>=> {
        if(!tableHttpReload){
            setTableHttpReload(true)
            return Promise.resolve({
                data: tableListDataSource,
                success: true,
            });
        }
        const query = {
            page:params.current,
            pageSize:params.pageSize,
            keyword:searchWord,
            partition:JSON.stringify(partition)
        }
        setTableListStatus({})
        setTableStatusLoading(true)
        return Promise.allSettled(
            [
                fetchData<BasicResponse<DynamicTableConfig>>(`dynamic/${moduleId}/list`,{method:'GET',eoParams:query,eoTransformKeys:['pageSize']}).then((res)=>{
                    if(res.code === STATUS_CODE.SUCCESS) getConfig(res.data) ;message.destroy(); return res}),
                fetchData<BasicResponse<{[k:string]:{[key:string]:string}}>>(`dynamic/${moduleId}/status`,{method:'GET',eoParams:query,eoTransformKeys:['pageSize']}).then((res)=>{
                    if(res.code === STATUS_CODE.SUCCESS) setTableListStatus(res.data) ; return res})
                ])
                .then(([resultA, 
                    resultB
                ])=>{
                    // 检查两个请求是否都成功
                    const isSuccessA = resultA.status === 'fulfilled' && resultA.value.code === STATUS_CODE.SUCCESS;
                    const isSuccessB = resultB.status === 'fulfilled' && resultB.value.code === STATUS_CODE.SUCCESS;
                    // 根据请求结果更新状态
                    if(!isSuccessB){
                        message.error(resultB?.value?.msg || '操作失败')
                    }
                    if (isSuccessA) {
                        setTableStatusLoading(false)
                        setTableStatusError(!isSuccessB)
                        const fullTableListData =isSuccessB ? resultA.value.data.list?.map((x:DynamicTableItem)=>{return {...x,...resultB.value.data.list[x.id]}}):resultA.value.data.list
                        setTableListDataSource(fullTableListData);
                        getTableConfig(resultA.value.data, false,!isSuccessB)
                        return { data: fullTableListData, success: true,total:resultA.value.data.total };
                    } else {
                        setTableListDataSource([]);
                        return { data: [], success: false };
                    }
        }).catch((e)=>{console.warn(e)})
    }

    const getConfig = (data:DynamicTableConfig)=>{
        const {basic,list } = data
        const {title,drivers} = basic

        setBreadcrumb([
            {title:location.includes('resourcesettings') ? '资源配置': '日志配置'},
            {
                title
            }
        ])

        setPluginName(title)
        if(!tableListStatus || !Object.keys(tableListStatus).length){
            setTableListDataSource(list)
            getTableConfig(data,true,false) // 获取列表配置
        }
        setDriverOptions(drivers?.map((driver:DynamicDriverData) => {
            return { label: driver.title, value: driver.name }
        }) || [])
    }

    const getTableConfig = (data:DynamicTableConfig,tableStatusLoading:boolean, tableStatusError:boolean)=>{
        const {basic,list } = data
        const {title,drivers,fields} = basic
        let statusColFlag:boolean = true
        const newColumn : ProColumns<DynamicTableItem>[] = fields?.filter(
            (x:DynamicTableField)=>{
                if(x.attr === 'status' && (tableStatusLoading || tableStatusError) && statusColFlag){
                    statusColFlag = false
                    return true
                }
                    return !(x.attr === 'status' && (tableStatusLoading || tableStatusError) && !statusColFlag)
            }
            ).map(
                (x,index)=>{
            // 当状态list还未返回时，页面显示一个状态为loading的table
                    if (x.attr === 'status' && (tableStatusLoading || tableStatusError)) {
                        return {
                            title:(tableStatusLoading || tableStatusError) ? "状态" : x.title,
                            dataIndex:x.name,
                            render: ()=>tableStatusLoading ? <LoadingOutlined spin /> : '-',
                            rowSpan:list.length || 1} as ProColumns<DynamicTableItem>
                    }

            return {
                title:x.title,
                fixed:index === 0 ? 'left' : false,
                width:index === 0 ? 150 : undefined,
                ...(x.enum?.length > 0 ?{
                    onFilter: (value: string, record) => record[x.name].indexOf(value) === 0,
                    filters:x.enum?.map((x:string)=>{return {text:x, value:x}}),
                    render:(dom, entity)=> {
                        return <span className={StatusColorClass[entity[x.name]]}>{(entity[x.name] as string)}</span>                        
                    },
                }:{}),
                dataIndex:x.name,
                ellipsis:true
            } as ProColumns<DynamicTableItem>
        })
        setColumns(newColumn)
    }

    const getRender = ()=>{
        return fetchData<BasicResponse<DynamicRender>>(`dynamic/${moduleId}/render`,{method:'GET'}).then((resp) => {
            if (resp.code === STATUS_CODE.SUCCESS) {
                setRenderSchema(resp.data.render)
                return Promise.resolve(resp.data.render)
            }
            return Promise.reject(resp.msg || '操作失败')
        })
    }

    
    const getPartitionList = ()=>{
        setPartitionOption([])
        return fetchData<BasicResponse<{partitions:PartitionItem[]}>>('simple/partitions',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setPartitionOption(data.partitions?.map((x:PartitionItem)=>{return {
                    label:x.name,value:x.id
                }}))
                return Promise.resolve(data.partitions)
            }else{
                message.error(msg || '操作失败')
                return Promise.reject(msg || '操作失败')
            }
        })
    }

    const operation:ProColumns<DynamicTableItem>[] =[
        {
            title: '操作',
            key: 'option',
            width: 178,
            fixed:'right',
            valueType: 'option',
            render: (_: React.ReactNode, entity: DynamicTableItem) => [
                <TableBtnWithPermission  access="" key="publish" onClick={()=>{openModal('publish',entity)}} btnTitle="发布管理"/>,
                <Divider type="vertical" className="mx-0"  key="div1"/>,
                <TableBtnWithPermission  access="" key="edit" onClick={()=>{openDrawer('edit',entity)}} btnTitle="查看"/>,
                <Divider type="vertical" className="mx-0"  key="div2"/>,
                <TableBtnWithPermission  access="" key="delete" onClick={()=>{openModal('delete',entity)}} btnTitle="删除"/>,
            ],
        }
    ]
    const handleClusterChange = (e:string[])=>{
        setCluster(e)
        setTableHttpReload(true)
        pageListRef.current?.reload()
    }

    const manualReloadTable = () => {
        setTableHttpReload(true); // 表格数据需要从后端接口获取
        pageListRef.current?.reload()
    };

    const deleteInstance = (entity:DynamicTableItem)=>{
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>(`dynamic/${moduleId}/batch`,{method:'DELETE',eoParams:{ids:JSON.stringify([entity!.id])}}).then(response=>{
                const {code,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功！')
                    resolve(true)
                }else{
                    message.error(msg || '操作失败')
                    reject(msg || '操作失败')
                }
            })
        })
    }

    const openDrawer = async (type:'add'|'edit', entity?:DynamicTableItem)=>{
        switch (type){
            case 'add':
                setCurDetail({driver:driverOptions[0].value || '',config:{'c3ebd745-f7d5-45cd-8d3e-e0e43099d20e':{scopes:[]},'550e2537-8436-48e4-ab84-f9f58faf1b18':{scopes:[]}}})
                break;
            case 'edit':{
                setDrawerLoading(true)
                fetchData<BasicResponse<{info:DynamicTableItem}>>(
                    `dynamic/${moduleId}/info`,
                    {method:'GET',eoParams:{id:entity!.id}}).then((res)=>{
                        const {code, data, msg } = res
                        if(code === STATUS_CODE.SUCCESS){
                            if(data.info.config){
                                for (const tab in data.info.config) {
                                    data.info.config[tab]._apinto_show = true
                                }
                            }
                            setCurDetail(data.info)
                        }else{
                            message.error(msg || '操作失败')
                        }
                    }).finally(()=>setDrawerLoading(false))
                break;
            }
        }
        setDrawerType(type)
        setDrawerOpen(true)
    }

    const openModal = async (type:'publish'|'delete', entity?:DynamicTableItem)=>{
        let title:string = ''
        let content:string|React.ReactNode = ''
        //console.log(renderSchema,driverOptions,entity)
        switch (type){
            case 'publish':{
                title=`${pluginName}发布管理`
                message.loading('正在加载数据')
                const {code,data,msg} = await fetchData<BasicResponse<{data:DynamicPublishData}>>(`dynamic/${moduleId}`,{method:'GET',eoParams:{id:entity!.id},eoTransformKeys:['update_time']})
                message.destroy()
                if(code === STATUS_CODE.SUCCESS){
                    content=<IntelligentPluginPublish ref={publishRef} entity={data.info} moduleId={moduleId!}/>
                }else{
                    message.error(msg || '操作失败')
                    return
                }
                break;}
            case 'delete':
                title='删除'
                content=<span>确定删除成员<span className="text-status_fail"></span>？此操作无法恢复，确认操作？</span>
                break;
        }

        const modalInst = modal.confirm({
            title,
            content,
            onOk:()=>{
                switch (type){ // case 'publish':
                    //     return editRef.current?.save().then((res)=>{if(res === true) manualReloadTable()})
                    case 'delete':
                        return deleteInstance(entity!).then((res)=>{if(res === true) manualReloadTable()})
                }
            },
            width: type === 'delete'? 600 : 900,
            okText:'确认',
            okButtonProps:{
                disabled:false
            },
            cancelText:'取消',
            closable:true,
            icon:<></>,
            footer:(_, { OkBtn, CancelBtn }) =>{//console.log(_,OkBtn,CancelBtn);
                return (
                    <>
                        {type === 'publish' ? <>
                        <WithPermission access="system.partition.self.edit"><CancelBtn/></WithPermission>
                        <WithPermission access="system.partition.self.edit"><Button type="primary" danger loading={publishBtnLoading} onClick={()=>{setPublishBtnLoading(true); publishRef?.current?.offline().then((res)=>{if(res === true) { modalInst.destroy(); manualReloadTable()}}).finally(()=>{setPublishBtnLoading(false); })}}>下线</Button></WithPermission>
                        <WithPermission access="system.partition.self.edit"><Button type="primary"  loading={publishBtnLoading} onClick={()=>{setPublishBtnLoading(true); publishRef?.current?.online().then((res)=>{if(res === true) {  modalInst.destroy(); manualReloadTable()}}).finally(()=>{setPublishBtnLoading(false); })}}>上线</Button></WithPermission>
                            </> :
                            <>
                                <WithPermission access="system.partition.self.edit"><CancelBtn/></WithPermission>
                                <WithPermission access="system.partition.self.edit"><OkBtn/></WithPermission>
                                </>
                        }
                    </>
                );
            },
        })
    }

    // 渲染配置页时需要用到环境数据，在此合并数据
    const getFinalRender = ()=>{
        Promise.all([getRender(),getPartitionList()]).then(([render, partitions])=>{
            if(!partitions || partitions.length === 0) return
            if(partitions.length === 1){
                setRenderSchema({[partitions[0].id]:{
                    properties:render,
                    type:'object'
                }})
            }else{
                setRenderSchema(partitions.map((p:EntityItem)=>({
                    [p.id]:{
                        type: 'void',
                        'x-component': 'FormTab.TabPane',
                        'x-component-props': {
                            tab: p.name,
                        },
                        properties:render,
                    }
                })))
            }
        })
    }

    useEffect(() => {
        getRender()
        getPartitionList()
        pageListRef.current?.reload()
    }, [moduleId]);


    return (<>
    <PageList
            ref={pageListRef}
            columns = {[...columns,...operation]}
            request={(params)=>getIntelligentPluginTableList(params)}
            addNewBtnTitle={`添加${pluginName}`}
            beforeSearchNode={[ <Select key="zoneSelect"
                                        mode="multiple"
                                        allowClear
                                        style={{ width: '100%' }}
                                        placeholder="所有环境"
                                        value={partition}
                                        onChange={handleClusterChange}
                                        options={partitionOptions}
            />]}
            searchPlaceholder={`搜索${pluginName}名称`}
            onChange={() => {
                setTableHttpReload(false)
            }}
            onAddNewBtnClick={()=>{openDrawer('add')}}
            onSearchWordChange={(e)=>{setSearchWord(e.target.value);setTableHttpReload(true);setTableHttpReload(true)}}
        />
        
        <DrawerWithFooter title={`${drawerType === 'add' ? '添加' : '编辑'}${pluginName }`} open={drawerOpen} onClose={()=>{setCurDetail(undefined);setDrawerOpen(false)}} onSubmit={()=>drawerFormRef.current?.save()?.then((res)=>{res && manualReloadTable();return res})}  submitAccess='system.partition.self.edit'>
            <Spin indicator={<LoadingOutlined style={{ fontSize: 24 }} spin/>} spinning={drawerLoading}>
                <IntelligentPluginConfig 
                    ref={drawerFormRef!} 
                    type={drawerType} 
                    renderSchema={renderSchema} 
                    tabData={partitionOptions}
                    moduleId={moduleId!} 
                    driverSelectionOptions={driverOptions}  
                    initFormValue={curDetail as { [k: string]: unknown; }} />
            </Spin>
        </DrawerWithFooter>
    </>)
}