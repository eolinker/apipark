/*
 * @Date: 2023-11-27 18:13:40
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-07 11:01:44
 * @FilePath: \frontend\packages\common\src\components\aoplatform\BasicLayout.tsx
 */
import { Button, Layout } from 'antd';
import Logo from '@common/assets/logo.png'
import Navigation from "./Navigation";
import TopBreadcrumb from "./Breadcrumb";
import {Outlet, useLocation, useNavigate} from "react-router-dom";
import {useEffect, useState} from "react";
import UserAvatar from "./UserAvatar.tsx";
import ErrorBoundary from './ErrorBoundary.tsx';
import { ShopOutlined } from '@ant-design/icons';
const { Header, Sider, Content } = Layout;

 function BasicLayout({project = 'core'}:{project:string}){
     const navigator = useNavigate()
     const location = useLocation()
     const currentUrl = location.pathname
    const query =new URLSearchParams(useLocation().search)
    const [isServiceHub,setIsServiceHub] = useState<boolean>(false)
    const [,setIsRentMng] = useState<boolean>(false)
    const [isMng,setIsMng] = useState<boolean>(false)
     useEffect(() => {
        setIsServiceHub(currentUrl.includes('/serviceHub'))
        setIsRentMng(currentUrl.includes('/tenantManagement'))
         setIsMng(project === 'core' && !currentUrl.includes('/serviceHub') && !currentUrl.includes('/tenantManagement'))
         if(currentUrl === '/'){
             navigator(  project === 'core' ?'/system/list':'/serviceHub/list')
         }
         
     }, [currentUrl]);

     const openServiceHub =()=>{
        isMng ? window.open(`/serviceHub/list`,'_blank') :  navigator('/serviceHub/list')
     }

     const backToPage =()=>{
        // const backUrl = query.get('callbackUrl') 
        // navigator(backUrl&& backUrl !== 'null' ?backUrl : '/')
        isServiceHub?window.open(`/tenantManagement/list`,'_blank') :  navigator('/serviceHub/list')
     }

    return(
        <Layout className="h-full w-full overflow-hidden">
            <Header className="border-0 border-b border-solid border-b-BORDER flex items-center px-[20px]">
                <div className="w-[175px] flex items-center">
                  <img
                    className="h-[32px]"
                    src={Logo}
                  />
                </div>
                <div className="flex justify-between items-center w-[calc(100%-175px)]">
                    <TopBreadcrumb />
                    <div className="flex justify-between items-center ">
                      <Button className="mr-btnbase" onClick={()=>isServiceHub ?backToPage() : openServiceHub()}>
                          {isServiceHub && <span className='flex items-center'><span className="mr-[4px] flex items-center"><iconpark-icon  className="" name="auto-generate-api"></iconpark-icon></span>租户管理</span>}
                          {!isServiceHub && <span className='flex items-center'><ShopOutlined className="mr-[4px]" />服务市场</span>}
                      </Button>
                      {/* <Button  className="mr-[20px]">
                        <span className='flex items-center'><QuestionCircleOutlined className="mr-[4px]" />帮助文档</span>
                      </Button> */}
                      <UserAvatar />
                    </div>
                </div>
            </Header>
            <Layout hasSider={isMng}>
                {!isMng? undefined :<Sider width={192} className="border-r border-r-BORDER overflow-hidden hover:overflow-y-auto " >
                  <Navigation></Navigation>
                </Sider>}
                <Content className="h-full block">
                    <ErrorBoundary>
                        <Outlet />
                    </ErrorBoundary>
                </Content>
            </Layout>
        </Layout>
    )
}
export default BasicLayout