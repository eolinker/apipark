/*
 * @Date: 2024-02-02 14:37:42
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-05-29 18:14:33
 * @FilePath: \frontend\packages\core\src\utils\navigation.tsx
 */

import { MenuItem } from "@common/components/aoplatform/Navigation";

export function getItem(
    label: React.ReactNode,
    key: React.Key,
    icon?: React.ReactNode,
    children?: MenuItem[],
    type?: 'group',
    access?:string[] | string
  ): MenuItem {
    return {
      key,
      icon,
      children,
      label,
      type,
      access
    } as MenuItem;
  }

  export function getTabItem(
    label: React.ReactNode,
    key: React.Key,
    children?: MenuItem[],
    type?: 'group',
    access?:string
  ) {
    return {
      key,
      label,
      access
    } 
  }