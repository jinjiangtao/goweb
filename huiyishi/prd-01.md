开发会议室预订系统 - 管理后台和后端

用 Go + Gin + SQLite 做后端，Vue3 + Element Plus 做管理后台，深色科技感风格。

一、后端接口（server 文件夹）

数据库三张表：

· 会议室表：id, name, capacity, devices（逗号分隔）, status, created_at
· 预订表：id, room_id, name, phone, date, start_time, end_time, purpose, status, created_at
· 管理员表：id, username, password（bcrypt）, nickname, created_at

初始化数据：默认管理员 admin / 123456

接口列表：

· POST /api/admin/login - 登录，返回 JWT token
· 会议室 CRUD（需 JWT）：
  · GET /api/admin/rooms - 列表（支持搜索）
  · POST /api/admin/rooms - 新增
  · PUT /api/admin/rooms/:id - 编辑
  · DELETE /api/admin/rooms/:id - 删除（检查是否有预订记录）
· 预订管理（需 JWT）：
  · GET /api/admin/bookings - 列表，支持按日期、会议室、状态筛选，分页
  · PUT /api/admin/bookings/:id/cancel - 管理员取消预订，需传取消原因
· 统计（需 JWT）：
  · GET /api/admin/stats - 返回今日预订总数、本周各会议室预订数量、本周每天预订数量

JWT 中间件：所有 /api/admin/* 接口都要验证 token。

二、管理后台（admin 文件夹）

深色主题，顶部导航+侧边栏布局。

登录页：背景模糊效果，账号 admin 密码 123456，用 Pinia 存储 token，路由守卫拦截未登录访问。

功能模块：

1. 会议室管理
   · 表格列：名称、容纳人数、设备（标签形式）、状态（开关/标签）、操作
   · 搜索框：按名称搜索
   · 新增/编辑弹窗：名称、容纳人数（数字）、设备（多选框：投影、电视、白板、会议电话）、状态（启用/禁用）
   · 删除：二次确认，如果该会议室有预订记录则提示无法删除
2. 预订管理
   · 顶部筛选：日期选择器、会议室下拉、状态下拉（全部/已预订/已取消）
   · 表格展示：会议室名称、预订人、手机号、日期、时间段、用途、状态、取消原因（如有）
   · 每条记录操作按钮：取消预订（仅状态为已预订时可点），弹窗输入取消原因
   · 分页
3. 统计看板
   · 顶部卡片：今日总预订数、启用会议室数量
   · 柱状图：本周每天预订数量（用 ECharts）
   · 水平条形图或表格：各会议室本周使用时长排行（按预订数量或时间段计算）

三、额外要求

· 管理后台所有接口请求统一封装，自动带 token
· 退出登录清除 token 和 Pinia 状态
· 风格统一深色系，侧边栏半透明玻璃效果，表格和卡片用暗色背景

