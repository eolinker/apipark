# 部署

## 代码同步
    packages目录下，部分子项目为企业版独有，不要同步到开源版：
    packages/businessEntry, packages/dashboard, packages/openApi, packages/systemRunning, README.pro.md

## 安装依赖
    建议使用pnpm
    `npm install -g pnpm`
    使用pnpm安装依赖
    `pnpm install`
    
## 编译
### 开源版本
        仅编译管理后台（打包目录为dist）：`pnpm run build`
        仅编译租户端（打包目录为tenant_dist）：`pnpm run build:tenant`
        同时编译管理后台和租户端:`pnpm run build:all`
### 企业版本
        仅编译管理后台（打包目录为dist）：`pnpm run build:pro`
        仅编译租户端（租户端暂时不区分企业版和开源版，打包目录为tenant_dist）：`pnpm run build:tenant`
        同时编译管理后台和租户端:`pnpm run build:pro:all`