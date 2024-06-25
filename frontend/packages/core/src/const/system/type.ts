/*
 * @Date: 2024-02-04 16:37:27
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 11:18:44
 * @FilePath: \frontend\packages\core\src\const\system\type.ts
 */
import { FormInstance, UploadFile } from "antd";
import { HeaderParamsType, BodyParamsType, QueryParamsType, RestParamsType, ResultListType, ApiBodyType } from "@common/const/api-detail";
import { EntityItem, MatchItem, PartitionItem } from "@common/const/type";
import { SubscribeEnum, SubscribeFromEnum } from "./const";
import { HTTPMethod, Protocol } from "@common/components/postcat/api/RequestMethod";

export type SystemTableListItem = {
    id:string;
    name: string;
    organization:EntityItem;
    team: EntityItem;
    apiNum: number;
    serviceNum: number,
    description:string;
    master:EntityItem;
    createTime:string;
};

export type SystemConfigFieldType = {
    name?: string;
    id?: string;
    prefix?:string;
    description?: string;
    team?:string;
    master?:string;
    partition?:string[]
};

export type SystemSubServiceTableListItem = {
    id:string;
    applyStatus:SubscribeEnum;
    partition:EntityItem[];
    project:EntityItem;
    team:EntityItem
    service:EntityItem
    applier:EntityItem
    from:SubscribeFromEnum
    createTime:string
};


export type SimpleMemberItem = {
    id:string
    name:string
    email:string
    department:string
    avatar:string
}

export type SystemSubscriberTableListItem = {
    id:string
    service:EntityItem
    partition:EntityItem[];
    applyStatus:SubscribeEnum
    project:EntityItem
    team:EntityItem;
    applier:EntityItem
    approver:EntityItem;
    from:SubscribeFromEnum
    applyTime:string
};

export type SystemSubscriberConfigFieldType = {
    service:string
    subscriber:string
    applier:string
    partition:string
};

export type SystemSubscriberConfigProps = {
    systemId:string
    partitionList:PartitionItem[]
}

export type SystemSubscriberConfigHandle = {
    save:()=>Promise<boolean|string>
}

export type SystemMemberTableListItem = {
    user: EntityItem;
    email:string;
    roles:Array<EntityItem>;
    canDelete:boolean
};

export type SystemApiDetail = {
    id:string
    name:string
    description:string
    protocol:Protocol
    method:HTTPMethod
    path:string
    creator:EntityItem
    createTime:string
    updater:EntityItem
    updateTime:string
    match?:MatchItem[]
    proxy?:SystemApiProxyType
    doc?:{
        encoding: string,
        tag: string,
        requestParams: {
            headerParams: HeaderParamsType[],
            bodyParams: BodyParamsType[],
            queryParams: QueryParamsType[],
            restParams: RestParamsType[]
        },
        resultList: ResultListType[],
        responseList: [{
            id: number,
            responseUuid: string,
            apiUuid: string,
            oldId: number,
            name: string,
            httpCode: string,
            contentType: ApiBodyType,
            isDefault: number,
            updateUserId: number,
            createUserId: number,
            createTime: number,
            updateTime: number,
            responseParams: {
                headerParams: HeaderParamsType[],
                bodyParams: BodyParamsType[]
                queryParams: QueryParamsType[],
                restParams: RestParamsType[]
            }
        }]
    }
}


export type SystemApiProxyType = {
    path:string
    timeout:number
    retry:number
    headers:Array<ProxyHeaderItem>

}
export type SystemApiProxyFieldType = {
    name: string;
    id:string;
    description?:string;
    path:string;
    method:string;
    match:MatchItem[]
    isDisable?: boolean;
    service?:string;
    proxy:SystemApiProxyType
};

export type SystemInsideApiCreateProps = {
    type?:'copy'
    entity?:SystemApiProxyFieldType &{systemId:string}
    modalApiPrefix?:string
    modalPrefixForce?:boolean
}

export type SystemInsideApiCreateHandle = {
    copy:()=>Promise<boolean|string>;
    save:()=>Promise<boolean|string>;
}


export type SystemApiTableListItem = {
    id:string;
    name: string;
    method:string;
    requestPath:string;
    creator:EntityItem;
    createTime:string;
    updater:EntityItem
    updateTime:string
    canDelete:boolean
};


export type EditAuthFieldType  = {
    id?:string
    name: string
    driver: string
    hideCredential: boolean
    expireTime: number
    position: string
    tokenName: string
    config: {
        userName?: string
        password?: string
        apikey?: string
        ak?: string
        sk?: string
        iss?: string
        algorithm?: string
        secret?: string
        publicKey?: string
        user?: string
        userPath?: string
        claimsToVerify?: string[]
        signatureIsBase64?: boolean
    }
}


export type SystemAuthorityConfigProps = {
    type:'add'|'edit'
    data?:EditAuthFieldType
    systemId:string
}

export type SystemAuthorityConfigHandle = {
    save:()=>Promise<boolean|string>
}

export type SystemAuthorityViewProps = {
    entity:Array<{key:string, value:string}>
}

