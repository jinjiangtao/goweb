# 在线拖拽式流程审批系统

针对企业轻量化办公审批场景，搭建可视化流程配置平台，支持自定义审批流程、单据提交、节点流转、状态跟进，实现全流程线上自动化审批。

## 功能特性

### 流程设计器
- ✅ 可视化拖拽画布，支持拖拽节点、连线搭建审批流程
- ✅ 多种节点类型：开始节点、审批节点、条件节点、结束节点
- ✅ 节点属性配置：审批人员、流转条件、节点名称
- ✅ 表单字段自定义：支持文本、多行文本、数字、下拉框、日期、时间、开关、单选、多选等类型
- ✅ 流程模板保存、复用与编辑

### 审批管理
- ✅ 请假、报销、通用申请单据线上提交
- ✅ 审批进度可视化，节点状态实时更新
- ✅ 审批操作：同意、驳回、撤回、加签、转交
- ✅ 审批记录时间线，完整记录审批历史

### 系统管理
- ✅ 用户管理，支持多用户权限
- ✅ JWT 认证，接口权限校验
- ✅ 首页数据统计看板

## 技术栈

### 后端
- **Go 1.20+** - 编程语言
- **Gin v1.12.0** - Web 框架
- **GORM v1.31.1** - ORM 框架
- **SQLite** - 数据库
- **golang-jwt/v5** - JWT 认证
- **CORS** - 跨域处理

### 前端
- **Vue 3** - 前端框架
- **Vite** - 构建工具
- **Element Plus** - UI 组件库
- **Vue Router** - 路由管理
- **Pinia** - 状态管理
- **Axios** - HTTP 客户端
- **Sass** - CSS 预处理器

## 目录结构

```
shenpi/
├── server/                 # 后端代码
│   ├── config/             # 配置文件
│   ├── controllers/        # API 控制器
│   ├── middleware/         # 中间件
│   ├── models/             # 数据模型
│   ├── routes/             # 路由注册
│   ├── services/           # 业务逻辑
│   ├── utils/              # 工具函数
│   ├── main.go             # 程序入口
│   ├── go.mod              # 依赖管理
│   └── config.yaml         # 配置文件
├── web/                    # 前端代码
│   ├── src/
│   │   ├── api/            # API 接口封装
│   │   ├── components/     # 公共组件
│   │   ├── layout/         # 布局组件
│   │   ├── router/         # 路由配置
│   │   ├── store/          # 状态管理
│   │   ├── utils/          # 工具函数
│   │   ├── views/          # 页面视图
│   │   ├── App.vue         # 根组件
│   │   └── main.js         # 入口文件
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
├── start-all.bat           # 一键启动脚本
├── start-server.bat        # 后端启动脚本
├── start-web.bat           # 前端启动脚本
└── README.md               # 说明文档
```

## 环境要求

- **Go** >= 1.20
- **Node.js** >= 16.0
- **npm** >= 8.0

## 快速开始

### 方式一：一键启动（推荐）

双击运行 `start-all.bat`，脚本会自动启动后端和前端服务。

### 方式二：分别启动

1. **启动后端服务**
   ```bash
   双击 start-server.bat
   # 或手动执行
   cd server
   go mod download
   go run main.go
   ```
   后端服务启动在 http://localhost:8080

2. **启动前端服务**
   ```bash
   双击 start-web.bat
   # 或手动执行
   cd web
   npm install
   npm run dev
   ```
   前端服务启动在 http://localhost:5173

### 访问系统

在浏览器中打开 http://localhost:5173

## 默认测试账号

系统初始化时会自动创建以下测试用户：

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | 123456 | 管理员 |
| zhangsan | 123456 | 普通用户 |
| lisi | 123456 | 普通用户 |
| wangwu | 123456 | 普通用户 |
| zhaoliu | 123456 | 普通用户 |

## API 接口说明

### 认证接口（公开）
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

