
# 抽奖转盘后台管理系统

一个基于 Go + Gin + SQLite 后端和 Vue3 + Element Plus 前端的抽奖转盘后台管理系统。

## 项目结构

```
choujiang/
├── server/       # 后端 Go 项目
├── admin/        # 前端 Vue3 项目
└── README.md
```

## 后端 (Server)

### 技术栈
- Go 1.21+
- Gin Web Framework
- GORM (ORM)
- SQLite (纯 Go 驱动，无需 CGO)
- JWT 认证

### 初始化项目

1. 进入 server 目录：
```bash
cd server
```

2. 初始化 Go 模块并安装依赖：
```bash
go mod tidy
```

### 运行后端

```bash
go run main.go
```

后端将在 `http://localhost:8080` 启动。

## 前端 (Admin)

### 技术栈
- Vue 3
- Vue Router
- Element Plus
- Axios
- ECharts
- Vite

### 安装依赖

1. 进入 admin 目录：
```bash
cd admin
```

2. 安装 npm 依赖：
```bash
npm install
```

### 运行前端

```bash
npm run dev
```

前端将在 `http://localhost:3000` 启动。

## 默认账号

- 账号：admin
- 密码：123456

## 功能特性

### 1. 奖品管理
- 奖品增删改查
- 设置中奖概率
- 库存管理
- 启用/禁用奖品
- 概率总和检查（不超过100%）

### 2. 抽奖记录
- 查看所有抽奖记录
- 按姓名、手机号、奖品、中奖状态搜索
- 标记领取
- 导出 Excel

### 3. 统计看板
- 总抽奖次数
- 中奖次数和中奖率
- 待领取数量
- 奖品中奖饼图
- 近7天抽奖趋势图

### 4. 登录认证
- JWT 令牌认证
- 路由守卫

### 5. 抽奖核心逻辑
- 基于概率的抽奖算法
- 库存耗尽自动排除
- 互斥锁防止超卖

## API 接口

### 公共接口
- POST `/api/login` - 用户登录
- POST `/api/lottery` - 执行抽奖

### 需要认证的接口
- GET `/api/prizes` - 获取奖品列表
- POST `/api/prizes` - 添加奖品
- PUT `/api/prizes/:id` - 更新奖品
- DELETE `/api/prizes/:id` - 删除奖品
- PUT `/api/prizes/:id/toggle` - 切换奖品启用状态
- GET `/api/records` - 获取抽奖记录
- PUT `/api/records/:id/claim` - 标记领取
- GET `/api/records/export` - 导出记录
- GET `/api/stats` - 获取统计数据