export type SystemUpstreamTableListItem = {
    name: string;
    id:string;
    driver:string;
    creator:EntityItem;
    updater:EntityItem;
    createTime:string;
    updateTime:string;
    canDelete:boolean
};

export type ProxyHeaderItem = {
    key:string
    value:string
    optType:string
    id?:string
}

export type GlobalNodeItem = {
    address:string
    weight:number
}

export type NodeItem = Partial<GlobalNodeItem> & {
    cluster:string
    clusterName?:string
    _id?:string }

export type DiscoverItem = {
    cluster:string
    service:string
    discover:string
}

export type ServiceUpstreamFieldType = {
    _apinto_show?:boolean
    driver:string
    nodes:GlobalNodeItem[],
    discover?:DiscoverItem
    timeout:number;
    retry?:number;
    limitPeerSecond?:number;
    scheme:string,
    passHost:string,
    upstreamHost:string,
    balance:string;
    proxyHeaders:ProxyHeaderItem[]
};


export type MyServiceFieldType = {
    name?: string;
    id?: string;
    description?: string;
    logo?:string;
    logoFile?:UploadFile;
    tags?:Array<string>;
    serviceType?:'public'|'inner';
    team?:string;
    project?:string;
    partition?:Array<string>;
    group?:string | string[];
    status?:'off'|'on'
};

export type SimpleSystemItem = {
    id:string
    name:string
    team:EntityItem
    organization:EntityItem
    partition:EntityItem[]
}

export type ServiceApiTableListItem = {
    id:string;
    name: string;
    method:string;
    path:string;
    description:string;
};

export type SimpleApiItem = {
    id:string
    name:string
    method:string
    requestPath:string
}

export type SystemAuthorityTableListItem = {
    id:string
    name: string;
    driver:string;
    hideCredential:boolean;
    expireTime:number;
    creator:EntityItem;
    updater:EntityItem;
    createTime:string;
    updateTime:string
};

export type MyServiceTableListItem = {
    id:string;
    name: string;
    partition:EntityItem[];
    partitionId:string;
    serviceType:'public'|'inner';
    apiNum:number;
    status:string;
    createTime:string;
    updateTime:string;
};


export type SystemInsideApiDetailProps = {
    systemId:string;
    apiId:string;
}


export type SystemInsideApiDocumentHandle  = {
    save:()=>Promise<boolean|string>|undefined
}

export type SystemInsideApiDocumentProps = {
    systemId:string
    apiId:string
}


export type SystemInsideApiProxyProps = {
    className?:string
    systemId:string
    initProxyValue?:SystemApiProxyType
    value?:SystemApiProxyType
    onChange?: (newConfigItems: SystemApiProxyType) => void; // 当配置项变化时，外部传入的回调函数
}

export type SystemInsideApiProxyHandle = {
    validate:()=>Promise<void>
}


export interface MyServiceInsideConfigHandle {
    save:()=>Promise<boolean|string>
}

export interface MyServiceInsideConfigProps {
    systemId:string,
    teamId:string
    serviceId?:string
    closeDrawer?:() => void
}


export type SubSubscribeApprovalModalProps = {
    type:'reApply'|'view'
    data?:SystemSubServiceTableListItem
    systemId?:string
}

export type SubSubscribeApprovalModalHandle = {
    reApply:() =>Promise<boolean|string>
}

export type SubSubscribeApprovalModalFieldType = {
    partitions?:string[]
    reason?:string;
    opinion?:string;
};

export type SystemInsideUpstreamConfigProps = {
    upstreamNameForm:FormInstance
    partitionId:string
    setLoading:(loading:boolean) => void
}

export type SystemInsideUpstreamConfigHandle = {
    save:()=>Promise<boolean|string>|undefined
}

export type SystemInsideUpstreamContentHandle = {
    save:()=>Promise<boolean|string>|undefined
}


export type SystemConfigHandle = {
    save:()=>Promise<string|boolean>|undefined
}


export type SystemTopologyServiceItem = EntityItem & {
    project:string
}

export interface SystemTopologySubscriber {
    project: EntityItem;
    services: EntityItem[];
  }
  
  export interface SystemTopologyInvoke {
    project: EntityItem;
    services: EntityItem[];
  }
  
  
  // 接口返回的数据格式
  export interface SystemTopologyResponse {
    services: SystemTopologyServiceItem[];
    subscribers: SystemTopologySubscriber[];
    invoke: SystemTopologyInvoke[];
  }

export enum SystemReleaseStatus {
    '正常' = 0,
    '未设置' = 1,
    '缺失' = 2
}

  export type SystemPublishReleaseItem = {
    api: Array<{
        name: string,
        method: string,
        path: string,
        upstream: string,
        change: string,
        status: {
            upstreamStatus: SystemReleaseStatus,
            docStatus: SystemReleaseStatus,
            proxyStatus: SystemReleaseStatus
        }
    }>
    upstream: Array<{
        name: "",
        partition: EntityItem,
        cluster: EntityItem,
        type: "",
        addr: [],
        status: ""
    }>
  }