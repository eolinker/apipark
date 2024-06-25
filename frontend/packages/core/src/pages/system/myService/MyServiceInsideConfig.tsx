import  { forwardRef, useEffect, useImperativeHandle, useState} from "react";
import { App,  Button, Checkbox, Divider,  Form,  Input, Row, Select, Spin, TagType, TreeSelect, Upload,  UploadFile, UploadProps} from "antd";
import {LoadingOutlined, PlusOutlined} from "@ant-design/icons";
import {RcFile, UploadChangeParam} from "antd/es/upload";
import Radio from "antd/es/radio";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {DefaultOptionType} from "antd/es/cascader";
import {useSystemContext} from "../../../contexts/SystemContext.tsx";
import {v4 as uuidv4} from 'uuid'
import { visualizations } from "../../../const/system/const.tsx";
import { SimpleSystemItem, MyServiceFieldType, MyServiceInsideConfigHandle, MyServiceInsideConfigProps } from "../../../const/system/type.ts";
import { EntityItem, SimpleTeamItem } from "@common/const/type.ts";
import WithPermission from "@common/components/aoplatform/WithPermission.tsx";
import { getImgBase64 } from "@common/utils/dataTransfer.ts";
import { CategorizesType } from "@market/const/serviceHub/type.ts";

const MAX_SIZE = 2 * 1024; // 1KB

