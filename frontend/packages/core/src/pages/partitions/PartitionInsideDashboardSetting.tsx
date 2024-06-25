import { App, Form, Input, Select, Button, Divider, Row } from "antd";
import { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import { RouterParams } from "@core/components/aoplatform/RenderRoutes";
import WithPermission from "@common/components/aoplatform/WithPermission";
import { BasicResponse, STATUS_CODE } from "@common/const/const";
import { PartitionDashboardConfigFieldType } from "../../const/partitions/types";
import { useBreadcrumb } from "@common/contexts/BreadcrumbContext";
import { useFetch } from "@common/hooks/http";
import { DASHBOARD_SETTING_DRIVER_OPTION_LIST } from "../../const/partitions/const";

export default function PartitionInsideDashboardSetting(){

    const { message } = App.useApp()
    const { partitionId } = useParams<RouterParams>();
    const [ form ] = Form.useForm();
    const { setBreadcrumb } = useBreadcrumb()
    const { fetchData} = useFetch()
    const [, forceUpdate] = useState({});
    const onFinish = () => {
        form.validateFields().then((value)=>{
            fetchData<BasicResponse<{info: PartitionDashboardConfigFieldType}>>('partition/monitor',{method: 'POST',eoBody:(value),eoParams:{partition:partitionId}}).then(response=>{
                const {code,data,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    const config = data.info.config
                    form.setFieldsValue({...data.info,config:{addr:config.addr || '', org:config.org || '', token:config.token || ''}})
                    message.success(msg || '操作成功！')
                }else{
                    message.error(msg || '操作失败')
                }
            })
        })
    }

    // 获取表单默认值
    const getDashboardSetting = () => {
        fetchData<BasicResponse<{ info: PartitionDashboardConfigFieldType }>>('partition/monitor',{method:'GET',eoParams:{partition:partitionId}},).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                const config = data.info.config
                form.setFieldsValue({...data.info,config:{addr:config.addr || '', org:config.org || '', token:config.token || ''}})
                forceUpdate({})
            }else{
                message.error(msg || '操作失败')
            }
        })
    };

    const resetDashboardConfig = ()=>{
            fetchData<BasicResponse<null>>('partition/monitor',{method:'DELETE',eoParams:{partition:partitionId}}).then(response=>{
                const {code,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    message.success(msg || '操作成功')
                    form.resetFields()
                    getDashboardSetting()
                }else{
                    message.error(msg || '操作失败')
                }
            })
    }

    useEffect(() => {
        getDashboardSetting();

        setBreadcrumb([
            {title:<Link to="/partition/list">部署管理</Link>},
            {title:'监控配置'}
        ])
        return (form.setFieldsValue({}))
    }, [partitionId]);

    return (
        <>
                <div className="h-full">
                <WithPermission access={'system.partition.self.edit'} >
                    <Form
                        layout='vertical'
                        labelAlign='left'
                        scrollToFirstError
                        form={form}
                        name="paritionInsideDashboardSetting"
                        className="mx-auto pb-[20px]  flex flex-col justify-between h-full"
                        // labelCol={{ offset:1, span: 4}}
                        // wrapperCol={{ span: 19}}
                        onFinish={onFinish}
                        autoComplete="off"
                    >
                        {/* <Form.Item<PartitionDashboardConfigFieldType>
                            label="分区名称"
                            name="name"
                            rules={[{ required: true, message: '必填项' }]}
                        >
                            <Input className="w-INPUT_NORMAL"  placeholder="请输入分区名称"/>
                        </Form.Item> */}
                        <div>
                            <Form.Item<PartitionDashboardConfigFieldType>
                                label="数据源类型"
                                name="driver"
                                rules={[{ required: true, message: '必填项' }]}
                            >
                                <Select className="w-INPUT_NORMAL" placeholder="请选择数据源类型" options={[...DASHBOARD_SETTING_DRIVER_OPTION_LIST]}/>
                            </Form.Item>

                            <Form.Item<PartitionDashboardConfigFieldType>
                                label="数据源地址"
                                name={['config','addr']}
                                rules={[{ required: true, message: '必填项',whitespace:true  }]}
                            >
                                <Input className="w-INPUT_NORMAL"  placeholder="请输入数据源地址"/>
                            </Form.Item>

                            <Form.Item<PartitionDashboardConfigFieldType>
                                label="Organization"
                                name={['config','org']}
                                rules={[{ required: true, message: '必填项',whitespace:true  }]}
                            >
                                <Input className="w-INPUT_NORMAL"  placeholder="请输入 Organization"/>
                            </Form.Item>

                            <Form.Item<PartitionDashboardConfigFieldType>
                                label="鉴权 Token"
                                name={['config','token']}
                            >
                                <Input className="w-INPUT_NORMAL"  placeholder="请输入鉴权 Token"/>
                            </Form.Item>

                            <Row className="mb-[10px]" 
                            // wrapperCol={{ offset: 6, span: 16 }}
                            >
                                <WithPermission access='system.partition.self.edit'>
                                    <Button type="primary" htmlType="submit">
                                        保存
                                    </Button>
                                </WithPermission>
                            </Row>
                            </div>
                            <div>
                                <Divider />
                                <p className="text-center">重置监控：重置操作不可恢复，请谨慎操作！！</p>
                                <div className="text-center">
                                    <WithPermission access="system.partition.self.delete"><Button className="m-auto mt-[16px] mb-[20px]" type="default" onClick={resetDashboardConfig}>重置</Button></WithPermission>
                                </div>
                            </div>
                    </Form>
                </WithPermission>
                </div>
        </>
    )
}