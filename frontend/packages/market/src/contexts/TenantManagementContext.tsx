/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-06 14:28:34
 * @FilePath: \frontend\packages\market\src\contexts\TenantManagementContext.tsx
 */
import  { createContext, useContext, useState, ReactNode, FC } from 'react';

interface TenantManagementContextProps {
    appName:string|undefined
    setAppName:React.Dispatch<React.SetStateAction<string|undefined>>;
}

const TenantManagementContext = createContext<TenantManagementContextProps | undefined>(undefined);

export const useTenantManagementContext = () => {
    const context = useContext(TenantManagementContext);
    if (!context) {
        throw new Error('useArray must be used within a ArrayProvider');
    }
    return context;
};

export const TenantManagementProvider: FC<{ children: ReactNode }> = ({ children }) => {
    const [appName, setAppName] = useState<string>()

    
    return <TenantManagementContext.Provider value={{ appName, setAppName }}>{children}</TenantManagementContext.Provider>;
};