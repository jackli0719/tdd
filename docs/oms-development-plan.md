# OMS 订单管理系统 - 开发文档

## 项目概述

- **项目名称**: OMS (Order Management System)
- **技术栈**: Go + Gin + GORM + MySQL
- **目标**: 全功能版订单管理系统

---

## 一、项目初始化

### 1.1 初始化 Go 项目
```bash
mkdir -p oms && cd oms
go mod init oms
```

### 1.2 安装依赖
```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/mysql
```

### 1.3 创建目录结构
```
oms/
├── cmd/server/main.go
├── internal/
│   ├── config/
│   ├── handler/
│   ├── model/
│   ├── repository/
│   ├── service/
│   └── router/
├── pkg/response/
└── docs/
```

---

## 二、数据库 Schema

### 2.1 统一 Schema (SQL)

```sql
-- OMS 数据库 Schema
CREATE DATABASE IF NOT EXISTS oms DEFAULT CHARSET utf8mb4;
USE oms;

-- ----------------------------
-- 用户表
-- ----------------------------
CREATE TABLE `users` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(50) NOT NULL COMMENT '用户名',
    `email` varchar(100) NOT NULL COMMENT '邮箱',
    `phone` varchar(20) NOT NULL COMMENT '手机号',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- ----------------------------
-- 产品表
-- ----------------------------
CREATE TABLE `products` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(100) NOT NULL COMMENT '产品名称',
    `price` decimal(10,2) NOT NULL COMMENT '单价',
    `stock` int NOT NULL DEFAULT 0 COMMENT '库存',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品表';

-- ----------------------------
-- 订单表
-- ----------------------------
CREATE TABLE `orders` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `order_no` varchar(32) NOT NULL COMMENT '订单号',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `total_amount` decimal(10,2) NOT NULL COMMENT '总金额',
    `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/confirmed/in_service/completed/cancelled',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_no` (`order_no`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_orders_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- ----------------------------
-- 订单明细表
-- ----------------------------
CREATE TABLE `order_items` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `order_id` bigint unsigned NOT NULL COMMENT '订单ID',
    `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
    `price` decimal(10,2) NOT NULL COMMENT '下单时单价',
    `quantity` int NOT NULL COMMENT '数量',
    `subtotal` decimal(10,2) NOT NULL COMMENT '小计',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_product_id` (`product_id`),
    CONSTRAINT `fk_order_items_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`),
    CONSTRAINT `fk_order_items_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单明细表';
```

### 2.2 ER 关系

```
users (1) ----< orders (N)
orders (1) ----< order_items (N)
products (1) ----< order_items (N)
```

### 2.3 字段规范

| 规范 | 说明 |
|------|------|
| 主键 | `bigint unsigned`, 自增 |
| 时间 | `datetime`, `created_at`/`updated_at` |
| 金额 | `decimal(10,2)` |
| 状态 | `varchar(20)`, 枚举值 |
| 外键 | 带 `fk_` 前缀约束名 |
| 索引 | `idx_` 前缀, `uk_` 唯一索引 |

---

## 三、任务清单

> 每任务遵循：**契约先行 → 测试用例 → 实现代码 → 提交**

### Phase 1: 项目基础

#### 1.1 创建项目目录结构
- [ ] **契约**: 定义目录结构规范
- [ ] 创建 `oms/{cmd,internal, pkg, docs}/`
- [ ] **测试**: 验证目录存在
- [ ] 提交

#### 1.2 初始化 go.mod
- [ ] **契约**: `go.mod` 内容规范
- [ ] `go mod init oms`
- [ ] **测试**: `go build` 无报错
- [ ] 提交

#### 1.3 创建 config 配置模块
- [ ] **契约**: `Config` 结构，包含 `Port`, `DSN`
- [ ] 创建 `internal/config/config.go`
- [ ] **测试**: 配置加载测试
- [ ] 提交

#### 1.4 创建数据库连接
- [ ] **契约**: `DB` 实例，`gorm.Open` 返回 `*gorm.DB`
- [ ] 创建 `internal/config/database.go`
- [ ] **测试**: Ping 测试
- [ ] 提交

#### 1.5 创建统一响应 pkg
- [ ] **契约**: `Response{code, message, data}` 结构
- [ ] 创建 `pkg/response/response.go`
- [ ] **测试**: `Success(ctx, data)`, `Error(ctx, code, msg)` 测试
- [ ] 提交

---

### Phase 2: 用户模块

#### 2.1 创建 user model
- [ ] **契约**: `User` 结构体字段（参照 Schema）
- [ ] 创建 `internal/model/user.go`
- [ ] **测试**: JSON 序列化/反序列化测试
- [ ] 提交

#### 2.2 创建 user repository (CRUD)
- [ ] **契约**: `Create`, `GetByID`, `Update`, `Delete`, `List` 接口
- [ ] 创建 `internal/repository/user.go`
- [ ] **测试**: CRUD 单元测试
- [ ] 提交

#### 2.3 创建 user service (业务逻辑)
- [ ] **契约**: 业务方法签名
- [ ] 创建 `internal/service/user.go`
- [ ] **测试**: 业务逻辑测试（含参数校验）
- [ ] 提交

#### 2.4 创建 user handler (HTTP)
- [ ] **契约**: API 路由、请求/响应结构
  ```go
  POST   /api/users        → CreateUser
  GET    /api/users        → ListUsers
  GET    /api/users/:id    → GetUser
  PUT    /api/users/:id    → UpdateUser
  DELETE /api/users/:id    → DeleteUser
  ```
- [ ] 创建 `internal/handler/user.go`
- [ ] **测试**: HTTP handler 测试
- [ ] 提交

#### 2.5 注册 user 路由
- [ ] **契约**: 路由组 `/api/users`
- [ ] 在 router 中注册
- [ ] **测试**: 路由注册测试
- [ ] 提交

---

### Phase 3: 产品模块

#### 3.1 创建 product model
- [ ] **契约**: `Product` 结构体字段（参照 Schema）
- [ ] 创建 `internal/model/product.go`
- [ ] **测试**: JSON 序列化/反序列化测试
- [ ] 提交

#### 3.2 创建 product repository (CRUD)
- [ ] **契约**: `Create`, `GetByID`, `Update`, `Delete`, `List` 接口
- [ ] 创建 `internal/repository/product.go`
- [ ] **测试**: CRUD 单元测试
- [ ] 提交

#### 3.3 创建 product service (业务逻辑)
- [ ] **契约**: 业务方法签名（含库存校验）
- [ ] 创建 `internal/service/product.go`
- [ ] **测试**: 业务逻辑测试（含库存扣减校验）
- [ ] 提交

#### 3.4 创建 product handler (HTTP)
- [ ] **契约**: API 路由
  ```go
  POST   /api/products
  GET    /api/products
  GET    /api/products/:id
  PUT    /api/products/:id
  DELETE /api/products/:id
  ```
- [ ] 创建 `internal/handler/product.go`
- [ ] **测试**: HTTP handler 测试
- [ ] 提交

#### 3.5 注册 product 路由
- [ ] **契约**: 路由组 `/api/products`
- [ ] 在 router 中注册
- [ ] **测试**: 路由注册测试
- [ ] 提交

---

### Phase 4: 订单模块

#### 4.1 创建 order model
- [ ] **契约**: `Order` 结构体字段（参照 Schema）
- [ ] 创建 `internal/model/order.go`
- [ ] **测试**: JSON 序列化测试
- [ ] 提交

#### 4.2 创建 order_item model
- [ ] **契约**: `OrderItem` 结构体字段
- [ ] 创建 `internal/model/order_item.go`
- [ ] **测试**: JSON 序列化测试
- [ ] 提交

#### 4.3 创建 order repository (CRUD)
- [ ] **契约**: `Create`, `GetByID`, `Update`, `Delete`, `List`, `GetByOrderNo` 接口
- [ ] 创建 `internal/repository/order.go`
- [ ] **测试**: CRUD 单元测试
- [ ] 提交

