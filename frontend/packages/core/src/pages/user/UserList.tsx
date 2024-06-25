/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-07 17:36:52
 * @FilePath: \frontend\packages\core\src\pages\user\UserList.tsx
 */
import {App, Button, Modal} from "antd";
import PageList from "@common/components/aoplatform/PageList.tsx";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {TransferTableHandle} from "@common/components/aoplatform/TransferTable.tsx";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import MemberTransfer from "@common/components/aoplatform/MemberTransfer.tsx";
import  {useEffect, useMemo, useRef, useState} from "react";
import {useOutletContext, useParams} from "react-router-dom";
import {ActionType, ProColumns} from "@ant-design/pro-components";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {EntityItem, MemberItem} from "@common/const/type.ts";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import { USER_LIST_COLUMNS } from "../../const/user/const.tsx";
import { DepartmentListItem } from "../../const/member/type.ts";
import { handleDepartmentListToFilter } from "@common/utils/dataTransfer.ts";
import { checkAccess } from "@common/utils/permission.ts";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";
import { ColumnFilterItem } from "antd/es/table/interface";
import {v4 as uuidv4} from 'uuid'
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";

type DepartmentWithMemberItem = 
    DepartmentListItem & {type:'member'|'department', children?: Array< DepartmentListItem & {type:'member'|'department'} | {
        id:string,
        name:string,
        avatar:string,
        email:string} & {type:'member'|'department'}>
}


export const getDepartmentWithMember = (department:(DepartmentListItem & {type?:'department'|'member'})[],departmentMap:Map<string, (MemberItem & {type:'department'|'member'})[]>) : (DepartmentWithMemberItem | undefined)[] =>{
    return department.map((x:DepartmentListItem & {type?:'department'|'member'})=>{
        const res =  ({
            ...x,
            key:x.id,
            title:x.name,
            type: x.type || 'department',
            children:((x.type === 'member' || (!x.children||x.children.length === 0 )&& (!departmentMap.get(x.id) || departmentMap.get(x.id)!.length === 0))? undefined : [...(x.children && x.children.length > 0 ? getDepartmentWithMember(x.children,departmentMap) : []),...departmentMap.get(x.id) || []])
        });
        return res}).filter(node=>node.type === 'member' ||( node.children && node.children.length > 0))
}

export const addMemberToDepartment = (departmentMap: Map<string, (MemberItem & {type:'department'|'member'})[]>, departmentId: string, member: MemberItem) => {
    const members = departmentMap.get(departmentId) || [];
    members.push({...member, type: 'member'});
    departmentMap.set(departmentId, members);
  }

