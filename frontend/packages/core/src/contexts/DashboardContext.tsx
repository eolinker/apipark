/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-04-01 16:50:55
 * @FilePath: \frontend\packages\core\src\contexts\DashboardContext.tsx
 */
import  { createContext, useContext, useState, ReactNode, FC } from 'react';
import {EntityItem, PartitionItem} from "@common/const/type.ts";

interface DashboardContextProps {
    partitionList: PartitionItem[];
    setPartitionList: React.Dispatch<React.SetStateAction<PartitionItem[]>>;
    currentClusterList:EntityItem[]; 
    setCurrentClusterList: React.Dispatch<React.SetStateAction<EntityItem[]>>;
}

const DashboardContext = createContext<DashboardContextProps | undefined>(undefined);

export const useDashboardContext = () => {
    const context = useContext(DashboardContext);
    if (!context) {
        throw new Error('useArray must be used within a ArrayProvider');
    }
    return context;
};

export const DashboardProvider: FC<{ children: ReactNode }> = ({ children }) => {
    const [partitionList, setPartitionList] = useState<PartitionItem[]>([]);
    const [currentClusterList, setCurrentClusterList] = useState<EntityItem[]>([])
    return <DashboardContext.Provider value={{ partitionList, setPartitionList,currentClusterList, setCurrentClusterList }}>{children}</DashboardContext.Provider>;
};