#### 4.4 创建 order service (业务逻辑 + 状态流转)
- [ ] **契约**: 状态机接口
  ```go
  Confirm()   // pending → confirmed
  Start()     // confirmed → in_service
  Complete()  // in_service → completed
  Cancel()    // pending/confirmed → cancelled
  ```
- [ ] 创建 `internal/service/order.go`
- [ ] **测试**: 状态流转测试（含非法流转校验）
- [ ] 提交

#### 4.5 创建 order handler (HTTP)
- [ ] **契约**: API 路由
  ```go
  POST   /api/orders
  GET    /api/orders
  GET    /api/orders/:id
  PUT    /api/orders/:id
  DELETE /api/orders/:id
  POST   /api/orders/:id/confirm
  POST   /api/orders/:id/start
  POST   /api/orders/:id/complete
  POST   /api/orders/:id/cancel
  ```
- [ ] 创建 `internal/handler/order.go`
- [ ] **测试**: HTTP handler 测试
- [ ] 提交

#### 4.6 注册 order 路由
- [ ] **契约**: 路由组 `/api/orders`
- [ ] 在 router 中注册
- [ ] **测试**: 路由注册测试
- [ ] 提交

---

### Phase 5: 统计模块

#### 5.1 订单统计 API
- [ ] **契约**: `GET /api/stats/orders` → 按状态分组的订单数
- [ ] 创建 `internal/handler/stats.go` → `OrderStats`
- [ ] **测试**: 统计准确性测试
- [ ] 提交

#### 5.2 营收统计 API
- [ ] **契约**: `GET /api/stats/revenue` → 总营收、日营收
- [ ] 创建 `internal/service/stats.go` → `RevenueStats`
- [ ] **测试**: 营收计算测试
- [ ] 提交

---

### Phase 6: 路由整合

#### 6.1 创建 router 模块
- [ ] **契约**: `SetupRouter() *gin.Engine`
- [ ] 创建 `internal/router/router.go`
- [ ] **测试**: 路由注册完整性测试
- [ ] 提交

#### 6.2 整合所有路由
- [ ] **契约**: 所有模块路由已注册
- [ ] 整合 user/product/order/stats 路由组
- [ ] **测试**: 路由覆盖测试
- [ ] 提交

#### 6.3 创建 main.go 入口
- [ ] **契约**: 启动服务，监听 `:8080`
- [ ] 创建 `cmd/server/main.go`
- [ ] **测试**: 服务启动测试
- [ ] 提交

---

### Phase 7: 测试验证

#### 7.1 编译检查
- [ ] **契约**: `go build ./...` 无错误
- [ ] 执行 `go build ./...`
- [ ] **测试**: 编译成功
- [ ] 提交

#### 7.2 启动服务测试
- [ ] **契约**: 服务正常启动，API 可响应
- [ ] `go run cmd/server/main.go`
- [ ] **测试**: `curl http://localhost:8080/api/users` 返回正确响应
- [ ] 提交

---

## 四、API 列表

### 4.1 用户 API
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/users | 用户列表 |
| GET | /api/users/:id | 用户详情 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/:id | 更新用户 |
| DELETE | /api/users/:id | 删除用户 |

### 4.2 产品 API
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/products | 产品列表 |
| GET | /api/products/:id | 产品详情 |
| POST | /api/products | 创建产品 |
| PUT | /api/products/:id | 更新产品 |
| DELETE | /api/products/:id | 删除产品 |

### 4.3 订单 API
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/orders | 订单列表 |
| GET | /api/orders/:id | 订单详情 |
| POST | /api/orders | 创建订单 |
| PUT | /api/orders/:id | 更新订单 |
| DELETE | /api/orders/:id | 删除订单 |
| POST | /api/orders/:id/confirm | 确认订单 |
| POST | /api/orders/:id/start | 开始服务 |
| POST | /api/orders/:id/complete | 完成服务 |
| POST | /api/orders/:id/cancel | 取消订单 |

### 4.4 统计 API
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/stats/orders | 订单统计 |
| GET | /api/stats/revenue | 营收统计 |

---

## 五、订单状态流转（上门服务模式）

```
pending (待确认)
    ├── confirmed (已确认)
    │       ├── in_service (服务中)
    │       │       └── completed (已完成)
    │       └── cancelled (已取消)
    └── cancelled (已取消)
```

---

## 六、Agent 分工参考

| Agent | 负责模块 | 任务编号 |
|-------|----------|----------|
| Agent-User | 用户模块 | 2.1 - 2.5 |
| Agent-Product | 产品模块 | 3.1 - 3.5 |
| Agent-Order | 订单模块 | 4.1 - 4.6 |
| Agent-Stats | 统计模块 | 5.1 - 5.2 |
| Agent-Integration | 整合 | 1.1 - 1.5, 6.1 - 6.3, 7.1 - 7.2 |

---

## 七、验证命令

```bash
# 编译
go build ./...

# 运行
go run cmd/server/main.go

# 测试
curl http://localhost:8080/api/users
curl http://localhost:8080/api/products
curl http://localhost:8080/api/orders
```

---

## 八、数据库集成 (SQLite)

> 每任务遵循：**契约先行 → 测试验证 → 实现提交**

### Phase 8: 数据库集成

#### 8.1 安装 SQLite 驱动
- [ ] **契约**: `gorm.io/driver/sqlite`
- [ ] `go get gorm.io/driver/sqlite`
- [ ] **测试**: `go build ./...` 无报错
- [ ] 提交

#### 8.2 修改 database.go 支持 SQLite
- [ ] **契约**: DSN 支持 `sqlite://oms.db`
- [ ] 修改 `internal/config/database.go` 支持 sqlite driver
- [ ] **测试**: DB 连接成功
- [ ] 提交

#### 8.3 表结构迁移
- [ ] **契约**: AutoMigrate 所有表到 SQLite
- [ ] users, products, orders, order_items
- [ ] **测试**: E2E CRUD API 测试
- [ ] 提交

---

## 九、前端开发 (Vue 3 + Element Plus)

> 每任务遵循：**契约先行 → 测试验证 → 实现提交**
>
> **API 契约**: 参照本文档第四节的 API 列表

### Phase 9: 前端项目初始化

#### 9.1 创建 Vue 3 项目
- [ ] **契约**: `npm create vue@latest frontend`
- [ ] 创建项目，安装依赖
- [ ] **测试**: `npm run dev` 启动成功
- [ ] 提交

#### 9.2 安装前端依赖
- [ ] **契约**: Element Plus, Axios, Vue Router
- [ ] `npm install element-plus axios vue-router`
- [ ] **测试**: 编译无报错
- [ ] 提交

#### 9.3 配置 Vite 代理
- [ ] **契约**: 开发环境 API 代理到 `:8080`
- [ ] 修改 `vite.config.ts`
- [ ] **测试**: 前端访问 `/api/*` 代理到后端
- [ ] 提交

---

### Phase 10: 前端公共层

#### 10.1 创建 Layout 组件
- [ ] **契约**: 侧边栏 + 头部布局
- [ ] `src/components/Layout.vue`
- [ ] **测试**: 页面渲染正确
- [ ] 提交

#### 10.2 封装 API 请求
- [ ] **契约**: 统一响应格式 `{code, message, data}`
- [ ] `src/api/*.ts` - user, product, order, stats
- [ ] **测试**: API 调用成功
- [ ] 提交

#### 10.3 配置路由
- [ ] **契约**: Vue Router 配置
- [ ] `/` - 仪表盘
- [ ] `/users` - 用户管理
- [ ] `/products` - 产品管理
- [ ] `/orders` - 订单管理
- [ ] `/stats` - 统计页面
- [ ] **测试**: 路由跳转正常
- [ ] 提交

---

### Phase 11: 用户管理页面

#### 11.1 用户列表页
- [ ] **契约**: `GET /api/users` → 列表展示
- [ ] `src/views/user/UserList.vue`
- [ ] **测试**: 分页、搜索正常
- [ ] 提交

#### 11.2 用户表单弹窗
- [ ] **契约**: `POST/PUT /api/users` → 新建/编辑
- [ ] `src/views/user/UserForm.vue`
- [ ] **测试**: 新建、编辑、删除正常
- [ ] 提交

---

