import {ActionType, ProColumns} from "@ant-design/pro-components";
import  { useRef, useState} from "react";
import PageList from "@common/components/aoplatform/PageList.tsx";
import {App, Button, Modal} from "antd";
import TransferTable, {TransferTableHandle} from "@common/components/aoplatform/TransferTable.tsx";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import { apiModalColumn, SYSTEM_MYSERVICE_API_TABLE_COLUMNS } from "../../../const/system/const.tsx";
import { ServiceApiTableListItem, SimpleApiItem } from "../../../const/system/type.ts";
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";
import { checkAccess } from "@common/utils/permission.ts";

const MyServiceInsideApi = ({systemId,serviceId}:{systemId:string, serviceId:string})=>{
    const { modal,message } = App.useApp()
    const [open, setOpen] = useState(false);
    const {fetchData} = useFetch()
    const [init, setInit] = useState<boolean>(true)
    const pageListRef = useRef<ActionType>(null);
    const [addApiBtnLoading, setAddApiBtnLoading] = useState<boolean>()
    const [addApiBtnDisabled, setAddApiBtnDisabled] = useState<boolean>(true)
    const addRef = useRef<TransferTableHandle<ServiceApiTableListItem>>(null)
    const [apiIds, setApiIds] = useState<string[]>([])
    const {accessData} = useGlobalContext()

    const getServiceApiList = ()=>{
        return fetchData<BasicResponse<{apis:ServiceApiTableListItem[]}>>('project/service/apis',{method:'GET',eoParams:{service:serviceId,project:systemId}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setInit((prev)=>prev ? false : prev)
                setApiIds(data.apis?.map((x:ServiceApiTableListItem)=>x.id) || [])
                return  {data:data.apis, success: true}
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        }).catch(() => {
            return {data:[], success:false}
        })
    }

    const handleDragSortEnd = (beforeIndex: number, afterIndex: number, newDataSource: ServiceApiTableListItem[]) => {
        fetchData<BasicResponse<{apis:ServiceApiTableListItem}>>('project/service/api/sort',{method:'PUT',eoParams:{service:serviceId,project:systemId},eoBody:({apis:newDataSource?.map(x=>x.id)})|| []}).then(response=>{
            const {code,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                message.success(msg || '操作成功')
                manualReloadTable()
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        })
    };

    const manualReloadTable = () => {
        pageListRef.current?.reload()
    };

    const operation:ProColumns<ServiceApiTableListItem>[] =[
        {
            title: '操作',
            key: 'option',
            width: 62,
            fixed:'right',
            valueType: 'option',
            render: (_: React.ReactNode, entity: ServiceApiTableListItem) => [
                <TableBtnWithPermission  access="project.mySystem.service.edit" key="delete"  onClick={()=>{openModal('delete',entity)}} btnTitle="删除"/>,
            ],
        }
    ]

    const deleteApi = (entity:ServiceApiTableListItem) =>{
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>('project/service/unbind',{method:'DELETE',eoParams:{service:serviceId,project:systemId, apis:JSON.stringify([entity.id])}}).then(response=>{
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

    const addApi = ()=>{
        setAddApiBtnLoading(true)
        fetchData<BasicResponse<null>>('project/service/bind',{method:'POST' ,eoBody:({apis:addRef.current?.selectedRowKeys().filter(x=>apiIds.indexOf(x) === -1)}),eoParams:{service:serviceId,project:systemId}}).then(response=>{
            const {code,msg} = response
            setAddApiBtnLoading(false)
            if(code === STATUS_CODE.SUCCESS){
                setOpen(false)
                message.success(msg || '操作成功！')
                manualReloadTable()
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const openModal = async (type:'add'|'delete',entity?:ServiceApiTableListItem)=>{
        switch(type){
            case 'add':
                setOpen(true)
                break;
            case 'delete':
                modal.confirm({
                    title:'删除',
                    content:'该数据删除后将无法找回，请确认是否删除？',
                    onOk:()=> deleteApi(entity!).then((res)=>{if(res === true) manualReloadTable()}),
                    width:600,
                    okText:'确认',
                    okButtonProps:{
                        disabled : !checkAccess( 'project.mySystem.service.delete', accessData)
                    },
                    cancelText:'取消',
                    closable:true,
                    icon:<></>
                })
        }
    }

    const getSelectableApiList = (keyword?:string)=>{
        return fetchData<BasicResponse<{apis:SimpleApiItem[]}>>('project/apis/simple',{method:'GET',eoParams:{ keyword,project:systemId},eoTransformKeys:['request_path']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                return  {data:data.apis, success: true}
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        }).catch(() => {
            return {data:[], success:false}
        })
    }

    return (
            <div className="h-full">
                <PageList
                    id="global_system_myService_api"
                    ref={pageListRef}
                    columns={[...SYSTEM_MYSERVICE_API_TABLE_COLUMNS, ...operation]}
                    request={()=>getServiceApiList()}
                    addNewBtnTitle="添加 API"
                    addNewBtnAccess="project.mySystem.service.edit"
                    onAddNewBtnClick={() => {
                      openModal('add')
                    }}
                    dragSortKey="id"
                    onDragSortEnd={handleDragSortEnd}
                />
                <Modal
                    title="添加 API"
                    open={open}
                    width={600}
                    destroyOnClose={true}
                    wrapClassName="height-fixed-modal"
                    maskClosable={false}
                    onCancel={() => setOpen(false)}
                    footer={[
                        <Button key="back" onClick={() => setOpen(false)}>
                            取消
                        </Button>,
                        <WithPermission access="project.mySystem.service.edit"><Button
                            key="submit"
                            type="primary"
                            disabled={addApiBtnDisabled}
                            loading={addApiBtnLoading}
                            onClick={addApi}
                        >
                            确认
                        </Button></WithPermission>,
                    ]}
                >
                    <TransferTable
                        ref={addRef}
                        request={getSelectableApiList}
                        columns={apiModalColumn}
                        primaryKey="id"
                        disabledData={apiIds}
                        onSelect={(selectedData: SimpleApiItem[]) => {
                            setAddApiBtnDisabled(!(selectedData.length > 0));
                        }}
                    />
                </Modal>
            </div>
    )
}
export default MyServiceInsideApi