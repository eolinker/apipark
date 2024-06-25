/*
 * @Date: 2024-06-04 13:57:45
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-05 16:03:51
 * @FilePath: \frontend\packages\common\src\hooks\useInitializeMonaco.ts
 */
import { useEffect,useState } from 'react';
import { loader } from '@monaco-editor/react';
import { monaco } from '../monacoConfig';

const useInitializeMonaco = () => {
    const [initialized, setInitialized] = useState(false);

    useEffect(() => {
        if (!initialized) {
            loader.config({ monaco });
            loader.init().then(() => {
                setInitialized(true);
            });
        }
    }, [initialized]);
};
export default useInitializeMonaco;