### Phase 12: 产品管理页面

#### 12.1 产品列表页
- [ ] **契约**: `GET /api/products` → 列表展示
- [ ] `src/views/product/ProductList.vue`
- [ ] **测试**: 分页、搜索正常
- [ ] 提交

#### 12.2 产品表单弹窗
- [ ] **契约**: `POST/PUT /api/products` → 新建/编辑
- [ ] `src/views/product/ProductForm.vue`
- [ ] **测试**: 新建、编辑、删除正常
- [ ] 提交

---

### Phase 13: 订单管理页面

#### 13.1 订单列表页
- [ ] **契约**: `GET /api/orders` → 列表展示
- [ ] `src/views/order/OrderList.vue`
- [ ] **测试**: 分页、状态筛选正常
- [ ] 提交

#### 13.2 订单状态操作
- [ ] **契约**: `POST /api/orders/:id/{confirm,start,complete,cancel}`
- [ ] 订单状态按钮操作
- [ ] **测试**: 状态流转正常
- [ ] 提交

---

### Phase 14: 统计页面

#### 14.1 仪表盘
- [ ] **契约**: 展示关键指标
- [ ] `src/views/dashboard/Dashboard.vue`
- [ ] **测试**: 数据展示正常
- [ ] 提交

#### 14.2 订单统计
- [ ] **契约**: `GET /api/stats/orders`
- [ ] 订单数量按状态分组
- [ ] **测试**: 统计数据正确
- [ ] 提交

#### 14.3 营收统计
- [ ] **契约**: `GET /api/stats/revenue`
- [ ] 营收金额统计
- [ ] **测试**: 营收数据正确
- [ ] 提交

---

## 十、任务统计

| 分支 | 模块 | Phase | 主任务 | 子任务 |
|------|------|-------|--------|--------|
| phase1-7 | 后端 | 1-7 | 28 | 112 |
| phase8 | 数据库 | 8 | 3 | 12 |
| phase9-14 | 前端 | 9-14 | 10 | 40 |
| phase15-17 | 完善项 | 15-17 | 14 | 56 |
| phase18 | CI/CD | 18 | 3 | 12 |
| phase19 | 开发清单 | 19 | 2 | 8 |
| phase20 | 部署 | 20 | 3 | 12 |
| phase21 | 仪表盘增强 | 21 | 1 | 6 |
| phase22 | 品类管理 | 22 | 6 | 24 |
| phase23 | 上门服务订单 | 23 | 9 | 36 |
| **待开发** | | | | | |
| phase24 | 登录认证 | P0 | 14 | 56 |
| phase27 | 服务人员 | P1 | 14 | 56 |
| phase31 | 系统管理-角色权限 | P1 | 12 | 48 |
| phase25 | 预约时间 | P2 | 13 | 52 |
| phase26 | 服务地址 | P3 | 13 | 52 |
| phase31 | 系统管理-操作日志 | P3 | 4 | 16 |
| phase28 | 评价反馈 | P4 | 13 | 52 |
| phase29 | 商家管理 | P5 | 13 | 52 |
| phase30 | 服务者管理 | P6 | 13 | 52 |
| phase31 | 系统管理-配置项 | P6 | 3 | 12 |
| phase32 | 财务管理 | P5 | 17 | 68 |
| phase33 | 价格策略配置 | P5 | 16 | 64 |
| **总计** | | | **168** | **696** |

---

## 十一、Agent 分工参考

| Agent | 负责模块 | Phase |
|-------|----------|-------|
| Agent-DB | 数据库集成 | 8.1 - 8.3 |
| Agent-Frontend | 前端开发 | 9 - 14 |
| Agent-Test | 单元测试覆盖率 | 15.1 - 15.3 |
| Agent-Error | 错误处理 | 16.1 - 16.2 |
| Agent-Log | 日志记录 | 17.1 - 17.2 |
| Agent-CI | CI/CD 配置 | 18.1 - 18.3 |
| Agent-Checklist | 开发检查清单 | 19.1 - 19.2 |
| Agent-Deploy | 部署文档 | 20.1 - 20.3 |
| Agent-Category | 品类管理 | 22.1 - 22.6 |

---

## 十二、前端技术栈

- Vue 3 (Composition API)
- Element Plus
- Axios
- Vue Router
- Vite

---

## 十三、单元测试覆盖率 (Phase 15)

> 每任务遵循：**契约先行 → 补充测试 → 验证覆盖率 → 提交**

### Phase 15: 单元测试覆盖率提升

#### 15.1 后端单元测试覆盖
- [ ] **契约**: 测试覆盖率目标 70%+
- [ ] `go test ./... -cover` 查看当前覆盖率
- [ ] 补充 handler/service/repository 测试
- [ ] **测试**: 覆盖率达标
- [ ] 提交

#### 15.2 前端组件测试
- [ ] **契约**: Vitest + Vue Test Utils
- [ ] `npm install -D vitest @vue/test-utils`
- [ ] 编写组件基础测试
- [ ] **测试**: `npm run test` 通过
- [ ] 提交

#### 15.3 API 集成测试
- [ ] **契约**: 使用 JSON Schema 验证 API 响应
- [ ] 编写 API 响应格式测试
- [ ] **测试**: E2E API 测试通过
- [ ] 提交

---

## 十四、错误处理完善 (Phase 16)

> 每任务遵循：**契约先行 → 定义错误码 → 实现处理 → 提交**

### Phase 16: 错误处理完善

#### 16.1 统一错误码定义
- [ ] **契约**: 错误码规范
  ```go
  const (
      ErrCodeSuccess = 0
      ErrCodeParam = 400
      ErrCodeUnauthorized = 401
      ErrCodeForbidden = 403
      ErrCodeNotFound = 404
      ErrCodeInternal = 500
  )
  ```
- [ ] 创建 `pkg/errors/errors.go`
- [ ] **测试**: 错误码一致性
- [ ] 提交

#### 16.2 全局错误处理中间件
- [ ] **契约**: Gin 中间件统一捕获 panic
- [ ] 创建 `internal/middleware/recovery.go`
- [ ] 统一日志记录错误
- [ ] **测试**: 模拟 panic 被捕获
- [ ] 提交

---

## 十五、日志记录 (Phase 17)

> 每任务遵循：**契约先行 → 定义日志格式 → 实现记录 → 提交**

### Phase 17: 日志记录

#### 17.1 日志库集成
- [ ] **契约**: 结构化日志 (JSON 格式)
- [ ] `go get go.uber.org/zap`
- [ ] 创建 `pkg/logger/logger.go`
- [ ] **测试**: 日志输出正常
- [ ] 提交

#### 17.2 日志中间件
- [ ] **契约**: HTTP 请求日志
- [ ] 创建 `internal/middleware/logger.go`
- [ ] 记录请求耗时、状态码、路径
- [ ] **测试**: 请求日志正确输出
- [ ] 提交

---

## 十六、CI/CD 配置 (Phase 18)

> CI 包含：依赖安装、编译检查、单元测试、Lint、E2E 测试、预提交汇总

### Phase 18: CI/CD 配置

#### 18.1 CI workflow 配置
- [ ] **契约**: `.github/workflows/ci.yml`
- [ ] Jobs: deps, build, test-unit, lint, test-e2e, pre-commit
- [ ] **测试**: push PR 时 CI 通过
- [ ] 提交

#### 18.2 Lint 门禁
- [ ] **契约**: golangci-lint + ESLint
- [ ] 本地无 golangci-lint 时降级到 gofmt
- [ ] **测试**: CI lint job 通过
- [ ] 提交

#### 18.3 测试缓存和报告
- [ ] **契约**: 测试结果和 E2E 报告上传到 artifacts
- [ ] node_modules 缓存
- [ ] **测试**: 报告可下载
- [ ] 提交

---

## 十七、开发检查清单 (每次提交前必查)

### 提交前检查项
- [ ] `go build ./...` 通过
- [ ] `go test ./...` 通过
- [ ] `npm run build` 通过
- [ ] `npm test` 通过
- [ ] API 字段名与文档一致
- [ ] git status 无意外修改

