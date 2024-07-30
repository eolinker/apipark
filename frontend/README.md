<!--
 * @Date: 2024-06-05 16:00:58
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-07-12 20:15:05
 * @FilePath: \frontend\README.md
-->
# 部署

## 安装依赖
    建议使用pnpm
    `npm install -g pnpm`
    使用pnpm安装依赖
    `pnpm install`
    
## 编译
    仅编译管理后台（打包目录为dist）：`pnpm run build`
    仅编译租户端（打包目录为tenant_dist）：`pnpm run build:tenant`
    同时编译管理后台和租户端:`pnpm run build:all`