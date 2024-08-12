import {ActionType, ProColumns} from "@ant-design/pro-components";
import  {FC, forwardRef, useEffect, useImperativeHandle, useMemo, useRef, useState} from "react";
import {Link, useParams} from "react-router-dom";
import {App, Checkbox, Form, Select,TreeSelect} from "antd";
import {useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {useFetch} from "@common/hooks/http.ts";
import { RouterParams } from "@core/components/aoplatform/RenderRoutes.tsx";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import PageList from "@common/components/aoplatform/PageList.tsx";
import {useSystemContext} from "../../contexts/SystemContext.tsx";
import {DefaultOptionType} from "antd/es/cascader";
import { SYSTEM_SUBSCRIBER_TABLE_COLUMNS } from "../../const/system/const.tsx";
import { SystemSubscriberTableListItem, SystemSubscriberConfigFieldType, SystemSubscriberConfigHandle, SystemSubscriberConfigProps, SimpleSystemItem } from "../../const/system/type.ts";
import { EntityItem, SimpleMemberItem } from "@common/const/type.ts";
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";
import TableBtnWithPermission from "@common/components/aoplatform/TableBtnWithPermission.tsx";
import { useGlobalContext } from "@common/contexts/GlobalStateContext.tsx";
import { checkAccess } from "@common/utils/permission.ts";

const SystemInsideSubscriber:FC = ()=>{
    const { setBreadcrumb } = useBreadcrumb()
    const { modal,message } = App.useApp()
    // const [confirmLoading, setConfirmLoading] = useState(false);
    const {fetchData} = useFetch()
    const [init, setInit] = useState<boolean>(true)
    // const [tableListDataSource, setTableListDataSource] = useState<SystemSubscriberTableListItem[]>([]);
    // const [tableHttpReload, setTableHttpReload] = useState(true);
    const {systemId} = useParams<RouterParams>()
    const addRef = useRef<SystemSubscriberConfigHandle>(null)
    const pageListRef = useRef<ActionType>(null);
    const [memberValueEnum, setMemberValueEnum] = useState<{[k:string]:{text:string}}>({})
    const {partitionList} = useSystemContext()
    const {accessData} = useGlobalContext()
    const getSystemSubscriber = ()=>{
        return fetchData<BasicResponse<{subscribers:SystemSubscriberTableListItem[]}>>('project/subscribers',{method:'GET',eoParams:{project:systemId},eoTransformKeys:['apply_time']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setInit((prev)=>prev ? false : prev)
                return  {data:data.subscribers, success: true}
            }else{
                message.error(msg || '操作失败')
                return {data:[], success:false}
            }
        }).catch(() => {
            return {data:[], success:false}
        })
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

    const manualReloadTable = () => {
        // setTableHttpReload(true); // 表格数据需要从后端接口获取
        pageListRef.current?.reload()
    };

    const deleteSubscriber = (entity:SystemSubscriberTableListItem)=>{
        return new Promise((resolve, reject)=>{
            fetchData<BasicResponse<null>>('project/subscriber',{method:'DELETE',eoParams:{application:entity!.id,service:entity!.service.id,project:systemId}}).then(response=>{
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

    const openModal =async (type:'delete'|'add',entity?:SystemSubscriberTableListItem)=>{
        let title:string = ''
        let content:string|React.ReactNode = ''
        switch (type){
            case 'add':
                title='新增订阅方'
                content=<SystemSubscriberConfig ref={addRef} systemId={systemId!} partitionList={partitionList}/>
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
                    case 'delete':
                        return deleteSubscriber(entity!).then((res)=>{if(res === true) manualReloadTable()})
                }
            },
            width:600,
            okText:'确认',
            okButtonProps:{
                disabled : !checkAccess( `project.mySystem.subscriber.${type}`, accessData)
            },
            cancelText:'取消',
            closable:true,
            icon:<></>,
        })
    }

    const operation:ProColumns<SystemSubscriberTableListItem>[] =[
        {
            title: '操作',
            key: 'option',
            width: 62,
            fixed:'right',
            valueType: 'option',
            render: (_: React.ReactNode, entity: SystemSubscriberTableListItem) => [
                <TableBtnWithPermission  access="project.mySystem.subscriber.delete" key="delete" onClick={()=>{openModal('delete',entity)}} btnTitle="删除"/>,
            ],
        }
    ]

    useEffect(() => {
        setBreadcrumb([
            {
                title:<Link to={`/system/list`}>内部数据服务</Link>
            },
            {
                title:'订阅方管理'
            }
        ])
        getMemberList()
        manualReloadTable()
    }, [systemId]);

    const columns = useMemo(()=>{
        return SYSTEM_SUBSCRIBER_TABLE_COLUMNS.map(x=>{if(x.filters &&((x.dataIndex as string[])?.indexOf('applier') !== -1 || (x.dataIndex as string[])?.indexOf('approver') !== -1) ){x.valueEnum = memberValueEnum} return x})
    },[memberValueEnum])

    return (
        <PageList
            id="global_system_subscriber"
            ref={pageListRef}
            columns = {[...columns,...operation]}
            request={()=>getSystemSubscriber()}
            // dataSource={tableListDataSource}
            showPagination={false}
            addNewBtnTitle="新增订阅方"
            onAddNewBtnClick={()=>{openModal('add')}}
            addNewBtnAccess="project.mySystem.subscriber.add"
        />
    )
}

export default SystemInsideSubscriber


export const SystemSubscriberConfig = forwardRef<SystemSubscriberConfigHandle,SystemSubscriberConfigProps>((props, ref) => {
    const { message } = App.useApp()
    const { systemId,partitionList} = props
    const [form] = Form.useForm();
    const {fetchData} = useFetch()
    const [myServiceOptionList, setMyServiceOptionList] = useState<DefaultOptionType[]>()
    const [systemOptionList, setSystemOptionList] = useState<DefaultOptionType[]>()
    const [memberOptionList, setMemberOptionList] = useState<DefaultOptionType[]>()
    // const [avaliablePartitions, setAvaPartitions] = useState<Array<string>>([])
    const [subscriberSystemId, setSubscriberSystemId] = useState<string>()
    // const [partitionsList, setPartitionsList] = useState<DefaultOptionType[]>([])
    const save:()=>Promise<boolean | string> =  ()=>{
        return new Promise((resolve, reject)=>{
            form.validateFields().then((value)=>{
                fetchData<BasicResponse<null>>('project/subscriber',{method:'POST',eoBody:({...value,project:systemId}), eoParams:{project:systemId}}).then(response=>{
                    const {code,msg} = response
                    if(code === STATUS_CODE.SUCCESS){
                        message.success(msg || '操作成功！')
                        resolve(true)
                    }else{
                        message.error(msg || '操作失败')
                        reject(msg || '操作失败')
                    }
                })
            }).catch((errorInfo)=> reject(errorInfo))
        })
    }

    useImperativeHandle(ref, ()=>({
            save
        })
    )

    const getMyServiceList = ()=>{
        setMyServiceOptionList([])
        fetchData<BasicResponse<{ services: EntityItem[] }>>('simple/project/services',{method:'GET',eoParams:{project:systemId}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setMyServiceOptionList(data.services?.map((x:EntityItem)=>{return {
                    label:x.name, value:x.id
                }}))
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const getSystemList = ()=>{
        setSystemOptionList([])
        fetchData<BasicResponse<{ projects: SimpleSystemItem[] }>>('simple/apps/mine',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                const teamMap = new Map<string, unknown>();

                data.projects
                    .filter((x:SimpleSystemItem)=>x.id !== systemId)
                    .forEach((item:SimpleSystemItem) => {
                    if (!teamMap.has(item.team.id)) {
                        teamMap.set(item.team.id, {
                            title: item.team.name,
                            value: item.team.id,
                            key: item.team.id,
                            children: [],
                            selectable: false, // 第一级不可选
                            disabled:true
                                      });
                    }

                    teamMap.get(item.team.id)!.children!.push({
                        title: item.name,
                        value: item.id,
                        key: item.id,
                        selectable: true, // 子级可选
                        // partition:item.partition?.map((x:EntityItem)=>x.id) || []
                    });
                });
                setSystemOptionList(Array.from(teamMap.values()))
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    useEffect(()=>{
        subscriberSystemId && getMemberList()
        form.setFieldValue('applier',null)
    },[subscriberSystemId])

    const getMemberList = ()=>{
        setMemberOptionList([])
        fetchData<BasicResponse<{ members: SimpleMemberItem[] }>>('simple/project/members',{method:'GET',eoParams:{project:subscriberSystemId}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setMemberOptionList(data.members?.map((x:SimpleMemberItem)=>{return {
                    label:x.name, value:x.id
                }}))
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    // const getPartitionList = (applicationId:string)=>{
    //     fetchData<BasicResponse<{partitions:(EntityItem &{ serviceNum:number})[]}>>('simple/application/partitions',{method:'GET',eoParams:{application:applicationId},eoTransformKeys:['service_num']}).then(response=>{
    //         const {code,data,msg} = response
    //         if(code === STATUS_CODE.SUCCESS){
    //             setPartitionsList(data.partitions?.map((x:SimpleTeamItem)=>({label:x.name, value:x.id})))
    //         }else{
    //             message.error(msg || '操作失败')
    //         }
    //     })
    // }

    useEffect(() => {
        getMyServiceList()
        getSystemList()
    }, [systemId]);

    

    // const partitionsList = useMemo(()=>{
    //     const newList = partitionList?.filter((x:EntityItem)=> avaliablePartitions && avaliablePartitions?.indexOf(x.id) !== -1)?.map((x:EntityItem)=>({label:x.name, value:x.id})) || []
    //     if(newList.length === 1){
    //         form.setFieldValue('partition',[newList[0].value])
    //     }
    //     return newList
    // },[partitionList, avaliablePartitions])

    return  (<WithPermission access="project.mySystem.subscriber.add">
        <Form
            layout='vertical'
            labelAlign='left'
            scrollToFirstError
            form={form}
            className="mx-auto  "
            name="systemInsideSubscriber"
            // labelCol={{ offset:1, span: 4 }}
            // wrapperCol={{ span: 19}}
            autoComplete="off"
        >
            <Form.Item<SystemSubscriberConfigFieldType>
                label="订阅服务"
                name="service"
                rules={[{ required: true, message: '必填项' }]}
            >
                <Select className="w-INPUT_NORMAL" options={myServiceOptionList}  placeholder="请选择"/>
            </Form.Item>

            <Form.Item<SystemSubscriberConfigFieldType>
                label="订阅方"
                name="subscriber"
                rules={[{ required: true, message: '必填项' }]}
            >
                <TreeSelect
                    className="w-INPUT_NORMAL" 
                    dropdownStyle={{ maxHeight: 400, overflow: 'auto' }}
                    treeData={systemOptionList}
                    placeholder="请选择"
                    treeDefaultExpandAll
                    // onSelect={(_:unknown,selectedItemInfo)=>{setAvaPartitions(selectedItemInfo.partition);setSubscriberSystemId(_)}}
                    onSelect={(_:unknown)=>{setSubscriberSystemId(_)}}
                />
            </Form.Item>

            <Form.Item
                label="申请人"
                name="applier"
                rules={[{ required: true, message: '必填项' }]}
            >
                <Select className="w-INPUT_NORMAL"  options={memberOptionList}  placeholder="请选择"/>
            </Form.Item>

            <Form.Item
                label="可用环境"
                name="partition"
                rules={[{ required: true, message: '必填项' }]}
            >
                <Checkbox.Group className="flex flex-col gap-[8px] mt-[5px]" options={partitionList.map((x:EntityItem)=>({label:x.name, value:x.id}))} />
            </Form.Item>
        </Form>
    </WithPermission>)
})