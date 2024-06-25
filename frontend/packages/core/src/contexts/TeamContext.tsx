/*
 * @Date: 2024-02-06 16:44:50
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-04-22 17:30:40
 * @FilePath: \frontend\packages\core\src\contexts\TeamContext.tsx
 */
import  { FC, createContext, useContext, useState, ReactNode } from 'react';
import { TeamConfigType } from '../const/team/type';

interface TeamContextProps {
    teamInfo?:TeamConfigType
    setTeamInfo?: React.Dispatch<React.SetStateAction<TeamConfigType|undefined>>;
    // partitionList: PartitionItem[];
    // setPartitionList: React.Dispatch<React.SetStateAction<PartitionItem[]>>;
    // apiPrefix:string;
    // setApiPrefix:React.Dispatch<React.SetStateAction<string>>;
    // prefixForce:boolean;
    // setPrefixForce:React.Dispatch<React.SetStateAction<boolean>>;
}

const TeamContext = createContext<TeamContextProps | undefined>(undefined);

export const useTeamContext = () => {
    const context = useContext(TeamContext);
    if (!context) {
        throw new Error('useArray must be used within a ArrayProvider');
    }
    return context;
};

export const TeamProvider: FC<{ children: ReactNode }> = ({ children }) => {
    // const [partitionList, setPartitionList] = useState<PartitionItem[]>([]);
    // const [apiPrefix, setApiPrefix] = useState<string>('');
    // const [prefixForce, setPrefixForce] = useState<boolean>(false);
    const [teamInfo, setTeamInfo] = useState<TeamConfigType>()
    
    return <TeamContext.Provider value={{ teamInfo, setTeamInfo }}>{children}</TeamContext.Provider>;
};