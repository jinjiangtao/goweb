# ERP 管理系统前端

基于 Vue 3 + Element Plus + Vite 构建的企业资源管理系统前端。

## 功能特性

- 用户登录与 JWT Token 认证
- 动态侧边栏菜单
- 用户管理（新增、编辑、删除、禁用）
- 角色管理（新增、编辑、删除、分配菜单权限）
- 菜单管理（树形结构、新增、编辑、删除）
- 产品管理（列表、搜索、新增、编辑、删除）
- 仪表盘数据统计
- 路由权限控制

## 技术栈

- Vue 3 (Composition API)
- Vite
- Element Plus
- Vue Router
- Pinia
- Axios

## 项目结构

```
frontend/
├── package.json
├── vite.config.js
├── index.html
└── src/
    ├── main.js          # 应用入口
    ├── App.vue          # 根组件
    ├── router/          # 路由配置
    │   └── index.js
    ├── store/           # Pinia 状态管理
    │   └── user.js
    ├── api/             # API 接口封装
    │   ├── auth.js
    │   ├── user.js
    │   ├── role.js
    │   ├── menu.js
    │   └── product.js
    ├── utils/           # 工具函数
    │   └── request.js   # Axios 封装
    ├── views/           # 页面组件
    │   ├── Login.vue
    │   ├── Dashboard.vue
    │   ├── UserManagement.vue
    │   ├── RoleManagement.vue
    │   ├── MenuManagement.vue
    │   └── ProductManagement.vue
    ├── layout/          # 布局组件
    │   └── index.vue
    ├── components/      # 通用组件
    └── assets/          # 静态资源
```

## 快速开始

### 安装依赖

```bash
cd frontend
npm install
```

### 启动开发服务器

```bash
npm run dev
```

开发服务器将在 http://localhost:3000 启动。

### 构建生产版本

```bash
npm run build
```

## 后端接口

后端 API 运行在 http://localhost:8080，前端通过代理转发请求。

## 说明

- 确保后端服务已启动
- 默认管理员账号可通过后端初始化创建
- Token 存储在 localStorage 中
