
# 会议室预订系统 - 管理后台和后端

## 项目结构

```
huiyishi/
├── server/          # Go 后端服务
│   ├── main.go
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   └── utils/
├── admin/           # Vue3 管理后台
│   ├── src/
│   │   ├── views/
│   │   ├── stores/
│   │   ├── router/
│   │   └── main.js
│   └── package.json
└── prd-01.md
```

## 功能特性

### 后端 (Go + Gin + SQLite)
- JWT 认证登录
- 会议室 CRUD 管理
- 预订管理（取消预订等）
- 统计数据接口

### 前端 (Vue3 + Element Plus)
- 深色科技感主题
- 登录页面
- 会议室管理
- 预订管理
- 统计看板（ECharts 图表）

## 如何运行

### 1. 后端运行

**注意**：在 Windows 上运行需要 GCC 编译器（安装 MinGW-w64 或 TDM-GCC），或者我们可以使用纯 Go 的 SQLite 驱动。

#### 方法 1：使用 CGO（需要 GCC）

```bash
cd server
set CGO_ENABLED=1
go mod tidy
go run main.go
```

#### 方法 2：使用纯 Go SQLite 驱动

更新 `database/database.go` 导入：
```go
import (
    _ "modernc.org/sqlite"
    "gorm.io/driver/sqlite"
    // ...
)
```

然后运行：
```bash
cd server
go mod tidy
go run main.go
```

后端服务默认运行在 `http://localhost:8080`

### 2. 前端运行

```bash
cd admin
npm install
npm run dev
```

前端默认运行在 `http://localhost:3000`

## 默认账号

- 用户名：`admin`
- 密码：`123456`

## 技术栈

### 后端
- Go 1.x
- Gin Web Framework
- GORM ORM
- SQLite 数据库
- JWT 认证
- CORS 跨域支持

### 前端
- Vue 3
- Element Plus UI 组件库
- Pinia 状态管理
- Vue Router 路由
- Axios HTTP 客户端
- ECharts 图表库
