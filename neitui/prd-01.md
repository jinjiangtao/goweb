产品需求文档：内推招聘系统（管理后台 + 后端）

一、项目背景

为公司内部开发一套员工内推招聘管理系统，实现职位发布、候选人推荐、进度跟踪和数据统计功能。后端使用 Go + Gin + SQLite，管理后台使用 Vue3 + Element Plus。

二、角色与权限

角色 权限说明
员工 发布职位、推荐候选人、查看自己的内推记录和进度
HR 查看所有内推记录、更新候选人状态、查看统计看板
管理员 拥有HR权限，外加用户管理（启用/禁用员工账号）

三、后端接口设计（server 文件夹）

3.1 数据库表结构（SQLite）

· 用户表：id, username, password(bcrypt), role(employee/hr/admin), real_name, status(启用/禁用), created_at
· 职位表：id, title, requirement, salary_range, location, status(发布中/已关闭), created_by(员工id), created_at
· 内推记录表：id, job_id, candidate_name, candidate_phone, resume_path(本地文件路径), status(初筛中/已面试/发offer/已入职/淘汰), reject_reason, employee_id(推荐人), hr_remark, created_at, updated_at

3.2 文件存储

· 简历文件上传后保存在本地 server/uploads/resumes/ 目录
· 文件名规则：时间戳_原文件名，避免重名
· 接口返回相对路径或访问URL

3.3 默认初始化数据

· 启动时自动创建管理员账号：admin / 123456（角色为admin）
· 可选创建测试HR账号：hr / 123456（角色为hr）

3.4 接口列表

需登录（JWT token）：

· POST /api/user/login - 登录，返回 token 和用户信息

员工接口：

· POST /api/jobs - 发布职位
· GET /api/jobs/my - 获取自己发布的职位列表
· PUT /api/jobs/:id/status - 关闭/开启职位
· POST /api/referrals - 推荐候选人（含简历上传，multipart/form-data）
· GET /api/referrals/my - 获取自己的内推记录列表（支持分页、按状态筛选）
· GET /api/referrals/my/stats - 自己的内推统计（总数、各状态数量）

HR/管理员接口（需角色校验）：

· GET /api/admin/referrals - 获取所有内推记录（支持分页、按状态筛选、按职位筛选、按推荐人筛选）
· PUT /api/admin/referrals/:id/status - 更新候选人状态（含淘汰原因、HR备注）
· GET /api/admin/stats - 统计看板数据（总内推数、已入职数、面试通过率、各员工内推成功排行、近30天内推趋势）
· GET /api/admin/jobs - 获取所有职位列表
· GET /api/admin/referrals/export - 导出当前筛选结果为Excel

管理员专属：

· GET /api/admin/users - 获取员工列表
· POST /api/admin/users - 添加员工账号
· PUT /api/admin/users/:id/status - 启用/禁用员工账号
· PUT /api/admin/users/:id/reset-password - 重置密码

3.5 JWT中间件

· 除 /api/user/login 外，其他接口均需验证 token
· 员工接口需验证角色为 employee/hr/admin
· HR接口需验证角色为 hr/admin
· 管理员接口需验证角色为 admin

四、管理后台功能（admin 文件夹）

4.1 登录页面

· 表单：用户名、密码
· 登录成功后根据角色展示不同菜单

4.2 员工端

菜单：职位管理、我的推荐

职位管理：

· 表格展示自己发布的职位（标题、要求、薪资、地点、状态、创建时间）
· 按钮：发布新职位（弹窗表单）
· 操作：关闭/开启职位

我的推荐：

· 表格展示自己推荐的候选人（姓名、手机号、职位、状态、HR备注、提交时间）
· 支持按状态筛选、按姓名或手机号搜索、分页
· 卡片显示自己的内推统计

4.3 HR端/管理员端

菜单：内推列表、职位列表、统计看板（管理员多一个用户管理）

内推列表：

· 表格展示所有内推记录（推荐人、候选人姓名、手机号、职位、状态、提交时间、更新时间）
· 筛选：按状态、按职位、按推荐人
· 搜索：候选人姓名或手机号
· 分页、导出Excel
· 操作：更新状态（弹窗选择状态、填写淘汰原因、HR备注）

职位列表：

· 表格展示所有职位，支持按发布人筛选

统计看板：

· 卡片：总内推数、已入职数、面试通过率
· 柱状图：各员工内推成功人数排行（ECharts）
· 折线图：近30天内推趋势

用户管理（仅管理员）：

· 表格展示员工账号
· 添加员工、启用/禁用、重置密码

4.4 通用要求

· 路由守卫、Pinia存储token、axios统一封装
· 上传文件用 multipart/form-data

五、技术栈

· 后端：Go 1.21+，Gin，SQLite，JWT，bcrypt
· 前端：Vue3，Element Plus，Pinia，Vue Router，axios，ECharts
· 简历存储：本地文件系统

六、项目结构

```
server/
├── main.go
├── models/        # 数据库模型
├── handlers/      # 接口处理
├── middleware/    # JWT、角色校验
├── routes/        # 路由
├── utils/         # 工具函数
└── uploads/resumes/  # 简历存储目录
admin/             # Vue3管理后台
```

