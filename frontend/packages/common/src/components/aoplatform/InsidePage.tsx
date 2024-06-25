/*
 * @Date: 2023-11-28 11:54:17
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-05-15 15:03:46
 * @FilePath: \frontend\packages\core\src\components\InsidePage.tsx
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
    backUrl?:string = '/'
    btnAccess?:string
}

const InsidePage:FC<InsidePageProps> = ({showBanner=true,pageTitle,tagList,showBtn,btnTitle,btnAccess,description,children,onBtnClick,backUrl})=>{
    const navigate = useNavigate();

    const goBack = () => {
        navigate(backUrl || '/');
    };
    return (
        // <div className="h-full flex flex-col flex-1 overflow-hidden bg-[#f7f8fa]">
        <div className="h-full flex flex-col flex-1 overflow-hidden  ">
            { showBanner &&  <div className="p-btnbase  mx-[4px] border-[0px] border-b-[1px] border-solid border-BORDER">
                {backUrl && <p className="mb-[16px] h-[20px] "><a className="leading-[20px] inline-block "onClick={goBack}> {"<  返回"}</a></p>}
                <div className="flex justify-between">
                    <div className="flex items-center">
                        <span className="text-[20px] mr-[8px] leading-[32px]">{pageTitle}</span>
                        {tagList && tagList?.length > 0 && tagList?.map((tag)=>{
                            return ( <Tag className="" key={tag.label as string} bordered={false} >{tag.label}</Tag>)
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

export default InsidePage