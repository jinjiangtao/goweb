# 在线协同白板系统

## 项目简介

基于 Go Gin + WebSocket + Vue3 + Canvas 开发的多人实时协同在线白板系统，支持跨用户实时绘图、标注、互动协作，适用于线上会议、远程协作、教学演示等场景。

## 技术栈

### 后端
- **Go 1.21+** - 编程语言
- **Gin** - Web 框架
- **Gorilla WebSocket** - WebSocket 通信
- **GORM + SQLite** - 数据库 ORM 和存储
- **godotenv** - 环境变量管理

### 前端
- **Vue 3** - 前端框架
- **Vite** - 构建工具
- **Pinia** - 状态管理
- **Vue Router** - 路由管理
- **Element Plus** - UI 组件库
- **Canvas API** - 绘图引擎

## 功能特性

### 绘图工具
- ✏️ **画笔** - 自由绘制，支持贝塞尔曲线平滑
- 📏 **直线** - 绘制直线
- 🔲 **矩形** - 绘制矩形
- ⭕ **圆形** - 绘制圆形
- 🔵 **椭圆** - 绘制椭圆
- 🔼 **三角形** - 绘制三角形
- ➡️ **箭头** - 绘制带箭头的直线
- 📝 **文字** - 添加文字标注
- 🧹 **橡皮擦** - 擦除内容
- 🖼️ **图片** - 支持粘贴图片（Ctrl+V）

### 样式配置
- 🎨 **颜色选择** - 线条颜色和填充颜色
- 📏 **粗细调节** - 1-50px 线条粗细
- 👁️ **透明度** - 0.1-1.0 透明度调节
- 🎨 **背景切换** - 自定义白板背景色

### 实时协同
- 👥 **多人在线** - 支持多人同时在线编辑
- ⚡ **毫秒级同步** - WebSocket 实时广播操作
- 🖱️ **光标显示** - 显示其他用户的光标位置和名称
- 📋 **成员列表** - 查看在线成员列表
- 🔄 **自动重连** - 网络异常时自动重连

### 白板管理
- 💾 **云端存储** - 数据持久化到 SQLite
- 📜 **历史记录** - 保存操作历史，支持回溯恢复
- 🗑️ **一键清空** - 快速清空白板内容
- ✂️ **撤销/重做** - 支持操作撤销和重做
- 🔍 **缩放平移** - 支持 Ctrl+滚轮缩放，抓手工具平移
- 📚 **图层管理** - 图层顺序调整、删除

### 快捷键
- `Ctrl+S` - 保存白板
- `Ctrl+Z` - 撤销
- `Ctrl+Y` - 重做
- `Ctrl+V` - 粘贴图片
- `Delete/Backspace` - 删除选中图层
- `Esc` - 切换回画笔工具
- `Ctrl+滚轮` - 缩放画布

## 快速开始

### 环境要求

#### 后端
- Go 1.21+
- GCC (TDM-GCC 或 MinGW-w64，用于 CGO 编译 SQLite)
  - 下载: https://jmeubank.github.io/tdm-gcc/
  - 安装后请确保添加到 PATH 环境变量

#### 前端
- Node.js 18+
- npm 或 pnpm

### 一键启动

1. **安装 GCC** (如果尚未安装)
   - Windows: 下载安装 TDM-GCC
   - 重启终端使环境变量生效

2. **编译并启动**
   ```bash
   # 方式一：使用启动脚本（推荐）
   start.bat
   
   # 方式二：手动启动
   # 1. 编译并启动后端
   cd backend
   go mod tidy
   set CGO_ENABLED=1
   go build -o whiteboard.exe ./cmd/main.go
   whiteboard.exe
   
   # 2. 启动前端（新终端）
   cd ../frontend
   npm install
   npm run dev
   ```

3. **访问应用**
   - 前端地址: http://localhost:3000
   - 后端地址: http://localhost:8080
   - API 文档: http://localhost:8080/health

### 编译部署

```bash
# 编译后端和前端
build.bat

# 停止所有服务
stop.bat
```

## 项目结构

```
whiteboard/
├── backend/                    # 后端 Go 项目
│   ├── cmd/
│   │   └── main.go            # 应用入口
│   ├── internal/
│   │   ├── server/            # 服务器初始化
│   │   │   └── server.go
│   │   ├── ws/                # WebSocket 模块
│   │   │   └── hub.go         # Hub 用户管理、消息广播
│   │   ├── models/            # 数据模型
│   │   │   ├── models.go      # 数据库模型定义
│   │   │   └── database.go    # 数据库初始化
│   │   ├── handlers/          # API 处理器
│   │   │   └── handlers.go    # REST API 和 WebSocket 处理
│   │   └── middleware/        # 中间件
│   │       └── middleware.go  # CORS、日志、异常恢复
│   ├── data/                  # 数据库文件目录
│   ├── .env                   # 环境变量
│   ├── go.mod
│   └── whiteboard.exe         # 编译后的可执行文件
├── frontend/                   # 前端 Vue3 项目
│   ├── src/
│   │   ├── components/         # 组件
│   │   │   ├── Toolbar.vue      # 工具栏
│   │   │   ├── UserList.vue     # 用户列表
│   │   │   ├── LayerPanel.vue   # 图层面板
│   │   │   └── HistoryPanel.vue # 历史记录面板
│   │   ├── views/              # 页面
│   │   │   ├── Home.vue        # 首页 - 白板列表
│   │   │   └── Board.vue       # 白板页面
│   │   ├── composables/        # 组合式函数
│   │   │   ├── useCanvas.js    # Canvas 绘图引擎
│   │   │   └── useWebSocket.js # WebSocket 客户端
│   │   ├── store/              # 状态管理
│   │   │   └── board.js        # 白板状态
│   │   ├── utils/              # 工具函数
│   │   │   ├── request.js      # HTTP 请求封装
│   │   │   └── common.js       # 通用工具
│   │   ├── types/              # 类型定义 (JSDoc)
│   │   │   └── index.js
│   │   ├── styles/             # 样式
│   │   │   └── index.scss
│   │   ├── router/             # 路由
│   │   │   └── index.js
│   │   ├── App.vue
│   │   └── main.js
│   ├── public/
│   ├── dist/                   # 编译后的静态文件
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
├── start.bat                   # 一键启动脚本
├── build.bat                   # 编译脚本
├── stop.bat                    # 停止脚本
└── prd-01.md                   # 需求文档
```

