/*
 * @Date: 2024-01-31 15:00:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-06-04 11:44:31
 * @FilePath: \frontend\packages\common\src\components\postcat\api\PreviewTable\index.tsx
 */
import { GridActionsCellItem } from '@mui/x-data-grid-pro'
import type { ReactNode } from 'react'
import type { SxProps, Theme } from '@mui/material'
import { IconButton } from '../IconButton'
import { commonTableSx } from '@common/const/api-detail'

interface PreviewGridActionsCellItemProps {
  icon: IconParkIconElement['name']
  label: string
  onClick?: () => void
}

export function PreviewGridActionsCellItem({ icon, label, onClick }: PreviewGridActionsCellItemProps): ReactNode {
  return (
    <GridActionsCellItem
      onClick={onClick}
      disableRipple
      className="table-actions"
      sx={{
        visibility: 'hidden'
      }}
      component="div"
      icon={<IconButton name={icon} />}
      label={label}
    />
  )
}

export function previewTableHoverSx(): SxProps<Theme> {
  return {...commonTableSx
  }
}

export function collapseTableSx(borderRadius: string | number): SxProps<Theme> {
  return {
    border: 'none',
    borderRadius: `0 0 ${borderRadius} ${borderRadius}`,
    overflow: 'hidden'
  }
}
