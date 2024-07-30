/*
 * @Date: 2024-05-30 12:07:53
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-06 17:01:16
 * @FilePath: \frontend\packages\market\src\pages\serviceHub\ServiceHubGroup.tsx
 */
import {debounce} from "lodash-es";
import {SearchOutlined} from "@ant-design/icons";
import {App, Divider, Input, TreeDataNode} from "antd";
import  {useCallback, useEffect, useState} from "react";
import Tree, {DataNode, TreeProps} from "antd/es/tree";
import {BasicResponse, STATUS_CODE} from "@common/const/const.ts";
import {useFetch} from "@common/hooks/http.ts";
import { CategorizesType, TagType } from "../../const/serviceHub/type.ts";
import { PartitionItem } from "@common/const/type.ts";
import { filterServiceList, initialServiceHubListState, SERVICE_HUB_LIST_ACTIONS, ServiceHubListActionType } from "./ServiceHubList.tsx";

type ServiceHubGroup = {
    children:JSX.Element
    filterOption:typeof initialServiceHubListState
    dispatch:React.Dispatch<ServiceHubListActionType>
}

export const ServiceHubGroup = ({children,filterOption,dispatch}:ServiceHubGroup)=>{
    const {message} = App.useApp()
    const {fetchData} = useFetch()
    // const [treeHeight, setTreeHeight] = useState<number>(Math.ceil((window.innerHeight - 50 - 20 * 2 - (32 + 4 + 20) - ( 12 * 2 + 1 ) * 2 - ( 25 + 15 ) * 3) /3) )
    
    useEffect(() => {
        getTagAndServiceClassifyList()
        getPartitionList()
        // const handleResize = () => {
        //     setTreeHeight(Math.ceil((window.innerHeight - 50 - 20 * 2 - (32 + 4 + 20) - ( 12 * 2 + 1 ) * 2 - ( 25 + 15 ) * 3) /3))
        // };
    
        // const debouncedHandleResize = debounce(handleResize, 200);
    
        // // 监听窗口大小变化
        // window.addEventListener('resize', debouncedHandleResize);
        // handleResize();
        // return () => {
        // window.removeEventListener('resize', debouncedHandleResize);
        // };
    }, []);

    const onSearchWordChange = (e:string)=>{
        dispatch({type:SERVICE_HUB_LIST_ACTIONS.SET_KEYWORD,payload:e})
        dispatch({type:SERVICE_HUB_LIST_ACTIONS.LIST_LOADING,payload:true})
        dispatch({type:SERVICE_HUB_LIST_ACTIONS.SET_SERVICES,payload: filterServiceList({...filterOption,keyword:e})})
        dispatch({type:SERVICE_HUB_LIST_ACTIONS.LIST_LOADING,payload:false})
    }

    const getTagAndServiceClassifyList = ()=>{
        fetchData<BasicResponse<{ catalogues:CategorizesType[],tags:TagType[]}>>('catalogues',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                dispatch({type:SERVICE_HUB_LIST_ACTIONS.GET_CATEGORIES,payload:data.catalogues})
                dispatch({type:SERVICE_HUB_LIST_ACTIONS.GET_TAGS,payload:[...data.tags,{id:'empty',name:'(空标签)'}]})
                dispatch({type:SERVICE_HUB_LIST_ACTIONS.SET_SELECTED_CATE,payload:[...data.catalogues.map((x:CategorizesType)=>x.id)]})
                dispatch({type:SERVICE_HUB_LIST_ACTIONS.SET_SELECTED_TAG,payload:[...data.tags.map((x:TagType)=>x.id),'empty']})
            }else{
                message.error(msg || '操作失败')
            }
        })
    }
    
    const getPartitionList = ()=>{
        return fetchData<BasicResponse<{partitions:PartitionItem[]}>>('simple/partitions',{method:'GET'}).then(response=>{
            const {code,data,msg} = response
            if(code === STATUS_CODE.SUCCESS){
                dispatch({type:SERVICE_HUB_LIST_ACTIONS.GET_PARTITIONS,payload:data.partitions})
                dispatch({type:SERVICE_HUB_LIST_ACTIONS.SET_SELECTED_PARTITION,payload:data.partitions.map((x:PartitionItem)=>x.id)})
                return Promise.resolve(data.partitions)
            }else{
                message.error(msg || '操作失败')
                return Promise.reject(msg || '操作失败')
            }
        })
    }

    const transferToTreeData = useCallback((data:CategorizesType[] | TagType[] | PartitionItem[]):TreeDataNode[]=>{
        const loop = (data: CategorizesType[] | TagType[] | PartitionItem[]): DataNode[] =>
            data?.map((item) => {
                if ((item as CategorizesType).children) {
                    return {
                        title:item.name,
                        key: item.id, children: loop((item as CategorizesType).children)
                    };
                }
                return {
                    title:item.name,
                    key: item.id,
                };
            });
        return loop(data || [])
    },[])

    const onCheckHandler = (type: 'SET_SELECTED_CATE' | 'SET_SELECTED_TAG' | 'SET_SELECTED_PARTITION') => (checkedKeys:string[]) => {
        dispatch({ type: SERVICE_HUB_LIST_ACTIONS[type], payload: checkedKeys });
        dispatch({type:SERVICE_HUB_LIST_ACTIONS.LIST_LOADING,payload:true})
        
        dispatch({type:SERVICE_HUB_LIST_ACTIONS.SET_SERVICES,payload: filterServiceList({...filterOption,[(type === 'SET_SELECTED_CATE' ? 'selectedCate' : type === 'SET_SELECTED_TAG' ? 'selectedTag' : 'selectedPartition' ) as keyof typeof filterOption]: checkedKeys })})
        dispatch({type:SERVICE_HUB_LIST_ACTIONS.LIST_LOADING,payload:false})
    };


    return (
        <div className="flex flex-1 h-full">
            <div className="w-[250px] border-0 border-solid border-r-[1px] border-r-BORDER">
            <div className=" h-full">
                <Input className="rounded-SEARCH_RADIUS m-[10px] h-[40px] bg-[#f8f8f8] w-[230px]" onChange={(e) => debounce(onSearchWordChange, 500)(e.target.value)}
                    allowClear placeholder="搜索服务"
                    prefix={<SearchOutlined className="cursor-pointer"/>}/>
                    <div className="h-[calc(100%-60px)] overflow-auto">
                        <div className="mt-[20px] ml-[20px] pr-[10px] ">
                            <p className="text-[18px] h-[25px] leading-[25px] font-bold mb-[15px]">分类</p>
                            <Tree
                                className={`no-selected-tree ${transferToTreeData(filterOption.categoriesList).filter(x=>x.children && x.children.length > 0).length > 0 ? '' : 'no-first-switch-tree'}`}
                                checkable
                                blockNode={true}
                                checkedKeys={filterOption.selectedCate}
                                onCheck={onCheckHandler('SET_SELECTED_CATE')}
                                treeData={transferToTreeData(filterOption.categoriesList)}
                                showIcon={false}
                                selectable={false}
                                />
                        </div>
                        <Divider  className="my-[20px]" />
                        <div className="ml-[20px] pr-[10px]">
                        <p className="text-[18px] h-[25px] leading-[25px] font-bold mb-[15px]">标签</p>
                            <Tree
                                className="no-first-switch-tree no-selected-tree"
                                checkable
                                blockNode={true}
                                checkedKeys={filterOption.selectedTag}
                                onCheck={onCheckHandler('SET_SELECTED_TAG')}
                                treeData={transferToTreeData(filterOption.tagsList)}
                                showLine={false}
                                showIcon={false}
                                selectable={false}
                                />
                        </div>
                        <Divider className="my-[20px]" />
                        <div className="ml-[20px] pr-[10px]">
                        <p className="text-[18px] h-[25px] leading-[25px] font-bold mb-[15px]">环境</p>
                            <Tree
                                className="no-first-switch-tree no-selected-tree"
                                checkable
                                blockNode={true}
                                checkedKeys={filterOption.selectedPartition}
                                onCheck={onCheckHandler('SET_SELECTED_PARTITION')}
                                treeData={transferToTreeData(filterOption.partitionList)}
                                showLine={false}
                                showIcon={false}
                                selectable={false}
                                />
                        </div>
                </div>
            </div>
        </div>
        <div className="w-[calc(100%-224px)]">
          {children}
        </div>
    </div>);
}

export default ServiceHubGroup