const UserList = ()=>{
    const { userGroupId }  = useParams<RouterParams>();
    const [searchWord, setSearchWord] = useState<string>('')
    const { modal,message} = App.useApp()
    // const [confirmLoading, setConfirmLoading] = useState(false);
    const addRef = useRef<TransferTableHandle<MemberItem>>(null)
    const [addMemberBtnDisabled, setAddMemberBtnDisabled] = useState<boolean>(true)
    const [init, setInit] = useState<boolean>(true)
    const [tableHttpReload, setTableHttpReload] = useState(true);
    const [tableListDataSource, setTableListDataSource] = useState<MemberItem[]>([]);
    const {fetchData} = useFetch()
    const [allMemberIds,setAllMemberIds] = useState<string[]>([])
    const [allMemberSelectedDepartIds, setAllMemberSelectedDepartIds] = useState<string[]>([])
    const {refreshGroup} = useOutletContext<{refreshGroup:()=>void}>()
    const [departmentValueEnum,setDepartmentValueEnum] = useState<ColumnFilterItem[] >([])
    const [selectableMemberIds,setSelectableMemberIds] = useState<Set<string>>(new Set())
    const [addMemberBtnLoading, setAddMemberBtnLoading] = useState<boolean>(false)
    const [modalVisible, setModalVisible] = useState<boolean>(false)
    const {accessData} = useGlobalContext()
    const getUserList = ()=>{
        if(!tableHttpReload){
            setTableHttpReload(true)
            return Promise.resolve({
                data: tableListDataSource,
                success: true,
            });
        }
        return fetchData<BasicResponse<{members:MemberItem[]}>>('user/group/members',{method:'GET',eoParams:{userGroup:userGroupId,keyword:searchWord},eoTransformKeys:['userGroup']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setTableListDataSource(data.members)
                setInit((prev)=>prev ? false : prev)
                setTableHttpReload(false)
                if(!searchWord) setAllMemberIds(data.members?.map((x:MemberItem)=>x.id) || [])
                return  {data:data.members, success: true}
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        }).catch(() => {
            return {data:[], success:false}
        })
    }


    const getDepartmentMemberList = () => {
        const topDepartmentId:string = uuidv4()
        return Promise.all([
          fetchData<BasicResponse<{department:DepartmentListItem}>>('simple/departments', {method:'GET'}),
          fetchData<BasicResponse<{members:MemberItem}>>('simple/member', {method:'GET', eoParams:{}, eoTransformKeys:[]})
        ]).then(([departmentResponse, memberResponse])=>{
            const departmentMap = new Map<string, (MemberItem & {type:'department'|'member'})[]>();
            memberResponse.data.members.forEach((member: MemberItem) => {
                setSelectableMemberIds((pre)=>{pre.add(member.id);return pre})
                member = {...member, title:member.name, key:member.id}
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
      

    const operation:ProColumns<MemberItem>[] =[
        {
            title: '操作',
            key: 'option',
            width: 62,
            fixed:'right',
            valueType: 'option',
            render: (_: React.ReactNode, entity: MemberItem) => [
                <TableBtnWithPermission  access="system.user.member.delete" key="delete" onClick={()=>{openModal('delete',entity!)}} btnTitle="移除"/>
            ],
        }
    ]

    const treeDisabledData = useMemo(()=>{ return [...allMemberIds,...allMemberSelectedDepartIds]},[allMemberIds,allMemberSelectedDepartIds])

    const { setBreadcrumb } = useBreadcrumb()
    const pageListRef = useRef<ActionType>(null);

    const manualReloadTable = () => {
        setTableHttpReload(true); // 表格数据需要从后端接口获取
        pageListRef.current?.reload()
    };

    const addMember = (selectableMemberIds:Set<string>)=>{
        setAddMemberBtnLoading(true)
        const keyFromModal = addRef.current?.selectedRowKeys()
        const memberKeyFromModal = keyFromModal?.filter(x => allMemberIds.indexOf(x as string) === -1 && selectableMemberIds.has(x)) || [];
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>('user/group/member',{method:'POST',eoParams:{userGroup:userGroupId},eoBody:({ids:Array.from(memberKeyFromModal) || []}),eoTransformKeys:['userGroup']}).then(response=>{
                const {code,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功！')
                    manualReloadTable()
                    refreshGroup && refreshGroup()
                    resolve(true)
                }else{
                    message.error(msg || '操作失败')
                    reject(msg || '操作失败')
                }
            }).catch((errorInfo)=> reject(errorInfo)).finally(()=>setAddMemberBtnLoading(false))
        })
    }

    const deleteUser = (entity:MemberItem)=>{
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>('user/group/member',{method:'DELETE',eoParams:{userGroup:userGroupId,member:entity!.id},eoTransformKeys:['userGroup']}).then(response=>{
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


    const openModal = (type:'addMember'|'delete',entity?:MemberItem)=>{
        let title:string = ''
        let content:string|React.ReactNode = ''
        switch (type){
            case 'addMember':
                setModalVisible(true)
                setAddMemberBtnDisabled(true)
                setAddMemberBtnLoading(false)
                return;
            case 'delete':
                title='移除成员'
                content=<span>确定删除成员<span className="text-status_fail"></span>？此操作无法恢复，确认操作？</span>
                break;
        }

        modal.confirm({
            title,
            content,
            onOk:()=>{
                    return deleteUser(entity!).then((res)=>{if(res === true) {refreshGroup && refreshGroup() ; manualReloadTable()}})
            },
            width:600,
            okText:'确认',
            cancelText:'取消',
            okButtonProps:{
                disabled : !checkAccess( `system.user.member.delete`, accessData)
            },
            closable:true,
            icon:<></>,
        })
    }

    useEffect(() => {
        !init && manualReloadTable()
    }, [userGroupId]);

    useEffect(()=>{
        setBreadcrumb([{ title: '用户组'}])
        getDepartmentList()
    },[])

    
    const cleanModalData = ()=>{
        setModalVisible(false);setAddMemberBtnDisabled(true);setAddMemberBtnLoading(false)
    }

    const getDepartmentList = async ()=>{
        setDepartmentValueEnum([])
        const {code,data,msg}  = await fetchData<BasicResponse<{ department: DepartmentListItem }>>('simple/departments',{method:'GET'})
        if(code === STATUS_CODE.SUCCESS){
            const tmpValueEnum:ColumnFilterItem[]   = [{text:data.department.name, value:data.department.id,children:handleDepartmentListToFilter(data.department.children)}]
            setDepartmentValueEnum(tmpValueEnum)
        }else{
            message.error(msg || '操作失败')
        }
    }

    const columns = useMemo(()=>{
        return USER_LIST_COLUMNS.map(x=>{if((x.dataIndex as string[])?.indexOf('department') !== -1 ){
            x.filters = departmentValueEnum
            x.onFilter = (value: string, record) => {
                return value ? record.department?.filter(x=>x.id === value).length > 0 : true
            }
        } return x})
    },[departmentValueEnum])

    return (<><PageList
        id="global_user"
        ref={pageListRef}
        columns={[...columns, ...operation]}
        request={()=>getUserList()}
        addNewBtnTitle="添加成员"
        searchPlaceholder="输入用户名、邮箱查找成员"
        onAddNewBtnClick={() => {
            openModal('addMember')
        }}
        addNewBtnAccess="system.user.member.add"
        onSearchWordChange={(e) => {
            setSearchWord(e.target.value)
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
                    <WithPermission access="system.user.member.add"><Button
                        key="submit"
                        type="primary"
                        disabled={addMemberBtnDisabled}
                        loading={addMemberBtnLoading}
                        onClick={()=>addMember(selectableMemberIds as Set<string>)}
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
            </>)

}
export default UserList;