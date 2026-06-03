开发一个完整的后台管理系统，第一版只支持菜单和用户的管理。
技术栈：Go + Gin + SQLite（后端代码放在 server 文件夹），Vue3 + Element Plus + Pinia（前端代码放在 admin 文件夹）。要求实现以下功能：

一、后端接口开发（server 文件夹）

1. 数据库表设计

· 管理员用户表：id、用户名、密码（bcrypt加密）、昵称、角色（super/admin/operator）、状态、最后登录时间、创建时间
· 菜单表：id、父级id（支持三级菜单）、菜单名称、路由路径、图标、排序值、是否显示、创建时间
· 角色菜单权限表：id、角色、菜单id

2. 初始化超级管理员

· 程序启动时自动检测，若无管理员则创建超级管理员：用户名 admin，密码 123456

3. 认证接口

· POST /api/admin/login：登录，返回JWT token
· GET /api/admin/info：获取当前管理员信息及所拥有的菜单树
· POST /api/admin/logout：退出登录

4. 管理员管理接口

· GET /api/admin/users：分页查询管理员列表（支持用户名搜索）
· POST /api/admin/users：创建管理员（超级管理员可操作）
· PUT /api/admin/users/:id：编辑管理员
· DELETE /api/admin/users/:id：删除管理员
· PUT /api/admin/users/:id/status：启用/禁用管理员
· PUT /api/admin/users/:id/password：重置密码

5. 菜单管理接口

· GET /api/admin/menus：获取完整菜单树（支持多级嵌套，最多三级）
· POST /api/admin/menus：创建菜单
· PUT /api/admin/menus/:id：编辑菜单
· DELETE /api/admin/menus/:id：删除菜单（同时删除子菜单）

6. 角色权限接口

· GET /api/admin/roles：获取角色列表
· GET /api/admin/roles/:role/menus：获取角色拥有的菜单权限ID列表
· PUT /api/admin/roles/:role/menus：设置角色的菜单权限

7. 项目结构要求

· server/handlers/：处理函数
· server/models/：数据模型
· server/middleware/：中间件（JWT、管理员权限）
· server/routes/：路由注册
· server/utils/：工具函数

二、前端页面开发（admin 文件夹）

1. 登录页

· 用户名、密码表单，Element Plus表单验证
· 登录成功跳转后台首页，存储token和用户信息

2. 后台布局

· 左侧可折叠菜单栏：根据后端返回的菜单树动态渲染（支持三级菜单，使用 el-sub-menu 嵌套）
· 右侧顶部：显示管理员昵称、退出按钮，面包屑导航
· 右侧内容区域：展示当前选中菜单对应的页面

3. 管理员管理页

· 表格展示：用户名、昵称、角色（super/admin/operator）、状态、最后登录时间
· 顶部搜索：按用户名搜索
· 操作按钮：新增、编辑、删除、启用/禁用、重置密码
· 新增/编辑弹窗：用户名、昵称、密码、角色下拉选择
· 超级管理员不可被删除或修改角色

4. 菜单管理页

· 树形表格展示菜单：名称、路由路径、图标、排序、是否显示
· 支持展开/折叠
· 操作按钮：新增（可指定父级）、编辑、删除
· 新增/编辑弹窗：父级菜单下拉选择（最多三级）、菜单名称、路由路径、图标选择器、排序数字、是否显示开关

5. 角色权限页

· 角色切换：Tab或下拉切换角色（super/admin/operator）
· 权限分配：左侧展示完整菜单树（带复选框），勾选后保存该角色的菜单权限
· 超级管理员角色默认拥有所有权限且不可修改

6. 项目结构要求

· admin/src/api/：API请求封装
· admin/src/views/：页面组件
· admin/src/router/：路由配置
· admin/src/stores/：Pinia状态管理
· admin/src/layout/：布局组件

三、权限控制

后端：JWT中间件验证token，每个接口校验当前角色是否有权限访问该菜单对应的路由

前端：

· 路由守卫：未登录跳转登录页
· 动态路由：根据后端返回的菜单树动态生成可访问的路由
· 按钮级权限：根据角色控制页面上新增/编辑/删除等按钮的显示

四、开发要求

· 后端代码全部放在 server 文件夹，前端代码全部放在 admin 文件夹
· 菜单最多支持三级，使用递归组件渲染
· 初始化后使用 admin / 123456 登录
· 数据库文件使用 server/ecommerce.db

