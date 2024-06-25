/*
 * @Date: 2024-04-08 15:05:34
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-05 15:17:55
 * @FilePath: \frontend\packages\core\src\pages\user\UserPage.tsx
 */
import {DataNode} from "antd/es/tree";
import {Outlet, useNavigate, useParams} from "react-router-dom";
import  {forwardRef, useEffect, useImperativeHandle, useMemo, useRef, useState} from "react";
import {App, Button, Empty, Form, Input} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import styles from '../system/SystemList.module.css'
import {useFetch} from "@common/hooks/http.ts";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import { UserGroupModalHandle, UserGroupModalProps, FieldType, UserGroupItem } from "../../const/user/types.ts";
import { RouterParams } from "@core/components/aoplatform/RenderRoutes.tsx";
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";
import { PERMISSION_DEFINITION } from "@common/const/permissions.ts";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";
import { checkAccess } from "@common/utils/permission.ts";
import GroupTree, { GroupTreeHandle } from "@common/components/aoplatform/GroupTree.tsx";
import { useBreadcrumb } from "@common/contexts/BreadcrumbContext.tsx";

export const UserGroupModal = forwardRef<UserGroupModalHandle,UserGroupModalProps>((props,ref)=>{
    const { message } = App.useApp()
    const [form] = Form.useForm();
    const {type,entity} = props
    const {fetchData} = useFetch()

    const save:()=>Promise<boolean | string> =  ()=>{
        //console.log(type)
        const url:string = 'user/group'
        let method:string
        switch (type){
            case 'add':
                method = 'POST'
                break;
            case 'edit':
                method = 'PUT'
                break;
        }
        return new Promise((resolve, reject)=>{
            form.validateFields().then((value)=>{
                fetchData<BasicResponse<null>>(url,{method,eoBody:({name:value.name}), ...(type ==='edit' ? {eoParams:{userGroup:value.id}}:{}),eoTransformKeys:['userGroup']}).then(response=>{
                    const {code,msg} = response
                    if(code === STATUS_CODE.SUCCESS){
                        message.success(msg || '操作成功！')
                        resolve(true)
                    }else{
                        message.error(msg || '操作失败')
                        reject(msg || '操作失败')
                    }
                }).catch((errorInfo)=> reject(errorInfo))
            }).catch((errorInfo)=> reject(errorInfo))
        })
    }

    useImperativeHandle(ref, ()=>({
            save
        })
    )

    useEffect(() => {
        
        if(type === 'edit'){
            form.setFieldsValue(entity)
        }
    }, []);

    return (<WithPermission access={type === 'add' ? 'system.user.group.add' : 'system.user.group.edit'}>
        <Form
            layout='vertical'
            labelAlign='left'
            scrollToFirstError
            form={form}
            className="mx-auto mt-mbase "
            name="userPage"
            // labelCol={{ span: 7 }}
            // wrapperCol={{ span: 17}}
            autoComplete="off"
        >

            {type === 'edit' &&
                <Form.Item<FieldType>
                    label="ID"
                    name="id"
                    hidden
                    rules={[{ required: true, message: '必填项',whitespace:true  }]}
                >
                    <Input className="w-INPUT_NORMAL" placeholder="ID"/>
                </Form.Item>
            }
                <Form.Item<FieldType>
                    label="用户组名称"
                    name="name"
                    rules={[{ required: true, message: '必填项' ,whitespace:true }]}
                >
                    <Input className="w-INPUT_NORMAL" placeholder="请输入"/>
                </Form.Item>
        </Form>
    </WithPermission>)
})

