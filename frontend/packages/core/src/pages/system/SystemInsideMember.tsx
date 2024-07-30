import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {App, Button, Modal, Select} from "antd";
import PageList from "@common/components/aoplatform/PageList.tsx";
import  {useEffect, useMemo, useRef, useState} from "react";
import {Link, useParams} from "react-router-dom";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {ActionType, ProColumns} from "@ant-design/pro-components";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {TransferTableHandle} from "@common/components/aoplatform/TransferTable.tsx";
import {useFetch} from "@common/hooks/http.ts";
import {EntityItem, MemberItem, TeamSimpleMemberItem} from "@common/const/type.ts";
import { SYSTEM_MEMBER_TABLE_COLUMN } from "../../const/system/const.tsx";
import { SystemMemberTableListItem } from "../../const/system/type.ts";
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";
import { checkAccess } from "@common/utils/permission.ts";
import MemberTransfer from "@common/components/aoplatform/MemberTransfer.tsx";
import { DepartmentListItem } from "../../const/member/type.ts";
import { addMemberToDepartment, getDepartmentWithMember } from "../user/UserList.tsx";
import {v4 as uuidv4} from 'uuid'

export default function SystemInsideMember(){
    const {systemId, teamId} = useParams<RouterParams>();
    const [searchWord, setSearchWord] = useState<string>('')
    const {modal, message} = App.useApp()
    const {setBreadcrumb} = useBreadcrumb()
    // const [confirmLoading, setConfirmLoading] = useState(false);
    const {fetchData} = useFetch()
    const [init, setInit] = useState<boolean>(true)
    const [tableListDataSource, setTableListDataSource] = useState<SystemMemberTableListItem[]>([]);
    const [tableHttpReload, setTableHttpReload] = useState(true);
    const pageListRef = useRef<ActionType>(null);
    const addRef = useRef<TransferTableHandle<SystemMemberTableListItem>>(null)
    const [columns,setColumns] = useState<ProColumns<SystemMemberTableListItem>[]>([])
    const [allMemberIds,setAllMemberIds] = useState<string[]>([])
    const {accessData} = useGlobalContext()
    const [selectableMemberIds,setSelectableMemberIds] = useState<Set<string>>(new Set())
    const [addMemberBtnLoading, setAddMemberBtnLoading] = useState<boolean>(false)
    const [modalVisible, setModalVisible] = useState<boolean>(false)
    const [addMemberBtnDisabled, setAddMemberBtnDisabled] = useState<boolean>(true)
    const [allMemberSelectedDepartIds, setAllMemberSelectedDepartIds] = useState<string[]>([])

    const operation: ProColumns<SystemMemberTableListItem>[] = [
        {
            title: '操作',
            key: 'option',
            width: 62,
            fixed:'right',
            valueType: 'option',
            render: (_: React.ReactNode, entity: SystemMemberTableListItem) => [
                <TableBtnWithPermission  access="project.mySystem.member.edit" key="delete" disabled={!entity.canDelete} tooltip="暂无权限" onClick={()=>{openModal('delete',entity)}} btnTitle="删除"/>,
            ],
        }
    ]

    const getMemberList = (): Promise<{ data: SystemMemberTableListItem[], success: boolean }> => {
        if (!tableHttpReload) {
            setTableHttpReload(true)
            return Promise.resolve({
                data: tableListDataSource,
                success: true,
            });
        }
        return fetchData<BasicResponse<{ members: SystemMemberTableListItem }>>('project/members', {method: 'GET', eoParams:{project:systemId},eoTransformKeys:['can_delete']}).then(response => {
            const {code, data, msg} = response
            if (code === STATUS_CODE.SUCCESS) {
                if(!searchWord) setAllMemberIds(data.members?.map((x:SystemMemberTableListItem)=>x.user.id) || [])
                setTableListDataSource(data.members)
                setInit((prev) => prev ? false : prev)
                setTableHttpReload(false)
                return {data: data.members, success: true}
            } else {
                message.error(msg || '操作失败')
                return {data: [], success: false}
            }
        }).catch(() => {
            return {data: [], success: false}
        })
    }

    const getDepartmentMemberList = () => {
        const topDepartmentId:string = uuidv4()
        return Promise.all([
          fetchData<BasicResponse<{department:DepartmentListItem}>>('simple/departments', {method:'GET'}),
          fetchData<BasicResponse<{teams:TeamSimpleMemberItem[]}>>('team/members/simple',{method:'GET',eoParams:{team:teamId},eoTransformKeys:['user_group','attach_time','user_id']})
        ]).then(([departmentResponse, memberResponse])=>{
            const departmentMap = new Map<string, (MemberItem & {type:'department'|'member'})[]>();
            memberResponse.data.teams.forEach((member: (TeamSimpleMemberItem | MemberItem) & {title?:string, key?:string}) => {
                setSelectableMemberIds((pre)=>{pre.add((member as TeamSimpleMemberItem).user?.id);return pre})
                member = {department:(member as TeamSimpleMemberItem).department ? (member as TeamSimpleMemberItem).department : undefined, email:(member as TeamSimpleMemberItem).mail, ...(member as TeamSimpleMemberItem).user, title:(member as TeamSimpleMemberItem).user.name, key:(member as TeamSimpleMemberItem).user.id}
              if (member.department) {
                member.department.forEach((department: EntityItem) => {
                  addMemberToDepartment(departmentMap, department.id, member);
                });
              } else {
                addMemberToDepartment(departmentMap, '_withoutDepartment', member);
              }
            });
            const finalData = departmentResponse.data.department 
              ? [
                  {
                    id: topDepartmentId, 
                    key:topDepartmentId,
                    name: departmentResponse.data.department.name, 
                    title:departmentResponse.data.department.name, 
                    children: [
                      ...getDepartmentWithMember(departmentResponse.data.department.children, departmentMap),
                      ...departmentMap.get('_withoutDepartment') || []
                    ]
                  }
                ] 
              : [...departmentMap.get('_withoutDepartment') || []];
          
              for(const [k,v] of departmentMap){
                if(k !== '_withoutDepartment' && allMemberIds.length > 0 ){
                     // 筛选出部门内没被勾选的用户，如果不存在没勾选用户，需要将部门id放入ids中
                     if(v.filter(m => allMemberIds.indexOf(m.id) === -1).length  === 0){
                         setAllMemberSelectedDepartIds((pre)=>[...pre, k])
                     }
                }
             }

             if(!finalData[0].children || finalData[0].children.filter(m => allMemberIds.indexOf(m.id) === -1).length  === 0){
                 setAllMemberSelectedDepartIds((pre)=>[...pre, topDepartmentId])
             }
             
              return  {data:finalData, success: true}
        }).catch(()=>({data:[], success:false}))
      }
      

    const openModal = (type: 'add' | 'edit' | 'delete', entity?: SystemMemberTableListItem) => {
        let title: string = ''
        let content: string | React.ReactNode = ''
        switch (type) {
            case 'add':
                setModalVisible(true)
                setAddMemberBtnDisabled(true)
                setAddMemberBtnLoading(false)
                return;
            case 'delete':
                title = '删除'
                content = '该数据删除后将无法找回，请确认是否删除？'
                break;
        }

        modal.confirm({
            title,
            content,
            onOk:()=> removeMember(entity!).then((res)=>{if(res === true) manualReloadTable()}),
            width: 600,
            okText: '确认',
            okButtonProps:{
                disabled : !checkAccess( `project.mySystem.member.edit`, accessData)
            },
            cancelText: '取消',
            closable: true,
            icon: <></>,
        })
    }

    const addMember = () => {
        setAddMemberBtnLoading(true)
        const keyFromModal = addRef.current?.selectedRowKeys()
        const memberKeyFromModal = keyFromModal?.filter(x => allMemberIds.indexOf(x as string) === -1 && selectableMemberIds.has(x)) || [];
       fetchData<BasicResponse<null>>('project/member', {method: 'POST', eoBody: ({users: memberKeyFromModal}), eoParams: {project: systemId}}).then(response => {
            const {code, msg} = response
            setAddMemberBtnLoading(false)
            if (code === STATUS_CODE.SUCCESS) {
                setModalVisible(false)
                message.success(msg || '操作成功！')
                manualReloadTable()
            } else {
                message.error(msg || '操作失败')
            }
        }).finally(()=>setAddMemberBtnLoading(false))
    }

    const manualReloadTable = () => {
        setTableHttpReload(true)
        pageListRef.current?.reload()
    };

    const removeMember = (entity: SystemMemberTableListItem) => {
        return new Promise((resolve, reject) => {
            fetchData<BasicResponse<null>>(`project/member`, {method: 'DELETE', eoParams: {project: systemId, user:entity.user.id}}).then(response => {
                const {code, msg} = response
                if (code === STATUS_CODE.SUCCESS) {
                    message.success(msg || '操作成功！')
                    resolve(true)
                } else {
                    message.error(msg || '操作失败')
                    reject(msg || '操作失败')
                }
            }).catch((errorInfo)=> reject(errorInfo))
        })
    }

    const changeMemberInfo = (value:string[],entity:SystemMemberTableListItem )=>{
        //console.log(value)
        return new Promise((resolve, reject) => {
            fetchData<BasicResponse<null>>(`project/member`, {method: 'PUT',eoBody:({roles:value}), eoParams: {project: systemId, user:entity.user.id}}).then(response => {
                const {code, msg} = response
                if (code === STATUS_CODE.SUCCESS) {
                    message.success(msg || '操作成功！')
                    resolve(true)
                } else {
                    message.error(msg || '操作失败')
                    reject(msg || '操作失败')
                }
            }).catch((errorInfo)=> reject(errorInfo))
        })
    }

    const getRoleList = ()=>{
        fetchData<BasicResponse<{roles:Array<{id:string,name:string}>}>>('simple/roles', {method: 'GET'}).then(response => {
            const {code, data,msg} = response
            if (code === STATUS_CODE.SUCCESS) {

                const newCol = [...SYSTEM_MEMBER_TABLE_COLUMN]
                for(const col of newCol){
                    //console.log(col)
                    if(col.dataIndex instanceof Array && col.dataIndex[0] === 'roles'){
                        col.render = (_,entity)=>(
                            <WithPermission access="project.mySystem.member.edit">
                                <Select
                                    className="w-full"
                                    mode="multiple"
                                    value={entity.roles?.map((x:EntityItem)=>x.id)}
                                    options={data.roles?.map((x:{id:string,name:string})=>({label:x.name, value:x.id}))}
                                    onChange={(value)=>{
                                        changeMemberInfo(value,entity ).then((res)=>{
                                            if(res) manualReloadTable()
                                        })
                                    }}
                                />
                            </WithPermission>
                        )
                        col.filters = data.roles?.map((x:{id:string,name:string})=>({text:x.name, value:x.id}))
                        col.onFilter = (value: unknown, record:SystemMemberTableListItem) =>{
                            return record.roles ? record.roles?.map((x)=>x.id).indexOf(value as string) !== -1 : false;}
                        setColumns(newCol)
                        return
                    }
                }
            } else {
                message.error(msg || '操作失败')
            }
        })
    }

    
    const cleanModalData = ()=>{
        setModalVisible(false);setAddMemberBtnDisabled(true);setAddMemberBtnLoading(false)
    }

    const treeDisabledData = useMemo(()=>{ return [...allMemberIds,...allMemberSelectedDepartIds]},[allMemberIds,allMemberSelectedDepartIds])


    useEffect(() => {
        getRoleList()
        setBreadcrumb([
            {
                title: <Link to={`/system/list`}>内部数据服务</Link>
            },
            {
                title: '成员'
            }])
        manualReloadTable()
    }, [systemId]);

    return (
        <>
            <PageList
                id="global_system_member"
                ref={pageListRef}
                columns={[...columns, ...operation]}
                request={() => getMemberList()}
                dataSource={tableListDataSource}
                showPagination={false}
                addNewBtnTitle="添加成员"
                addNewBtnAccess="project.mySystem.member.add"
                onAddNewBtnClick={() => {
                    openModal('add')
                }}
                onChange={() => {
                    setTableHttpReload(false)
                }}
            />
            <Modal
                title="添加成员"
                open={modalVisible}
                destroyOnClose={true}
                width={900}
                onCancel={() => {cleanModalData()}}
                maskClosable={false}
                footer={[
                    <Button key="back" onClick={() => cleanModalData()}>
                        取消
                    </Button>,
                    <WithPermission access="project.mySystem.member.add"><Button
                        key="submit"
                        type="primary"
                        disabled={addMemberBtnDisabled}
                        loading={addMemberBtnLoading}
                        onClick={addMember}
                    >
                        确认
                    </Button></WithPermission>,
                ]}
            >
                <MemberTransfer
                    ref={addRef}
                    primaryKey="id"
                    disabledData={treeDisabledData}
                    request={()=>getDepartmentMemberList()}
                    onSelect={(selectedData: Set<string>) => {
                        const memberKeyFromModal = Array.from(selectedData)?.filter(x => allMemberIds.indexOf(x) === -1 &&selectableMemberIds.has(x)) || [];
                        setAddMemberBtnDisabled((memberKeyFromModal.length === 0));
                    }}
                    searchPlaceholder="搜索用户名、邮箱"
                 />
            </Modal>
        </>
    )
}