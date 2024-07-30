import  {forwardRef, useEffect, useImperativeHandle, useState} from "react";
import {App, Button, Checkbox, CheckboxOptionType, Divider, Form, Input, Row, Select} from "antd";
import { Link, useNavigate, useParams} from "react-router-dom";
import {RouterParams} from "@core/components/aoplatform/RenderRoutes.tsx";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {DefaultOptionType} from "antd/es/cascader";
import { MemberItem, SimpleTeamItem, EntityItem, TeamSimpleMemberItem} from "@common/const/type.ts";
import { v4 as uuidv4 } from 'uuid'
import { SystemConfigFieldType, SystemConfigHandle } from "../../const/system/type.ts";
import { validateUrlSlash } from "@common/utils/validate.ts";
// import WithPermission from "@common/components/aoplatform/WithPermission.tsx";
import { useBreadcrumb } from "@common/contexts/BreadcrumbContext.tsx";
import { useSystemContext } from "../../contexts/SystemContext.tsx";

const SystemConfig = forwardRef<SystemConfigHandle>((_,ref) => {
    const { message,modal } = App.useApp()
    const { teamId, systemId } = useParams<RouterParams>();
    const [onEdit, setOnEdit] = useState<boolean>(!!teamId)
    const [form] = Form.useForm();
    const {fetchData} = useFetch()
    const [teamOptionList, setTeamOptionList] = useState<DefaultOptionType[]>()
    const [memberOptionList, setMemberOptionList] = useState<DefaultOptionType[]>()
    const navigate = useNavigate();
    const [partitionOption,setPartitionOption] = useState<CheckboxOptionType[]>([])
    const [currentTeamId, setCurrentTeamId] = useState<string>(teamId || '')
    const {setBreadcrumb} = useBreadcrumb()
    const { setSystemInfo,setPartitionList} = useSystemContext()

    useImperativeHandle(ref, () => ({
        save:onFinish
    }));

    // 获取表单默认值
    const getSystemInfo = () => {
        fetchData<BasicResponse<{ project: SystemConfigFieldType }>>('project/info',{method:'GET',eoParams:{project:systemId},eoTransformKeys:['team_id']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setTimeout(()=>{form.setFieldsValue({...data.project,organization:data.project.organization.id,team:data.project.team.id,master:data.project.master.id,partition:data.project.partition?.map((x:EntityItem)=>x.id)})},0)
                setCurrentTeamId(data.project.team.id)
            }else{
                message.error(msg || '操作失败')
            }
        })
    };

    useEffect(()=>{
        const newPartitions = (teamOptionList as (Array<(DefaultOptionType & { id: string; availablePartitions: EntityItem[] })>) )?.find(x => x.id === currentTeamId)?.availablePartitions?.map((p: EntityItem) => ({ label: p.name, value: p.id })) || []
        setPartitionOption(newPartitions)
        if(!newPartitions || newPartitions.length === 0){
            form.setFieldValue('partition',[])
            return
        }
        const selectedPartitions =form.getFieldValue('partition')
        if(selectedPartitions && selectedPartitions?.length > 0 ){
            form.setFieldValue('partition',newPartitions.filter(x=> selectedPartitions.indexOf(x.value) !== -1).map(x=>x.value))
        }
    },[currentTeamId,teamOptionList])

    const onFinish:()=>Promise<boolean|string> = () => {
        return form.validateFields().then((value)=>{
            return fetchData<BasicResponse<{project:{id:string}}>>(systemId === undefined? 'team/project':'project/info',{method:systemId === undefined? 'POST' : 'PUT',eoParams: {...(systemId === undefined ? {team:value.team} :{project:systemId})},eoBody:({...value,prefix:value.prefix?.trim()})}).then(response=>{
                const {code,data,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功！')
                    setSystemInfo(data.project)
                    setPartitionList(data.project.partition)
                    return Promise.resolve(true)
                }else{
                    message.error(msg || '操作失败')
                    return Promise.reject(msg || '操作失败')
                }
            }).catch((errorInfo)=>{
                return Promise.reject(errorInfo)
            })
        })
    };

    useEffect(()=>{
        currentTeamId && getMemberOptionList()
    },[currentTeamId])

    const getMemberOptionList = ()=>{
        setMemberOptionList([])
        fetchData<BasicResponse<{ teams: TeamSimpleMemberItem[] }>>('team/members/simple',{method:'GET',eoParams:{team:currentTeamId}}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setMemberOptionList(data.teams?.map((x:TeamSimpleMemberItem)=>{return {
                    label:x.user.name, value:x.user.id
                }}))
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const getTeamOptionList = ()=>{
        setTeamOptionList([])
        fetchData<BasicResponse<{ teams: SimpleTeamItem[] }>>('simple/teams/mine',{method:'GET',eoTransformKeys:['available_partitions']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setTeamOptionList(data.teams?.map((x:MemberItem)=>{return {...x,
                    label:x.name, value:x.id
                }}))
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const deleteSystem = ()=>{
        fetchData<BasicResponse<null>>('team/project',{method:'DELETE',eoParams:{team:teamId,project:systemId}}).then(response=>{
            const {code,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                message.success(msg || '操作成功！')
                navigate(`/system/list`)
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    useEffect(() => {
        
        // getMemberOptionList()
        getTeamOptionList()
        // getPartitionList()

        if (systemId !== undefined) {
            setOnEdit(true);
            getSystemInfo();
            
            setBreadcrumb([
                {
                    title: <Link to={`/system/list`}>内部数据服务</Link>
                },
                {
                    title: '设置'
                }])

        } else {
            setOnEdit(false);
            form.setFieldValue('id',uuidv4()); // 清空 initialValues
            form.setFieldValue('team',teamId); // 清空 initialValues
        }
        return (form.setFieldsValue({}))
    }, [systemId]);


    // const getPartitionList = ()=>{
    //     setPartitionOption([])
    //     fetchData<BasicResponse<{partitions:PartitionItem[]}>>('simple/organization/partitions',{method:'GET',eoParams:{organization:orgId}}).then(response=>{
    //         const {code,data,msg} = response
    //         if(code === STATUS_CODE.SUCCESS){
    //             setPartitionOption(data.partitions?.map((x:PartitionItem)=>{return {
    //                 label:x.name,value:x.id
    //             }}))
    //             if(systemId === undefined && data.partitions?.length === 1){
    //                 form.setFieldValue('partition',[data.partitions[0].id])
    //             }
    //         }else{
    //             message.error(msg || '操作失败')
    //         }
    //     })
    // }

    
    const deleteSystemModal = async ()=>{
        modal.confirm({
            title:'删除',
            content:'该数据删除后将无法找回，请确认是否删除？',
            onOk:()=> {
                return deleteSystem()
            },
            width:600,
            okText:'确认',
            okButtonProps:{
                danger:true
            },
            cancelText:'取消',
            closable:true,
            icon:<></>
        })
    }


    return (
        <>
            <div className={`h-full min-w-[570px]`}>
                {/* <WithPermission access={onEdit ? 'team.mySystem.self.edit' :'team.mySystem.self.add'}> */}
                <Form
                    layout='vertical'
                    labelAlign='left'
                    scrollToFirstError
                    form={form}
                    className="mx-auto  flex flex-col justify-between h-full"
                    name="systemConfig"
                    // labelCol={{ offset:1, span: 4 }}
                    // wrapperCol={{ span: 19}}
                    onFinish={onFinish}
                    autoComplete="off"
                >
                    <div>
                        <Form.Item<SystemConfigFieldType>
                            label="服务名称"
                            name="name"
                            rules={[{ required: true, message: '必填项' ,whitespace:true }]}
                        >
                            <Input className="w-INPUT_NORMAL" placeholder="请输入服务名称"/>
                        </Form.Item>

                        <Form.Item<SystemConfigFieldType>
                            label="服务ID"
                            name="id"
                            extra="服务ID（sys_id）可用于检索服务或日志"
                            rules={[{ required: true, message: '必填项' ,whitespace:true }]}
                        >
                            <Input className="w-INPUT_NORMAL" disabled={onEdit} placeholder="请输入服务ID"/>
                        </Form.Item>

                        <Form.Item<SystemConfigFieldType>
                            label="API 调用前缀"
                            name="prefix"
                            extra="选填，作为服务内所有服务的API的前缀，比如host/{sys_name}/{service_name}/{api_path}，一旦保存无法修改"
                            rules={[
                            {
                            validator: validateUrlSlash,
                            }]}
                        >
                            <Input prefix={onEdit ? '' : '/'} className="w-INPUT_NORMAL" disabled={onEdit} placeholder="请输入 API 调用前缀"/>
                        </Form.Item>

                        <Form.Item<SystemConfigFieldType>
                            label="描述"
                            name="description"
                        >
                            <Input.TextArea className="w-INPUT_NORMAL" placeholder="请输入描述"/>
                        </Form.Item>

                        {!onEdit && <Form.Item<SystemConfigFieldType>
                            label="所属团队"
                            name="team"
                            rules={[{ required: true, message: '必填项' }]}
                        >
                            <Select className="w-INPUT_NORMAL" disabled={onEdit} placeholder="请选择" options={teamOptionList} onChange={(x)=>{setCurrentTeamId(x)}}>
                            </Select>
                        </Form.Item>}

                        <Form.Item<SystemConfigFieldType>
                            label="负责人"
                            name="master"
                            extra="负责人对服务内的服务、服务、成员有管理权限"
                            rules={[{required: true, message: '必填项'}]}
                        >
                            <Select className="w-INPUT_NORMAL" placeholder="请选择负责人" options={memberOptionList}>
                            </Select>
                        </Form.Item>

                        <Form.Item<SystemConfigFieldType>
                            label="环境权限"
                            name="partition"
                            rules={[{ required: true, message: '必填项' }]}
                        >
                            <Checkbox.Group className="flex flex-col gap-[8px] mt-[5px]" options={partitionOption} />{partitionOption.length === 0 && <span className="text-status_fail block h-[27px]">暂无可选环境，请检查所属团队可用环境</span>}
                        </Form.Item>
                        {onEdit && <>
                        <Row className="mb-[10px]"
                            // wrapperCol={{ offset: 5, span: 19 }}
                            >
                        {/* <WithPermission access={onEdit ? 'team.mySystem.self.edit' :'team.mySystem.self.add'}> */}
                            <Button type="primary" htmlType="submit">
                                保存
                            </Button>
                            {/* </WithPermission> */}
                        </Row></>}
                    </div>
                    {onEdit && <>
                        <div>
                            <Divider />
                            <p className="text-center">删除服务前，需要先删除所有服务内的服务，删除服务之后将无法找回，请谨慎操作！</p>
                            <div className="text-center">
                                {/* <WithPermission access="project.mySystem.self.delete"> */}
                                    <Button className="m-auto mt-[16px] mb-[20px]" type="default" danger={true} onClick={deleteSystemModal}>删除服务</Button>
                                    {/* </WithPermission> */}
                            </div>
                        </div>
                    </>}
                </Form>
                {/* </WithPermission> */}
                </div>
        </>
    )
})
export default SystemConfig