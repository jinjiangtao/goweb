# 内推招聘管理系统

一个完整的员工内推招聘管理系统，包含后端 API 和管理后台。

## 技术栈

### 后端
- Go 1.21+
- Gin Web Framework
- SQLite 数据库
- JWT 认证

### 前端
- Vue 3
- Element Plus
- Pinia 状态管理
- Vue Router
- Axios
- ECharts 图表

## 项目结构

```
neitui/
├── server/                 # 后端服务
│   ├── main.go            # 主入口
│   ├── models/            # 数据模型
│   ├── handlers/          # 接口处理
│   ├── middleware/        # 中间件
│   ├── routes/            # 路由配置
│   ├── utils/             # 工具函数
│   └── uploads/resumes/   # 简历存储
├── admin/                 # 前端管理后台
│   ├── src/
│   │   ├── views/         # 页面组件
│   │   ├── stores/        # Pinia 状态
│   │   ├── router/        # 路由配置
│   │   └── api/           # API 封装
│   └── package.json
└── README.md
```

## 快速开始

### 1. 启动后端服务

```bash
cd server
go mod tidy
go run main.go
```

后端服务将在 http://localhost:8080 启动

### 2. 启动前端服务

```bash
cd admin
npm install
npm run dev
```

前端服务将在 http://localhost:5173 启动

## 默认账号

系统启动时会自动创建以下账号：

- **管理员**: admin / 123456
- **HR**: hr / 123456

## 功能说明

### 员工角色
- 发布职位
- 查看自己发布的职位
- 推荐候选人（支持上传简历）
- 查看自己的内推记录和统计

### HR 角色
- 查看所有内推记录
- 更新候选人状态
- 查看所有职位
- 查看统计看板（内推数、入职数、排行榜、趋势图）

### 管理员角色
- 拥有 HR 的所有权限
- 用户管理（添加、启用/禁用、重置密码）

## 数据库

数据存储在 `server/neitui.db`（SQLite 数据库），启动时自动创建。
