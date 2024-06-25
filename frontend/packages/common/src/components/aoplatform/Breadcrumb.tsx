/*
 * @Date: 2023-11-28 11:41:17
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-05-29 18:30:51
 * @FilePath: \frontend\packages\common\src\components\aoplatform\Breadcrumb.tsx
 */
import { Breadcrumb } from "antd"
import { useBreadcrumb} from "@common/contexts/BreadcrumbContext.tsx";
import {FC,useEffect} from "react";


const TopBreadcrumb: FC = () => {
     const { breadcrumb } = useBreadcrumb()
    useEffect(() => {
    }, [breadcrumb]);
    return (
        <Breadcrumb items={breadcrumb} />
    )
}

export default TopBreadcrumb