import {forwardRef, useEffect, useImperativeHandle, useState} from "react";
import { action } from '@formily/reactive'
import {
    FormItem,
    Space,
    ArrayItems,
    DatePicker,
    Editable,
    FormButtonGroup,
    Input,
    Radio,
    Select,
    Submit,
    Cascader,
    Form,
    FormGrid,
    FormLayout,
    Upload,
    ArrayCollapse,
    ArrayTable,
    ArrayTabs,
    Checkbox,
    FormCollapse,
    FormDialog,
    FormDrawer,
    FormStep,
    FormTab,
    NumberPicker,
    Password,
    PreviewText,
    Reset,
    SelectTable,
    Switch,
    TimePicker,
    Transfer,
    TreeSelect,
    ArrayCards
} from '@formily/antd-v5'
import { createForm } from '@formily/core'
import {CustomCodeboxComponent} from "@common/components/aoplatform/formily2-customize/CustomCodeboxComponent.tsx";
import {SimpleMapComponent} from "@common/components/aoplatform/formily2-customize/SimpleMapComponent.tsx";
import {CustomDialogComponent} from "@common/components/aoplatform/formily2-customize/CustomDialogComponent.tsx";
import {ArrayItemBlankComponent} from "@common/components/aoplatform/formily2-customize/ArrayItemBlankComponent.tsx";
import {DefaultOptionType} from "antd/es/cascader";
import {createSchemaField, FormProvider} from "@formily/react";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import {App} from "antd";
import { cloneDeep } from "lodash-es";

export type IntelligentPluginConfigProps = {
    type:'add'|'edit'
    renderSchema:unknown
    tabData:DefaultOptionType[]
    moduleId:string
    driverSelectionOptions:DefaultOptionType[]
    entityId?:string
    initFormValue:{[k:string]:unknown}
}

export type IntelligentPluginConfigHandle = {
    save:()=>Promise<boolean | string>
}

