## 角色：全栈开发工程师

## 任务：为现有 ERP 系统同时添加「客户管理」和「供应商管理」两个模块

## 现有系统情况：
- 已有功能：用户登录/注册、用户组管理、菜单管理、产品管理
- 技术栈：Go + Gin + GORM + SQLite，Vue3 + Element Plus + Pinia
- 认证方式：JWT Token，需要登录后才能访问

## 功能需求：

### 模块一：客户管理
1. 客户列表页，表格展示所有客户
2. 支持新增客户（弹窗表单）
3. 支持编辑客户信息
4. 支持删除客户（二次确认）
5. 支持按客户名称模糊搜索
6. 支持分页展示

**客户字段**：客户名称（必填）、联系人、联系电话、地址、备注

### 模块二：供应商管理
1. 供应商列表页，表格展示所有供应商
2. 支持新增供应商（弹窗表单）
3. 支持编辑供应商信息
4. 支持删除供应商（二次确认）
5. 支持按供应商名称模糊搜索
6. 支持分页展示

**供应商字段**：供应商名称（必填）、联系人、联系电话、地址、备注

## 后端要求：
- 使用 Go + Gin + GORM
- 每个模块独立创建 RESTful API：
  - 客户：POST /api/customers、GET /api/customers、GET /api/customers/:id、PUT /api/customers/:id、DELETE /api/customers/:id
  - 供应商：POST /api/suppliers、GET /api/suppliers、GET /api/suppliers/:id、PUT /api/suppliers/:id、DELETE /api/suppliers/:id
- 列表接口支持分页和按名称搜索
- 自动迁移新增的两个数据表

## 前端要求：
- 使用 Vue3 + Element Plus
- 客户管理页面：搜索区（客户名称 + 搜索/重置/新增按钮）、表格区（名称、联系人、电话、地址、备注、操作）、分页区
- 供应商管理页面：结构同上，字段改为供应商相关
- 新增/编辑使用 el-dialog 弹窗表单
- 在菜单管理中添加两个菜单：客户管理（/customers）、供应商管理（/suppliers）

## 输出要求：
请输出以下代码：
1. 后端：models/customer.go、models/supplier.go
2. 后端：handlers/customer.go、handlers/supplier.go
3. 后端：路由注册代码片段
4. 前端：views/erp/customer/index.vue
5. 前端：views/erp/supplier/index.vue
6. 前端：api/customer.js、api/supplier.js
7. 前端：路由配置代码片段

## 特别说明：
- 代码风格与现有产品管理模块保持一致
- 每个文件标注文件路径
```