### 本地 lint 检查
```bash
# Go lint (无 golangci-lint 时用 gofmt)
make lint
# 或
cd oms && gofmt -l .

# Frontend lint
cd frontend && npm run lint
```

---

## 十八、历史错误记录

> 参照 `memory/oms-common-errors.md` 避免重复犯错

常见错误类型：
1. **API 响应解析**: 前端用 `res.users` 但后端返回 `res.data.users`
2. **字段名不匹配**: `total_price` vs `total_amount`
3. **缺少事务**: 订单创建和库存扣减未在同一事务
4. **RowsAffected 未检查**: 并发下库存可能变负
5. **cancelled 订单计入营收**: 需要跳过 cancelled 状态

---

## 十九、部署文档 (Phase 20)

> 每任务遵循：**契约先行 → 编写文档 → 验证部署 → 提交**

### Phase 20: 部署文档

#### 20.1 Docker 配置
- [ ] **契约**: Dockerfile + docker-compose.yml
- [ ] 创建 `Dockerfile` (Go 构建)
- [ ] 创建 `docker-compose.yml` (Go + Vue + SQLite)
- [ ] **测试**: `docker-compose up` 启动成功
- [ ] 提交

#### 20.2 环境变量配置
- [ ] **契约**: `.env.example` 模板
- [ ] 创建 `oms/.env.example`
- [ ] 创建 `frontend/.env.example`
- [ ] **测试**: 按模板配置可正常运行
- [ ] 提交

#### 20.3 README 部署说明
- [ ] **契约**: 部署步骤清晰、可执行
- [ ] 编写 `docs/DEPLOYMENT.md`
- [ ] 包含：环境要求、本地部署、Docker 部署
- [ ] **测试**: 文档可读性检查
- [ ] 提交

---

## 二十一、仪表盘增强 (Phase 21)

> 不需要后端 API，纯前端功能

### Phase 21: 仪表盘添加当前时间显示

#### 21.1 添加日期时间显示
- [ ] **契约**: 纯前端展示，不需要 API
- [ ] 前端: 在 `Dashboard.vue` 顶部添加时间卡片
- [ ] 样式: 36px 大字体，视觉清晰
- [ ] 功能: 每秒更新，显示年/月/日 时:分:秒
- [ ] **验证**: `make pre-commit` 通过
- [ ] 提交

---

## 二十二、品类管理 (Phase 22)

> 产品管理升级为品类管理，一个品类可包含多个产品
> 例如：品类"家政"下有"4小时保洁"、"深度清洗"等产品

### Phase 22: 品类管理模块

#### 22.1 后端 - 品类 Model
- [ ] **契约**: `Category` 结构体
- [ ] 创建 `internal/model/category.go`
- [ ] **字段**: id, name, description, created_at, updated_at
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 22.2 后端 - 品类 CRUD
- [ ] **契约**: repository/service/handler/router
- [ ] 创建 `internal/repository/category.go`
- [ ] 创建 `internal/service/category.go`
- [ ] 创建 `internal/handler/category.go`
- [ ] 添加路由 `/api/categories`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 22.3 后端 - 产品关联品类
- [ ] **契约**: products 表添加 category_id 外键
- [ ] 修改 `internal/model/product.go` 添加 CategoryID
- [ ] 修改 API 契约文档
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 22.4 前端 - 品类管理页面
- [ ] **契约**: 品类 CRUD 页面
- [ ] 创建 `src/views/category/CategoryList.vue`
- [ ] 创建 `src/views/category/CategoryForm.vue`
- [ ] 创建 `src/api/category.js`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 22.5 前端 - 产品管理关联品类
- [ ] **契约**: 产品列表/表单显示所属品类
- [ ] 修改 `ProductList.vue` 和 `ProductForm.vue`
- [ ] 修改 `src/api/product.js`
- [ ] **测试**: `npm test` + `npx playwright test`
- [ ] 提交

#### 22.6 E2E 测试
- [ ] **契约**: 品类和产品 E2E 测试
- [ ] 添加 `e2e/category.spec.js`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

---

## 二十三、订单状态改为上门服务模式 (Phase 23) ✅ 完成

> 服务行业（家政、上门保洁等）订单状态不同于电商

### Phase 23: 上门服务订单状态 ✅

#### 23.1 后端 - 修改订单状态常量 ✅
- [x] **契约**: 新状态 pending/confirmed/in_service/completed/cancelled
- [x] 修改 `internal/model/order.go` - OrderStatus 常量
- [x] **测试**: `go build ./...`
- [x] 提交

#### 23.2 后端 - 修改状态流转逻辑 ✅
- [x] **契约**: pending → confirmed → in_service → completed
- [x] 修改 `internal/service/order.go` - 状态转换方法
- [x] 移除 paid/shipped 相关方法，改用 confirmed/in_service
- [x] **测试**: `go test ./...`
- [x] 提交

#### 23.3 后端 - 修改 Handler 路由 ✅
- [x] **契约**: /orders/:id/confirm, /orders/:id/start, /orders/:id/complete
- [x] 修改 `internal/handler/order.go`
- [x] **测试**: `go test ./...`
- [x] 提交

#### 23.4 后端 - 修改统计 API ✅
- [x] **契约**: 统计按新状态分组
- [x] 修改 `internal/handler/stats.go` 和 `internal/service/order.go`
- [x] **测试**: `go test ./...`
- [x] 提交

#### 23.5 前端 - 修改订单列表状态显示 ✅
- [x] **契约**: 新状态映射
- [x] pending→待确认, confirmed→已确认, in_service→服务中, completed→已完成, cancelled→已取消
- [x] 修改 `OrderList.vue`
- [x] **测试**: `npm test`
- [x] 提交

#### 23.6 前端 - 修改订单表单状态操作 ✅
- [x] **契约**: 新状态按钮操作
- [x] 修改 `OrderForm.vue` - 确认、开始服务、完成操作
- [x] **测试**: `npm test` + `npx playwright test`
- [x] 提交

#### 23.7 API 契约文档更新 ✅
- [x] **契约**: 更新 `docs/api-contract.md`
- [x] 更新订单状态说明和状态流转图
- [x] 提交

#### 23.8 E2E 测试 ✅
- [x] **契约**: 新状态 E2E 测试
- [x] 修改 `e2e/oms.spec.js` 订单相关测试
- [x] **测试**: `npx playwright test`
- [x] 提交

#### 23.9 验证 ✅
- [x] `make pre-commit` 通过
- [x] 所有单元测试通过
- [x] 所有 E2E 测试通过

---

## 二十四、登录认证 (Phase 24)

> 用户登录、Token 验证

### Phase 24: 登录认证模块

#### 24.1 后端 - 用户 Model 添加密码字段
- [ ] **契约**: users 表添加 password 字段，密码需 bcrypt 加密存储
- [ ] 修改 `internal/model/user.go` 添加 Password 字段
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 24.2 后端 - 用户 Repository 密码验证
- [ ] **契约**: FindByUsername 方法需返回 password 字段
- [ ] 修改 `internal/repository/user.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 24.3 后端 - Auth Service 登录逻辑
- [ ] **契约**: 登录成功返回 JWT Token，失败返回错误
- [ ] 创建 `internal/service/auth.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 24.4 后端 - Auth Handler 登录接口
- [ ] **契约**: `POST /api/auth/login` → `{username, password}` → `{token}`
- [ ] 创建 `internal/handler/auth.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 24.5 后端 - Auth Handler 注册接口
- [ ] **契约**: `POST /api/auth/register` → `{username, password, email, phone}` → `{user}`
- [ ] 添加注册方法到 auth handler
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 24.6 后端 - JWT 中间件
- [ ] **契约**: 请求头 `Authorization: Bearer <token>` 验证
- [ ] 创建 `internal/middleware/auth.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 24.7 后端 - 路由注册
- [ ] **契约**: 路由组 `/api/auth`
- [ ] 在 router 中注册 auth 路由
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 24.8 前端 - Auth API 封装
- [ ] **契约**: `POST /api/auth/login`, `POST /api/auth/register`
- [ ] 创建 `src/api/auth.js`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 24.9 前端 - 登录页面
- [ ] **契约**: 用户名密码登录表单
- [ ] 创建 `src/views/auth/Login.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 24.10 前端 - 注册页面
- [ ] **契约**: 用户注册表单
- [ ] 创建 `src/views/auth/Register.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 24.11 前端 - 路由守卫
- [ ] **契约**: 未登录跳转到登录页
- [ ] 创建 `src/router/guards.js`
- [ ] 修改路由配置添加导航守卫
- [ ] **测试**: `npm run build`
- [ ] 提交

