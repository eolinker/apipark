import {FC, useEffect, useMemo, useState} from 'react';
import type { MenuProps } from 'antd';
import { Menu } from 'antd';
import { useLocation, useNavigate} from "react-router-dom";
import { getItem } from '@common/utils/navigation';
import { PERMISSION_DEFINITION } from '@common/const/permissions';
import { useGlobalContext } from '@common/contexts/GlobalStateContext';
import { DashboardOutlined, DeploymentUnitOutlined, HddOutlined, TeamOutlined } from '@ant-design/icons';

export type MenuItem = Required<MenuProps>['items'][number];

// type NavigationItemType = {
//   title: string,
//   iconType: string,
//   icon: string,
//   router: string,
//   access:string [],
//   children: NavigationItemType[]
// }
const APP_MODE = import.meta.env.VITE_APP_MODE;

const routerKeyMap = new Map([
  ['assets','/assets'],
  ['dashboard','/dashboard'],
  ['systemrunning','/systemrunning'],
  ['system','/system/list'],
  ['servicecategories','/servicecategories'],
  ['organization','/organization/list'],
  ['team','/team/list'],
  ['member','/member/list'],
  ['user','/user/list'],
  ['role','/role'],
  ['access','/access'],
  ['partition','/partition/list'],
  ['openapi','/openapi'],
  ['logsettings','/logsettings'],
  ['resourcesettings','/resourcesettings']])

  
const TOTAL_MENU_ITEMS: MenuProps['items'] = [
  APP_MODE === 'pro' ? getItem('仪表盘', 'mainPage', <DashboardOutlined />,[
    // getItem(<a >资产视图</a>, 'assets',null,undefined,undefined,''),
    getItem(<a >运行视图</a>, 'dashboard',null,undefined,undefined,''),
    getItem(<a >系统拓扑图</a>, 'systemrunning',null,undefined,undefined,''),
    // getItem((<Link to="/approval"  className="flex items-center"><span className='mr-[4px]'>审批</span><Badge size="small" count={2} /></Link>), 'approval', null),
  ]):null,

  getItem('数据服务资产', 'dataAssets',<DeploymentUnitOutlined />, [
    getItem(<a>内部数据服务</a>, 'system',null,undefined,undefined,''),
    getItem(<a>服务分类管理</a>, 'servicecategories',null,undefined,undefined,''),
  ]),

  getItem('组织架构', 'operationCenter',<TeamOutlined />, [
    getItem(<a>成员与部门</a>, 'member',null,undefined,undefined,'system.member.self.view'),
    getItem(<a>组织</a>, 'organization',null,undefined,undefined,'system.organization.self.view'),
    getItem(<a>团队</a>, 'team',null,undefined,undefined,''),
    getItem(<a>用户组</a>, 'user',null,undefined,undefined,'system.user.self.view'),
    getItem(<a>自定义角色</a>, 'role',null,undefined,undefined,'system.role.self.view'),
    getItem(<a>权限配置</a>, 'access',null,undefined,undefined,'system.access.self.view'),
  ]),

  getItem('运维与集成', 'maintenanceCenter', <HddOutlined />, [
    getItem(<a>部署管理</a>, 'partition',null,undefined,undefined,'system.partition.self.view'),
    getItem(<a>日志配置</a>, 'logsettings',null,undefined,undefined,'system.partition.self.view'),
    APP_MODE === 'pro' ? getItem(<a>资源配置</a>, 'resourcesettings',null,undefined,undefined,'system.partition.self.view'):null,
    // getItem(<Link to="/email">邮箱设置</Link>, 'email'),
    APP_MODE === 'pro' ? getItem(<a>Open API</a>, 'openapi',null,undefined,undefined,'system.openapi.self.view'):null,
    // getItem(<Link to="/webhook">Webhook</Link>, 'webhook'),
    // getItem(<Link to="/template/httplog">HTTP 日志配置</Link>, 'httplog'),
    // getItem(<Link to="/logretrieval">日志检索</Link>, 'logretrieval',null,undefined,undefined,'system.logRetrieval.self.view'),
    // getItem(<Link to="/auditlog">审计日志</Link>, 'auditlog'),
  ]),
];

const Navigation: FC = () => {
  // const { message } = App.useApp()
  const location = useLocation()
  const [selectedKeys, setSelectedKeys] = useState<string>('')
  const currentUrl = location.pathname
  const navigateTo = useNavigate()
  // const {fetchData} = useFetch()
  // const [navigationItemList, setNavigationItemList] = useState<NavigationItemType[]>([])
  const { accessData,checkPermission} = useGlobalContext()

  const onClick: MenuProps['onClick'] = (e) => {
    if(location.pathname.split('/')[1] === e.key) return
    const newUrl = routerKeyMap.get(e.key)
    newUrl && navigateTo(newUrl)
  };

  // useEffect(()=>{
  //   const newUrl = routerKeyMap.get(selectedKeys)
  //   console.log(newUrl)
  //   newUrl && navigateTo(newUrl)
  // },[selectedKeys])

  // 插件相关，暂时隐藏
  // const getNavigationData = ()=>{
  //   fetchData<BasicResponse<{ navigation: NavigationItemType[] }>>('navigation',{method:'GET',eoApiPrefix:'_system/'},).then(response=>{
  //         const {code,data,msg} = response
  //         if(code === STATUS_CODE.SUCCESS){
  //           setNavigationItemList(data.navigation)
  //         }else{
  //             message.error(msg || '操作失败')
  //         }
  //     })
  // }

  // const navigations = useMemo(()=>{
  //  const getNavTitle = (data:NavigationItemType,root?:boolean)=>{
  //       if(root){
  //         return data.title
  //       }
  //       return <Link to={data.router} >{data.title}</Link>
  //   }

  //   const handleNavigationData:(data:NavigationItemType[],root?:boolean)=>MenuItem[] = (data:NavigationItemType[],root?:boolean)=>{
  //     return data?.map(x=>getItem(
  //       getNavTitle(x,root),
  //       x.router.split('/')[0] === 'template' ? x.router.split('/')[1] : x.router.split('/')[0],
  //       x.icon,
  //       x.children && handleNavigationData(x.children)))
  //   }

  //   return handleNavigationData(navigationItemList,true)
  // },[navigationItemList])

  
  const menuData = useMemo(()=>{
    const filterMenu = (menu:Array<{[k:string]:unknown}>)=>{
        return menu.filter(x=> x && (x.access ? checkPermission(x.access as keyof typeof PERMISSION_DEFINITION[0]): true))
    }
    return TOTAL_MENU_ITEMS!.filter(x=>x).map((x)=> ( x.children ? {...x, children:filterMenu(x.children)} : x))?.filter(x=> x.key === 'system' || (x.children && x.children?.length > 0))
},[accessData])

  useEffect(() => {
    setSelectedKeys(currentUrl.split('/')[1] === 'template' ? currentUrl.split('/')[2] : currentUrl.split('/')[1])
    // getNavigationData()
  }, [currentUrl]);

  return (
    <Menu
      onClick={onClick}
      style={{height:'100%' }}
      selectedKeys={[selectedKeys]}
      defaultOpenKeys={['mainPage','dataAssets','operationCenter','maintenanceCenter']}
      mode="inline"
      items={[...menuData]}
    />
  );
};

export default Navigation;