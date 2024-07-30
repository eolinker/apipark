import { getItem } from "@common/utils/navigation";
import { ProColumns } from "@ant-design/pro-components";
import { PartitionCertTableListItem, PartitionClusterNodeModalTableListItem, PartitionClusterNodeTableListItem, PartitionClusterTableListItem, PartitionTableListItem } from "./types";
import { ColumnType } from "antd/es/table";
import CopyAddrList from "@common/components/aoplatform/CopyAddrList";
import { MenuProps } from "antd";
import { Link } from "react-router-dom";

const APP_MODE = import.meta.env.VITE_APP_MODE;

export const PARTITIONS_INNER_MENU: MenuProps['items'] = [
    getItem('管理', 'grp', null,
        [getItem(<Link to="cluster">集群</Link>, 'cluster',undefined,undefined,undefined,'system.partition.cluster.view'),
            getItem(<Link to="cert">证书管理</Link>, 'cert',undefined,undefined,undefined,'system.partition.cert.view'),
            APP_MODE ==='pro' ? getItem(<Link to="dashboard_setting">监控配置</Link>, 'dashboard_setting',undefined,undefined,undefined,'system.partition.self.view'):null,
            getItem(<Link to="setting">环境设置</Link>, 'setting',undefined,undefined,undefined,'system.partition.self.view')],
        'group'),
];


export const PARTITION_CERT_TABLE_COLUMNS: ProColumns<PartitionCertTableListItem>[] = [
    {
        title: '证书',
        dataIndex: 'name',
        copyable: true,
        ellipsis:true,
        width:160,
        fixed:'left',
        sorter: (a,b)=> {
            return a.name.localeCompare(b.name)
        },
    },
    {
        title: '绑定域名',
        dataIndex: 'domains',
        renderText:(_,entity) =>(
            entity.domains.join(',')
        ),
        copyable: true,
        ellipsis:true
    },
    {
        title: '证书有效期',
        ellipsis: true,
        dataIndex: 'notAfter',
        copyable: true,
        renderText: (value:string,entity:PartitionCertTableListItem) => {
            return `${entity.notBefore} - ${entity.notAfter}`
        },
    },
    {
        title: '更新者',
        dataIndex: ['updater','name'],
        ellipsis: true,
        filters: true,
        onFilter: true,
        valueType: 'select',
        filterSearch: true
    },
    {
        title: '更新时间',
        key: 'updateTime',
        dataIndex: 'updateTime',
        ellipsis:true,
        width:182,
        sorter: (a,b)=> {
            return a.updateTime.localeCompare(b.updateTime)
        },
    },
];

export const PARTITION_CLUSTER_TABLE_COLUMNS : ProColumns<PartitionClusterTableListItem>[] = [
    {
        title: '集群名称',
        dataIndex: 'name',
        copyable: true,
        ellipsis:true,
        width:160,
        fixed:'left',
        sorter: (a,b)=> {
            return a.name.localeCompare(b.name)
        },
    },
    {
        title: '集群 ID',
        dataIndex: 'id',
        width: 140,
        copyable: true,
        ellipsis:true
    },
    {
        title: '状态',
        dataIndex: 'status',
        ellipsis:true,
        valueType: 'select',
        filters: true,
        onFilter: true,
        valueEnum: new Map([
            [0, <span className="text-status_fail">异常</span>],
            [1,<span className="text-status_success">正常</span>],
        ])
    },
    {
        title: '描述',
        dataIndex: 'description',
        copyable: true,
        ellipsis:true
    }
];


export const PARTITION_CLUSTER_NODE_COLUMNS: ProColumns<PartitionClusterNodeTableListItem>[] = [
    {
        title: '节点名称',
        dataIndex: 'name',
        copyable: true,
        ellipsis:true,
        fixed:'left',
        sorter: (a,b)=> {
            return a.name.localeCompare(b.name)
        },
    },
    {
        title: '管理地址',
        dataIndex: 'managerAddress',
        ellipsis:true,
        width:200,
        render:(_,entity)=>(<CopyAddrList keyName="managerAddress" addrItem={entity} />)
    },
    {
        title: '服务地址',
        dataIndex: 'serviceAddress',
        ellipsis:true,
        width:230,
        render:(_,entity)=>(<CopyAddrList keyName="serviceAddress" addrItem={entity} />)
    },
    {
        title: '集群同步地址',
        dataIndex: 'peerAddress',
        ellipsis:true,
        width:230,
        render:(_,entity)=>(<CopyAddrList keyName="peerAddress" addrItem={entity} />)
    },
    {
        title: '状态',
        dataIndex: 'status',
        ellipsis:true,
        width:86,
        valueType: 'select',
        filters: true,
        onFilter: true,
        valueEnum: new Map([
            [0, <span className="text-status_fail">异常</span>],
            [1,<span className="text-status_success">正常</span>],
        ])
    },
];

export const NODE_MODAL_COLUMNS:ColumnType<PartitionClusterNodeModalTableListItem>[] = [
    {title:'名称', dataIndex:'name',width:200,
    ellipsis:true,
    fixed:'left'},
    {title:'管理地址', dataIndex:'managerAddress',width:240,ellipsis:true,render:(_,entity)=>(<CopyAddrList keyName="managerAddress" addrItem={entity} />)},
    {title:'服务地址', dataIndex:'serviceAddress',width:240,ellipsis:true,render:(_,entity)=>(<CopyAddrList keyName="serviceAddress" addrItem={entity} />)},
    {title:'状态', dataIndex:'status',
    render:(text)=>(
        <span className={text === 0 ? 'text-status_fail' : 'text-status_success'}>{ClusterStatusEnum[text]}</span>
    )}
]

export const PARTITION_LIST_COLUMNS: ProColumns<PartitionTableListItem>[] = [
    {
        title: '环境名称',
        dataIndex: 'name',
        copyable: true,
        ellipsis:true,
        fixed:'left',
        sorter: (a,b)=> {
            return a.name.localeCompare(b.name)
        },
    },
    {
        title: 'ID',
        dataIndex: 'id',
        copyable: true,
        ellipsis:true,
        width:140,
    },
    // {
    //     title: '集群数量',
    //     dataIndex: 'clusterNum',
    //     sorter: (a,b)=> {
    //         return a.clusterNum - b.clusterNum
    //     },
    // },
    {
        title: '更新者',
        dataIndex: ['updater','name'],
        ellipsis: true,
        filters: true,
        onFilter: true,
        width:100,
        valueType: 'select',
        filterSearch: true
    },
    {
        title: '更新时间',
        dataIndex: 'updateTime',
        ellipsis:true,
        width:182,
        sorter: (a,b)=> {
            return a.updateTime.localeCompare(b.updateTime)
        },
    },
];

export enum ClusterStatusEnum {
    '异常',
    '正常'
}

export const DASHBOARD_SETTING_DRIVER_OPTION_LIST = [
    {
        label:'influxdb',
        value:'influxdb-v2'
    }
]