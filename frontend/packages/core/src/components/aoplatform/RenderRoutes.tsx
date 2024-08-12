/*
 * @Date: 2024-06-04 08:54:24
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-07-12 20:08:54
 * @FilePath: \frontend\packages\core\src\components\aoplatform\RenderRoutes.tsx
 */
import { BrowserRouter as Router, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import Login from "@core/pages/Login.tsx"
import BasicLayout from '@common/components/aoplatform/BasicLayout';
import {createElement, ReactElement,ReactNode,Suspense} from 'react';
import { v4 as uuidv4 } from 'uuid'
import {App, Skeleton} from "antd";
import ApprovalPage from "@core/pages/approval/ApprovalPage.tsx";
import {SystemProvider} from "@core/contexts/SystemContext.tsx";
import {useGlobalContext} from "@common/contexts/GlobalStateContext.tsx";
import {FC,lazy} from 'react';
import { TeamProvider } from '@core/contexts/TeamContext.tsx';
import SystemOutlet from '@core/pages/system/SystemOutlet.tsx';
import { DashboardProvider } from '@core/contexts/DashboardContext.tsx';
import { PartitionProvider } from '@core/contexts/PartitionContext.tsx';
import { TenantManagementProvider } from '@market/contexts/TenantManagementContext.tsx';

type RouteConfig = {
    path:string
    component?:ReactElement
    children?:(RouteConfig|false)[]
    key:string
    provider?:FC<{ children: ReactNode; }>
    lazy?:unknown
}
const APP_MODE = import.meta.env.VITE_APP_MODE;

export type RouterParams  = {
    orgId:string
    teamId:string
    systemId:string
    apiId:string
    serviceId:string
    partitionId:string
    clusterId:string;
    memberGroupId:string
    userGroupId:string
    pluginName:string
    moduleId:string
    accessType:'project'|'team'|'system'
    categoryId:string
    tagId:string
    dashboardType:string
    dashboardDetailId:string
    topologyId:string
    appId:string
}

const PUBLIC_ROUTES:RouteConfig[] = [
    {
        path:'/',
        component:<Login/>,
        key: uuidv4(),
    },
    {
        path:'/login',
        component:<Login/>,
        key: uuidv4()
    },
    {
        path:'/',
        component:<ProtectedRoute/>,
        key: uuidv4(),
        children:[
            {
                path:'approval/*',
                component:<ApprovalPage />,
                key:uuidv4()
            },
            { 
                path:'organization/*',
                component:<Outlet/>,
                key: uuidv4(),
                children:[
                    { 
                        path:'list',
                        key: uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/organization/OrganizationList.tsx'))
                    }
                ]
            },
            {
                path:'team/*',
                component:<Outlet/>,
                key: uuidv4(),
                provider: TeamProvider,
                children:[
                    {
                        path:'list',
                        key: uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/team/TeamList.tsx'))
                    },
                    {
                        path:'inside/:orgId/:teamId',
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/team/TeamInsidePage.tsx')),
                        key: uuidv4(),
                        children:[
                            {
                                path:'member',
                                key: uuidv4(),
                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/team/TeamInsideMember.tsx')),
                            },
                            {
                                path:'access',
                                key:uuidv4(),
                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/team/TeamInsideAccess.tsx')),
                            },
                            {
                                path:'setting',
                                key: uuidv4(),
                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/team/TeamConfig.tsx')),
                            },
                        ]
                    }
                ]
            },
            {
                path:'system/*',
                component:<SystemOutlet />,
                key: uuidv4(),
                provider: SystemProvider,
                children:[
                    {
                        path:'list',
                        key: uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/SystemList.tsx')),
                    },
                    {
                        path:'list/:orgId/:teamId',
                        key: uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/SystemList.tsx')),
                    },
                    {
                        path:':orgId/:teamId',
                        component:<Outlet/>,
                        key: uuidv4(),
                        children:[
                            {
                                path:'inside/:systemId',
                                key: uuidv4(),
                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/SystemInsidePage.tsx')),
                                children:[
                                    {
                                        path:'api',
                                        key: uuidv4(),
                                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/api/SystemInsideApiList.tsx')),
                                    },
                                    {
                                        path:'upstream',
                                        key: uuidv4(),
                                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/upstream/SystemInsideUpstreamContent.tsx')),
                                    },
                                    {
                                        path:'myService',
                                        key: uuidv4(),
                                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/myService/SystemInsideMyService.tsx')),
                                    },
                                    {
                                        path:'subService',
                                        key: uuidv4(),
                                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/subSubscribe/SystemInsideSubService.tsx')),
                                        children:[

                                        ]
                                    },
                                    {
                                        path:'subscriber',
                                        key: uuidv4(),
                                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/SystemInsideSubscriber.tsx')),
                                        children:[

                                        ]
                                    },
                                    {
                                        path:'approval',
                                        key: uuidv4(),
                                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/approval/SystemInsideApproval.tsx')),
                                        children:[
                                            {
                                                path:'*',
                                                key: uuidv4(),
                                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/approval/SystemInsideApprovalList.tsx')),
                                            }
                                        ]
                                    },
                                    {
                                        path:'topology',
                                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/SystemTopology.tsx')),
                                        key: uuidv4(),
                                        children:[
                                        ]
                                    },
                                    {
                                        path:'authority',
                                        key: uuidv4(),
                                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/authority/SystemInsideAuthority.tsx')),
                                        children:[

                                        ]
                                    },
                                    {
                                        path:'publish',
                                        key: uuidv4(),
                                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/publish/SystemInsidePublish.tsx')),
                                        children:[
                                            {
                                                path:'*',
                                                key: uuidv4(),
                                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/publish/SystemInsidePublishList.tsx')),
                                            }
                                        ]
                                    },
                                    {
                                        path:'access',
                                        key: uuidv4(),
                                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/SystemInsideAccess.tsx')),
                                    },
                                    {
                                        path:'member',
                                        key: uuidv4(),
                                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/SystemInsideMember.tsx')),
                                    },
                                    {
                                        path:'setting',
                                        key: uuidv4(),
                                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/system/SystemConfig.tsx')),
                                        children:[

                                        ]
                                    },
                                ]
                            }
                        ]
                    }
                ]
            },
            {
                path:'partition/*',
                component:<Outlet/>,
                key: uuidv4(),
                children:[
                    {
                        path:'list',
                        key: uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/partitions/PartitionList.tsx')),
                    },
                    {
                        path:'inside/:partitionId',
                        key: uuidv4(),
                        provider:PartitionProvider,
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/partitions/PartitionInsidePage.tsx')),
                        children:[
                            {
                                path:'cluster',
                                key: uuidv4(),
                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/partitions/PartitionInsideCluster.tsx')),
                            },
                            {
                                path:'cert',
                                key: uuidv4(),
                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/partitions/PartitionInsideCert.tsx')),
                            },
                            APP_MODE === 'pro' &&{
                                path:'dashboard_setting',
                                key: uuidv4(),
                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/partitions/PartitionInsideDashboardSetting.tsx')),
                            },
                            {
                                path:'setting',
                                key: uuidv4(),
                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/partitions/PartitionConfig.tsx')),
                            },
                            {
                                path:'template/:moduleId',
                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@common/components/aoplatform/intelligent-plugin/IntelligentPluginList.tsx')),
                                key:uuidv4()
                            }
                        ]
                    }
                ]
            },
            {
                path:'serviceHub',
                component:<Outlet />,
                key:uuidv4(),
                children:[
                    {
                        path:'list',
                        key:uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@market/pages/serviceHub/ServiceHubList.tsx')),
                    },
                    {
                        path:'detail/:serviceId',
                        key:uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@market/pages/serviceHub/ServiceHubDetail.tsx')),
                    }]
            },
            {
                path:'servicecategories',
                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/serviceCategory/ServiceCategory.tsx')),
                key:uuidv4(),
            },
            {
                path:'tenantManagement',
                component:<Outlet />,
                provider:TenantManagementProvider,
                key:uuidv4(),
                children:[
                    {
                        path:':teamId/inside/:appId',
                        key:uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@market/pages/serviceHub/management/ManagementInsidePage.tsx')),
                        children:[
                            {
                                path:'service/:partitionId',
                                key:uuidv4(),
                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@market/pages/serviceHub/management/ManagementInsideService.tsx')),
                            },
                            {
                                path:'authorization',
                                key:uuidv4(),
                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@market/pages/serviceHub/management/ManagementInsideAuth.tsx')),
                            },
                            {
                                path:'setting',
                                key:uuidv4(),
                                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@market/pages/serviceHub/management/ManagementAppSetting.tsx')),
                            },
                        ]
                    },
                    {
                        path:'list',
                        key:uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@market/pages/serviceHub/management/ServiceHubManagement.tsx')),
                    },
                    {
                        path:'list/:teamId',
                        key:uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@market/pages/serviceHub/management/ServiceHubManagement.tsx')),
                    },
                ]
            },
            {
                path:'member/*',
                key:uuidv4(),
                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/member/MemberPage.tsx')),
                children:[
                    {
                        path:'list',
                        key:uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/member/MemberList.tsx')),
                    },
                    {
                        path:'list/:memberGroupId',
                        key:uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/member/MemberList.tsx')),
                    }
                ]
            },
            {
                path:'user/*',
                key:uuidv4(),
                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/user/UserPage.tsx')),
                children:[
                    {
                        path:'list',
                        key:uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/user/UserList.tsx')),
                    },
                    {
                        path:'list/:userGroupId',
                        key:uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/user/UserList.tsx')),
                    }
                ]
            },
            {
                path:'role/*',
                key:uuidv4(),
                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/role/RoleList.tsx')),
            },
            {
                path:'access',
                key:uuidv4(),
                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/access/AccessPage.tsx')),
                children:[
                    {
                        path:':accessType',
                        key:uuidv4(),
                        lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/access/AccessList.tsx')),
                    },
                ]
            },
            {
                path:'logretrieval',
                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/logRetrieval/LogRetrieval.tsx')),
                key:uuidv4(),
            },
            {
                path:'auditlog',
                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/auditLog/AuditLog.tsx')),
                key:uuidv4(),
            },
            {
                path:'assets',
                component:<p>设计中</p>,
                key:uuidv4()
            },
            {
                path:'template/:moduleId',
                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@common/components/aoplatform/intelligent-plugin/IntelligentPluginList.tsx')),
                key:uuidv4()
            },
            {
                path:'logsettings/*',
                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/logsettings/LogSettings.tsx')),
                key: uuidv4(),
                children:[{
                    path:'template/:moduleId',
                    lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@common/components/aoplatform/intelligent-plugin/IntelligentPluginList.tsx')),
                    key:uuidv4()
                }]
                
            },
            APP_MODE ==='pro' && {
                path:'resourcesettings/*',
                lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@core/pages/resourcesettings/ResourceSettings.tsx')),
                key: uuidv4(),
                children:[{
                    path:'template/:moduleId',
                    lazy:lazy(() => import(/* webpackChunkName: "[request]" */ '@common/components/aoplatform/intelligent-plugin/IntelligentPluginList.tsx')),
                    key:uuidv4()
                }]
                
            }
        ]
    },
]

const RenderRoutes = ()=> {
    return (
        <App className="h-full" message={{ maxCount: 1 }}>
            <Router>
                <Routes>
                    {generateRoutes(PUBLIC_ROUTES)}
                    </Routes>
            </Router>
        </App>
        )
}

const generateRoutes = (routerConfig: RouteConfig[]) => {
    return routerConfig?.map((route: RouteConfig) => {
            let routeElement;
            if (route.lazy) {
                const LazyComponent = route.lazy as React.ExoticComponent<unknown>;

                routeElement = (
                    <Suspense fallback={ <div className=''><Skeleton className='m-btnbase w-[calc(100%-20px)]' active /></div>}>
                        {route.provider ? (
                            createElement(route.provider, {}, <LazyComponent  />)
                        ) : (
                            <LazyComponent />
                        )}
                    </Suspense>
                );
            } else {
                routeElement = route.provider ? (
                    createElement(route.provider, {}, route.component)
                ) : (
                    route.component
                );
            }

                return (
                  <Route
                    key={route.key}
                    path={route.path}
                    element={routeElement}
                  >
                    {route.children && generateRoutes(route.children as RouteConfig[])}
                  </Route>
                );
              }
        )
}

// 保护的路由组件
function ProtectedRoute() {
    const {state} = useGlobalContext()
    return state.isAuthenticated? <BasicLayout project="core" /> : <Navigate to="/login" />;
  }

export default RenderRoutes