### 用户接口（需认证）
- `GET /api/user/profile` - 获取当前用户信息
- `GET /api/user/list` - 获取用户列表
- `GET /api/user/:id` - 获取指定用户信息

### 模板接口（需认证）
- `POST /api/template` - 创建流程模板
- `POST /api/template/default` - 创建默认模板
- `GET /api/template` - 获取模板列表
- `GET /api/template/:id` - 获取模板详情
- `PUT /api/template/:id` - 更新模板
- `DELETE /api/template/:id` - 删除模板
- `GET /api/template/:id/fields` - 获取模板表单字段

### 审批接口（需认证）
- `POST /api/approval` - 提交审批申请
- `GET /api/approval/my` - 获取我的申请列表
- `GET /api/approval/todo` - 获取待我审批列表
- `GET /api/approval/stats` - 获取统计数据
- `GET /api/approval/:id` - 获取审批详情
- `GET /api/approval/:id/records` - 获取审批记录
- `GET /api/approval/:id/progress` - 获取审批进度
- `POST /api/approval/:id/approve` - 同意审批
- `POST /api/approval/:id/reject` - 驳回审批
- `POST /api/approval/:id/revoke` - 撤回申请
- `POST /api/approval/:id/addsign` - 加签
- `POST /api/approval/:id/transfer` - 转交审批

## 使用说明

### 1. 创建流程模板

1. 登录系统后，点击左侧菜单"流程模板"
2. 点击"新建模板"或"创建默认模板"（请假/报销/通用申请）
3. 在流程设计器中：
   - 从左侧节点面板拖拽节点到画布
   - 点击节点右侧圆点拖拽到目标节点创建连线
   - 选中节点，在右侧面板配置节点属性（审批人、条件等）
   - 在底部配置表单字段
4. 点击"保存模板"

### 2. 提交审批申请

1. 点击左侧菜单"我的申请"
2. 点击"提交申请"
3. 选择流程模板
4. 填写表单内容
5. 点击"提交"

### 3. 处理审批

1. 点击左侧菜单"待我审批"
2. 点击"审批"查看详情
3. 查看审批进度和表单内容
4. 执行审批操作：同意、驳回、加签、转交

### 4. 查看审批进度

1. 在"我的申请"或"待我审批"列表中点击"详情"
2. 查看流程图（节点颜色表示状态）
   - 灰色：未开始
   - 蓝色：进行中
   - 绿色：已通过
   - 红色：已驳回
3. 查看审批记录时间线

## 核心功能演示

### 流程引擎特性
- 支持串行审批、并行审批
- 支持条件分支（基于表单字段值自动判断走向）
- 支持多人会签（需所有审批人同意）
- 支持或签（任一审批人同意即可）
- 支持加签（增加审批人员）
- 支持转交（转由他人审批）
- 支持撤回（申请人撤回申请）

### 条件表达式
支持以下运算符：
- `==` 等于
- `!=` 不等于
- `>` 大于
- `>=` 大于等于
- `<` 小于
- `<=` 小于等于
- `contains` 包含

示例：
- `amount > 1000` - 金额大于1000
- `type == 'leave'` - 类型等于请假
- `days <= 3` - 天数小于等于3
- `dept contains '技术'` - 部门包含"技术"

## 数据库

系统使用 SQLite 数据库，数据库文件位于 `server/data/shenpi.db`，首次运行时会自动创建并初始化测试数据。

## 注意事项

1. 启动前请确保已安装 Go 和 Node.js 环境
2. 确保 8080 和 5173 端口未被占用
3. 后端服务必须先启动，前端才能正常调用 API
4. 生产环境请修改 `server/config/config.yaml` 中的 JWT 密钥

## 开发说明

### 后端开发
```bash
cd server
go mod tidy    # 整理依赖
go build       # 编译
go run main.go # 运行
```

### 前端开发
```bash
cd web
npm install    # 安装依赖
npm run dev    # 开发模式
npm run build  # 生产构建
```

## License

MIT