const MyServiceInsideConfig = forwardRef<MyServiceInsideConfigHandle,MyServiceInsideConfigProps>((props,ref)=>{
    const {teamId, systemId, serviceId,closeDrawer} = props
    const { message,modal } = App.useApp()
    const [onEdit, setOnEdit] = useState<boolean>(!!serviceId)
    const [form] = Form.useForm();
    const {fetchData} = useFetch()
    const [, setImageUrl] = useState<string>();
    const [teamOptionList, setTeamOptionList] = useState<DefaultOptionType[]>([])
    const [systemOptionList, setSystemOptionList] = useState<DefaultOptionType[]>([])
    const [tagOptionList, setTagOptionList] = useState<DefaultOptionType[]>([])
    const {partitionList} = useSystemContext()
    const [serviceClassifyOptionList, setServiceClassifyOptionList] = useState<DefaultOptionType[]>()
    const [showClassify, setShowClassify] = useState<boolean>()
    const [imageBase64, setImageBase64] = useState<string | null>(null);
    const [status,setStatus] = useState<'off'|'on'>('on')
    const [loading, setLoading] = useState<boolean>(false)
    const [uploadLoading, setUploadLoading] = useState<boolean>(false)
    const [startBtnLoading, setStartBtnLoading] = useState<boolean>(false)
    const [delBtnLoading, setDelBtnLoading] = useState<boolean>(false)

    const beforeUpload = async (file: RcFile) => {
    if (!['image/png', 'image/jpeg', 'image/svg+xml'].includes(file.type)) {
        alert('只允许上传PNG、JPG或SVG格式的图片');
        return false;
      }
  
      if (file.size > MAX_SIZE) {
        try {
          const compressedBase64 = await compressImage(file, MAX_SIZE);
          setImageBase64(`data:${file.type};base64,${compressedBase64}`);
          form.setFieldValue('logo', `data:${file.type};base64,${compressedBase64}`);
        } catch (error) {
          console.error('压缩图片时出错', error);
        }
      } else {
        const reader = new FileReader();
        reader.onload = (e: ProgressEvent<FileReader>) => {
        
          setImageBase64(e.target?.result as string);
          form.setFieldValue('logo', e.target?.result);
        };
        reader.readAsDataURL(file);
      }
        return false;
    };


    const compressImage = (file: RcFile, maxSize: number): Promise<string> => {
        const img = document.createElement('img');
        const canvas = document.createElement('canvas');
        const reader = new FileReader();
    
        return new Promise((resolve, reject) => {
          reader.onload = (e) => {
            img.src = e.target.result as string;
            img.onload = () => {
              let quality = 0.9;
              let width = img.width;
              let height = img.height;
              
              const ctx = canvas.getContext('2d');
              
              const compress = () => {
                canvas.width = width;
                canvas.height = height;
                ctx.clearRect(0, 0, width, height);
                ctx.drawImage(img, 0, 0, width, height);
    
                const dataUrl = canvas.toDataURL(file.type, quality);
                const base64 = dataUrl.split(',')[1];
                return { base64, size: base64.length * 0.75 };
              };
    
              let { base64, size } = compress();
    
              while (size > maxSize && quality > 0.1) {
                quality -= 0.1;
                ({ base64, size } = compress());
              }
    
              while (size > maxSize && (width > 50 || height > 50)) {
                width *= 0.9;
                height *= 0.9;
                ({ base64, size } = compress());
              }
    
              resolve(base64);
            };
          };
          reader.onerror = (e) => reject(e);
          reader.readAsDataURL(file);
        });
      };
    // 获取表单默认值
    const getServiceInfo = () => {
        setLoading(true)
        fetchData<BasicResponse<{ service: MyServiceFieldType }>>('project/service/info',{method:'GET',eoParams:{service:serviceId,project:systemId},eoTransformKeys:['service_type']}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setTimeout(()=>{form.setFieldsValue({...data.service, logoFile:[
                        {
                            uid: '-1', // 文件唯一标识
                            name: 'image.png', // 文件名
                            status: 'done', // 状态有：uploading, done, error, removed
                            url: data.service.logo, // 图片 Base64 数据
                        }
                    ],team:data.service.team.id, project:data.service.project.id,tags:data.service.tags?.map((x:EntityItem)=>x.name)||[],partition:data.service.partition?.map((x:EntityItem)=>x.id)||[],group:data.service.group.id})},0)
                setImageBase64(data.service.logo)
                setStatus( data.service.status)
                setShowClassify(data.service.serviceType === 'public')
            }else{
                message.error(msg || '操作失败')
            }
        }).finally(()=>setLoading(false))

    }
    const onFinish = () => {
        return form.validateFields().then((value)=>{
            return fetchData<BasicResponse<{service:MyServiceFieldType}>>(serviceId === undefined?'project/service':'project/service/info',{method:serviceId === undefined? 'POST' : 'PUT',eoBody:({...value,tags:value.tags?.map((x:string|EntityItem)=>typeof x === 'string' ? x : x.name),partition:value.partition?.map((x:string)=>x.toString())}), eoParams:serviceId === undefined ? {project:systemId}:{service:serviceId,project:systemId},eoTransformKeys:['serviceType']}).then(response=>{
                const {code,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    // setServiceInfo(data.service)
                    message.success(msg || '操作成功！')
                    return Promise.resolve(true)
                    // if(serviceId === undefined) navigate(`/system/${orgId}/${teamId}/inside/${systemId}/myService/inside/${data.service.id}/api`)
                }else{
                    message.error(msg || '操作失败')
                    return Promise.reject(msg || '操作失败')
                }
            }).catch((errInfo)=>Promise.reject(errInfo))
        }).catch((err) => {
            form.scrollToField(err.errorFields[0].name[0]);
            return Promise.reject(msg || '操作失败')
          });
    };
    

    useImperativeHandle(ref,()=>({
        save:onFinish
    })
)

    const handleChange: UploadProps['onChange'] = (info: UploadChangeParam<UploadFile>) => {
        if (info.file.status === 'uploading') {
            setUploadLoading(true);
            return;
        }
        if (info.file.status === 'done') {
            // Get this url from response in real world.
            getImgBase64(info.file.originFileObj as RcFile, (url) => {
                setUploadLoading(false);
                setImageUrl(url);
            });
        }
        if (info.fileList.length === 0) {
            // 如果文件被移除，清除 logo 字段
            form.setFieldValue( "logo", null );
        }
    };

    const uploadButton = ( 
    <div>
        {uploadLoading ? <LoadingOutlined /> : <PlusOutlined />}
    </div>
    );

    const getTeamList = ()=>{
        setTeamOptionList([])
        fetchData<BasicResponse<{ teams: SimpleTeamItem[] }>>('simple/teams/mine',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setTeamOptionList(data.teams?.map((x:SimpleTeamItem)=>{return {
                    label:x.name, value:x.id
                }})||[])
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const getSystemList = ()=>{
        setSystemOptionList([])
        fetchData<BasicResponse<{ projects: SimpleSystemItem[] }>>('simple/projects/mine',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setSystemOptionList(data.projects?.map((x:SimpleSystemItem)=>{return {
                    label:x.name, value:x.id, partition:x.partition?.map((x:EntityItem)=>x.id)
                }}) || [])
            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const getTagAndServiceClassifyList = ()=>{
        setTagOptionList([])
        setServiceClassifyOptionList([])
        fetchData<BasicResponse<{ catalogues:CategorizesType[],tags:TagType[]}>>('catalogues',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                setTagOptionList(data.tags?.map((x:TagType)=>{return {
                    label:x.name, value:x.name
                }})||[])
                setServiceClassifyOptionList(data.catalogues)

            }else{
                message.error(msg || '操作失败')
            }
        })
    }

    const deleteService = ()=>{
        fetchData<BasicResponse<null>>('project/service',{method:'DELETE',eoParams:{service:serviceId,project:systemId}}).then(response=>{
            const {code,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                message.success(msg || '操作成功，即将返回列表页')
                closeDrawer?.()
            }else{
                message.error(msg || '操作失败')
            }
        }).finally(()=>setDelBtnLoading(false))
    }

    const enabledService =()=>{
        setStartBtnLoading(true)
        fetchData<BasicResponse<null>>(`project/service/${status==='off'?'enable':'disable'}`,{method:'PUT',eoParams:{service:serviceId,project:systemId}}).then(response=>{
            const {code,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                message.success(msg || '操作成功')
                setStatus((prevState)=>prevState === 'off' ? 'on' :'off')
            }else{
                message.error(msg || '操作失败')
            }
        }).finally(()=>setStartBtnLoading(false))
    }

    
const normFile = (e: unknown) => {
    if (Array.isArray(e)) {
      return e;
    }
    return( e as {fileList:unknown} )?.fileList;
  };

  
  const deleteSystemModal = async ()=>{
    setDelBtnLoading(true)
    modal.confirm({
        title:'删除',
        content:'该数据删除后将无法找回，请确认是否删除？',
        onOk:()=> {
            return deleteService()
        },
        width:600,
        okText:'确认',
        okButtonProps:{
            danger:true
        },
        onCancel:()=>{
            setDelBtnLoading(false)
        },
        cancelText:'取消',
        closable:true,
        icon:<></>
    })
}


    useEffect(() => {
        getTeamList()
        getSystemList()
        getTagAndServiceClassifyList()
        if (serviceId !== undefined) {
            setOnEdit(true);
            getServiceInfo();
        } else {
            setOnEdit(false);
            form.setFieldsValue({id:uuidv4(), serviceType:'inner',partition:partitionList?.length === 1 ? [partitionList[0].id] : [], team:teamId,project:systemId}); // 清空 initialValues
        }

        return (form.setFieldsValue({}))
    }, []);

    return (
        <Spin indicator={<LoadingOutlined style={{ fontSize: 24 }} spin />} spinning={loading}>
            <div className="h-full overflow-y-auto pb-btnybase">
                <WithPermission access={onEdit ? 'project.mySystem.service.edit':'project.mySystem.service.add'}>
                    <Form
                        labelAlign='left'
                        layout='vertical'
                        scrollToFirstError
                        form={form}
                        className="mx-auto  "
                        name="myServiceInsideConfigForm"
                        // labelCol={{ offset:1, span:4 }}
                        // wrapperCol={{ span: 19}}
                        onFinish={onFinish}
                        autoComplete="off"
                    >
                        <Form.Item<MyServiceFieldType>
                            label="服务名称"
                            name="name"
                            rules={[{ required: true, message: '必填项',whitespace:true }]}
                        >
                            <Input className="w-INPUT_NORMAL" placeholder="请输入服务名称"/>
                        </Form.Item>

                        <Form.Item<MyServiceFieldType>
                            label="服务ID"
                            name="id"
                            extra="服务ID（service_id）可用于检索服务或日志"
                            rules={[{ required: true, message: '必填项',whitespace:true  }]}
                        >
                            <Input className="w-INPUT_NORMAL" disabled={onEdit} placeholder="请输入服务ID"/>
                        </Form.Item>

                        <Form.Item<MyServiceFieldType>
                            label="图标"
                            name="logoFile"
                            extra="仅支持 .png .jpg .jpeg .svg 格式的图片文件, 大于 1KB 的文件将被压缩"
                            valuePropName="fileList" getValueFromEvent={normFile}
                        >
                            <Upload
                                listType="picture"
                                beforeUpload={beforeUpload}
                                onChange={handleChange}
                                showUploadList={false}
                                maxCount={1}
                                accept=".png, .jpg, .jpeg, .svg"
                            >
                                <div className="h-[68px] w-[68px] border-[1px] border-dashed border-BORDER flex items-center justify-center rounded bg-bar-theme cursor-pointer" style={{ marginTop: 8 }}>
                                    {imageBase64 ? <img src={imageBase64} alt="Logo" style={{  maxWidth: '200px', width:'68px',height:'68px'}} /> : uploadButton}
                                </div>
                            </Upload>

                        </Form.Item>


                        <Form.Item<MyServiceFieldType>
                            label="Logo"
                            name="logo"
                            hidden
                        >
                        </Form.Item>

                        <Form.Item<MyServiceFieldType>
                            label="所属团队"
                            name="team"
                            rules={[{ required: true, message: '必填项' }]}
                        >
                            <Select className="w-INPUT_NORMAL" disabled={true} placeholder="请选择" options={teamOptionList}>
                            </Select>
                        </Form.Item>

                        <Form.Item<MyServiceFieldType>
                            label="所属服务"
                            name="project"
                            rules={[{ required: true, message: '必填项' }]}
                        >
                            <Select className="w-INPUT_NORMAL" disabled={true} placeholder="请选择" options={systemOptionList}>
                            </Select>
                        </Form.Item>

                        <Form.Item<MyServiceFieldType>
                            label="标签"
                            name="tags"
                        >
                            <Select 
                                className="w-INPUT_NORMAL" 
                                mode="tags" 
                                placeholder="请选择" 
                                options={tagOptionList}>
                            </Select>
                            
                        </Form.Item>

                        <Form.Item<MyServiceFieldType>
                            label="描述"
                            name="description"
                        >
                            <Input.TextArea className="w-INPUT_NORMAL" placeholder="请输入描述"/>
                        </Form.Item>

                        <Form.Item<MyServiceFieldType>
                            label="可访问分区"
                            name="partition"
                            rules={[{required: true, message: '必填项'}]}
                        >
                            <Checkbox.Group disabled={onEdit}  className="flex flex-col gap-[8px] mt-[5px]" options={partitionList?.map((x)=>({label:x.name, value:x.id})) || []} />
                        </Form.Item>

                        <Form.Item<MyServiceFieldType>
                            label="可见性"
                            name="serviceType"
                            rules={[{required: true, message: '必填项'}]}
                        >
                            <Radio.Group className="flex flex-col" options={visualizations} onChange={(e)=>{setShowClassify(e.target.value === 'public')}} />
                        </Form.Item>

                        {showClassify &&
                        <Form.Item<MyServiceFieldType>
                            label="所属服务分类"
                            name="group"
                            extra="设置服务展示在服务市场中的哪个分类下"
                            rules={[{required: true, message: '必填项'}]}
                        >
                            <TreeSelect
                                className="w-INPUT_NORMAL"
                                fieldNames={{label:'name',value:'id',children:'children'}}
                                showSearch
                                dropdownStyle={{ maxHeight: 400, overflow: 'auto' }}
                                placeholder="请选择"
                                allowClear
                                treeDefaultExpandAll
                                treeData={serviceClassifyOptionList}
                            />
                        </Form.Item>
                        }
                        {serviceId !== undefined && 
                        <Row className="mb-[10px]"
                        // wrapperCol={{ offset: 5, span: 16 }}
                        >
                        <WithPermission access={onEdit ? 'project.mySystem.service.edit':'project.mySystem.service.add'}><Button className="mr-btnbase" type="primary" htmlType="submit">
                                完成
                            </Button></WithPermission>
                        </Row>}
                        {serviceId !== undefined && <>
                            <Divider className=""/>
                                <p className="text-center">停用服务后，服务将从服务市场中下线，所有请求将无法通过该服务访问 API，请谨慎操作！</p>
                                <div className="text-center"><WithPermission access="project.mySystem.service.running"><Button className="mt-[16px] m-auto"  onClick={enabledService} loading={startBtnLoading}>{status === 'off'? '启用':'停用'}</Button></WithPermission></div>

                            <Divider className="" />
                                <p className="text-center">删除操作不可恢复，请谨慎操作！</p>
                                <div className="text-center"><WithPermission access="project.mySystem.service.delete"><Button className="mt-[16px]" onClick={deleteSystemModal} loading={delBtnLoading}>删除</Button></WithPermission></div>
                        </>}
                    </Form>
                    </WithPermission>
            </div>
        </Spin>
    )
})
export default MyServiceInsideConfig