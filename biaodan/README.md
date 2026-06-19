# 低代码在线表单设计平台

基于低代码可视化开发理念的在线表单设计平台，支持拖拽式表单搭建、自定义属性配置与数据收集，无需代码即可快速制作各类业务表单，适配多场景数据采集需求。

## 功能特性

### 表单设计器
- **丰富组件库**：内置单行文本、多行文本、数字、邮箱、手机号、单选、多选、下拉选择、日期、时间、文件上传、评分、开关、滑块、分割线、说明文字等 16 种常用组件
- **拖拽式布局**：组件可自由拖拽排序，支持实时调整位置
- **属性配置**：可视化配置组件标签、占位符、默认值、必填校验、选项等属性
- **实时预览**：一键预览表单实际效果，支持模拟填写
- **模板保存**：表单设计可保存复用，随时编辑修改

### 表单填写
- **响应式页面**：适配PC、移动端，美观的渐变风格界面
- **表单校验**：前端 + 后端双重校验（必填、邮箱、手机号格式等）
- **文件上传**：支持单/多文件上传，大小限制可配置

### 数据管理后台
- **数据汇总**：展示总提交数、今日新增、近7天提交趋势图表
- **数据检索**：关键词全文搜索，支持分页浏览
- **数据导出**：一键导出 CSV 格式（Excel可直接打开）
- **数据管理**：支持单条数据删除

## 技术架构

| 层级 | 技术栈 | 目录 |
|------|--------|------|
| 后端 | Go 1.21 + Gin + GORM + SQLite | `server/` |
| 前端 | Vue 3 + Vite 5 + Element Plus + Pinia | `web/` |

### 后端结构
```
server/
├── main.go           # 入口文件，路由配置
├── go.mod            # Go模块依赖
├── database/
│   └── database.go   # SQLite数据库初始化
├── models/
│   └── models.go     # 数据模型（表单模板、提交数据）
└── handlers/
    ├── handlers.go   # 核心API处理
    └── upload.go     # 文件上传处理
```

### 前端结构
```
web/
├── index.html
├── vite.config.js
├── package.json
└── src/
    ├── main.js           # 入口文件
    ├── App.vue
    ├── router/index.js   # 路由配置
    ├── api/index.js      # API接口封装
    ├── stores/form.js    # Pinia状态管理
    ├── utils/request.js  # Axios封装
    ├── styles/global.css # 全局样式
    ├── components/
    │   ├── FormRenderer.vue     # 表单渲染器
    │   ├── PropertyConfig.vue   # 属性配置面板
    │   └── DraggableContainer.vue # 拖拽容器
    └── views/
        ├── Home.vue         # 模板列表首页
        ├── Designer.vue     # 表单设计器
        ├── FillForm.vue     # 表单填写页
        └── DataManage.vue   # 数据管理后台
```

## 快速开始

### 环境要求
- Go >= 1.21
- Node.js >= 16
- 包管理器：npm / yarn / pnpm

### 一、启动后端服务

```bash
cd server

# 1. 安装依赖（首次运行）
go mod tidy

# 2. 编译运行
go run main.go

# 或编译成可执行文件：
go build -o biaodan.exe .
./biaodan.exe
```

服务启动后监听 `http://localhost:8080`

首次运行会自动创建：
- `data/biaodan.db` - SQLite 数据库文件
- `uploads/` - 文件上传存储目录

### 二、启动前端（开发模式）

```bash
cd web

# 1. 安装依赖（首次运行）
npm install

# 2. 启动开发服务器
npm run dev
```

开发服务器启动后访问 `http://localhost:5173`，Vite已配置代理到后端8080端口。

### 三、生产部署

前端已构建，直接运行后端即可访问：

```bash
# 1. 构建前端（已经构建过可跳过）
cd web && npm run build && cd ..

# 2. 运行后端（会自动托管前端dist目录）
cd server && go run main.go
```

浏览器打开 `http://localhost:8080` 即可使用完整功能。

## API接口列表

| 方法 | 路径 | 说明 |
|------|------|------|
| GET    | `/api/templates`               | 获取模板列表，支持keyword搜索 |
| GET    | `/api/templates/:id`           | 获取单个模板详情 |
| POST   | `/api/templates`               | 创建新模板 |
| PUT    | `/api/templates/:id`           | 更新模板 |
| DELETE | `/api/templates/:id`           | 删除模板（级联删除数据） |
| POST   | `/api/submissions`             | 提交表单数据（含服务端校验） |
| GET    | `/api/submissions`             | 分页查询提交数据 |
| GET    | `/api/submissions/export`      | 导出CSV |
| DELETE | `/api/submissions/:id`         | 删除单条数据 |
| POST   | `/api/upload`                  | 文件上传，multipart/form-data |
| GET    | `/api/stats/:templateId`       | 获取模板统计数据 |

## 使用场景

1. **活动报名**：姓名、手机、人数、备注等字段快速搭建报名表单
2. **问卷调查**：单选、多选、评分、文本组合实现调研问卷
3. **信息登记**：访客/员工/学生信息登记，支持文件上传附件
4. **意见反馈**：多类型组件组合收集用户反馈

## 数据安全说明

- 表单提交数据存储于本地 SQLite 数据库 `server/data/biaodan.db`
- 上传文件存放于 `server/uploads/` 目录，通过 `/uploads/*` 路径访问
- 建议生产环境定期备份数据库文件与上传目录
