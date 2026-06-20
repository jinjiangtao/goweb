# 在线简历编辑系统

面向用户简历快速制作需求的可视化在线简历编辑系统，依托模板套用、拖拽排版、实时预览功能，实现零基础快速制作、导出专业简历。

## 功能特性

### 📝 核心功能
- **模板管理**：内置3套通用专业简历模板（简约商务、创意设计、技术工程师）
- **模块拖拽**：支持拖拽调整模块排序，自定义增删简历模块
- **内容编辑**：在线编辑文字内容，支持教育经历、工作经历、项目经验、专业技能等多种模块
- **样式调整**：自定义主色调、字体大小、行高、页面边距等样式
- **实时预览**：页面实时渲染预览制作效果，所见即所得
- **版本管理**：自动保存多版本简历，支持版本自由切换和恢复
- **PDF导出**：一键导出PDF格式简历文件

### 🎨 支持模块类型
- 基本信息（姓名、电话、邮箱、求职意向等）
- 教育背景（学校、学历、专业、时间、描述）
- 工作经历（公司、职位、时间、工作描述、工作亮点）
- 项目经验（项目名称、角色、时间、项目描述、项目亮点）
- 专业技能（技能分类、技能内容）
- 个人简介/自我评价

## 技术架构

### 后端技术栈
- **Web框架**：Go Gin v1.9.1
- **ORM**：GORM v1.31.1
- **数据库**：SQLite (modernc.org/sqlite 纯Go实现，无需CGO)
- **PDF生成**：gofpdf/v2 v2.17.3

### 前端技术栈
- **框架**：Vue 3.4 + Vite 5.0
- **状态管理**：Pinia 2.1
- **路由**：Vue Router 4.2
- **UI组件库**：Element Plus 2.4
- **拖拽功能**：vuedraggable 4.1
- **HTTP客户端**：Axios 1.6
- **工具库**：VueUse 10.7
- **样式**：Sass 1.69

### 项目结构
```
jianli/
├── server/                    # 后端 Go Gin 项目
│   ├── main.go               # 应用入口
│   ├── go.mod
│   ├── config/               # 配置管理
│   ├── controllers/          # 控制器
│   │   ├── template.go       # 模板管理接口
│   │   ├── resume.go         # 简历管理接口
│   │   ├── version.go        # 版本管理接口
│   │   └── pdf.go            # PDF导出接口
│   ├── models/               # 数据模型
│   ├── routes/               # 路由配置
│   ├── middleware/           # 中间件
│   ├── db/                   # 数据库层
│   ├── utils/                # 工具函数
│   ├── fonts/                # 中文字体目录
│   └── jianli.db             # SQLite数据库文件（运行时生成）
│
├── web/                      # 前端 Vue3 项目
│   ├── package.json
│   ├── vite.config.js
│   ├── index.html
│   └── src/
│       ├── main.js           # 应用入口
│       ├── App.vue           # 根组件
│       ├── api/              # API接口
│       ├── views/            # 页面组件
│       │   ├── Home.vue          # 首页/简历列表
│       │   ├── TemplateSelect.vue # 模板选择页
│       │   └── Editor.vue        # 简历编辑器
│       ├── components/       # 公共组件
│       │   ├── ModuleList.vue    # 模块列表（拖拽排序）
│       │   ├── ContentEditor.vue # 内容编辑器
│       │   ├── StylePanel.vue    # 样式调整面板
│       │   ├── ResumePreview.vue # 实时预览组件
│       │   └── VersionPanel.vue  # 版本历史面板
│       ├── stores/           # Pinia 状态管理
│       ├── router/           # Vue Router 配置
│       ├── utils/            # 工具函数
│       ├── styles/           # 全局样式
│       └── assets/           # 静态资源
│
├── start.bat               # Windows一键启动脚本
└── README.md
```

## API 接口文档

### 模板接口
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/templates | 获取所有模板列表 |
| GET | /api/templates/:id | 获取单个模板详情 |

### 简历接口
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/resumes?user_identity=xxx | 获取用户简历列表 |
| GET | /api/resumes/:id | 获取简历详情 |
| POST | /api/resumes | 创建新简历 |
| PUT | /api/resumes/:id | 更新简历（自动保存新版本） |
| DELETE | /api/resumes/:id | 删除简历 |

### 版本接口
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/resumes/:id/versions | 获取简历版本列表 |
| GET | /api/versions/:id | 获取版本详情 |
| POST | /api/versions/:id/restore | 恢复到指定版本 |

### PDF接口
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/pdf/export | 导出PDF文件 |

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- npm 9+

### 方式一：一键启动（Windows）
```bash
# 直接运行启动脚本
start.bat
```

### 方式二：手动启动

#### 1. 启动后端服务
```bash
cd server
go build -o jianli-server.exe .
jianli-server.exe
```
后端服务将在 `http://localhost:8080` 启动

#### 2. 启动前端服务
```bash
cd web
npm install
npm run dev
```
前端服务将在 `http://localhost:5173` 启动

### 访问应用
打开浏览器访问 `http://localhost:5173`

## 中文字体支持

PDF导出功能需要中文字体支持，请按以下步骤配置：

1. 下载 Noto Sans SC 字体：
   - 访问 https://fonts.google.com/noto/specimen/Noto+Sans+SC
   - 下载 `NotoSansSC-Regular.ttf` 和 `NotoSansSC-Bold.ttf`

2. 将字体文件放入 `server/fonts/` 目录

如果未配置中文字体，PDF导出功能将使用系统默认字体，中文可能显示为乱码。

## 数据库说明

系统使用 SQLite 作为数据库，数据库文件 `jianli.db` 将在首次运行时自动创建在 `server/` 目录下。

首次启动时会自动初始化3套简历模板数据。

## 使用流程

1. **首页**：查看已有简历列表，可编辑、删除、导出PDF
2. **创建简历**：点击"新建简历"选择模板，输入标题创建
3. **编辑简历**：
   - 左侧：拖拽调整模块顺序，显示/隐藏模块，添加新模块，调整样式
   - 中间：编辑选中模块的内容
   - 右侧：实时预览简历效果
4. **版本管理**：点击"版本历史"查看历史版本，可恢复到任意版本
5. **导出PDF**：编辑完成后点击"导出PDF"下载简历文件

## 开发说明

### 后端开发
```bash
cd server
go mod tidy
go run main.go
```

### 前端开发
```bash
cd web
npm install
npm run dev    # 开发模式
npm run build  # 生产构建
```

## 特色功能说明

### 🔄 自动保存
- 编辑器内容每3秒自动保存一次
- 每次保存自动创建新版本
- 支持自定义版本快照名称

### 🎯 拖拽排版
- 使用 vuedraggable 实现流畅的拖拽体验
- 拖拽时实时更新预览
- 支持模块显示/隐藏控制

### 📱 响应式设计
- 编辑器三栏布局，空间利用合理
- 预览区 A4 纸张比例，所见即所得
- 打印样式支持，可直接浏览器打印

### 🎨 主题定制
- 主色调自定义
- 强调色自定义
- 字体大小、行高、边距全方位可调

## License

MIT