const UserPage = ()=>{
    const {modal,message} = App.useApp()
    const navigate = useNavigate()
    const {fetchData} = useFetch()
    const {userGroupId} = useParams<RouterParams>()
    const [userGroup, setUserGroup] = useState<UserGroupItem[]>()
    const [selectedGroupIds, setSelectedGroupIds] = useState<string[]>([userGroupId!])
    const {accessData} = useGlobalContext()
    const groupRef = useRef<GroupTreeHandle>(null)
    const {setBreadcrumb} = useBreadcrumb()


    const dropdownMenu = (entity:UserGroupItem) => [
        // {
        //     key: 'add',
        //     label: (
        //         <WithPermission access="system.user.group.add"><Button className="h-[32px] border-none p-0 flex items-center bg-transparent "  onClick={()=>{console.log(entity.id);groupRef.current?.startAdd('addPeer',entity)}}>
        //             添加用户组
        //         </Button></WithPermission>
        //     ),
        // },
        {
            key: 'edit',
            label: (
                <WithPermission access="system.user.group.edit"><Button className=" h-[32px] border-none p-0 flex items-center bg-transparent "  onClick={()=>{groupRef.current?.startEdit(entity.id)}}>
                    编辑用户组
                </Button></WithPermission>
            ),
        },
        {
            key: 'delete',
            label: (
                <WithPermission access="system.user.group.delete"><Button className="h-[32px] border-none p-0 flex items-center bg-transparent "  onClick={()=>openModal('delete',entity)}>
                    删除
                </Button></WithPermission>
            ),
        },
    ];

    const handleEditUserGroup = (type:'rename'|'addChild'|'addPeer',entity:UserGroupItem, name:string)=>{
        const url:string = 'user/group'
        let method:string
        switch (type){
            case 'rename':
                method = 'PUT'
                break;
            case 'addPeer':
                method = 'POST'
                break;
        }

        return new Promise((resolve, reject)=>{
            return fetchData<BasicResponse<null>>(url,{method,eoBody:({name:name}), ...(type ==='rename' ? {eoParams:{userGroup:entity.id}}:{}),eoTransformKeys:['userGroup']}).then(response=>{
                const {code,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功！')
                    getUserGroup()
                    return resolve(true)
                }else{
                    message.error(msg || '操作失败')
                    return reject(msg || '操作失败')
                }
            }).catch((errorInfo)=> reject(errorInfo))
        })
    }

    const treeData = useMemo(() => {
        const loop = (data: UserGroupItem[]): (DataNode & UserGroupItem)[] =>
            data?.map((item) => {
                return {
                    ...item,
                    title:item.name,
                    key: item.id,
                };
            });
        return userGroup ? loop(userGroup) :[];
    }, [userGroup]);


    const getUserGroup = ()=>{
        setUserGroup([])
        fetchData<BasicResponse<{userGroups:UserGroupItem[]}>>('user/groups',{method:'GET',eoTransformKeys:['user_groups']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setUserGroup(data.userGroups)
                if(userGroup === undefined) {navigate(`/user/list/${data.userGroups[0].id}`); return}
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        })
    }

    const deleteUserGroup = (id:string)=>{
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>('user/group',{method:'DELETE',eoParams:{userGroup:id},eoTransformKeys:['userGroup']}).then(response=>{
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

    const openModal = (type:'add'|'edit'|'delete',entity?:{id:string, name:string})=>{
        modal.confirm({
            title:'删除',
            content:'该数据删除后将无法找回，请确认是否删除？',
            onOk:()=>deleteUserGroup(entity!.id!).then((res)=>{if(res === true) getUserGroup()}),
            width:600,
            okText:'确认',
            okButtonProps:{
                disabled : !checkAccess(`system.user.group.${type}` as keyof typeof PERMISSION_DEFINITION[0], accessData)
            },
            cancelText:'取消',
            closable:true,
            icon:<></>,
        })
    }

    useEffect(() => {
        setBreadcrumb([
            {title:'用户组'}
        ])
       getUserGroup()
    }, []);

    useEffect(()=>{
        userGroupId && setSelectedGroupIds([userGroupId])
    },[userGroupId])

    return (
        <div className="flex flex-1 h-full">
            <div className={styles['system-tree'] + ` w-[200px] border-0 border-solid border-r-[1px] border-r-BORDER`}>
                    <GroupTree
                        ref={groupRef}
                        icon={<></>}
                        groupData={treeData}
                        selectedKeys={selectedGroupIds}
                        addBtnName={<><PlusOutlined className='mr-[2px]' />添加用户组</>}
                        addBtnAccess="system.user.group.add"
                        onSelect={(selectedKeys) => {
                            navigate(`/user/list/${selectedKeys[0]}`)
                        }}
                        treeNameSuffixKey='usage'
                        dropdownMenu={dropdownMenu}
                        onEditGroup={handleEditUserGroup}
                        withMore
                        placeholder="搜索用户组"
                    />
                  </div>
            {treeData?.length > 0 ? <div className="w-[calc(100%-200px)]">
                <Outlet context={{refreshGroup:()=>getUserGroup()}} />
            </div> :<div className="block w-full h-full align-middle"><Empty className="mt-[20%]" image={Empty.PRESENTED_IMAGE_SIMPLE}/></div>}
        </div>)
}
export default UserPage;