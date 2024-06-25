import {Alert, App, Divider, Form, Input} from "antd";
import PageList from "@common/components/aoplatform/PageList.tsx";
import  {forwardRef, useEffect, useImperativeHandle, useRef, useState} from "react";
import {ActionType, ProColumns} from "@ant-design/pro-components";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import { ROLE_TABLE_COLUMNS } from "../../const/role/const.tsx";
import { RoleModalContentHandle, RoleModalContentProps, RoleTableListItem } from "../../const/role/type.ts";
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import { PERMISSION_DEFINITION } from "@common/const/permissions.ts";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";
import { checkAccess } from "@common/utils/permission.ts";


const RoleModalContent = forwardRef<RoleModalContentHandle,RoleModalContentProps>((props, ref)=>{
    const { message } = App.useApp()
    const {type,entity} = props
    const [form] = Form.useForm();
    const {fetchData} = useFetch()
    const save:()=>Promise<boolean | string> =  ()=>{
        return new Promise((resolve, reject)=>{
            form.validateFields().then((value)=>{
                fetchData<BasicResponse<null>>(type === 'add'? 'manage/role':'manage/role',{method:type === 'add'? 'POST' : 'PUT',eoBody:({name:value.name}),eoParams:{id:value!.id}}).then(response=>{
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

    return (
        <WithPermission access={type === 'edit' ? 'system.role.self.edit':'system.role.self.add'}>
        <Form
            layout='vertical'
            labelAlign='left'
            scrollToFirstError
            form={form}
            className="mx-auto "
            name="RoleList"
            // labelCol={{ offset:1, span: 4 }}
            // wrapperCol={{ span: 19}}
            autoComplete="off"
        >

            {type === 'edit' && <Form.Item<RoleTableListItem>
                    label="ID"
                    name="id"
                    hidden
                    rules={[{ required: true, message: '必填项' ,whitespace:true }]}
                >
                    <Input className="w-INPUT_NORMAL" placeholder="ID"/>
                </Form.Item>}

                <Form.Item<RoleTableListItem>
                    label="角色名称"
                    name="name"
                    rules={[{ required: true, message: '必填项' ,whitespace:true }]}
                >
                    <Input className="w-INPUT_NORMAL" placeholder="请输入角色名称"/>
                </Form.Item>
        </Form></WithPermission>)
})

const RoleList = ()=>{
    // const [searchWord, setSearchWord] = useState<string>('')
    const { modal,message } = App.useApp()
    const { setBreadcrumb } = useBreadcrumb()
    // const [confirmLoading, setConfirmLoading] = useState(false);
    const {fetchData} = useFetch()
    const [init, setInit] = useState<boolean>(true)
    const pageListRef = useRef<ActionType>(null);
    const addRef = useRef<RoleModalContentHandle>(null)
    const editRef = useRef<RoleModalContentHandle>(null)
    const {accessData} = useGlobalContext()

    const operation:ProColumns<RoleTableListItem>[] =[
        {
            title: '操作',
            key: 'option',
            width: 93,
            fixed:'right',
            valueType: 'option',
            render: (_: React.ReactNode, entity: RoleTableListItem) => [
                <TableBtnWithPermission  access="system.role.self.edit" key="edit" onClick={()=>{openModal('edit',entity)}} btnTitle="编辑"/>,
                <Divider type="vertical" className="mx-0"  key="div1" />,
                <TableBtnWithPermission  access="system.role.self.delete" key="delete" onClick={()=>{openModal('delete',entity)}} btnTitle="删除"/>,
            ],
        }
    ]

    const getRoleList = ()=>{
        return fetchData<BasicResponse<{roles:RoleTableListItem[]}>>('manage/roles',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setInit((prev)=>prev ? false : prev)
                return  {data:data.roles, success: true}
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        }).catch(() => {
            return {data:[], success:false}
        })
    }

    const deleteRole = (entity:RoleTableListItem)=>{
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>(`manage/role`,{method:'DELETE',eoParams:{id:entity.id}}).then(response=>{
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

    const manualReloadTable = () => {
        pageListRef.current?.reload()
    };

    const isActionAllowed = (type:'add'|'edit'|'delete') => {
        
        const permission = `system.role.self.${type}` as keyof typeof PERMISSION_DEFINITION[0] ;
        
        return !checkAccess(permission, accessData);
        };

    const openModal = (type:'add'|'edit'|'delete',entity?:RoleTableListItem)=>{
        let title:string = ''
        let content:string|React.ReactNode = ''
        switch (type){
            case 'add':
                title='添加角色'
                content=<RoleModalContent ref={addRef} type={type} />
                break;
            case 'edit':
                title='编辑角色'
                content=<RoleModalContent ref={editRef} type={type} entity={entity}/>
                break;
            case 'delete':
                title='删除'
                content='该数据删除后将无法找回，请确认是否删除？'
                break;
        }

        modal.confirm({
            title,
            content,
            onOk:()=>{
                switch (type){
                    case 'add':
                        return addRef.current?.save().then((res)=>{if(res === true) manualReloadTable()})
                    case 'edit':
                        return editRef.current?.save().then((res)=>{if(res === true) manualReloadTable()})
                    case 'delete':
                        return deleteRole(entity!).then((res)=>{if(res === true) manualReloadTable()})
                }
            },
            width:600,
            okText:'确认',
            okButtonProps:{
                disabled:isActionAllowed(type)
            },
            cancelText:'取消',
            closable:true,
            icon:<></>,
        })
    }

    useEffect(() => {
        setBreadcrumb([
            {
                title: '自定义角色'}])
    }, []);

    return (<>
        <Alert showIcon banner={true}  className="b-none m-btnbase mb-0 rounded" type="info" message="自定义角色用于配置系统下需要参与的角色组成，以下列表为平台可用的角色列表，各个系统的各个角色以及成员组成请到自定义角色入口进行配置。"/>
        <div className="h-[calc(100%-40px)]">
        <PageList
            id="global_role"
            ref={pageListRef}
            columns={[...ROLE_TABLE_COLUMNS, ...operation]}
            request={()=>getRoleList()}
            addNewBtnTitle="添加角色"
            // searchPlaceholder="输入角色名称进行搜索"
            onAddNewBtnClick={() => {
               openModal('add')
            }}
            addNewBtnAccess="system.role.self.add"
            // onSearchWordChange={(e) => {
            //     setSearchWord(e.target.value)
            // }}
            onRowClick={(row:RoleTableListItem)=>openModal('edit',row)}
            tableClickAccess="system.role.self.edit"
        />
        </div>
    </>)
}
export default RoleList;