#### 24.12 前端 - Token 存储
- [ ] **契约**: Token 存储在 localStorage
- [ ] 封装 auth utility: getToken, setToken, removeToken
- [ ] Axios 请求拦截器添加 Token
- [ ] **测试**: `npm test`
- [ ] 提交

#### 24.13 E2E 测试
- [ ] **契约**: 登录注册 E2E 测试
- [ ] 添加 `e2e/auth.spec.js`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 24.14 验证
- [ ] `make pre-commit` 通过
- [ ] 所有单元测试通过
- [ ] 所有 E2E 测试通过

---

## 二十五、预约时间 (Phase 25)

> 客户选择上门服务时间段

### Phase 25: 预约时间模块

#### 25.1 后端 - Order Model 添加预约时间字段
- [ ] **契约**: orders 表添加 appointment_time 字段 (datetime)
- [ ] 修改 `internal/model/order.go`
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 25.2 后端 - 数据库迁移
- [ ] **契约**: AutoMigrate 添加 appointment_time 列
- [ ] 修改 database.go 或创建 migration
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 25.3 后端 - 更新创建订单接口
- [ ] **契约**: `POST /api/orders` 请求体添加 appointment_time
- [ ] 修改 CreateOrder 处理函数
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 25.4 后端 - 可用时间段 Service
- [ ] **契约**: 根据日期返回可用时间段 (9:00-18:00，每小时一段)
- [ ] 创建 `internal/service/slot.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 25.5 后端 - 可用时间段 Handler
- [ ] **契约**: `GET /api/slots?date=2026-05-26` → 可用时间段列表
- [ ] 创建 `internal/handler/slot.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 25.6 后端 - 路由注册
- [ ] **契约**: 路由组 `/api/slots`
- [ ] 在 router 中注册 slot 路由
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 25.7 前端 - Slot API 封装
- [ ] **契约**: `GET /api/slots?date=xxx`
- [ ] 创建 `src/api/slot.js`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 25.8 前端 - 预约时间选择组件
- [ ] **契约**: 日期选择 + 时段选择
- [ ] 创建 `src/components/TimeSlotPicker.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 25.9 前端 - 订单表单集成预约时间
- [ ] **契约**: OrderForm.vue 添加预约时间选择
- [ ] 修改 `src/views/order/OrderForm.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 25.10 前端 - 订单列表显示预约时间
- [ ] **契约**: OrderList.vue 列显示 appointment_time
- [ ] 修改 `src/views/order/OrderList.vue`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 25.11 前端 - 订单详情弹窗显示预约时间
- [ ] **契约**: 订单详情弹窗显示预约时间
- [ ] 修改 OrderDetail 弹窗组件
- [ ] **测试**: `npm run build`
- [ ] 提交

#### 25.12 E2E 测试
- [ ] **契约**: 预约时间 E2E 测试
- [ ] 添加 `e2e/slot.spec.js`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 25.13 验证
- [ ] `make pre-commit` 通过
- [ ] 所有单元测试通过
- [ ] 所有 E2E 测试通过

---

## 二十六、服务地址 (Phase 26)

> 客户地址管理

### Phase 26: 服务地址模块

#### 26.1 后端 - Address Model
- [ ] **契约**: `Address` 结构体 (id, user_id, name, phone, province, city, district, detail, created_at, updated_at)
- [ ] 创建 `internal/model/address.go`
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 26.2 后端 - Address Repository CRUD
- [ ] **契约**: Create, GetByID, Update, Delete, ListByUserID
- [ ] 创建 `internal/repository/address.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 26.3 后端 - Address Service
- [ ] **契约**: 业务逻辑：用户地址列表、默认地址设置
- [ ] 创建 `internal/service/address.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 26.4 后端 - Address Handler
- [ ] **契约**:
  ```
  GET    /api/addresses           # 用户地址列表
  GET    /api/addresses/:id       # 地址详情
  POST   /api/addresses           # 新增地址
  PUT    /api/addresses/:id       # 更新地址
  DELETE /api/addresses/:id       # 删除地址
  PUT    /api/addresses/:id/default # 设为默认地址
  ```
- [ ] 创建 `internal/handler/address.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 26.5 后端 - 路由注册
- [ ] **契约**: 路由组 `/api/addresses`
- [ ] 在 router 中注册 address 路由
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 26.6 前端 - Address API 封装
- [ ] **契约**: CRUD API
- [ ] 创建 `src/api/address.js`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 26.7 前端 - 地址列表页
- [ ] **契约**: 地址列表展示
- [ ] 创建 `src/views/address/AddressList.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 26.8 前端 - 地址表单弹窗
- [ ] **契约**: 新增/编辑地址表单
- [ ] 创建 `src/views/address/AddressForm.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 26.9 前端 - 地址管理入口
- [ ] **契约**: 用户菜单添加地址管理入口
- [ ] 在 Layout.vue 或用户菜单添加地址管理链接
- [ ] **测试**: `npm run build`
- [ ] 提交

#### 26.10 前端 - 订单表单地址选择
- [ ] **契约**: OrderForm.vue 添加服务地址选择
- [ ] 修改 `src/views/order/OrderForm.vue`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 26.11 前端 - 地址回退到输入
- [ ] **契约**: 无地址时支持手动输入地址
- [ ] OrderForm.vue 添加手动输入模式
- [ ] **测试**: `npm test`
- [ ] 提交

#### 26.12 E2E 测试
- [ ] **契约**: 地址管理 E2E 测试
- [ ] 添加 `e2e/address.spec.js`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 26.13 验证
- [ ] `make pre-commit` 通过
- [ ] 所有单元测试通过
- [ ] 所有 E2E 测试通过

---

## 二十七、服务人员 (Phase 27)

> 服务人员列表、状态、分配

### Phase 27: 服务人员模块

#### 27.1 后端 - Staff Model
- [ ] **契约**: `Staff` 结构体 (id, name, phone, avatar, status, created_at, updated_at)
- [ ] 创建 `internal/model/staff.go`
- [ ] status: available(空闲), busy(忙碌), off(休息)
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 27.2 后端 - Staff Repository CRUD
- [ ] **契约**: Create, GetByID, Update, Delete, List, ListAvailable
- [ ] 创建 `internal/repository/staff.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 27.3 后端 - Staff Service
- [ ] **契约**: 业务逻辑：人员列表、可分配人员筛选
- [ ] 创建 `internal/service/staff.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 27.4 后端 - Staff Handler
- [ ] **契约**:
  ```
  GET    /api/staff              # 人员列表
  GET    /api/staff/:id          # 人员详情
  POST   /api/staff              # 新增人员
  PUT    /api/staff/:id          # 更新人员
  DELETE /api/staff/:id          # 删除人员
  PUT    /api/staff/:id/status   # 更新状态 (available/busy/off)
  ```
- [ ] 创建 `internal/handler/staff.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 27.5 后端 - 路由注册
- [ ] **契约**: 路由组 `/api/staff`
- [ ] 在 router 中注册 staff 路由
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 27.6 后端 - Order Model 添加 StaffID
- [ ] **契约**: orders 表添加 staff_id 字段 (bigint, nullable)
- [ ] 修改 `internal/model/order.go`
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 27.7 后端 - 订单分配人员
- [ ] **契约**: `PUT /api/orders/:id` 可更新 staff_id
- [ ] 修改 order handler 的 UpdateOrder
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 27.8 前端 - Staff API 封装
- [ ] **契约**: CRUD + 状态更新 API
- [ ] 创建 `src/api/staff.js`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 27.9 前端 - 人员列表页
- [ ] **契约**: 人员列表展示 + 状态筛选
- [ ] 创建 `src/views/staff/StaffList.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 27.10 前端 - 人员表单弹窗
- [ ] **契约**: 新增/编辑人员表单
- [ ] 创建 `src/views/staff/StaffForm.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 27.11 前端 - 订单分配人员选择
- [ ] **契约**: OrderForm.vue 或订单详情选择服务人员
- [ ] 修改 `src/views/order/OrderDetail.vue` 添加人员选择
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 27.12 前端 - 人员状态管理
- [ ] **契约**: 人员状态切换 (空闲/忙碌/休息)
- [ ] StaffList.vue 添加状态切换按钮
- [ ] **测试**: `npm run build`
- [ ] 提交