## API 接口

### 白板管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/boards` | 获取白板列表 |
| POST | `/api/boards` | 创建白板 |
| GET | `/api/boards/:id` | 获取白板详情 |
| PUT | `/api/boards/:id` | 更新白板信息 |
| DELETE | `/api/boards/:id` | 删除白板 |

### 元素管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/boards/:id/elements` | 获取白板元素 |
| POST | `/api/boards/:id/elements` | 保存白板元素 |
| POST | `/api/boards/:id/clear` | 清空白板 |

### 历史记录
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/boards/:id/history` | 获取历史记录 |
| POST | `/api/boards/:id/history/:hid/restore` | 恢复历史版本 |

### WebSocket
| 路径 | 说明 |
|------|------|
| `/ws/:id?name=xxx&color=xxx` | 白板实时通信 |

### WebSocket 消息类型
| 类型 | 说明 |
|------|------|
| `draw` | 绘图操作 |
| `cursor` | 光标位置 |
| `user_join` | 用户加入 |
| `user_leave` | 用户离开 |
| `users` | 用户列表同步 |
| `clear` | 清空白板 |
| `sync` | 全量数据同步 |
| `undo` | 撤销操作 |
| `redo` | 重做操作 |
| `ping/pong` | 心跳检测 |

## 数据库表结构

### boards - 白板表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | 主键 UUID |
| name | VARCHAR(100) | 白板名称 |
| background | VARCHAR(50) | 背景颜色 |
| thumbnail | TEXT | 缩略图 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### elements - 元素表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | 主键 UUID |
| board_id | VARCHAR(36) | 关联白板ID |
| type | VARCHAR(20) | 元素类型 |
| user_id | VARCHAR(36) | 用户ID |
| points | TEXT | 点坐标 (JSON) |
| style | TEXT | 样式 (JSON) |
| text | TEXT | 文字内容 |
| image_data | TEXT | 图片数据 (base64) |
| timestamp | BIGINT | 时间戳 |

### histories - 历史记录表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | 主键 UUID |
| board_id | VARCHAR(36) | 关联白板ID |
| action | VARCHAR(50) | 操作类型 |
| user_id | VARCHAR(36) | 用户ID |
| user_name | VARCHAR(50) | 用户名 |
| snapshot | TEXT | 快照数据 (JSON) |
| timestamp | DATETIME | 操作时间 |

### operation_logs - 操作日志表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | 主键 UUID |
| board_id | VARCHAR(36) | 关联白板ID |
| user_id | VARCHAR(36) | 用户ID |
| action | VARCHAR(50) | 操作类型 |
| data | TEXT | 操作数据 |
| timestamp | DATETIME | 操作时间 |

## 配置说明

### 后端配置 (backend/.env)
```env
SERVER_PORT=8080              # 服务端口
WEBSOCKET_PORT=8080           # WebSocket 端口
DATABASE_PATH=./data/whiteboard.db  # 数据库路径
ALLOW_ORIGINS=*               # 允许的跨域源
MAX_CONNECTIONS=1000          # 最大连接数
```

### 前端配置 (frontend/vite.config.js)
```javascript
export default {
  server: {
    port: 3000,
    proxy: {
      '/api': 'http://localhost:8080',
      '/ws': {
        target: 'ws://localhost:8080',
        ws: true
      }
    }
  }
}
```

## 常见问题

### 1. 后端编译失败，提示 "CGO_ENABLED=0"
- 原因: SQLite 驱动需要 CGO 支持
- 解决: 
  1. 安装 TDM-GCC 或 MinGW-w64
  2. 添加到 PATH 环境变量
  3. 设置 `set CGO_ENABLED=1`
  4. 重新编译

### 2. 前端启动后无法连接后端
- 检查后端服务是否正常启动
- 检查 vite.config.js 中的代理配置
- 检查防火墙设置

### 3. 多人协同不同步
- 检查 WebSocket 连接状态
- 确认使用相同的白板 ID
- 检查浏览器控制台是否有错误信息

## 开发计划

- [ ] 支持导出为图片/PDF
- [ ] 支持模板库
- [ ] 支持语音/视频通话
- [ ] 支持文件附件
- [ ] 支持权限管理
- [ ] 移动端适配

## 许可证

MIT License
