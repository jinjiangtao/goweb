# 设备运维工单管理系统

适配企业、校园设备运维场景的全流程工单系统，实现设备故障运维数字化管控。

## 功能特性

- 工单提交：用户在线提交报修工单，上传故障图片、描述问题，选择设备类型和报修区域
- 工单派单：管理员支持手动派单和智能自动派单（基于区域匹配和负载均衡算法）
- 维修处理：运维人员更新维修进度、填写维修记录、上传维修凭证
- 状态管控：待处理 → 已派单 → 维修中 → 已完成 / 已驳回 / 已取消，完整的状态流转
- 进度追踪：工单进度可视化展示，维修记录时间线
- 逾期预警：自动检测逾期工单，红色标记提醒
- 工单评价：完成后用户可评分和评价
- 数据统计：工单处理效率、故障率、区域分布、运维人员绩效等多维度统计分析
- 设备管理：设备台账管理
- 用户管理：多角色权限（管理员、运维人员、普通用户）

## 技术栈

- 后端：Go + Gin + GORM + SQLite
- 前端：Vue3 + Vite + Pinia + Vue Router + Element Plus + ECharts
- 认证：JWT Token

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+

### 一键启动（推荐）

**Windows:**
```bat
start.bat
```

**Linux/Mac:**
```bash
chmod +x start.sh
./start.sh
```

启动后访问 http://localhost:8080

### 开发模式

**Windows:**
```bat
dev.bat
```

**手动开发模式：**

```bash
# 后端 (端口 8080)
cd backend
go mod tidy
go run main.go

# 前端 (端口 5173)
cd frontend
npm install
npm run dev
```

### 默认账号

| 角色 | 用户名 | 密码 |
|------|--------|------|
| 管理员 | admin | 123456 |
| 普通用户 | user1 | 123456 |
| 运维人员 | tech1 | 123456 |
| 运维人员 | tech2 | 123456 |

## 项目结构

```
weixiu/
├── backend/              # Go 后端
│   ├── config/          # 配置（数据库等）
│   ├── handlers/        # API 处理器
│   │   ├── auth.go      # 用户认证
│   │   ├── device.go    # 设备管理
│   │   ├── order.go     # 工单 CRUD
│   │   ├── order_flow.go# 工单流转、派单、处理
│   │   ├── stats.go     # 数据统计
│   │   └── log.go       # 操作日志
│   ├── middleware/      # 中间件（JWT认证、权限）
│   ├── models/          # 数据库模型
│   ├── data/            # SQLite 数据库文件
│   ├── uploads/         # 上传图片存储
│   └── main.go          # 入口文件
├── frontend/            # Vue3 前端
│   ├── src/
│   │   ├── api/         # API 请求
│   │   ├── router/      # 路由
│   │   ├── stores/      # Pinia 状态管理
│   │   ├── utils/       # 工具函数
│   │   └── views/       # 页面组件
│   └── package.json
├── start.bat            # Windows 启动脚本
├── start.sh             # Linux/Mac 启动脚本
└── dev.bat              # Windows 开发模式脚本
```

## API 接口

### 认证
- `POST /api/auth/login` 用户登录
- `POST /api/auth/register` 用户注册
- `GET /api/auth/me` 获取当前用户信息

### 工单
- `POST /api/orders` 创建工单
- `GET /api/orders` 获取工单列表（支持分页、筛选）
- `GET /api/orders/:id` 获取工单详情（含维修日志、评价）
- `POST /api/orders/:id/assign` 手动派单
- `POST /api/orders/:id/auto-assign` 智能自动派单
- `POST /api/orders/:id/start` 开始处理
- `POST /api/orders/:id/process` 提交维修进度
- `POST /api/orders/:id/complete` 完成维修
- `POST /api/orders/:id/reject` 驳回工单
- `POST /api/orders/:id/cancel` 取消工单
- `POST /api/orders/:id/evaluate` 评价工单

### 设备
- `POST /api/devices` 创建设备
- `GET /api/devices` 获取设备列表
- `GET /api/devices/:id` 获取设备详情
- `PUT /api/devices/:id` 更新设备
- `DELETE /api/devices/:id` 删除设备

### 统计
- `GET /api/stats` 获取运维统计数据

### 其他
- `POST /api/upload` 文件上传
- `GET /api/users` 用户列表（管理员）
- `GET /api/technicians` 运维人员列表