#### 27.13 E2E 测试
- [ ] **契约**: 服务人员 E2E 测试
- [ ] 添加 `e2e/staff.spec.js`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 27.14 验证
- [ ] `make pre-commit` 通过
- [ ] 所有单元测试通过
- [ ] 所有 E2E 测试通过

---

## 二十八、评价反馈 (Phase 28)

> 服务完成后客户评价

### Phase 28: 评价反馈模块

#### 28.1 后端 - Review Model
- [ ] **契约**: `Review` 结构体 (id, order_id, user_id, staff_id, rating, comment, created_at)
- [ ] 创建 `internal/model/review.go`
- [ ] rating: 1-5 评分
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 28.2 后端 - Review Repository
- [ ] **契约**: Create, GetByID, GetByOrderID, ListByStaffID
- [ ] 创建 `internal/repository/review.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 28.3 后端 - Review Service
- [ ] **契约**: 业务逻辑：提交评价、查询评价
- [ ] 创建 `internal/service/review.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 28.4 后端 - Review Handler
- [ ] **契约**:
  ```
  POST   /api/reviews              # 提交评价
  GET    /api/reviews/:id         # 评价详情
  GET    /api/reviews?order_id=x  # 订单评价查询
  GET    /api/reviews?staff_id=x  # 服务人员评价查询
  ```
- [ ] 创建 `internal/handler/review.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 28.5 后端 - 路由注册
- [ ] **契约**: 路由组 `/api/reviews`
- [ ] 在 router 中注册 review 路由
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 28.6 后端 - 订单完成后可评价
- [ ] **契约**: completed 状态订单允许创建评价
- [ ] 修改 review service 添加订单状态校验
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 28.7 前端 - Review API 封装
- [ ] **契约**: 评价 CRUD API
- [ ] 创建 `src/api/review.js`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 28.8 前端 - 评价表单弹窗
- [ ] **契约**: 星级评分 + 文字评论
- [ ] 创建 `src/views/review/ReviewForm.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 28.9 前端 - 订单完成入口评价
- [ ] **契约**: OrderList.vue 完成状态订单添加"评价"按钮
- [ ] 修改 `src/views/order/OrderList.vue`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 28.10 前端 - 评价列表展示
- [ ] **契约**: 评价列表展示 (可按服务人员筛选)
- [ ] 创建 `src/views/review/ReviewList.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 28.11 前端 - 服务人员评分展示
- [ ] **契约**: StaffList.vue 显示平均评分
- [ ] 修改 `src/views/staff/StaffList.vue`
- [ ] **测试**: `npm run build`
- [ ] 提交

#### 28.12 E2E 测试
- [ ] **契约**: 评价功能 E2E 测试
- [ ] 添加 `e2e/review.spec.js`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 28.13 验证
- [ ] `make pre-commit` 通过
- [ ] 所有单元测试通过
- [ ] 所有 E2E 测试通过

---

## 二十九、商家/商户管理 (Phase 29)

> 商家/商户（平台入驻商家）管理

### Phase 29: 商家管理模块

#### 29.1 后端 - Merchant Model
- [ ] **契约**: `Merchant` 结构体 (id, name, logo, phone, province, city, district, detail, status, created_at, updated_at)
- [ ] 创建 `internal/model/merchant.go`
- [ ] status: pending(待审核), approved(已审核), rejected(已拒绝), suspended(已停用)
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 29.2 后端 - Merchant Repository CRUD
- [ ] **契约**: Create, GetByID, Update, Delete, List, ListApproved
- [ ] 创建 `internal/repository/merchant.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 29.3 后端 - Merchant Service
- [ ] **契约**: 业务逻辑：商家入驻审核通过/拒绝/停用
- [ ] 创建 `internal/service/merchant.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 29.4 后端 - Merchant Handler
- [ ] **契约**:
  ```
  GET    /api/merchants              # 商家列表
  GET    /api/merchants/:id          # 商家详情
  POST   /api/merchants              # 商家入驻申请
  PUT    /api/merchants/:id          # 更新商家信息
  DELETE /api/merchants/:id          # 删除商家
  PUT    /api/merchants/:id/approve  # 审核通过
  PUT    /api/merchants/:id/reject   # 审核拒绝
  PUT    /api/merchants/:id/suspend  # 停用商家
  ```
- [ ] 创建 `internal/handler/merchant.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 29.5 后端 - 路由注册
- [ ] **契约**: 路由组 `/api/merchants`
- [ ] 在 router 中注册 merchant 路由
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 29.6 后端 - Product 关联 Merchant
- [ ] **契约**: products 表添加 merchant_id 字段
- [ ] 修改 `internal/model/product.go`
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 29.7 前端 - Merchant API 封装
- [ ] **契约**: CRUD + 状态更新 API
- [ ] 创建 `src/api/merchant.js`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 29.8 前端 - 商家列表页
- [ ] **契约**: 商家列表展示 + 状态筛选
- [ ] 创建 `src/views/merchant/MerchantList.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 29.9 前端 - 商家表单弹窗
- [ ] **契约**: 新增/编辑商家表单
- [ ] 创建 `src/views/merchant/MerchantForm.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 29.10 前端 - 商家审核操作
- [ ] **契约**: MerchantList.vue 添加审核/停用按钮
- [ ] 修改 `src/views/merchant/MerchantList.vue`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 29.11 前端 - 产品关联商家
- [ ] **契约**: ProductForm.vue 添加商家选择
- [ ] 修改 `src/views/product/ProductForm.vue`
- [ ] **测试**: `npm run build`
- [ ] 提交

#### 29.12 E2E 测试
- [ ] **契约**: 商家管理 E2E 测试
- [ ] 添加 `e2e/merchant.spec.js`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 29.13 验证
- [ ] `make pre-commit` 通过
- [ ] 所有单元测试通过
- [ ] 所有 E2E 测试通过

---

## 三十、服务者管理 (Phase 30)

> 服务者（实际提供服务的员工/自由职业者）管理

### Phase 30: 服务者管理模块

#### 30.1 后端 - Provider Model
- [ ] **契约**: `Provider` 结构体 (id, merchant_id, name, phone, avatar, id_card, status, created_at, updated_at)
- [ ] 创建 `internal/model/provider.go`
- [ ] status: pending(待审核), approved(已审核), rejected(已拒绝), suspended(已停用)
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 30.2 后端 - Provider Repository CRUD
- [ ] **契约**: Create, GetByID, Update, Delete, List, ListByMerchantID, ListApproved
- [ ] 创建 `internal/repository/provider.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 30.3 后端 - Provider Service
- [ ] **契约**: 业务逻辑：服务者入驻审核、服务者列表
- [ ] 创建 `internal/service/provider.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 30.4 后端 - Provider Handler
- [ ] **契约**:
  ```
  GET    /api/providers              # 服务者列表
  GET    /api/providers/:id          # 服务者详情
  POST   /api/providers              # 服务者入驻申请
  PUT    /api/providers/:id          # 更新服务者信息
  DELETE /api/providers/:id          # 删除服务者
  PUT    /api/providers/:id/approve  # 审核通过
  PUT    /api/providers/:id/reject   # 审核拒绝
  PUT    /api/providers/:id/suspend   # 停用服务者
  ```
