import PageList from "@common/components/aoplatform/PageList.tsx"
import {ActionType, ProColumns} from "@ant-design/pro-components";
import  {FC, useEffect, useMemo, useRef, useState} from "react";
import {Link, useParams} from "react-router-dom";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {App, Button, Modal} from "antd";
import {TransferTableHandle} from "@common/components/aoplatform/TransferTable.tsx";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {EntityItem, MemberItem} from "@common/const/type.ts";
import { TeamMemberTableListItem } from "../../const/team/type.ts";
import {  TEAM_MEMBER_TABLE_COLUMNS } from "../../const/team/const.tsx";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import { checkAccess } from "@common/utils/permission.ts";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";
import MemberTransfer from "@common/components/aoplatform/MemberTransfer.tsx";
import { DepartmentListItem } from "../../const/member/type.ts";
import { addMemberToDepartment, getDepartmentWithMember } from "../user/UserList.tsx";
import {v4 as uuidv4} from 'uuid'
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";

const TeamInsideMember:FC = ()=>{
    const [searchWord, setSearchWord] = useState<string>('')
    const { setBreadcrumb} = useBreadcrumb()
    const { modal,message } = App.useApp()
    const {fetchData} = useFetch()
    const {teamId} = useParams<RouterParams>();
    const addRef = useRef<TransferTableHandle<TeamMemberTableListItem>>(null)
    const pageListRef = useRef<ActionType>(null);
    const [allMemberIds, setAllMemberIds] = useState<string[]>([])
    const {accessData} = useGlobalContext()
    const [selectableMemberIds,setSelectableMemberIds] = useState<Set<string>>(new Set())
    const [addMemberBtnLoading, setAddMemberBtnLoading] = useState<boolean>(false)
    const [modalVisible, setModalVisible] = useState<boolean>(false)
    const [addMemberBtnDisabled, setAddMemberBtnDisabled] = useState<boolean>(true)
    const [allMemberSelectedDepartIds, setAllMemberSelectedDepartIds] = useState<string[]>([])

    const operation:ProColumns<TeamMemberTableListItem>[] =[
        {
            title: '操作',
            key: 'option',
            width: 76,
            fixed:'right',
            valueType: 'option',
            render: (_: React.ReactNode, entity: TeamMemberTableListItem) => [
                !entity.role?.includes('负责人') &&<TableBtnWithPermission disabled={!entity.canDelete} tooltip="暂无权限" access="team.myTeam.member.edit" key="removeMember" onClick={()=>{openModal('remove',entity)}} btnTitle="移出团队"/>]
        }
    ]

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
      
    const getMemberList = ()=>{
        return fetchData<BasicResponse<{members:TeamMemberTableListItem}>>('team/members',{method:'GET',eoParams:{keyword:searchWord, team:teamId},eoTransformKeys:['user_group','attach_time','user_id','can_delete']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                if(!searchWord){
                    setAllMemberIds(data.members?.map((x:TeamMemberTableListItem)=>x.userId) || [])
                }
                return  {data:data.members, success: true}
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        }).catch(() => {
            return {data:[], success:false}
        })
    }

    const addMember = (selectableMemberIds:Set<string>)=>{
        setAddMemberBtnLoading(true)
        const keyFromModal = addRef.current?.selectedRowKeys()
        const memberKeyFromModal = keyFromModal?.filter(x => allMemberIds.indexOf(x as string) === -1 && selectableMemberIds.has(x)) || [];
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>('team/member',{method:'POST' ,eoBody:({users:memberKeyFromModal}),eoParams:{team:teamId}}).then(response=>{
            const {code,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                message.success(msg || '操作成功！')
                manualReloadTable()
                cleanModalData()
                resolve(true)
            }else{
                message.error(msg || '操作失败')
                reject(msg || '操作失败')
            }
        }).catch((errorInfo)=> reject(errorInfo)).finally(()=>setAddMemberBtnLoading(false))
    })
    }


    const removeMember = (entity:TeamMemberTableListItem) =>{
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>(`team/member`,{method:'DELETE',eoParams:{team:teamId,user:entity.userId}}).then(response=>{
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

    const cleanModalData = ()=>{
        setModalVisible(false);setAddMemberBtnDisabled(true);setAddMemberBtnLoading(false)
    }

    const openModal = async (type:'add'|'remove',entity?:TeamMemberTableListItem)=>{
        let title:string = ''
        let content:string|React.ReactNode = ''
        switch(type){
            case 'add':
                setModalVisible(true)
                setAddMemberBtnDisabled(true)
                setAddMemberBtnLoading(false)
                return
            case 'remove':
                title='移除成员'
                content=<span>确定删除成员<span className="text-status_fail"></span>？此操作无法恢复，确认操作？</span>
               break
        }
        
        modal.confirm({
            title,
            content,
            onOk:()=>{
                    return  removeMember(entity!).then((res)=>{if(res === true) manualReloadTable()})
            },
            width:600,
            okText:'确认',
            okButtonProps:{
                disabled: !checkAccess(`team.myTeam.member.edit`,accessData)
            },
            cancelText:'取消',
            closable:true,
            icon:<></>
        })
    }


    const manualReloadTable = () => {
        pageListRef.current?.reload()
    };

    useEffect(() => {
        setBreadcrumb([
            {title:<Link to="/team/list">团队</Link>},
            {title:'成员'}
        ])
        manualReloadTable()
    }, [teamId]);

    const treeDisabledData = useMemo(()=>{ return [...allMemberIds,...allMemberSelectedDepartIds]},[allMemberIds,allMemberSelectedDepartIds])
    
    return (
        <>
        <PageList
            id="global_team_member"
            ref={pageListRef}
            columns = {[...TEAM_MEMBER_TABLE_COLUMNS,...operation]}
            request={()=>getMemberList()}
            primaryKey="userId"
            addNewBtnTitle="添加成员"
            searchPlaceholder="输入姓名查找"
            onAddNewBtnClick={()=>{openModal('add')}}
            addNewBtnAccess="team.myTeam.member.add"
            tableClickAccess="team.myTeam.member.edit"
            onSearchWordChange={(e)=>{setSearchWord(e.target.value)}}
        />
        <Modal
                   title="添加成员"
                   open={modalVisible}
                   destroyOnClose={true}
                   width={900}
                   onCancel={() => cleanModalData()}
                   maskClosable={false}
                   footer={[
                       <Button key="back" onClick={() => cleanModalData()}>
                           取消
                       </Button>,
                       <WithPermission access="team.myTeam.member.add"><Button
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
        </>
    )

}
export default TeamInsideMember