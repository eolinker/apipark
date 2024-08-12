import PageList from "@common/components/aoplatform/PageList.tsx"
import {ActionType, ProColumns} from "@ant-design/pro-components";
import  {FC, useEffect, useMemo, useRef, useState} from "react";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {App, Divider, Modal} from "antd";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import { OrganizationFieldType, OrganizationTableListItem } from "../../const/organization/type.ts";
import { ORGANIZATION_TABLE_COLUMNS } from "../../const/organization/const.tsx";
import { SimpleMemberItem } from "@common/const/type.ts";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";
import { checkAccess } from "@common/utils/permission.ts";
import OrganizationConfig, { OrganizationConfigHandle } from "./OrganizationConfig.tsx";
import { PartitionItem } from "@common/const/type.ts";

const OrganizationList:FC = ()=>{
    const [searchWord, setSearchWord] = useState<string>('')
    const { setBreadcrumb } = useBreadcrumb()
    const { modal,message } = App.useApp()
    const pageListRef = useRef<ActionType>(null);
    const [memberValueEnum, setMemberValueEnum] = useState<{[k:string]:{text:string}}>({})
    const {fetchData} = useFetch()
    const {accessData} = useGlobalContext()
    const orgConfigRef = useRef<OrganizationConfigHandle>(null)
    const [curOrg, setCurOrg] = useState<OrganizationFieldType>({})
    const [modalVisible, setModalVisible] = useState<boolean>(false)
    const [modalType, setModalType] = useState<'add'|'edit'>('add')

    const getOrganizationList = ()=>{
        //console.log('此处应该获取最新列表',searchWord)
        return fetchData<BasicResponse<{organizations:OrganizationTableListItem}>>('manager/organizations',{method:'GET',eoTransformKeys:['update_time','can_delete','create_time'],eoParams:{keyword:searchWord}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                return  {data:data.organizations, success: true}
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        }).catch(() => {
            return {data:[], success:false}
        })
    }

    const deleteOrganization = (id:string)=>{
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>(`manager/organization`,{method:'DELETE',eoParams:{id}}).then(response=>{
                const {code,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功！')
                    resolve(true)
                }else{
                    message.error(msg || '操作失败')
                    reject(msg || '操作失败')
                }
            }).catch((errorInfo)=>{
                reject(errorInfo || '操作失败')
            })
        })
    }

    const openModal = async (type:'add'|'edit'|'delete',entity?:OrganizationTableListItem)=>{
        let title:string = ''
        let content:string | React.ReactNode= ''
        switch (type){
            case 'add':{
                setModalType('add')
                setModalVisible(true)
                return;}
            case 'edit':{
                message.loading('正在加载数据')
                const {code,data,msg} = await fetchData<BasicResponse<{ organization: OrganizationFieldType }>>('manager/organization',{method:'GET',eoParams:{id:entity!.id},eoTransformKeys:['create_time','master_id','update_time']})
                message.destroy()
                if(code === STATUS_CODE.SUCCESS){
                    setCurOrg({...data.organization,partitions:data.organization.partitions?.map((x:PartitionItem)=>(x.id)),master:data.organization.master.id})
                    setModalVisible(true)
                }else{
                    message.error(msg || '操作失败')
                    return
                }
                setModalType('edit')
                return;}
            case 'delete':
                title='删除'
                content='该数据删除后将无法找回，请确认是否删除？'
                break;
        }

         modal.confirm({
            title,
            content,
            onOk:()=>{
                return deleteOrganization(entity!.id).then((res)=>{if(res === true) pageListRef.current?.reload()})
            },
            okText:'确认',
            okButtonProps:{
                disabled : !checkAccess('system.organization.self.delete', accessData)
            },
            cancelText:'取消',
            closable:true,
            icon:<></>
        })
    }
    const operation:ProColumns<OrganizationTableListItem>[] =[
            {
                title: '操作',
                key: 'option',
                width: 107,
                fixed:'right',
                valueType: 'option',
                render: (_: React.ReactNode, entity: OrganizationTableListItem) => [
                <TableBtnWithPermission  access="system.organization.self.edit" key="edit" onClick={()=>openModal('edit',entity)} btnTitle="编辑"/>,
                <Divider type="vertical" className="mx-0"  key="div1"/>,
                <TableBtnWithPermission  access="system.organization.self.delete" key="delete" disabled={!entity.canDelete} tooltip="团队数据清除后，方可删除" onClick={()=>{openModal('delete',entity)}} btnTitle="删除"/>,
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
        getMemberList()
        setBreadcrumb([
            {
                title: '组织'}
        ])
    }, []);

    
    const columns = useMemo(()=>{
        return ORGANIZATION_TABLE_COLUMNS.map(x=>{if(x.filters &&(x.dataIndex as string[])?.indexOf('master') !== -1 ){x.valueEnum = memberValueEnum} return x})
    },[memberValueEnum])
    
    
    const manualReloadTable = () => {
        pageListRef.current?.reload()
    };


    return (
        <>
            <PageList
                id="global_organization"
                ref={pageListRef}
                showPagination={false}
                columns = {[...columns,...operation]}
                request={()=> getOrganizationList()}
                addNewBtnTitle="添加组织"
                addNewBtnAccess="system.organization.self.add"
                tableClickAccess="system.organization.self.edit"
                searchPlaceholder="输入名称、ID、负责人查找组织"
                onAddNewBtnClick={()=>{openModal('add')}}
                onSearchWordChange={(e)=>{setSearchWord(e.target.value)}}
                onRowClick={(row:OrganizationTableListItem)=>openModal('edit',row)}
            />
            <Modal
                title={modalType === 'add' ? "添加组织" : "配置组织"}
                open={modalVisible}
                width={600}
                destroyOnClose={true}
                maskClosable={false}
                afterOpenChange={(open:boolean)=>{
                    if(!open){
                        setModalVisible(false)
                        setCurOrg({} as unknown as OrganizationFieldType)
                    }
                }}
                onCancel={() => {setModalVisible(false)}}
                okText="确认"
                okButtonProps={{disabled : !checkAccess( modalType === 'add' ? 'system.organization.self.add':'system.organization.self.edit', accessData)}}
                cancelText='取消'
                closable={true}
                onOk={()=>orgConfigRef.current?.save().then((res)=>{
                    if(res){
                        setModalVisible(false)
                        manualReloadTable()
                    }
                    return res})}
            >
                <OrganizationConfig ref={orgConfigRef}  entity={modalType === 'add' ? undefined : curOrg} />
            </Modal>
            {/* <DrawerWithFooter title={`${curOrgId=== undefined ? '添加' : '编辑'}组织`} open={open} onClose={onClose} onSubmit={()=>orgConfigRef.current?.save()?.then((res)=>{res && pageListRef.current?.reload();return res})} submitAccess={curOrgId === undefined ? } >
            </DrawerWithFooter> */}
        </>
    )

}
export default OrganizationList