const finalSchema = {
    "type": "object",
    "properties": {
        "layout": {
            "type": "void",
            "x-component": "FormLayout",
            "x-component-props": {
                "labelCol": 6,
                "wrapperCol": 10,
                "layout": "vertical"
            },
            "properties": {
                "id": {
                    "type": "string",
                    "title": "ID",
                    "required": true,
                    "pattern": {},
                    "x-decorator": "FormItem",
                    "x-decorator-props": {
                        "labelCol": 4,
                        "wrapperCol": 20,
                        "labelAlign": "left"
                    },
                    "x-component": "Input",
                    "x-component-props": {
                        "placeholder": "支持字母开头、英文数字中横线下划线组合"
                    },
                    "x-disabled": false
                },
                "title": {
                    "type": "string",
                    "title": "名称",
                    "required": true,
                    "x-decorator": "FormItem",
                    "x-decorator-props": {
                        "labelCol": 4,
                        "wrapperCol": 20,
                        "labelAlign": "left"
                    },
                    "x-component": "Input",
                    "x-component-props": {
                        "placeholder": "请输入名称"
                    }
                },
                "driver": {
                    "type": "string",
                    "title": "Driver",
                    "required": true,
                    "x-decorator": "FormItem",
                    "x-decorator-props": {
                        "labelCol": 4,
                        "wrapperCol": 20,
                        "labelAlign": "left"
                    },
                    "x-component": "Select",
                    "x-component-props": {
                        "disabled": false
                    },
                    "x-display": "hidden",
                    "enum": [
                        {
                            "label": "文件",
                            "value": "file"
                        }
                    ]
                },
                "description": {
                    "type": "string",
                    "title": "描述",
                    "x-decorator": "FormItem",
                    "x-decorator-props": {
                        "labelCol": 4,
                        "wrapperCol": 20,
                        "labelAlign": "left"
                    },
                    "x-component": "Input.TextArea",
                    "x-component-props": {
                        "placeholder": "请输入描述"
                    }
                },
                "config": {
                    "type": "object",
                    "x-component": "void",
                    "properties": {
                        "tabForm": {
                            "type": "void",
                            "x-component": "FormTab",
                            "x-component-props": {
                                "formTab": "{{formTab}}",
                                "centered": false
                            },
                            "properties": {
                                "0a3d7fcc-a59f-435f-8be0-9ff2cd0809c6": {
                                    "type": "object",
                                    "x-component": "FormTab.TabPane",
                                    "x-component-props": {
                                        "tab": "华南",
                                        "forceRender": true
                                    },
                                    "properties": {
                                        "_apinto_show": {
                                            "type": "boolean",
                                            "title": "启用",
                                            "x-decorator": "FormItem",
                                            "x-decorator-props": {
                                                "labelCol": 4,
                                                "wrapperCol": 20,
                                                "labelAlign": "left"
                                            },
                                            "x-component": "Switch",
                                            "x-index": 0
                                        },
                                        "_apinto_backend": {
                                            "type": "void",
                                            "x-component": "void",
                                            "x-reactions": {
                                                "dependencies": [
                                                    "._apinto_show"
                                                ],
                                                "fulfill": {
                                                    "state": {
                                                        "visible": "{{!!$deps[0]}}"
                                                    }
                                                }
                                            },
                                            "properties": {
                                                "dir": {
                                                    "name": "dir",
                                                    "required": true,
                                                    "title": "存放目录",
                                                    "type": "string",
                                                    "x-component": "Input",
                                                    "x-component-props": {},
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    },
                                                    "x-index": 2,
                                                    "x-validator": []
                                                },
                                                "expore": {
                                                    "default": "3",
                                                    "description": "单位：天",
                                                    "name": "expore",
                                                    "required": true,
                                                    "title": "过期时间",
                                                    "type": "number",
                                                    "x-component": "NumberPicker",
                                                    "x-component-props": {},
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    },
                                                    "x-index": 4,
                                                    "x-validator": "integer"
                                                },
                                                "file": {
                                                    "name": "file",
                                                    "required": true,
                                                    "title": "文件名称",
                                                    "type": "string",
                                                    "x-component": "Input",
                                                    "x-component-props": {},
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    },
                                                    "x-index": 1,
                                                    "x-validator": []
                                                },
                                                "formatter": {
                                                    "title": "格式化配置",
                                                    "type": "object",
                                                    "x-component": "CustomCodeboxComponent",
                                                    "x-component-props": {
                                                        "mode": "json"
                                                    },
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    }
                                                },
                                                "period": {
                                                    "default": "hour",
                                                    "enum": [
                                                        {
                                                            "children": [],
                                                            "label": "小时",
                                                            "value": "hour"
                                                        },
                                                        {
                                                            "children": [],
                                                            "label": "天",
                                                            "value": "day"
                                                        }
                                                    ],
                                                    "name": "period",
                                                    "required": true,
                                                    "title": "日志分割周期",
                                                    "x-component": "Select",
                                                    "x-component-props": {},
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    },
                                                    "x-index": 3,
                                                    "x-validator": []
                                                },
                                                "scopes": {
                                                    "items": {
                                                        "properties": {
                                                            "remove": {
                                                                "type": "void",
                                                                "x-component": "ArrayItems.Remove",
                                                                "x-decorator": "FormItem"
                                                            },
                                                            "select": {
                                                                "enum": [
                                                                    {
                                                                        "label": "Access日志",
                                                                        "value": "access_log"
                                                                    }
                                                                ],
                                                                "type": "string",
                                                                "x-component": "Select",
                                                                "x-decorator": "FormItem"
                                                            },
                                                            "sort": {
                                                                "type": "void",
                                                                "x-component": "ArrayItems.SortHandle",
                                                                "x-decorator": "FormItem"
                                                            }
                                                        },
                                                        "type": "void",
                                                        "x-component": "Space"
                                                    },
                                                    "name": "scopes",
                                                    "properties": {
                                                        "add": {
                                                            "title": "添加条目",
                                                            "type": "void",
                                                            "x-component": "ArrayItems.Addition",
                                                            "x-component-props": {
                                                                "defaultValue": "access_log"
                                                            }
                                                        }
                                                    },
                                                    "required": true,
                                                    "title": "作用范围",
                                                    "type": "array",
                                                    "x-component": "ArrayItems",
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    },
                                                    "x-index": 0
                                                },
                                                "type": {
                                                    "default": "line",
                                                    "enum": [
                                                        {
                                                            "children": [],
                                                            "label": "单行",
                                                            "value": "line"
                                                        },
                                                        {
                                                            "children": [],
                                                            "label": "Json",
                                                            "value": "json"
                                                        }
                                                    ],
                                                    "name": "type",
                                                    "required": true,
                                                    "title": "输出格式",
                                                    "x-component": "Select",
                                                    "x-component-props": {},
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    },
                                                    "x-index": 5,
                                                    "x-validator": []
                                                }
                                            }
                                        }
                                    }
                                },
                                "8bfbba9e-5931-40ac-8ac4-084d1ff36409": {
                                    "type": "object",
                                    "x-component": "FormTab.TabPane",
                                    "x-component-props": {
                                        "tab": "华北",
                                        "forceRender": true
                                    },
                                    "properties": {
                                        "_apinto_show": {
                                            "type": "boolean",
                                            "title": "启用",
                                            "x-decorator": "FormItem",
                                            "x-decorator-props": {
                                                "labelCol": 4,
                                                "wrapperCol": 20,
                                                "labelAlign": "left"
                                            },
                                            "x-component": "Switch",
                                            "x-index": 0
                                        },
                                        "_apinto_backend": {
                                            "type": "void",
                                            "x-component": "void",
                                            "x-reactions": {
                                                "dependencies": [
                                                    "._apinto_show"
                                                ],
                                                "fulfill": {
                                                    "state": {
                                                        "visible": "{{!!$deps[0]}}"
                                                    }
                                                }
                                            },
                                            "properties": {
                                                "dir": {
                                                    "name": "dir",
                                                    "required": true,
                                                    "title": "存放目录",
                                                    "type": "string",
                                                    "x-component": "Input",
                                                    "x-component-props": {},
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    },
                                                    "x-index": 2,
                                                    "x-validator": []
                                                },
                                                "expore": {
                                                    "default": "3",
                                                    "description": "单位：天",
                                                    "name": "expore",
                                                    "required": true,
                                                    "title": "过期时间",
                                                    "type": "number",
                                                    "x-component": "NumberPicker",
                                                    "x-component-props": {},
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    },
                                                    "x-index": 4,
                                                    "x-validator": "integer"
                                                },
                                                "file": {
                                                    "name": "file",
                                                    "required": true,
                                                    "title": "文件名称",
                                                    "type": "string",
                                                    "x-component": "Input",
                                                    "x-component-props": {},
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    },
                                                    "x-index": 1,
                                                    "x-validator": []
                                                },
                                                "formatter": {
                                                    "title": "格式化配置",
                                                    "type": "object",
                                                    "x-component": "CustomCodeboxComponent",
                                                    "x-component-props": {
                                                        "mode": "json"
                                                    },
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    }
                                                },
                                                "period": {
                                                    "default": "hour",
                                                    "enum": [
                                                        {
                                                            "children": [],
                                                            "label": "小时",
                                                            "value": "hour"
                                                        },
                                                        {
                                                            "children": [],
                                                            "label": "天",
                                                            "value": "day"
                                                        }
                                                    ],
                                                    "name": "period",
                                                    "required": true,
                                                    "title": "日志分割周期",
                                                    "x-component": "Select",
                                                    "x-component-props": {},
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    },
                                                    "x-index": 3,
                                                    "x-validator": []
                                                },
                                                "scopes": {
                                                    "items": {
                                                        "properties": {
                                                            "remove": {
                                                                "type": "void",
                                                                "x-component": "ArrayItems.Remove",
                                                                "x-decorator": "FormItem"
                                                            },
                                                            "select": {
                                                                "enum": [
                                                                    {
                                                                        "label": "Access日志",
                                                                        "value": "access_log"
                                                                    }
                                                                ],
                                                                "type": "string",
                                                                "x-component": "Select",
                                                                "x-decorator": "FormItem"
                                                            },
                                                            "sort": {
                                                                "type": "void",
                                                                "x-component": "ArrayItems.SortHandle",
                                                                "x-decorator": "FormItem"
                                                            }
                                                        },
                                                        "type": "void",
                                                        "x-component": "Space"
                                                    },
                                                    "name": "scopes",
                                                    "properties": {
                                                        "add": {
                                                            "title": "添加条目",
                                                            "type": "void",
                                                            "x-component": "ArrayItems.Addition",
                                                            "x-component-props": {
                                                                "defaultValue": "access_log"
                                                            }
                                                        }
                                                    },
                                                    "required": true,
                                                    "title": "作用范围",
                                                    "type": "array",
                                                    "x-component": "ArrayItems",
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    },
                                                    "x-index": 0
                                                },
                                                "type": {
                                                    "default": "line",
                                                    "enum": [
                                                        {
                                                            "children": [],
                                                            "label": "单行",
                                                            "value": "line"
                                                        },
                                                        {
                                                            "children": [],
                                                            "label": "Json",
                                                            "value": "json"
                                                        }
                                                    ],
                                                    "name": "type",
                                                    "required": true,
                                                    "title": "输出格式",
                                                    "x-component": "Select",
                                                    "x-component-props": {},
                                                    "x-decorator": "FormItem",
                                                    "x-decorator-props": {
                                                        "labelCol": 6,
                                                        "wrapperCol": 10
                                                    },
                                                    "x-index": 5,
                                                    "x-validator": []
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

const SchemaField = createSchemaField({
    components: {
        ArrayCards,
        ArrayCollapse,
        ArrayItems,
        ArrayTable,
        ArrayTabs,
        Cascader,
        Checkbox,
        DatePicker,
        Editable,
        Form,
        FormButtonGroup,
        FormCollapse,
        // @ts-ignore
        FormDialog,
        // @ts-ignore
        FormDrawer,
        FormGrid,
        FormItem,
        FormLayout,
        FormStep,
        FormTab,
        Input,
        NumberPicker,
        Password,
        PreviewText,
        Radio,
        Reset,
        Select,
        SelectTable,
        Space,
        Submit,
        Switch,
        TimePicker,
        Transfer,
        TreeSelect,
        Upload,
        CustomCodeboxComponent,
        SimpleMapComponent,
        CustomDialogComponent,
        ArrayItemBlankComponent
    }
})

export const IntelligentPluginConfig =  forwardRef<IntelligentPluginConfigHandle,IntelligentPluginConfigProps>((props,ref)=>{
    const { type,renderSchema,tabData,moduleId,driverSelectionOptions,initFormValue}  = props
    const { message } = App.useApp()
    const [schema, setSchema] = useState({})
    const {fetchData} = useFetch()
    const form = createForm({ validateFirst: type === 'edit' })
    form.setInitialValues(initFormValue || {})

    const pluginEditSchema = {
        type: 'object',
        properties: {
            layout: {
              type: 'void',
              'x-component': 'FormLayout',
              'x-component-props': {
                labelCol: 6,
                wrapperCol: 10,
                layout: 'vertical',
              },
              properties: {
                id: {
                    type: 'string',
                    title: 'ID',
                    required: true,
                    pattern: /^[a-zA-Z][a-zA-Z0-9-_]*$/,
                    'x-decorator': 'FormItem',
                    'x-decorator-props': {
                        labelCol:4,
                        wrapperCol: 20,
                        labelAlign:'left'
                    },
                    'x-component': 'Input',
                    'x-component-props': {
                        placeholder: '支持字母开头、英文数字中横线下划线组合',
                    },
                    'x-disabled': type === 'edit'
                },
                title: {
                    type: 'string',
                    title: '名称',
                    required: true,
                    'x-decorator': 'FormItem',
                    'x-decorator-props': {
                        labelCol:4,
                        wrapperCol: 20,
                        labelAlign:'left'
                    },
                    'x-component': 'Input',
                    'x-component-props': {
                        placeholder: '请输入名称',
                    }
                },
                driver: {
                    type: 'string',
                    title: 'Driver',
                    required: true,
                    'x-decorator': 'FormItem',
                    'x-decorator-props': {
                        labelCol:4,
                        wrapperCol: 20,
                        labelAlign:'left'
                    },
                    'x-component': 'Select',
                    'x-component-props': {
                        disabled: type === 'edit'
                    },
                    'x-display': driverSelectionOptions.length > 1 ? 'visible' : 'hidden',
                    enum: [...driverSelectionOptions]
                },
                description: {
                    type: 'string',
                    title: '描述',
                    'x-decorator': 'FormItem',
                    'x-decorator-props': {
                        labelCol:4,
                        wrapperCol: 20,
                        labelAlign:'left'
                    },
                    'x-component': 'Input.TextArea',
                    'x-component-props': {
                        placeholder: '请输入描述',
                    }
                }
            }
        }
    }
}

    const formTab = FormTab.createFormTab()

    const getNewSchema = ()=>{
        const newSchema:{[k:string]:unknown} = {}
        if(!tabData || tabData.length === 0) return
        for(const tab of tabData){
        newSchema[tab.value!] = {
            type: 'object',
            'x-component': 'FormTab.TabPane',
            'x-component-props': {
                tab: tab.label,
                forceRender:true
            },
            properties:{
                _apinto_show:{
                    type: 'boolean',
                    title: '启用',
                    'x-decorator': 'FormItem',
                    'x-decorator-props': {
                        labelCol:4,
                        wrapperCol: 20,
                        labelAlign:'left'
                    },
                    'x-component': 'Switch',
                    'x-index':0,
                },
                _apinto_backend: {
                    type: 'void',
                    'x-component': 'void',
                    'x-reactions': {
                        dependencies: ['._apinto_show'],
                        fulfill: {
                            state: {
                                 visible: '{{!!$deps[0]}}',
                            },
                        },
                    },
                    properties: renderSchema[form.values.driver]?.properties,
                } 
            }
        }
        }

        const newSchema2 = {...pluginEditSchema,
            properties:{
                ...pluginEditSchema.properties,
                layout:{...pluginEditSchema.properties.layout,
                    properties:{
                        ...pluginEditSchema.properties.layout.properties,
                        config: {
                            type: 'object',
                            'x-component': 'void',
                            properties: {
                            tabForm:{
                                type: 'void',
                                'x-component': 'FormTab',
                                'x-component-props': {
                                    formTab: '{{formTab}}',
                                    centered:false,
                                    },
                                properties: newSchema
                            }
                            }
                        }
                    }
                }
            }
        }

        setSchema(newSchema2)
    }

    useEffect(() => {
        getNewSchema()
    }, [tabData,renderSchema,form.values.driver]);
    
    const save :()=>Promise<boolean | string> = ()=>{
        return new Promise((resolve, reject)=>{
            form.validate().then(()=>{
                const res = form.values
                const newData = JSON.parse(JSON.stringify(res)); // 深拷贝对象，避免直接修改 Proxy
                const config = newData.config;
              
                for (const tab in config) {
                  if (config[tab]._apinto_show) {
                    delete config[tab]._apinto_show; // 删除 _apinto_show 属性
                  } else {
                    delete config[tab]; // 删除整个 tab
                  }
                }
              
                fetchData<BasicResponse<null>>(type === 'add'?`dynamic/${moduleId}`:`dynamic/${moduleId}/config`,{method:type === 'add'? 'POST' : 'PUT',eoBody:newData, eoParams:{...(type !== 'add' && {id:initFormValue.id})}}).then(response=>{
                    const {code,msg} = response
                    if(code === STATUS_CODE.SUCCESS){
                        message.success(msg || '操作成功！')
                        resolve(true)
                    }else{
                        message.error(msg || '操作失败')
                        reject(msg || '操作失败')
                    }
                }).catch((errorInfo)=> reject(errorInfo))
            }).catch((errorInfo:unknown)=> reject(errorInfo))
        })
    }

    useImperativeHandle(ref, ()=>({
        save
    })
    )


    const getSkillData = async (skill: string) => {
        return new Promise((resolve,reject) => {
            fetchData<BasicResponse<{[k:string]:Array<{name:string,title:string}>}>>(`api/common/provider/${skill}`,{method:'GET'}).then(response=>{
                const {code,data,msg} = response
                if(code === STATUS_CODE.SUCCESS){
                    resolve(data[skill]?.map((x:{name:string,title:string})=>{return{label:x.title, value:x.name}}) || [])
                }else{
                    message.error(msg || '操作失败')
                    reject(msg || '操作失败')
                }
            })
        })
    }

    const useAsyncDataSource =
        (service: unknown, skill: string) => (field: unknown) => {
            field.loading = true
            service(skill).then(
                action.bound &&
                action.bound((data: unknown) => {
                    field.dataSource = data
                    field.loading = false
                })
            )
        }
    return (
        <div  className="pl-[12px]">
            <FormProvider form={form}>
                <SchemaField
                    schema={schema}
                    scope={{ useAsyncDataSource, getSkillData,formTab,form }}
                />
            </FormProvider>
        </div>)
})