- [ ] 创建 `internal/handler/provider.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 30.5 后端 - 路由注册
- [ ] **契约**: 路由组 `/api/providers`
- [ ] 在 router 中注册 provider 路由
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 30.6 后端 - Staff 关联 Provider
- [ ] **契约**: staff 表添加 provider_id 字段
- [ ] 修改 `internal/model/staff.go`
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 30.7 前端 - Provider API 封装
- [ ] **契约**: CRUD + 状态更新 API
- [ ] 创建 `src/api/provider.js`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 30.8 前端 - 服务者列表页
- [ ] **契约**: 服务者列表展示 + 状态筛选 + 所属商家筛选
- [ ] 创建 `src/views/provider/ProviderList.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 30.9 前端 - 服务者表单弹窗
- [ ] **契约**: 新增/编辑服务者表单
- [ ] 创建 `src/views/provider/ProviderForm.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 30.10 前端 - 服务者审核操作
- [ ] **契约**: ProviderList.vue 添加审核/停用按钮
- [ ] 修改 `src/views/provider/ProviderList.vue`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 30.11 前端 - 商家菜单添加工者入口
- [ ] **契约**: 商家详情页显示旗下服务者
- [ ] MerchantList.vue 添加工者管理入口
- [ ] **测试**: `npm run build`
- [ ] 提交

#### 30.12 E2E 测试
- [ ] **契约**: 服务者管理 E2E 测试
- [ ] 添加 `e2e/provider.spec.js`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 30.13 验证
- [ ] `make pre-commit` 通过
- [ ] 所有单元测试通过
- [ ] 所有 E2E 测试通过

---

## 三十一、系统管理 (Phase 31)

> 系统管理：角色、权限、操作日志、系统配置

### Phase 31: 系统管理模块

#### 31.1 后端 - Role Model
- [ ] **契约**: `Role` 结构体 (id, name, code, description, created_at, updated_at)
- [ ] 创建 `internal/model/role.go`
- [ ] 预置角色: super_admin(超级管理员), admin(管理员), operator(运营人员)
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 31.2 后端 - Role Repository CRUD
- [ ] **契约**: Create, GetByID, Update, Delete, List, GetByCode
- [ ] 创建 `internal/repository/role.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 31.3 后端 - Permission Model
- [ ] **契约**: `Permission` 结构体 (id, name, code, group, created_at)
- [ ] 创建 `internal/model/permission.go`
- [ ] 权限码示例: user:read, user:write, order:read, order:write, staff:manage
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 31.4 后端 - RolePermission 关联表
- [ ] **契约**: 角色-权限多对多关联表 role_permissions
- [ ] 创建 `internal/model/role_permission.go`
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 31.5 后端 - Role Service
- [ ] **契约**: 角色CRUD、分配权限、查询角色权限
- [ ] 创建 `internal/service/role.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 31.6 后端 - Role Handler
- [ ] **契约**:
  ```
  GET    /api/roles              # 角色列表
  GET    /api/roles/:id          # 角色详情(含权限)
  POST   /api/roles              # 创建角色
  PUT    /api/roles/:id          # 更新角色
  DELETE /api/roles/:id          # 删除角色
  PUT    /api/roles/:id/permissions # 分配权限
  ```
- [ ] 创建 `internal/handler/role.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 31.7 后端 - 操作日志 Model
- [ ] **契约**: `OperationLog` 结构体 (id, user_id, username, action, resource, details, ip, created_at)
- [ ] 创建 `internal/model/operation_log.go`
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 31.8 后端 - 操作日志 Service + Handler
- [ ] **契约**: 记录日志、查询日志列表
- [ ] 创建 `internal/service/operation_log.go`
- [ ] 创建 `internal/handler/operation_log.go`
- [ ] 添加路由 `/api/operation-logs`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 31.9 后端 - 操作日志中间件
- [ ] **契约**: 关键操作自动记录日志
- [ ] 创建 `internal/middleware/audit.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 31.10 后端 - System Config Model
- [ ] **契约**: `SystemConfig` 结构体 (id, key, value, description, updated_at)
- [ ] 创建 `internal/model/system_config.go`
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 31.11 后端 - System Config Service + Handler
- [ ] **契约**: 系统配置CRUD
- [ ] 创建 `internal/service/system_config.go`
- [ ] 创建 `internal/handler/system_config.go`
- [ ] 添加路由 `/api/system-configs`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 31.12 后端 - 路由注册
- [ ] **契约**: 路由组 `/api/roles`, `/api/operation-logs`, `/api/system-configs`
- [ ] 在 router 中注册系统管理路由
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 31.13 前端 - Role API 封装
- [ ] **契约**: 角色 CRUD API
- [ ] 创建 `src/api/role.js`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 31.14 前端 - 角色管理页面
- [ ] **契约**: 角色列表、新增、编辑、删除、分配权限
- [ ] 创建 `src/views/role/RoleList.vue`
- [ ] 创建 `src/views/role/RoleForm.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 31.15 前端 - 操作日志页面
- [ ] **契约**: 操作日志列表，按用户/时间筛选
- [ ] 创建 `src/views/operation-log/OperationLogList.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 31.16 前端 - 系统配置页面
- [ ] **契约**: 系统配置列表、编辑配置项
- [ ] 创建 `src/views/system-config/SystemConfigList.vue`
- [ ] 创建 `src/views/system-config/SystemConfigForm.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 31.17 前端 - 路由守卫权限控制
- [ ] **契约**: 根据用户角色控制页面访问权限
- [ ] 修改 `src/router/guards.js`
- [ ] **测试**: `npm run build`
- [ ] 提交

#### 31.18 E2E 测试
- [ ] **契约**: 系统管理 E2E 测试
- [ ] 添加 `e2e/role.spec.js`, `e2e/operation-log.spec.js`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 31.19 验证
- [ ] `make pre-commit` 通过
- [ ] 所有单元测试通过
- [ ] 所有 E2E 测试通过

---

## 三十二、财务管理 (Phase 32)

> 财务管理：收款、对账、退款、发票

### Phase 32: 财务管理模块

#### 32.1 后端 - 收款记录 Model
- [ ] **契约**: `Payment` 结构体 (id, order_id, order_no, user_id, amount, payment_method, transaction_no, paid_at, created_at)
- [ ] 创建 `internal/model/payment.go`
- [ ] payment_method: wechat(微信), alipay(支付宝), cash(现金), card(银行卡)
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 32.2 后端 - 收款 Repository
- [ ] **契约**: Create, GetByID, GetByOrderID, ListByUserID, ListByDateRange
- [ ] 创建 `internal/repository/payment.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 32.3 后端 - 收款 Service + Handler
- [ ] **契约**: `POST /api/payments` 记录收款, `GET /api/payments` 收款列表
- [ ] 创建 `internal/service/payment.go`
- [ ] 创建 `internal/handler/payment.go`
- [ ] 添加路由 `/api/payments`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 32.4 后端 - 对账 Model
- [ ] **契约**: `Reconciliation` 结构体 (id, date, total_orders, total_amount, total_paid, total_refund, status, created_at)
- [ ] 创建 `internal/model/reconciliation.go`
- [ ] status: pending(待对账), completed(已完成)
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 32.5 后端 - 对账 Service + Handler
- [ ] **契约**: `GET /api/reconciliations` 对账列表, `POST /api/reconciliations/generate` 生成对账单
- [ ] 创建 `internal/service/reconciliation.go`
- [ ] 创建 `internal/handler/reconciliation.go`
- [ ] 添加路由 `/api/reconciliations`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 32.6 后端 - 退款 Model
- [ ] **契约**: `Refund` 结构体 (id, order_id, order_no, user_id, amount, reason, status, approved_by, processed_at, created_at)
- [ ] 创建 `internal/model/refund.go`
- [ ] status: pending(待审批), approved(已批准), rejected(已拒绝), completed(已完成)
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 32.7 后端 - 退款 Service + Handler
- [ ] **契约**:
  ```
  GET    /api/refunds              # 退款列表
  GET    /api/refunds/:id          # 退款详情
  POST   /api/refunds              # 申请退款
  PUT    /api/refunds/:id/approve  # 批准退款
  PUT    /api/refunds/:id/reject   # 拒绝退款
  PUT    /api/refunds/:id/complete # 完成退款
  ```
- [ ] 创建 `internal/service/refund.go`
- [ ] 创建 `internal/handler/refund.go`
- [ ] 添加路由 `/api/refunds`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 32.8 后端 - 发票 Model
- [ ] **契约**: `Invoice` 结构体 (id, order_id, user_id, title, tax_no, amount, status, issued_at, created_at)
- [ ] 创建 `internal/model/invoice.go`
- [ ] status: pending(待开票), issued(已开票), invalid(已作废)
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 32.9 后端 - 发票 Service + Handler
- [ ] **契约**: `GET /api/invoices` 发票列表, `POST /api/invoices` 申请开票, `PUT /api/invoices/:id/issue` 开票
- [ ] 创建 `internal/service/invoice.go`
- [ ] 创建 `internal/handler/invoice.go`
- [ ] 添加路由 `/api/invoices`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 32.10 后端 - 路由注册
- [ ] **契约**: 路由组 `/api/payments`, `/api/reconciliations`, `/api/refunds`, `/api/invoices`
- [ ] 在 router 中注册财务相关路由
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 32.11 前端 - Payment API 封装
- [ ] **契约**: 收款记录 CRUD
- [ ] 创建 `src/api/payment.js`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 32.12 前端 - 收款记录页面
- [ ] **契约**: 收款列表、按订单/用户/日期筛选
- [ ] 创建 `src/views/payment/PaymentList.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 32.13 前端 - 退款申请页面
- [ ] **契约**: 退款列表、申请退款、审批退款
- [ ] 创建 `src/views/refund/RefundList.vue`
- [ ] 创建 `src/views/refund/RefundForm.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 32.14 前端 - 对账页面
- [ ] **契约**: 对账单列表、按日期筛选
- [ ] 创建 `src/views/reconciliation/ReconciliationList.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 32.15 前端 - 发票页面
- [ ] **契约**: 发票列表、开票申请
- [ ] 创建 `src/views/invoice/InvoiceList.vue`
- [ ] 创建 `src/views/invoice/InvoiceForm.vue`
- [ ] **测试**: `npm run build`
- [ ] 提交

