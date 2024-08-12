
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