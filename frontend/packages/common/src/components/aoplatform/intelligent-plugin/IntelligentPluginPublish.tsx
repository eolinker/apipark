import  {forwardRef, useEffect, useImperativeHandle, useState} from "react";
import {DynamicPublishCluster, StatusColorClass} from "./IntelligentPluginList.tsx";
import {App, Col, Row, Table} from "antd";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";

export type IntelligentPluginPublishProps = {
    entity:{[k:string]:unknown}
    moduleId:string
}

export type IntelligentPluginPublishHandle = {
    offline:()=>Promise<boolean | string>
    online:()=>Promise<boolean | string>
}

export type DynamicPublish = {
    code:number,
    msg:string,
    data:{
        success:Array<string>,
        fail:Array<string>
    }
}

export const IntelligentPluginPublish = forwardRef<IntelligentPluginPublishHandle,IntelligentPluginPublishProps>((props,ref)=>{
    const { message } = App.useApp()
    const { entity, moduleId} = props
    const [showNoCluster, setShowNoCluster]= useState<boolean>(false)
    const [selectedCluster, setSelectedCluster] = useState<string[]>([])
    const [selectedClusterUuid, setSelectedClusterUuid] = useState<string[]>([])
    const [partitionDataSource, setPartitionDataSource] = useState<DynamicPublishCluster[]>(entity?.partitions)
    const {fetchData} = useFetch()
    const [startCheckClusterNum,setStartCCN] = useState<boolean>(false)
    const apiColumns = [
        {
            title:'环境',
            dataIndex:'title',
            copyable: true,
            ellipsis:true
        },
        {
            title:'状态',
            dataIndex:'status',
            render:(dom, entity)=> {
                return <span className={StatusColorClass[entity.status as keyof typeof StatusColorClass]}>{(entity.status as string)}</span>                        
            },
        },
        {
            title:'更新人',
            dataIndex:['updater','name'],
            width:88
        },
        {
            title:'更新时间',
            dataIndex:'updateTime',
            width:182,
        }
    ]

    const online :()=>Promise<boolean | string> = ()=>{
            setStartCCN(true)
            if (!selectedCluster.length) {
                setShowNoCluster(!selectedCluster.length)
                return Promise.reject('操作失败')
            }

            return fetchData<BasicResponse<DynamicPublish>>(`dynamic/${moduleId}/online`, {
                method: 'PUT',
                eoParams:{id:entity.id},
                eoBody: ({partitions:selectedClusterUuid}),
            }).then(response => {
                const {code, msg} = response
                if (code === STATUS_CODE.SUCCESS) {
                    message.success(msg || '操作成功！')
                    return Promise.resolve(true)
                } else {
                    message.error(msg || '操作失败')
                    return Promise.reject(msg || '操作失败')
                }
            }).catch((errorInfo)=> Promise.reject(errorInfo))
    }

    const offline :()=>Promise<boolean | string> = ()=>{
            setStartCCN(true)
            if (!selectedCluster.length) {
                setShowNoCluster(!selectedCluster.length)
                return Promise.reject('操作失败')
            }

            return fetchData<BasicResponse<DynamicPublish>>(`dynamic/${moduleId}/offline`, {
                method: 'PUT',
                eoParams:{id:entity.id},
                eoBody: ({partitions:selectedClusterUuid}),
            }).then(response => {
                const {code, msg} = response
                if (code === STATUS_CODE.SUCCESS) {
                    message.success(msg || '操作成功！')
                    return Promise.resolve(true)
                } else {
                    message.error(msg || '操作失败')
                    return Promise.reject(msg || '操作失败')
                }
            }).catch((errorInfo)=> Promise.reject(errorInfo))
    }

    useImperativeHandle(ref, ()=>({
        online,offline
        })
    )

    // rowSelection object indicates the need for row selection
    const rowSelection = {
        selectedRowKeys:selectedClusterUuid,
        onChange: (selectedRowKeys: React.Key[], selectedRows: {[k:string]:unknown}[]) => {
            setStartCCN(true)
            setSelectedCluster(selectedRows?.map((x)=>x.title))
            setSelectedClusterUuid(selectedRows?.map((x)=>x.name))
        }
    };

    const handleRowClick = (entity:DynamicPublishCluster)=>{
        setSelectedCluster(prev=>prev.indexOf(entity.title) === -1 ? [...prev, entity.title]:prev.filter(x=>x !== entity.title))
        setSelectedClusterUuid(prev=>prev.indexOf(entity.name) === -1 ? [...prev, entity.name]:prev.filter(x=>x !== entity.name))
    }

    useEffect(()=>{
        startCheckClusterNum && setShowNoCluster(selectedClusterUuid.length  === 0)
    },[selectedClusterUuid])

    return (
        <>
        <Row className="mb-[8px]">
            <Col className="text-left  ant-form-item-label w-[42px]"><span >名称：</span></Col>
            <Col >{entity.title ?? '-'}</Col>
        </Row>

        <Row className="mb-[8px]">
            <Col className="text-left  ant-form-item-label w-[42px]" ><span >ID：</span></Col>
            <Col >{entity.id ??'-'}</Col>
        </Row>

        <Row className="mb-[8px]">
            <Col className="text-left  ant-form-item-label w-[42px]" ><span >描述：</span></Col>
            <Col >{entity.description ??'-'}</Col>
        </Row>

        <Row className="">
                <Table
                    columns={apiColumns}
                    bordered={true}
                    rowKey="name"
                    size="small"
                    dataSource={partitionDataSource}
                    pagination={false}
                    rowClassName="cursor-pointer"
                    rowSelection={{
                        type: 'checkbox',
                        columnWidth: 40,
                        ...rowSelection,
                    }}
                    onRow={(record) => ({
                        onClick: () => {
                            handleRowClick(record);
                        }
                        })}
                />
        </Row>
            {showNoCluster && 
                <Row className="">
                    <Col offset={0} span={24}><p className="text-status_fail">请至少选中一个集群</p></Col>
                </Row>}
            
        <Row className="my-mbase">
            <Col  span={24}> <div  className="text-right"><span>已选择 {selectedCluster.length } 项{selectedCluster.length > 0 && <span>：{selectedCluster.join(',')}</span>}</span></div></Col>
        </Row>
        
        </>)
})