#### 32.16 E2E 测试
- [ ] **契约**: 财务管理 E2E 测试
- [ ] 添加 `e2e/payment.spec.js`, `e2e/refund.spec.js`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 32.17 验证
- [ ] `make pre-commit` 通过
- [ ] 所有单元测试通过
- [ ] 所有 E2E 测试通过

---

## 三十三、价格策略配置 (Phase 33)

> 价格策略：定价规则、折扣活动、优惠券

### Phase 33: 价格策略模块

#### 33.1 后端 - 定价规则 Model
- [ ] **契约**: `PriceRule` 结构体 (id, name, rule_type, product_id, category_id, base_price, unit_price, min_quantity, created_at, updated_at)
- [ ] 创建 `internal/model/price_rule.go`
- [ ] rule_type: fixed(固定价), time-based(按时计费), quantity-based(按量计费)
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 33.2 后端 - 定价规则 Repository CRUD
- [ ] **契约**: Create, GetByID, Update, Delete, List, ListByProductID, ListByCategoryID
- [ ] 创建 `internal/repository/price_rule.go`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 33.3 后端 - 定价规则 Service + Handler
- [ ] **契约**:
  ```
  GET    /api/price-rules              # 规则列表
  GET    /api/price-rules/:id           # 规则详情
  POST   /api/price-rules              # 创建规则
  PUT    /api/price-rules/:id           # 更新规则
  DELETE /api/price-rules/:id          # 删除规则
  GET    /api/price-rules/calculate    # 计算价格 (product_id, quantity, duration)
  ```
- [ ] 创建 `internal/service/price_rule.go`
- [ ] 创建 `internal/handler/price_rule.go`
- [ ] 添加路由 `/api/price-rules`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 33.4 后端 - 折扣活动 Model
- [ ] **契约**: `Discount` 结构体 (id, name, discount_type, discount_value, min_order_amount, start_at, end_at, status, created_at, updated_at)
- [ ] 创建 `internal/model/discount.go`
- [ ] discount_type: percentage(百分比), fixed(固定金额), gift(赠品)
- [ ] status: pending(待生效), active(进行中), expired(已过期), disabled(已禁用)
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 33.5 后端 - 折扣活动 Service + Handler
- [ ] **契约**:
  ```
  GET    /api/discounts              # 活动列表
  GET    /api/discounts/:id          # 活动详情
  POST   /api/discounts              # 创建活动
  PUT    /api/discounts/:id          # 更新活动
  DELETE /api/discounts/:id          # 删除活动
  PUT    /api/discounts/:id/enable   # 启用活动
  PUT    /api/discounts/:id/disable  # 禁用活动
  ```
- [ ] 创建 `internal/service/discount.go`
- [ ] 创建 `internal/handler/discount.go`
- [ ] 添加路由 `/api/discounts`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 33.6 后端 - 优惠券 Model
- [ ] **契约**: `Coupon` 结构体 (id, code, name, discount_type, discount_value, min_order_amount, total_count, remain_count, per_user_limit, start_at, end_at, status, created_at, updated_at)
- [ ] 创建 `internal/model/coupon.go`
- [ ] status: pending, active, expired, disabled
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 33.7 后端 - 优惠券 Service + Handler
- [ ] **契约**:
  ```
  GET    /api/coupons                # 优惠券列表
  GET    /api/coupons/:id            # 优惠券详情
  POST   /api/coupons                # 创建优惠券
  PUT    /api/coupons/:id            # 更新优惠券
  DELETE /api/coupons/:id            # 删除优惠券
  POST   /api/coupons/generate       # 生成优惠券码
  POST   /api/coupons/claim          # 用户领取优惠券
  GET    /api/coupons/my             # 我的优惠券
  POST   /api/coupons/validate      # 校验优惠券
  ```
- [ ] 创建 `internal/service/coupon.go`
- [ ] 创建 `internal/handler/coupon.go`
- [ ] 添加路由 `/api/coupons`
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 33.8 后端 - 价格计算中间件
- [ ] **契约**: 创建订单时自动计算最优价格(定价规则+折扣+优惠券)
- [ ] 修改 order service 添加价格计算逻辑
- [ ] **测试**: `go test ./...`
- [ ] 提交

#### 33.9 后端 - 路由注册
- [ ] **契约**: 路由组 `/api/price-rules`, `/api/discounts`, `/api/coupons`
- [ ] 在 router 中注册价格策略路由
- [ ] **测试**: `go build ./...`
- [ ] 提交

#### 33.10 前端 - PriceRule API 封装
- [ ] **契约**: 定价规则 CRUD
- [ ] 创建 `src/api/price-rule.js`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 33.11 前端 - 定价规则页面
- [ ] **契约**: 规则列表、新增、编辑、删除
- [ ] 创建 `src/views/price-rule/PriceRuleList.vue`
- [ ] 创建 `src/views/price-rule/PriceRuleForm.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 33.12 前端 - 折扣活动页面
- [ ] **契约**: 活动列表、新增、编辑、启用/禁用
- [ ] 创建 `src/views/discount/DiscountList.vue`
- [ ] 创建 `src/views/discount/DiscountForm.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 33.13 前端 - 优惠券页面
- [ ] **契约**: 优惠券列表、生成、领取、校验
- [ ] 创建 `src/views/coupon/CouponList.vue`
- [ ] 创建 `src/views/coupon/CouponForm.vue`
- [ ] **测试**: `npm test`
- [ ] 提交

#### 33.14 前端 - 下单价格计算
- [ ] **契约**: OrderForm.vue 添加价格计算展示(原价、折扣、优惠券、红包后价格)
- [ ] 修改 `src/views/order/OrderForm.vue`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 33.15 E2E 测试
- [ ] **契约**: 价格策略 E2E 测试
- [ ] 添加 `e2e/price-rule.spec.js`, `e2e/coupon.spec.js`
- [ ] **测试**: `npx playwright test`
- [ ] 提交

#### 33.16 验证
- [ ] `make pre-commit` 通过
- [ ] 所有单元测试通过
- [ ] 所有 E2E 测试通过
