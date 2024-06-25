/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 18:15:06
 * @FilePath: \frontend\packages\common\src\components\aoplatform\formily2-customize\CustomCodeboxComponent.tsx
 */
import {forwardRef, useImperativeHandle, useState} from 'react'
import { Codebox } from '@common/components/postcat/api/Codebox'

export const CustomCodeboxComponent = forwardRef(
  (props: { [k: string]: unknown }, ref) => {
    const {
      mode = 'yaml',
      theme = 'xcode',
      fontSize,
      height,
      width = '100%',
      onChange,
      value
    } = props
    const [code, setCode] = useState(
      mode === 'json' ? JSON.stringify(value) : value
    )
    useImperativeHandle(ref, () => ({}))
    const handleChange = (value: string) => {
      setCode(value)
      let res = value
      if (mode === 'json') {
        try {
          res = JSON.parse(value)
        } catch {
          console.warn(' 输入的json语句格式有误')
        }
      }
      onChange(res)
    }

    return (
      <div className=" mt-[4px] border-[1px] border-solid border-BORDER">
        <Codebox 
          value={code} 
          language={mode}
          enableToolbar={false}
          theme={theme}
          fontSize={fontSize}
          height={height ?? 500}
          width={width}
           onChange={handleChange} />
      </div>
    )
  }
)
