
帮我开发一个学生报名系统，技术栈是：
· 后端：Go + Gin + SQLite（代码放 server 文件夹）
· 管理后台：Vue3 + Element Plus（代码放 admin 文件夹）
· 报名页面：H5 单页，随便写，放在 web 文件夹

功能要求如下：

一、H5 报名页面（web 文件夹）

做一个简单的报名表单，字段包括：

· 姓名
· 手机号
· 年龄
. 户口地址
· 学校

不用太花哨，但要干净、适配手机。提交后调用后端接口 POST /api/signup，提交成功提示“报名成功”，失败提示错误。默认状态为“报名中”。

二、管理后台（admin 文件夹）

1. 登录页面：账号密码登录，内置一个管理员账号 admin，密码 123456（后端初始化时自动创建，密码用 bcrypt 加密存储到 SQLite 数据库）。
2. 登录成功后进入后台主页，左侧菜单有“报名列表”和“统计看板”。
3. 报名列表页：
   · 表格展示所有报名记录：姓名、手机号、年龄、户口地址、学校、状态、提交时间。
   · 支持按姓名或手机号搜索，支持分页。
   · 每条记录后面有操作按钮，可以修改状态：下拉选择“报名中”、“报名成功”、“报名失败”，点确认后立即更新。
4. 统计看板：
   · 用卡片显示当前“报名中”、“报名成功”、“报名失败”的数量。
   · 可选：画个简单的柱状图展示每天报名人数（如果后端有接口，没有就算了，不强求）。

三、后端接口（server 文件夹）

1. 数据库 SQLite，建两张表：
   · 报名表：id, name, phone, age, school, status（默认 'pending'）, created_at
   · 管理员表：id, username, password（bcrypt）, nickname, created_at
2. 启动时自动检查，如果不存在 admin 用户，则创建 username=admin, password=123456 加密后的记录。
3. 接口：
   · POST /api/signup：公开，接收 name, phone, age, school，存入报名表，返回成功。
   · POST /api/admin/login：接收 username, password，验证成功返回 JWT token。
   · GET /api/admin/signups：需要 JWT token，支持 query 参数 page, pageSize, keyword，返回分页列表。
   · PUT /api/admin/signups/:id/status：需要 JWT token，接收 body { status }，更新报名状态。
   · GET /api/admin/stats：需要 JWT token，返回 { pending, approved, rejected } 三个数量。
4. JWT 中间件：除了 /api/signup 和 /api/admin/login 外，其他 /api/admin/* 接口都要验证 token。

四、其他要求

· 后端项目结构清晰：handlers, models, middleware, routes, utils。
· 前端管理后台用 Pinia 存储 token，路由守卫拦截未登录访问。
· 报名页面独立，不需要登录。
· 所有代码直接给出，能跑起来。

