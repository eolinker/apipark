/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-06 17:54:34
 * @FilePath: \frontend\packages\core\src\pages\partitions\PartitionList.tsx
 */
import PageList from "@common/components/aoplatform/PageList.tsx"
import {ActionType, ProColumns} from "@ant-design/pro-components";
import  {FC, useEffect, useMemo, useRef, useState} from "react";
import { useNavigate} from "react-router-dom";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {App, Button, Modal} from "antd";
import { PARTITION_LIST_COLUMNS } from "../../const/partitions/const.tsx";
import { PartitionTableListItem } from "../../const/partitions/types.ts";
import { SimpleMemberItem } from "../../const/system/type.ts";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import PartitionConfig, { PartitionConfigHandle } from "./PartitionConfig.tsx";
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";


const PartitionList:FC = ()=>{
    const { message, } = App.useApp()
    const [searchWord, setSearchWord] = useState<string>('')
    const navigate = useNavigate();
    const { setBreadcrumb } = useBreadcrumb()
    const {fetchData} = useFetch()
    const [init, setInit] = useState<boolean>(true)
    const [memberValueEnum, setMemberValueEnum] = useState<{[k:string]:{text:string}}>({})
    const addPartitionRef = useRef<PartitionConfigHandle>(null)
    const pageListRef = useRef<ActionType>(null);
    const [modalVisible,setModalVisible] = useState<boolean>(false)
    const [addClusterStep,setAddClusterStep] = useState<'config'|'retry'|'result'>('config')
    const [configBtnLoading,setConfigBtnLoading] = useState<boolean>(false)
    const [retryBtnLoading,setRetryBtnLoading] = useState<boolean>(false)
    const operation:ProColumns<PartitionTableListItem>[] =[
        {
            title: '操作',
            key: 'option',
            width: 62,
            fixed:'right',
            valueType: 'option',
            render: (_: React.ReactNode, entity: PartitionTableListItem) => [
                <TableBtnWithPermission  access="system.partition.self.view" key="view" navigateTo={`../inside/${entity.id}/cluster`} btnTitle="查看"/>,
            ],
        }
    ]


    const getPartitionList =(): Promise<{ data: PartitionTableListItem[], success: boolean }>=> {
        return fetchData<BasicResponse<{partitions:PartitionTableListItem}>>('partitions',{method:'GET',eoParams:{keyword:searchWord},eoTransformKeys:['cluster_num','update_time']}).then(response=>{
            const {code,data,msg} = response
            //console.log(code, data,STATUS_CODE.SUCCESS,code === STATUS_CODE.SUCCESS)
            if(code === STATUS_CODE.SUCCESS){
                setInit((prev)=>prev ? false : prev)
                return  {data:data.partitions, success: true}
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
                title:'部署管理'
            },
        ])
        getMemberList()
    }, []);

    const openModal = ()=>{
        setModalVisible(true)
        setAddClusterStep('config')
        // modal.confirm({
        //     title:'添加分区',
        //     content:<PartitionConfig ref={addPartitionRef} />,
        //     onOk:()=> {
        //         return addPartitionRef.current?.save().then((res)=>{if(res === true) manualReloadTable()})
        //     },
        //     width:600,
        //     okText:'下一步，检查集群',
        //     cancelText:'取消',
        //     closable:true,
        //     icon:<></>,
        // })
    }

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
        return PARTITION_LIST_COLUMNS.map(x=>{if(x.filters &&((x.dataIndex as string[])?.indexOf('updater') !== -1 || (x.dataIndex as string[])?.indexOf('approver') !== -1) ){x.valueEnum = memberValueEnum} return x})
    },[memberValueEnum])

    const manualReloadTable = () => {
        pageListRef.current?.reload()
    };

    return (
        <>
        <PageList
            id="global_partition"
            ref={pageListRef}
            columns = {[...columns,...operation]}
            request={()=>getPartitionList()}
            addNewBtnTitle="添加分区"
            showPagination={false}
            searchPlaceholder="输入名称、ID 查找分区"
            onAddNewBtnClick={()=>{openModal()}}
            addNewBtnAccess="system.partition.self.add"
            onSearchWordChange={(e)=>{setSearchWord(e.target.value)}}
            onRowClick={(row:PartitionTableListItem)=>navigate(`../inside/${row.id}/cluster`)}
            tableClickAccess="system.partition.self.view"
            />
            <Modal
                title="添加分区"
                open={modalVisible}
                width={900}
                destroyOnClose={true}
                maskClosable={false}
                afterOpenChange={(open:boolean)=>{
                    !open && setModalVisible(false);setAddClusterStep('config');setConfigBtnLoading(false);setRetryBtnLoading(false)
                }}
                onCancel={() => {setModalVisible(false);setAddClusterStep('config');setConfigBtnLoading(false);setRetryBtnLoading(false)}}
                footer={[
                    <Button key="back" onClick={() => setModalVisible(false)}>
                        取消
                    </Button>,
                    <>{addClusterStep === 'result' && <>
                       <WithPermission access="system.partition.cluster.add" key="lastStepPermission"><Button key="lastStep" className="" onClick={() => setAddClusterStep('config')}>上一步</Button></WithPermission>
                       <WithPermission access="system.partition.cluster.add" key="finishPermission"><Button key="finish" type="primary" onClick={() => addPartitionRef.current?.save().then((res)=>{if(res === true) {manualReloadTable();setModalVisible(false)}})}>完成</Button></WithPermission>
                    </>}</>,
                    <>{addClusterStep === 'config' &&
                    <WithPermission access="system.partition.cluster.add" key="checkPermission"><Button key="check" loading={configBtnLoading} type="primary"  onClick={() =>{ setConfigBtnLoading(true);addPartitionRef.current?.check().then((res)=>{setConfigBtnLoading(false); setAddClusterStep(res === true ? 'result':'retry')}).finally(()=>{setConfigBtnLoading(false)})}}>下一步，检查集群</Button></WithPermission>
                    }</>,
                    <>{addClusterStep === 'retry' &&
                    <WithPermission access="system.partition.cluster.add" key="checkPermission"><Button key="check" loading={retryBtnLoading} type="primary" onClick={() => { setRetryBtnLoading(true); addPartitionRef.current?.check().then((res)=>{setRetryBtnLoading(false); setAddClusterStep(res === true ? 'result':'retry')}).finally(()=>{setRetryBtnLoading(false)})}}>重新检查</Button></WithPermission>
                    }</>
                ]}
            >
                <PartitionConfig ref={addPartitionRef}  mode={addClusterStep}/>
            </Modal>
        </>
    )

}
export default PartitionList