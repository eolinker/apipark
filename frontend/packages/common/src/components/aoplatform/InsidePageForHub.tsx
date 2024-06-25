/*
 * @Date: 2023-11-28 11:54:17
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-02 09:38:42
 * @FilePath: \frontend\packages\core\src\components\InsidePageForHub.tsx
 * @Description:内页（带顶部导航与描述、tag、右侧按钮-可选）
 */

import { Button, Tag } from "antd"
import {useNavigate} from "react-router-dom";
import WithPermission from "@common/components/aoplatform/WithPermission";
import { FC, ReactNode } from "react";


class InsidePageProps {
    showBanner?:boolean = true
    pageTitle:string = ''
    tagList?:Array<{label:string|ReactNode}> = []
    children:React.ReactNode
    showBtn?:boolean = false
    btnTitle?:string = ''
    description?:string = ''
    onBtnClick?:()=>void
    backUrl:string = '/'
    btnAccess?:string
}

const InsidePageForHub:FC<InsidePageProps> = ({showBanner=true,pageTitle,tagList,showBtn,btnTitle,btnAccess,description,children,onBtnClick,backUrl})=>{
    const navigate = useNavigate();

    const goBack = () => {
        // document.referrer 有误，隐藏此段逻辑
        // const currentPath = window.location.pathname; // 获取当前路径
        // // 注意：在使用 React Router 时，直接操作 window.location 可能不是最佳实践，但用于获取当前路径是可接受的
        // const referrerPath = document.referrer ? new URL(document.referrer).pathname : ''; // 获取上一步的路径
    
        // // 定义一个函数来提取路由的开始部分
        // const getRouteStart = (path:string) => {
        //   const parts = path.split('/').filter(Boolean); // 分割路径，并移除空字符串
        //   return parts.length > 0 ? parts[0] : null; // 返回路由的开始部分
        // };
    
        // const currentRouteStart = getRouteStart(currentPath);
        // const referrerRouteStart = getRouteStart(referrerPath);
        // // 比较当前路径和上一步路径的路由开始部分是否相同
        // // console.log(currentPath,document,document.referrer, referrerPath, currentRouteStart, referrerRouteStart)
        // if (currentRouteStart && referrerRouteStart && currentRouteStart === referrerRouteStart) {
        //   // 如果相同，跳转到`${currentRouteStart}/list`
        //   navigate(backUrl ?? `/${currentRouteStart}/list`);
        // } else {
        //   // 如果不同，使用navigate实现返回上一步
        //   navigate(-1);
        // }

        navigate(backUrl);
    };
    return (
        <div className="h-full flex flex-col flex-1 overflow-hidden max-w-[1500px] m-auto">
            { showBanner &&  <div className="p-btnbase  mx-[4px]">
                <p className="mb-[16px] h-[20px] "><a className="leading-[20px] inline-block "onClick={goBack}> {"<  返回列表"}</a></p>
                <div className="flex justify-between">
                    <div className="">
                        <span className="text-[20px] font-bold">{pageTitle}</span>
                        {tagList && tagList?.length > 0 && tagList?.map((tag)=>{
                            return ( <Tag key={tag.label as string} bordered={false}>{tag.label}</Tag>)
                        })}
                    </div>
                    {showBtn && <WithPermission access={btnAccess}><Button type="primary" onClick={()=> {
                        onBtnClick&&onBtnClick()
                    }}>{btnTitle}</Button></WithPermission>}
                </div>
                <p className="text-[14px] leading-[22px] text-[#999]">
                    {description}
                </p>
            </div>}
            <div className="h-full overflow-y-hidden">{children}</div>
        </div>
    )
}

export default InsidePageForHub