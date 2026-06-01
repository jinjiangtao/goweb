请帮我用 Go 和 Vue 实现一个完整的即时通讯系统，代码要可以直接运行。

后端要求：
- Go + Gin + WebSocket + SQLite，不需要 Redis（用 Go 本地缓存代替，比如一个简单的 sync.Map 或 go-cache 存在线状态和未读数）
- 实现 JWT 登录注册
- 单聊：文本消息、图片消息
- 群聊：创建群、拉人、踢人、发群消息
- 消息存 SQLite，在线状态和未读数用本地缓存
- 离线消息：上线后自动拉取
- 已读回执
- WebSocket 连接管理（Hub 模式），支持心跳检测
- 后端代码放到 server 目录

前端要求：
- Vue 3 + TypeScript + Vite + Pinia
- 页面：登录页、好友列表、聊天窗口
- WebSocket 封装：自动重连、心跳
- 聊天窗口显示消息气泡，支持发送文本和图片
- 有新消息时桌面通知
- 前端代码放到 web 目录

输出内容：
1. 完整的目录结构（server 和 web 分开）
2. 后端所有核心代码（server/main.go、handler、service、model、websocket、storage、cache）
3. 前端所有核心代码（web/src 下的 App、聊天组件、WebSocket composable、store）
4. 数据库建表 SQL（SQLite 语法：users、friends、groups、group_members、messages）
5. 启动说明（如何运行 server 和 web）
6. 接口都要跑下测试用例，前端也要跑下测试用例。

注意：不用 Redis，用 Go 本地缓存。不要用第三方 IM SDK，从零实现。代码要完整能跑起来。

请熟悉现在的im 聊天现在的功能， 在现在功能点的基础上完成下面的迭代。
1. im 用户注册的时候可以添加个人头像。注册后也可以编辑个人资料修改头像。
2. im 的个人头像展示在右侧的栏和聊天窗口中。 聊天窗口的头像替换之前的名词。
3. 发送消息的时候支持发送emoji表情和图片。