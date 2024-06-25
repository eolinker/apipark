/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-05-31 16:08:35
 * @FilePath: \frontend\packages\core\src\pages\system\authority\SystemAuthorityView.tsx
 */
import {Col, Row} from "antd";
import {useEffect, useState} from "react";
import { SystemAuthorityViewProps } from "../../../const/system/type";

export const SystemAuthorityView = ({entity}:SystemAuthorityViewProps)=>{
    const [detail,setDetail] = useState<Array<{key:string, value:string}>>(entity)

    useEffect(() => {
        setDetail(entity)
    }, [entity]);

    return (
        <div className="my-btnybase">{
            detail?.length > 0 && detail.map((k,i)=>(
                <Row className="leading-[32px]" key={i}>
                    <Col className="pr-[8px]" offset={1} span={5}>{k.key}:</Col>
                    <Col className="break-all" span={18}>{ k.value || '-'}</Col>
                </Row>
            ))
        }
        </div>
    )
}