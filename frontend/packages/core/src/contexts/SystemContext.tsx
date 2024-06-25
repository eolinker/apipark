/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 11:18:52
 * @FilePath: \frontend\packages\core\src\contexts\SystemContext.tsx
 */
import  {FC, createContext, useContext, useState, ReactNode, useEffect } from 'react';
import {PartitionItem} from "@common/const/type.ts";
import { SystemConfigFieldType } from '../const/system/type.ts';

interface SystemContextProps {
    partitionList: PartitionItem[];
    setPartitionList: React.Dispatch<React.SetStateAction<PartitionItem[]>>;
    apiPrefix:string;
    setApiPrefix:React.Dispatch<React.SetStateAction<string>>;
    prefixForce:boolean;
    setPrefixForce:React.Dispatch<React.SetStateAction<boolean>>;
    systemInfo:SystemConfigFieldType|undefined
    setSystemInfo:React.Dispatch<React.SetStateAction<SystemConfigFieldType|undefined>>;
}

const SystemContext = createContext<SystemContextProps | undefined>(undefined);

export const useSystemContext = () => {
    const context = useContext(SystemContext);
    if (!context) {
        throw new Error('useArray must be used within a ArrayProvider');
    }
    return context;
};

export const SystemProvider: FC<{ children: ReactNode }> = ({ children }) => {
    const [partitionList, setPartitionList] = useState<PartitionItem[]>([]);
    const [apiPrefix, setApiPrefix] = useState<string>('');
    const [prefixForce, setPrefixForce] = useState<boolean>(false);
    const [systemInfo, setSystemInfo] = useState<SystemConfigFieldType>()

    return <SystemContext.Provider value={{ partitionList, setPartitionList,apiPrefix,setApiPrefix,prefixForce,setPrefixForce,systemInfo, setSystemInfo }}>{children}</SystemContext.Provider>;
};