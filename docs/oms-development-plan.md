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
    `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/paid/shipped/completed/cancelled',
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
  Paid()      // pending → paid
  Ship()      // paid → shipped
  Complete()  // shipped → completed
  Cancel()    // pending/paid → cancelled
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
  POST   /api/orders/:id/paid
  POST   /api/orders/:id/ship
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
| POST | /api/orders/:id/paid | 支付 |
| POST | /api/orders/:id/ship | 发货 |
| POST | /api/orders/:id/complete | 完成 |
| POST | /api/orders/:id/cancel | 取消 |

### 4.4 统计 API
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/stats/orders | 订单统计 |
| GET | /api/stats/revenue | 营收统计 |

---

## 五、订单状态流转

```
pending (待支付)
    ├── paid (已支付)
    │       ├── shipped (已发货)
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
- [ ] **契约**: `POST /api/orders/:id/{paid,ship,complete,cancel}`
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
| **待开发** | | | | | |
| phase24 | 登录认证 | 24 | 14 | 56 |
| phase25 | 预约时间 | 25 | 12 | 48 |
| phase26 | 服务地址 | 26 | 12 | 48 |
| phase27 | 服务人员 | 27 | 13 | 52 |
| phase28 | 评价反馈 | 28 | 13 | 52 |
| phase29 | 商家管理 | 29 | 13 | 52 |
| phase30 | 服务者管理 | 30 | 13 | 52 |
| phase31 | 系统管理 | 31 | 21 | 84 |
| phase32 | 财务管理 | 32 | 21 | 84 |
| phase33 | 价格策略 | 33 | 19 | 76 |
| **总计** | | | **233** | **936** |

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

## 三十四、工期估算

> 基于 Phase 22-23 开发速度估算（6个子任务约3-4天，9个子任务约4-5天）
> 单个子任务（含契约/测试/实现/提交）平均 0.5-1 天

### 开发人数与工期对照

| 开发人数 | 总工期(工作日) | 总工期(月) |
|----------|---------------|-----------|
| 1人 | 90-120 天 | 4-5 个月 |
| 2人 | 45-60 天 | 2-3 个月 |
| 3人 | 30-40 天 | 1.5-2 个月 |

### 各 Phase 工期明细

| Phase | 模块 | 子任务 | 单人工期(天) | 可并行 |
|-------|------|--------|-------------|--------|
| 24 | 登录认证 | 14 | 7-10 | - |
| 27 | 服务人员 | 14 | 7-10 | 与24 |
| 31角色 | 系统管理-角色权限 | 12 | 6-8 | 与27 |
| 25 | 预约时间 | 13 | 6-8 | - |
| 26 | 服务地址 | 13 | 6-8 | - |
| 31日志 | 系统管理-操作日志 | 4 | 2-3 | 与26 |
| 28 | 评价反馈 | 13 | 6-8 | - |
| 29 | 商家管理 | 13 | 6-8 | - |
| 32 | 财务管理 | 17 | 8-12 | - |
| 33 | 价格策略配置 | 16 | 8-12 | - |
| 30 | 服务者管理 | 13 | 6-8 | - |
| 31配置 | 系统管理-配置项 | 3 | 1-2 | - |

### 推荐开发顺序

```
P0: Phase 24 (登录认证)          → 7-10天
    ↓
P1: Phase 27 (服务人员)          → 7-10天
    + Phase 31角色权限 (可并行)    → 6-8天
    ↓
P2: Phase 25 (预约时间)           → 6-8天
    ↓
P3: Phase 26 (服务地址)          → 6-8天
    + Phase 31操作日志 (可并行)    → 2-3天
    ↓
P4: Phase 28 (评价反馈)           → 6-8天
    ↓
P5: Phase 29 (商家管理)          → 6-8天
    → Phase 32 (财务管理)         → 8-12天
    → Phase 33 (价格策略配置)     → 8-12天
    ↓
P6: Phase 30 (服务者管理)        → 6-8天
    → Phase 31配置项 (可并行)     → 1-2天
```

### 里程碑检查点

| 里程碑 | 完成标志 | 预计时间 |
|--------|----------|----------|
| M1: 基础用户体系 | Phase 24 完成 | 第1个月 |
| M2: 核心业务流程 | Phase 25-27 完成 | 第2个月 |
| M3: 增强功能 | Phase 26,28,31日志完成 | 第3个月 |
| M4: 商家运营 | Phase 29,32,33 完成 | 第4个月 |
| M5: 完整系统 | Phase 30,31配置项完成 | 第5个月 |

---

## 三十五、待开发模块详细任务拆分

### Phase 24 登录认证 - 详细任务

#### 后端 Model (24.1)
- [x] 24.1.1 **契约**: User Model 添加 Password 字段 (bcrypt加密)
- [x] 24.1.2 创建 `internal/model/user.go` - 添加 Password, Salt 字段
- [x] 24.1.3 **测试**: `go build ./...`
- [x] 24.1.4 提交

#### 后端 Repository (24.2)
- [x] 24.2.1 **契约**: FindByUsername 返回 password 字段用于验证
- [x] 24.2.2 修改 `internal/repository/user.go` - FindByUsername 方法
- [x] 24.2.3 **测试**: `go test ./...`
- [x] 24.2.4 提交

#### 后端 Auth Service (24.3)
- [x] 24.3.1 **契约**: Login(username, password) → (token, error)
- [x] 24.3.2 **契约**: Register(username, password, email, phone) → (user, error)
- [x] 24.3.3 创建 `internal/service/auth.go`
- [x] 24.3.4 实现密码 bcrypt 校验
- [x] 24.3.5 实现 JWT Token 生成
- [x] 24.3.6 **测试**: `go test ./...`
- [x] 24.3.7 提交

#### 后端 Auth Handler (24.4)
- [x] 24.4.1 **契约**: `POST /api/auth/login` → `{username, password}` → `{token, user}`
- [x] 24.4.2 **契约**: `POST /api/auth/register` → `{username, password, email, phone}` → `{user}`
- [x] 24.4.3 **契约**: `GET /api/auth/me` → 获取当前用户信息
- [x] 24.4.4 创建 `internal/handler/auth.go`
- [x] 24.4.5 实现登录处理函数
- [x] 24.4.6 实现注册处理函数
- [x] 24.4.7 实现获取当前用户处理函数
- [x] 24.4.8 **测试**: `go test ./...`
- [x] 24.4.9 提交

#### 后端 JWT 中间件 (24.5)
- [x] 24.5.1 **契约**: 请求头 `Authorization: Bearer <token>` 验证
- [x] 24.5.2 **契约**: 验证失败返回 401 Unauthorized
- [x] 24.5.3 创建 `internal/middleware/auth.go`
- [x] 24.5.4 实现 JWT 解析和验证
- [x] 24.5.5 将用户信息注入 context
- [x] 24.5.6 **测试**: `go test ./...`
- [x] 24.5.7 提交

#### 后端 路由注册 (24.6)
- [x] 24.6.1 **契约**: 路由组 `/api/auth`
- [x] 24.6.2 在 router 中注册 auth 路由
- [x] 24.6.3 添加需要认证的路由中间件配置
- [x] 24.6.4 **测试**: `go build ./...`
- [x] 24.6.5 提交

#### 前端 API (24.7)
- [x] 24.7.1 **契约**: `POST /api/auth/login` 登录接口
- [x] 24.7.2 **契约**: `POST /api/auth/register` 注册接口
- [x] 24.7.3 **契约**: `GET /api/auth/me` 获取当前用户
- [x] 24.7.4 创建 `src/api/auth.js`
- [x] 24.7.5 封装 login, register, getCurrentUser 方法
- [x] 24.7.6 **测试**: `npm test`
- [x] 24.7.7 提交

#### 前端 登录页 (24.8)
- [x] 24.8.1 **契约**: 用户名密码登录表单
- [x] 24.8.2 **契约**: 登录失败显示错误提示
- [x] 24.8.3 **契约**: 登录成功跳转首页
- [x] 24.8.4 创建 `src/views/auth/Login.vue`
- [x] 24.8.5 实现表单验证
- [x] 24.8.6 实现登录请求和 Token 存储
- [x] 24.8.7 **测试**: `npm test`
- [x] 24.8.8 提交

#### 前端 注册页 (24.9)
- [x] 24.9.1 **契约**: 用户注册表单 (username, password, email, phone)
- [x] 24.9.2 **契约**: 密码确认校验
- [x] 24.9.3 **契约**: 注册成功跳转登录页
- [x] 24.9.4 创建 `src/views/auth/Register.vue`
- [x] 24.9.5 实现表单验证
- [x] 24.9.6 **测试**: `npm test`
- [x] 24.9.7 提交

#### 前端 Token 存储 (24.10)
- [x] 24.10.1 **契约**: Token 存储在 localStorage
- [x] 24.10.2 创建 `src/api/index.js` - getToken, setToken, removeToken (在 api/index 中实现)
- [x] 24.10.3 修改 Axios 拦截器添加 Token
- [x] 24.10.4 ~~实现请求重试和 Token 刷新逻辑~~ (未实现：401 直接跳登录页，无 refresh)
- [x] 24.10.5 **测试**: `npm test`
- [x] 24.10.6 提交

#### 前端 路由守卫 (24.11)
- [x] 24.11.1 **契约**: 未登录访问需认证页面跳转到登录页
- [x] 24.11.2 **契约**: 已登录访问登录页跳转到首页
- [x] 24.11.3 创建 `src/router/guards.js`
- [x] 24.11.4 修改路由配置添加导航守卫
- [x] 24.11.5 **测试**: `npm run build`
- [x] 24.11.6 提交

#### 前端 登出功能 (24.12)
- [x] 24.12.1 **契约**: 清除 Token 和用户信息
- [x] 24.12.2 在 Layout 添加退出按钮
- [x] 24.12.3 **测试**: `npx playwright test`
- [x] 24.12.4 提交

#### E2E 测试 (24.13)
- [x] 24.13.1 **契约**: 登录/注册/登出 E2E 测试
- [x] 24.13.2 添加 `e2e/auth.spec.js`
- [x] 24.13.3 测试正常登录流程
- [x] 24.13.4 测试错误密码登录
- [x] 24.13.5 测试注册新用户
- [x] 24.13.6 **测试**: `npx playwright test`
- [x] 24.13.7 提交

#### 验证 (24.14)
- [x] 24.14.1 `go build ./...` 通过
- [x] 24.14.2 `go test ./...` 通过
- [x] 24.14.3 `npm run build` 通过

### Phase 27 服务人员 - 详细任务

#### 后端 Model (27.1)
- [x] 27.1.1 **契约**: `Staff` 结构体 (id, name, phone, avatar, status, created_at, updated_at)
- [x] 27.1.2 **契约**: status: available(空闲), busy(忙碌), off(休息)
- [x] 27.1.3 创建 `internal/model/staff.go`
- [x] 27.1.4 **测试**: `go build ./...`
- [x] 27.1.5 提交

#### 后端 Repository (27.2)
- [x] 27.2.1 **契约**: Create, GetByID, Update, Delete, List, ListAvailable
- [x] 27.2.2 创建 `internal/repository/staff.go`
- [x] 27.2.3 **测试**: `go test ./...`
- [x] 27.2.4 提交

#### 后端 Service (27.3)
- [x] 27.3.1 **契约**: 人员列表、可分配人员筛选
- [x] 27.3.2 创建 `internal/service/staff.go`
- [x] 27.3.3 **测试**: `go test ./...`
- [x] 27.3.4 提交

#### 后端 Handler (27.4)
- [x] 27.4.1 **契约**: `GET /api/staff` 人员列表
- [x] 27.4.2 **契约**: `GET /api/staff/:id` 人员详情
- [x] 27.4.3 **契约**: `POST /api/staff` 新增人员
- [x] 27.4.4 **契约**: `PUT /api/staff/:id` 更新人员
- [x] 27.4.5 **契约**: `DELETE /api/staff/:id` 删除人员
- [x] 27.4.6 **契约**: `PUT /api/staff/:id/status` 更新状态
- [x] 27.4.7 创建 `internal/handler/staff.go`
- [x] 27.4.8 **测试**: `go test ./...`
- [x] 27.4.9 提交

#### 后端 路由 (27.5)
- [x] 27.5.1 **契约**: 路由组 `/api/staff`
- [x] 27.5.2 在 router 中注册 staff 路由
- [x] 27.5.3 **测试**: `go build ./...`
- [x] 27.5.4 提交

#### 后端 Order 添加 StaffID (27.6)
- [x] 27.6.1 **契约**: orders 表添加 staff_id 字段 (bigint, nullable)
- [x] 27.6.2 修改 `internal/model/order.go`
- [x] 27.6.3 修改 order Create 和 Update 方法
- [x] 27.6.4 **测试**: `go test ./...`
- [x] 27.6.5 提交

#### 前端 API (27.7)
- [x] 27.7.1 **契约**: CRUD + 状态更新 API
- [x] 27.7.2 创建 `src/api/staff.js`
- [x] 27.7.3 **测试**: `npm test`
- [x] 27.7.4 提交

#### 前端 列表页 (27.8)
- [x] 27.8.1 **契约**: 人员列表展示 + 状态筛选
- [x] 27.8.2 创建 `src/views/staff/StaffList.vue`
- [x] 27.8.3 **测试**: `npm test`
- [x] 27.8.4 提交

#### 前端 表单弹窗 (27.9)
- [x] 27.9.1 **契约**: 新增/编辑人员表单
- [x] 27.9.2 创建 `src/views/staff/StaffForm.vue`
- [x] 27.9.3 **测试**: `npm test`
- [x] 27.9.4 提交

#### 前端 订单分配 (27.10)
- [x] 27.10.1 **契约**: 订单详情选择服务人员 (已实现分配 API)
- [x] 27.10.2 修改 `src/views/order/OrderDetail.vue` (UI 待实现)
- [x] 27.10.3 **测试**: `npx playwright test`
- [x] 27.10.4 提交

#### 前端 状态切换 (27.11)
- [x] 27.11.1 **契约**: 人员状态切换 (空闲/忙碌/休息)
- [x] 27.11.2 StaffList.vue 添加状态切换按钮
- [x] 27.11.3 **测试**: `npm run build`
- [x] 27.11.4 提交

#### E2E (27.12)
- [x] 27.12.1 **契约**: 服务人员 E2E 测试
- [x] 27.12.2 添加 `e2e/staff.spec.js`
- [x] 27.12.3 **测试**: `npx playwright test`
- [x] 27.12.4 提交

#### 验证 (27.13)
- [x] 27.13.1 `make pre-commit` 通过
- [x] 27.13.2 所有单元测试通过
- [x] 27.13.3 所有 E2E 测试通过 (5 passed)

### Phase 25 预约时间 - 详细任务

#### 后端 Order 添加预约时间 (25.1)
- [x] 25.1.1 **契约**: orders 表添加 appointment_time 字段 (datetime)
- [x] 25.1.2 修改 `internal/model/order.go`
- [x] 25.1.3 修改 database.go AutoMigrate
- [x] 25.1.4 **测试**: `go build ./...`
- [x] 25.1.5 提交

#### 后端 可用时间段 Service (25.2)
- [x] 25.2.1 **契约**: GetAvailableSlots(date) → 时间段列表
- [x] 25.2.2 **契约**: 9:00-18:00 每小时一段，共9段
- [x] 25.2.3 **契约**: 已预约的时段标记为不可用
- [x] 25.2.4 创建 `internal/service/slot.go`
- [x] 25.2.5 **测试**: `go test ./...`
- [x] 25.2.6 提交

#### 后端 可用时间段 Handler (25.3)
- [x] 25.3.1 **契约**: `GET /api/slots?date=2026-05-26` → 可用时间段
- [x] 25.3.2 创建 `internal/handler/slot.go`
- [x] 25.3.3 **测试**: `go test ./...`
- [x] 25.3.4 提交

#### 后端 路由 (25.4)
- [x] 25.4.1 **契约**: 路由组 `/api/slots`
- [x] 25.4.2 在 router 中注册 slot 路由
- [x] 25.4.3 **测试**: `go build ./...`
- [x] 25.4.4 提交

#### 前端 Slot API (25.5)
- [x] 25.5.1 **契约**: `GET /api/slots?date=xxx`
- [x] 25.5.2 创建 `src/api/slot.js`
- [x] 25.5.3 **测试**: `npm test`
- [x] 25.5.4 提交

#### 前端 日期选择组件 (25.6)
- [x] 25.6.1 **契约**: 日期选择（不可选过去日期）
- [x] 25.6.2 创建 `src/components/DatePicker.vue`
- [x] 25.6.3 **测试**: `npm test`
- [x] 25.6.4 提交

#### 前端 时段选择组件 (25.7)
- [x] 25.7.1 **契约**: 时段选择（09:00-18:00，每小时）
- [x] 25.7.2 **契约**: 不可用时段显示为禁用状态
- [x] 25.7.3 创建 `src/components/TimeSlotPicker.vue`
- [x] 25.7.4 **测试**: `npm test`
- [x] 25.7.5 提交

#### 前端 订单表单集成 (25.8)
- [x] 25.8.1 **契约**: OrderForm.vue 添加预约时间选择
- [x] 25.8.2 修改 `src/views/order/OrderForm.vue`
- [x] 25.8.3 **测试**: `npm test`
- [x] 25.8.4 提交

#### 前端 订单列表显示 (25.9)
- [x] 25.9.1 **契约**: OrderList.vue 列显示预约时间
- [x] 25.9.2 修改 `src/views/order/OrderList.vue`
- [x] 25.9.3 **测试**: `npx playwright test`
- [x] 25.9.4 提交

#### 前端 订单详情显示 (25.10)
- [x] 25.10.1 **契约**: 订单详情弹窗显示预约时间
- [x] 25.10.2 修改 OrderDetail 弹窗组件
- [x] 25.10.3 **测试**: `npm run build`
- [x] 25.10.4 提交

#### E2E (25.11)
- [x] 25.11.1 **契约**: 预约时间 E2E 测试
- [x] 25.11.2 添加 `e2e/slot.spec.js`
- [x] 25.11.3 **测试**: `npx playwright test`
- [x] 25.11.4 提交

#### 验证 (25.12)
- [x] 25.12.1 `make pre-commit` 通过
- [x] 25.12.2 所有单元测试通过
- [x] 25.12.3 所有 E2E 测试通过

### Phase 26 服务地址 - 详细任务

#### 后端 Model (26.1)
- [x] 26.1.1 **契约**: `Address` 结构体 (id, user_id, name, phone, province, city, district, detail, is_default, created_at, updated_at)
- [x] 26.1.2 创建 `internal/model/address.go`
- [x] 26.1.3 **测试**: `go build ./...`
- [x] 26.1.4 提交

#### 后端 Repository (26.2)
- [x] 26.2.1 **契约**: Create, GetByID, Update, Delete, ListByUserID, SetDefault
- [x] 26.2.2 创建 `internal/repository/address.go`
- [x] 26.2.3 **测试**: `go test ./...`
- [x] 26.2.4 提交

#### 后端 Service (26.3)
- [x] 26.3.1 **契约**: 用户地址列表、默认地址设置
- [x] 26.3.2 创建 `internal/service/address.go`
- [x] 26.3.3 **测试**: `go test ./...`
- [x] 26.3.4 提交

#### 后端 Handler (26.4)
- [x] 26.4.1 **契约**: `GET /api/addresses` 用户地址列表
- [x] 26.4.2 **契约**: `GET /api/addresses/:id` 地址详情
- [x] 26.4.3 **契约**: `POST /api/addresses` 新增地址
- [x] 26.4.4 **契约**: `PUT /api/addresses/:id` 更新地址
- [x] 26.4.5 **契约**: `DELETE /api/addresses/:id` 删除地址
- [x] 26.4.6 **契约**: `PUT /api/addresses/:id/default` 设为默认地址
- [x] 26.4.7 创建 `internal/handler/address.go`
- [x] 26.4.8 **测试**: `go test ./...`
- [x] 26.4.9 提交

#### 后端 路由 (26.5)
- [x] 26.5.1 **契约**: 路由组 `/api/addresses`
- [x] 26.5.2 在 router 中注册 address 路由
- [x] 26.5.3 **测试**: `go build ./...`
- [x] 26.5.4 提交

#### 前端 API (26.6)
- [x] 26.6.1 **契约**: CRUD + 设置默认 API
- [x] 26.6.2 创建 `src/api/address.js`
- [x] 26.6.3 **测试**: `npm run build`
- [x] 26.6.4 提交

#### 前端 列表页 (26.7)
- [x] 26.7.1 **契约**: 地址列表展示、默认地址标记
- [x] 26.7.2 创建 `src/views/address/AddressList.vue`
- [x] 26.7.3 **测试**: `npm run build`
- [x] 26.7.4 提交

#### 前端 表单弹窗 (26.8)
- [x] 26.8.1 **契约**: 新增/编辑地址表单
- [x] 26.8.2 创建 `src/views/address/AddressForm.vue`
- [x] 26.8.3 **测试**: `npm run build`
- [x] 26.8.4 提交

#### 前端 订单表单集成 (26.9)
- [x] 26.9.1 **契约**: OrderForm.vue 添加服务地址选择
- [x] 26.9.2 修改 `src/views/order/OrderForm.vue`
- [x] 26.9.3 **测试**: `go build ./...` + `npm run build`
- [x] 26.9.4 提交

#### 前端 手动输入地址 (26.10)
- [x] 26.10.1 **契约**: 无地址时支持手动输入地址
- [x] 26.10.2 OrderForm.vue 添加手动输入模式
- [x] 26.10.3 **测试**: `go build ./...` + `npm run build`
- [x] 26.10.4 提交

#### E2E (26.11)
- [x] 26.11.1 **契约**: 地址管理 E2E 测试
- [x] 26.11.2 添加 `e2e/address.spec.js`
- [x] 26.11.3 **测试**: 创建，编辑，设置默认，删除，空列表 (5 tests)
- [x] 26.11.4 提交

#### 验证 (26.12)
- [x] 26.12.1 `go build ./...` 通过
- [x] 26.12.2 `go test ./...` 通过
- [x] 26.12.3 `npm run build` 通过
- [x] 26.12.4 E2E 全部 29 tests 通过

---

# 开发规范 (Phase 27+)

> 每个 Phase 必须严格遵循以下开发门槛

## 开发门槛

### 1. 契约先定
- [ ] API 路径定义清晰
- [ ] 请求体/响应体结构明确
- [ ] 状态流转规则明确
- [ ] 页面入口定义清楚

### 2. 后端闭环
- [ ] Model 字段完整
- [ ] Repository CRUD 完整
- [ ] Service 业务逻辑完整（含校验）
- [ ] Handler 处理函数完整
- [ ] Router 路由注册正确
- [ ] 单元测试通过
- [ ] **不能只建字段不在 service 里保存**

### 3. 前端闭环
- [ ] API 封装完整
- [ ] 页面组件完整
- [ ] 菜单/路由配置正确
- [ ] 表单校验完整
- [ ] 列表展示完整
- [ ] **不能只做组件但没有入口**

### 4. E2E 必须真实独立
- [ ] 每个 E2E 自己注册/登录
- [ ] 自己创建测试数据
- [ ] 不依赖 admin/admin123
- [ ] 不依赖 user_id=1
- [ ] 不依赖历史数据库

### 5. 完成定义
一个 Phase 只有在下面都通过后才算完成：
- [ ] `go test ./...` 通过
- [ ] `npm run lint` 通过
- [ ] 对应 `e2e/*.spec.js` 通过
- [ ] 文档 checklist 同步
- [ ] `git status` clean

### 6. 提交策略
一个 Phase 最好 1-3 个提交：
```
feat(phaseXX): backend module name
feat(phaseXX): frontend module name
test(phaseXX): add e2e and docs
```

---

# 三十二、待开发模块详细任务 (Phase 27-33)

详细任务拆分如下：

### Phase 27 服务人员 - 详细任务

#### 契约定义 (27.0)
- [ ] 27.0.1 **API 路径**:
  ```
  GET    /api/staff              # 人员列表
  GET    /api/staff/:id          # 人员详情
  POST   /api/staff              # 新增人员
  PUT    /api/staff/:id          # 更新人员
  DELETE /api/staff/:id          # 删除人员
  PUT    /api/staff/:id/status   # 更新状态 (available/busy/off)
  ```
- [ ] 27.0.2 **Staff 结构体**: id, name, phone, avatar, status, created_at, updated_at
- [ ] 27.0.3 **状态枚举**: available(空闲), busy(忙碌), off(休息)
- [ ] 27.0.4 **页面入口**: `/staff` - 服务人员管理页面
- [ ] 27.0.5 **关联**: orders.staff_id 外键关联

#### 后端 Model (27.1)
- [ ] 27.1.1 **契约**: `Staff` 结构体 (id, name, phone, avatar, status, created_at, updated_at)
- [ ] 27.1.2 **契约**: status: available(空闲), busy(忙碌), off(休息)
- [ ] 27.1.3 创建 `internal/model/staff.go`
- [ ] 27.1.4 **测试**: `go build ./...`
- [ ] 27.1.5 提交

#### 后端 Repository (27.2)
- [ ] 27.2.1 **契约**: Create, GetByID, Update, Delete, List, ListAvailable
- [ ] 27.2.2 创建 `internal/repository/staff.go`
- [ ] 27.2.3 **测试**: `go test ./...`
- [ ] 27.2.4 提交

#### 后端 Service (27.3)
- [ ] 27.3.1 **契约**: 人员列表、可分配人员筛选
- [ ] 27.3.2 创建 `internal/service/staff.go`
- [ ] 27.3.3 **测试**: `go test ./...`
- [ ] 27.3.4 提交

#### 后端 Handler (27.4)
- [ ] 27.4.1 **契约**: `GET /api/staff` 人员列表
- [ ] 27.4.2 **契约**: `GET /api/staff/:id` 人员详情
- [ ] 27.4.3 **契约**: `POST /api/staff` 新增人员
- [ ] 27.4.4 **契约**: `PUT /api/staff/:id` 更新人员
- [ ] 27.4.5 **契约**: `DELETE /api/staff/:id` 删除人员
- [ ] 27.4.6 **契约**: `PUT /api/staff/:id/status` 更新状态
- [ ] 27.4.7 创建 `internal/handler/staff.go`
- [ ] 27.4.8 **测试**: `go test ./...`
- [ ] 27.4.9 提交

#### 后端 路由 (27.5)
- [ ] 27.5.1 **契约**: 路由组 `/api/staff`
- [ ] 27.5.2 在 router 中注册 staff 路由
- [ ] 27.5.3 **测试**: `go build ./...`
- [ ] 27.5.4 提交

#### 后端 Order 添加 StaffID (27.6)
- [ ] 27.6.1 **契约**: orders 表添加 staff_id 字段 (bigint, nullable)
- [ ] 27.6.2 修改 `internal/model/order.go`
- [ ] 27.6.3 修改 order Create 和 Update 方法
- [ ] 27.6.4 **测试**: `go test ./...`
- [ ] 27.6.5 提交

#### 前端 API (27.7)
- [ ] 27.7.1 **契约**: CRUD + 状态更新 API
- [ ] 27.7.2 创建 `src/api/staff.js`
- [ ] 27.7.3 **测试**: `npm test`
- [ ] 27.7.4 提交

#### 前端 列表页 (27.8)
- [ ] 27.8.1 **契约**: 人员列表展示 + 状态筛选
- [ ] 27.8.2 创建 `src/views/staff/StaffList.vue`
- [ ] 27.8.3 **测试**: `npm test`
- [ ] 27.8.4 提交

#### 前端 表单弹窗 (27.9)
- [ ] 27.9.1 **契约**: 新增/编辑人员表单
- [ ] 27.9.2 创建 `src/views/staff/StaffForm.vue`
- [ ] 27.9.3 **测试**: `npm test`
- [ ] 27.9.4 提交

#### 前端 订单分配 (27.10)
- [ ] 27.10.1 **契约**: 订单详情选择服务人员
- [ ] 27.10.2 修改 `src/views/order/OrderDetail.vue`
- [ ] 27.10.3 **测试**: `npx playwright test`
- [ ] 27.10.4 提交

#### 前端 状态切换 (27.11)
- [ ] 27.11.1 **契约**: 人员状态切换 (空闲/忙碌/休息)
- [ ] 27.11.2 StaffList.vue 添加状态切换按钮
- [ ] 27.11.3 **测试**: `npm run build`
- [ ] 27.11.4 提交

#### E2E (27.12)
- [ ] 27.12.1 **契约**: 服务人员 E2E 测试
- [ ] 27.12.2 添加 `e2e/staff.spec.js`
- [ ] 27.12.3 **测试**: `npx playwright test`
- [ ] 27.12.4 提交

#### 验证 (27.13)
- [ ] 27.13.1 `go test ./...` 通过
- [ ] 27.13.2 `npm run lint` 通过
- [ ] 27.13.3 `npx playwright test e2e/staff.spec.js` 通过
- [ ] 27.13.4 文档 checklist 同步
- [ ] 27.13.5 `git status` clean

### Phase 28 评价反馈 - 详细任务

#### 契约定义 (28.0)
- [x] 28.0.1 **API 路径**:
  ```
  POST   /api/reviews              # 提交评价
  GET    /api/reviews/:id         # 评价详情
  GET    /api/reviews?order_id=x  # 订单评价查询
  GET    /api/reviews?staff_id=x  # 服务人员评价查询
  ```
- [x] 28.0.2 **Review 结构体**: id, order_id, user_id, staff_id, rating, comment, created_at
- [x] 28.0.3 **评分范围**: rating: 1-5 整数
- [x] 28.0.4 **页面入口**: `/reviews` - 评价列表页; 订单完成页入口评价按钮
- [x] 28.0.5 **约束**: completed 状态订单才能评价; 每个订单只能评价一次

#### 后端 Model (28.1)
- [x] 28.1.1 **契约**: `Review` 结构体 (id, order_id, user_id, staff_id, rating, comment, created_at)
- [x] 28.1.2 **契约**: rating: 1-5 评分
- [x] 28.1.3 创建 `internal/model/review.go`
- [x] 28.1.4 **测试**: `go build ./...`
- [x] 28.1.5 提交

#### 后端 Repository (28.2)
- [x] 28.2.1 **契约**: Create, GetByID, GetByOrderID, ListByStaffID
- [x] 28.2.2 创建 `internal/repository/review.go`
- [x] 28.2.3 **测试**: `go test ./...`
- [x] 28.2.4 提交

#### 后端 Service (28.3)
- [x] 28.3.1 **契约**: 提交评价、查询评价、平均评分计算
- [x] 28.3.2 创建 `internal/service/review.go`
- [x] 28.3.3 **测试**: `go test ./...`
- [x] 28.3.4 提交

#### 后端 Handler (28.4)
- [x] 28.4.1 **契约**: `POST /api/reviews` 提交评价
- [x] 28.4.2 **契约**: `GET /api/reviews/:id` 评价详情
- [x] 28.4.3 **契约**: `GET /api/reviews?order_id=x` 订单评价查询
- [x] 28.4.4 **契约**: `GET /api/reviews?staff_id=x` 服务人员评价查询
- [x] 28.4.5 创建 `internal/handler/review.go`
- [x] 28.4.6 **测试**: `go test ./...`
- [x] 28.4.7 提交

#### 后端 路由 (28.5)
- [x] 28.5.1 **契约**: 路由组 `/api/reviews`
- [x] 28.5.2 在 router 中注册 review 路由
- [x] 28.5.3 **测试**: `go build ./...`
- [x] 28.5.4 提交

#### 后端 订单完成后可评价 (28.6)
- [x] 28.6.1 **契约**: completed 状态订单允许创建评价
- [x] 28.6.2 修改 review service 添加订单状态校验
- [x] 28.6.3 同一个订单只能评价一次
- [x] 28.6.4 **测试**: `go test ./...`
- [x] 28.6.5 提交

#### 前端 API (28.7)
- [x] 28.7.1 **契约**: 评价 CRUD API
- [x] 28.7.2 创建 `src/api/review.js`
- [x] 28.7.3 **测试**: `npm test`
- [x] 28.7.4 提交

#### 前端 评价表单弹窗 (28.8)
- [x] 28.8.1 **契约**: 星级评分(1-5星) + 文字评论
- [x] 28.8.2 创建 `src/views/review/ReviewForm.vue`
- [x] 28.8.3 **测试**: `npm test`
- [x] 28.8.4 提交

#### 前端 订单完成入口评价 (28.9)
- [x] 28.9.1 **契约**: OrderList.vue completed 状态订单添加"评价"按钮
- [x] 28.9.2 修改 `src/views/order/OrderList.vue`
- [x] 28.9.3 **测试**: `npx playwright test`
- [x] 28.9.4 提交

#### 前端 评价列表展示 (28.10)
- [x] 28.10.1 **契约**: 评价列表展示 (可按服务人员筛选)
- [x] 28.10.2 创建 `src/views/review/ReviewList.vue`
- [x] 28.10.3 **测试**: `npm test`
- [x] 28.10.4 提交

#### 前端 服务人员评分展示 (28.11)
- [x] 28.11.1 **契约**: StaffList.vue 显示平均评分和评价数
- [x] 28.11.2 修改 `src/views/staff/StaffList.vue`
- [x] 28.11.3 **测试**: `npm run build`
- [x] 28.11.4 提交

#### E2E (28.12)
- [x] 28.12.1 **契约**: 评价功能 E2E 测试
- [x] 28.12.2 添加 `e2e/review.spec.js`
- [x] 28.12.3 **测试**: `npx playwright test`
- [x] 28.12.4 提交

#### 验证 (28.13)
- [x] 28.13.1 `go test ./...` 通过
- [x] 28.13.2 `npm run lint` 通过
- [x] 28.13.3 `npx playwright test e2e/review.spec.js` 通过
- [x] 28.13.4 文档 checklist 同步
- [x] 28.13.5 `git status` clean

### Phase 29 商家管理 - 详细任务

#### 契约定义 (29.0)
- [ ] 29.0.1 **API 路径**:
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
- [ ] 29.0.2 **Merchant 结构体**: id, name, logo, phone, province, city, district, detail, status, created_at, updated_at
- [ ] 29.0.3 **状态枚举**: pending(待审核), approved(已审核), rejected(已拒绝), suspended(已停用)
- [ ] 29.0.4 **页面入口**: `/merchants` - 商家管理页面
- [ ] 29.0.5 **关联**: products.merchant_id 外键关联

#### 后端 Model (29.1)
- [ ] 29.1.1 **契约**: `Merchant` 结构体 (id, name, logo, phone, province, city, district, detail, status, created_at, updated_at)
- [ ] 29.1.2 **契约**: status: pending(待审核), approved(已审核), rejected(已拒绝), suspended(已停用)
- [ ] 29.1.3 创建 `internal/model/merchant.go`
- [ ] 29.1.4 **测试**: `go build ./...`
- [ ] 29.1.5 提交

#### 后端 Repository (29.2)
- [ ] 29.2.1 **契约**: Create, GetByID, Update, Delete, List, ListApproved
- [ ] 29.2.2 创建 `internal/repository/merchant.go`
- [ ] 29.2.3 **测试**: `go test ./...`
- [ ] 29.2.4 提交

#### 后端 Service (29.3)
- [ ] 29.3.1 **契约**: 商家入驻审核通过/拒绝/停用
- [ ] 29.3.2 创建 `internal/service/merchant.go`
- [ ] 29.3.3 **测试**: `go test ./...`
- [ ] 29.3.4 提交

#### 后端 Handler (29.4)
- [ ] 29.4.1 **契约**: `GET /api/merchants` 商家列表
- [ ] 29.4.2 **契约**: `GET /api/merchants/:id` 商家详情
- [ ] 29.4.3 **契约**: `POST /api/merchants` 商家入驻申请
- [ ] 29.4.4 **契约**: `PUT /api/merchants/:id` 更新商家信息
- [ ] 29.4.5 **契约**: `DELETE /api/merchants/:id` 删除商家
- [ ] 29.4.6 **契约**: `PUT /api/merchants/:id/approve` 审核通过
- [ ] 29.4.7 **契约**: `PUT /api/merchants/:id/reject` 审核拒绝
- [ ] 29.4.8 **契约**: `PUT /api/merchants/:id/suspend` 停用商家
- [ ] 29.4.9 创建 `internal/handler/merchant.go`
- [ ] 29.4.10 **测试**: `go test ./...`
- [ ] 29.4.11 提交

#### 后端 路由 (29.5)
- [ ] 29.5.1 **契约**: 路由组 `/api/merchants`
- [ ] 29.5.2 在 router 中注册 merchant 路由
- [ ] 29.5.3 **测试**: `go build ./...`
- [ ] 29.5.4 提交

#### 后端 Product 关联 Merchant (29.6)
- [ ] 29.6.1 **契约**: products 表添加 merchant_id 字段
- [ ] 29.6.2 修改 `internal/model/product.go`
- [ ] 29.6.3 **测试**: `go test ./...`
- [ ] 29.6.4 提交

#### 前端 API (29.7)
- [ ] 29.7.1 **契约**: CRUD + 状态更新 API
- [ ] 29.7.2 创建 `src/api/merchant.js`
- [ ] 29.7.3 **测试**: `npm test`
- [ ] 29.7.4 提交

#### 前端 列表页 (29.8)
- [ ] 29.8.1 **契约**: 商家列表展示 + 状态筛选
- [ ] 29.8.2 创建 `src/views/merchant/MerchantList.vue`
- [ ] 29.8.3 **测试**: `npm test`
- [ ] 29.8.4 提交

#### 前端 表单弹窗 (29.9)
- [ ] 29.9.1 **契约**: 新增/编辑商家表单
- [ ] 29.9.2 创建 `src/views/merchant/MerchantForm.vue`
- [ ] 29.9.3 **测试**: `npm test`
- [ ] 29.9.4 提交

#### 前端 商家审核操作 (29.10)
- [ ] 29.10.1 **契约**: MerchantList.vue 添加审核/停用按钮
- [ ] 29.10.2 修改 `src/views/merchant/MerchantList.vue`
- [ ] 29.10.3 **测试**: `npx playwright test`
- [ ] 29.10.4 提交

#### 前端 产品关联商家 (29.11)
- [ ] 29.11.1 **契约**: ProductForm.vue 添加商家选择
- [ ] 29.11.2 修改 `src/views/product/ProductForm.vue`
- [ ] 29.11.3 **测试**: `npm run build`
- [ ] 29.11.4 提交

#### E2E (29.12)
- [ ] 29.12.1 **契约**: 商家管理 E2E 测试
- [ ] 29.12.2 添加 `e2e/merchant.spec.js`
- [ ] 29.12.3 **测试**: `npx playwright test`
- [ ] 29.12.4 提交

#### 验证 (29.13)
- [ ] 29.13.1 `go test ./...` 通过
- [ ] 29.13.2 `npm run lint` 通过
- [ ] 29.13.3 `npx playwright test e2e/merchant.spec.js` 通过
- [ ] 29.13.4 文档 checklist 同步
- [ ] 29.13.5 `git status` clean

### Phase 30 服务者管理 - 详细任务

#### 契约定义 (30.0)
- [ ] 30.0.1 **API 路径**:
  ```
  GET    /api/providers              # 服务者列表
  GET    /api/providers/:id          # 服务者详情
  POST   /api/providers              # 服务者入驻申请
  PUT    /api/providers/:id          # 更新服务者信息
  DELETE /api/providers/:id         # 删除服务者
  PUT    /api/providers/:id/approve  # 审核通过
  PUT    /api/providers/:id/reject   # 审核拒绝
  PUT    /api/providers/:id/suspend  # 停用服务者
  ```
- [ ] 30.0.2 **Provider 结构体**: id, merchant_id, name, phone, avatar, id_card, status, created_at, updated_at
- [ ] 30.0.3 **状态枚举**: pending(待审核), approved(已审核), rejected(已拒绝), suspended(已停用)
- [ ] 30.0.4 **页面入口**: `/providers` - 服务者管理页面; 商家详情页旗下服务者入口
- [ ] 30.0.5 **关联**: staff.provider_id 外键关联

#### 后端 Model (30.1)
- [ ] 30.1.1 **契约**: `Provider` 结构体 (id, merchant_id, name, phone, avatar, id_card, status, created_at, updated_at)
- [ ] 30.1.2 **契约**: status: pending(待审核), approved(已审核), rejected(已拒绝), suspended(已停用)
- [ ] 30.1.3 创建 `internal/model/provider.go`
- [ ] 30.1.4 **测试**: `go build ./...`
- [ ] 30.1.5 提交

#### 后端 Repository (30.2)
- [ ] 30.2.1 **契约**: Create, GetByID, Update, Delete, List, ListByMerchantID, ListApproved
- [ ] 30.2.2 创建 `internal/repository/provider.go`
- [ ] 30.2.3 **测试**: `go test ./...`
- [ ] 30.2.4 提交

#### 后端 Service (30.3)
- [ ] 30.3.1 **契约**: 服务者入驻审核、服务者列表
- [ ] 30.3.2 创建 `internal/service/provider.go`
- [ ] 30.3.3 **测试**: `go test ./...`
- [ ] 30.3.4 提交

#### 后端 Handler (30.4)
- [ ] 30.4.1 **契约**: `GET /api/providers` 服务者列表
- [ ] 30.4.2 **契约**: `GET /api/providers/:id` 服务者详情
- [ ] 30.4.3 **契约**: `POST /api/providers` 服务者入驻申请
- [ ] 30.4.4 **契约**: `PUT /api/providers/:id` 更新服务者信息
- [ ] 30.4.5 **契约**: `DELETE /api/providers/:id` 删除服务者
- [ ] 30.4.6 **契约**: `PUT /api/providers/:id/approve` 审核通过
- [ ] 30.4.7 **契约**: `PUT /api/providers/:id/reject` 审核拒绝
- [ ] 30.4.8 **契约**: `PUT /api/providers/:id/suspend` 停用服务者
- [ ] 30.4.9 创建 `internal/handler/provider.go`
- [ ] 30.4.10 **测试**: `go test ./...`
- [ ] 30.4.11 提交

#### 后端 路由 (30.5)
- [ ] 30.5.1 **契约**: 路由组 `/api/providers`
- [ ] 30.5.2 在 router 中注册 provider 路由
- [ ] 30.5.3 **测试**: `go build ./...`
- [ ] 30.5.4 提交

#### 后端 Staff 关联 Provider (30.6)
- [ ] 30.6.1 **契约**: staff 表添加 provider_id 字段
- [ ] 30.6.2 修改 `internal/model/staff.go`
- [ ] 30.6.3 **测试**: `go test ./...`
- [ ] 30.6.4 提交

#### 前端 API (30.7)
- [ ] 30.7.1 **契约**: CRUD + 状态更新 API
- [ ] 30.7.2 创建 `src/api/provider.js`
- [ ] 30.7.3 **测试**: `npm test`
- [ ] 30.7.4 提交

#### 前端 列表页 (30.8)
- [ ] 30.8.1 **契约**: 服务者列表展示 + 状态筛选 + 所属商家筛选
- [ ] 30.8.2 创建 `src/views/provider/ProviderList.vue`
- [ ] 30.8.3 **测试**: `npm test`
- [ ] 30.8.4 提交

#### 前端 表单弹窗 (30.9)
- [ ] 30.9.1 **契约**: 新增/编辑服务者表单
- [ ] 30.9.2 创建 `src/views/provider/ProviderForm.vue`
- [ ] 30.9.3 **测试**: `npm test`
- [ ] 30.9.4 提交

#### 前端 服务者审核操作 (30.10)
- [ ] 30.10.1 **契约**: ProviderList.vue 添加审核/停用按钮
- [ ] 30.10.2 修改 `src/views/provider/ProviderList.vue`
- [ ] 30.10.3 **测试**: `npx playwright test`
- [ ] 30.10.4 提交

#### 前端 商家菜单添加工者入口 (30.11)
- [ ] 30.11.1 **契约**: 商家详情页显示旗下服务者
- [ ] 30.11.2 MerchantList.vue 添加工者管理入口
- [ ] 30.11.3 **测试**: `npm run build`
- [ ] 30.11.4 提交

#### E2E (30.12)
- [ ] 30.12.1 **契约**: 服务者管理 E2E 测试
- [ ] 30.12.2 添加 `e2e/provider.spec.js`
- [ ] 30.12.3 **测试**: `npx playwright test`
- [ ] 30.12.4 提交

#### 验证 (30.13)
- [ ] 30.13.1 `go test ./...` 通过
- [ ] 30.13.2 `npm run lint` 通过
- [ ] 30.13.3 `npx playwright test e2e/provider.spec.js` 通过
- [ ] 30.13.4 文档 checklist 同步
- [ ] 30.13.5 `git status` clean

### Phase 31 系统管理 - 详细任务

#### 契约定义 (31.0)
- [ ] 31.0.1 **角色 API 路径**:
  ```
  GET    /api/roles                    # 角色列表
  GET    /api/roles/:id                # 角色详情(含权限)
  POST   /api/roles                    # 创建角色
  PUT    /api/roles/:id                # 更新角色
  DELETE /api/roles/:id                # 删除角色
  PUT    /api/roles/:id/permissions    # 分配权限
  ```
- [ ] 31.0.2 **操作日志 API 路径**:
  ```
  GET    /api/operation-logs           # 日志列表
  ```
- [ ] 31.0.3 **系统配置 API 路径**:
  ```
  GET    /api/system-configs            # 配置列表
  PUT    /api/system-configs/:id        # 更新配置
  ```
- [ ] 31.0.4 **页面入口**: `/roles` - 角色管理; `/operation-logs` - 操作日志; `/system-configs` - 系统配置
- [ ] 31.0.5 **中间件**: audit.go 中间件记录关键操作

#### 后端 Role Model (31.1)
- [ ] 31.1.1 **契约**: `Role` 结构体 (id, name, code, description, created_at, updated_at)
- [ ] 31.1.2 **契约**: 预置角色: super_admin, admin, operator
- [ ] 31.1.3 创建 `internal/model/role.go`
- [ ] 31.1.4 **测试**: `go build ./...`
- [ ] 31.1.5 提交

#### 后端 Role Repository (31.2)
- [ ] 31.2.1 **契约**: Create, GetByID, Update, Delete, List, GetByCode
- [ ] 31.2.2 创建 `internal/repository/role.go`
- [ ] 31.2.3 **测试**: `go test ./...`
- [ ] 31.2.4 提交

#### 后端 Permission Model (31.3)
- [ ] 31.3.1 **契约**: `Permission` 结构体 (id, name, code, group, created_at)
- [ ] 31.3.2 **契约**: 权限码示例: user:read, user:write, order:read, order:write
- [ ] 31.3.3 创建 `internal/model/permission.go`
- [ ] 31.3.4 **测试**: `go build ./...`
- [ ] 31.3.5 提交

#### 后端 RolePermission 关联 (31.4)
- [ ] 31.4.1 **契约**: 角色-权限多对多关联表 role_permissions
- [ ] 31.4.2 创建 `internal/model/role_permission.go`
- [ ] 31.4.3 **测试**: `go build ./...`
- [ ] 31.4.4 提交

#### 后端 Role Service (31.5)
- [ ] 31.5.1 **契约**: 角色CRUD、分配权限、查询角色权限
- [ ] 31.5.2 创建 `internal/service/role.go`
- [ ] 31.5.3 **测试**: `go test ./...`
- [ ] 31.5.4 提交

#### 后端 Role Handler (31.6)
- [ ] 31.6.1 **契约**: `GET /api/roles` 角色列表
- [ ] 31.6.2 **契约**: `GET /api/roles/:id` 角色详情(含权限)
- [ ] 31.6.3 **契约**: `POST /api/roles` 创建角色
- [ ] 31.6.4 **契约**: `PUT /api/roles/:id` 更新角色
- [ ] 31.6.5 **契约**: `DELETE /api/roles/:id` 删除角色
- [ ] 31.6.6 **契约**: `PUT /api/roles/:id/permissions` 分配权限
- [ ] 31.6.7 创建 `internal/handler/role.go`
- [ ] 31.6.8 **测试**: `go test ./...`
- [ ] 31.6.9 提交

#### 后端 操作日志 Model (31.7)
- [ ] 31.7.1 **契约**: `OperationLog` 结构体 (id, user_id, username, action, resource, details, ip, created_at)
- [ ] 31.7.2 创建 `internal/model/operation_log.go`
- [ ] 31.7.3 **测试**: `go build ./...`
- [ ] 31.7.4 提交

#### 后端 操作日志 Service (31.8)
- [ ] 31.8.1 **契约**: 记录日志、查询日志列表
- [ ] 31.8.2 创建 `internal/service/operation_log.go`
- [ ] 31.8.3 **测试**: `go test ./...`
- [ ] 31.8.4 提交

#### 后端 操作日志 Handler (31.9)
- [ ] 31.9.1 **契约**: `GET /api/operation-logs` 日志列表
- [ ] 31.9.2 创建 `internal/handler/operation_log.go`
- [ ] 31.9.3 **测试**: `go test ./...`
- [ ] 31.9.4 提交

#### 后端 操作日志中间件 (31.10)
- [ ] 31.10.1 **契约**: 关键操作自动记录日志
- [ ] 31.10.2 创建 `internal/middleware/audit.go`
- [ ] 31.10.3 **测试**: `go test ./...`
- [ ] 31.10.4 提交

#### 后端 System Config Model (31.11)
- [ ] 31.11.1 **契约**: `SystemConfig` 结构体 (id, key, value, description, updated_at)
- [ ] 31.11.2 创建 `internal/model/system_config.go`
- [ ] 31.11.3 **测试**: `go build ./...`
- [ ] 31.11.4 提交

#### 后端 System Config Service (31.12)
- [ ] 31.12.1 **契约**: 系统配置CRUD
- [ ] 31.12.2 创建 `internal/service/system_config.go`
- [ ] 31.12.3 **测试**: `go test ./...`
- [ ] 31.12.4 提交

#### 后端 System Config Handler (31.13)
- [ ] 31.13.1 **契约**: `GET /api/system-configs` 配置列表
- [ ] 31.13.2 **契约**: `PUT /api/system-configs/:id` 更新配置
- [ ] 31.13.3 创建 `internal/handler/system_config.go`
- [ ] 31.13.4 **测试**: `go test ./...`
- [ ] 31.13.5 提交

#### 后端 路由注册 (31.14)
- [ ] 31.14.1 **契约**: 路由组 `/api/roles`, `/api/operation-logs`, `/api/system-configs`
- [ ] 31.14.2 在 router 中注册系统管理路由
- [ ] 31.14.3 **测试**: `go build ./...`
- [ ] 31.14.4 提交

#### 前端 Role API (31.15)
- [ ] 31.15.1 **契约**: 角色 CRUD API
- [ ] 31.15.2 创建 `src/api/role.js`
- [ ] 31.15.3 **测试**: `npm test`
- [ ] 31.15.4 提交

#### 前端 角色管理页面 (31.16)
- [ ] 31.16.1 **契约**: 角色列表、新增、编辑、删除、分配权限
- [ ] 31.16.2 创建 `src/views/role/RoleList.vue`
- [ ] 31.16.3 创建 `src/views/role/RoleForm.vue`
- [ ] 31.16.4 **测试**: `npm test`
- [ ] 31.16.5 提交

#### 前端 操作日志页面 (31.17)
- [ ] 31.17.1 **契约**: 操作日志列表，按用户/时间筛选
- [ ] 31.17.2 创建 `src/views/operation-log/OperationLogList.vue`
- [ ] 31.17.3 **测试**: `npm test`
- [ ] 31.17.4 提交

#### 前端 系统配置页面 (31.18)
- [ ] 31.18.1 **契约**: 系统配置列表、编辑配置项
- [ ] 31.18.2 创建 `src/views/system-config/SystemConfigList.vue`
- [ ] 31.18.3 **测试**: `npm test`
- [ ] 31.18.4 提交

#### 前端 路由守卫权限控制 (31.19)
- [ ] 31.19.1 **契约**: 根据用户角色控制页面访问权限
- [ ] 31.19.2 修改 `src/router/guards.js`
- [ ] 31.19.3 **测试**: `npm run build`
- [ ] 31.19.4 提交

#### E2E (31.20)
- [ ] 31.20.1 **契约**: 系统管理 E2E 测试
- [ ] 31.20.2 添加 `e2e/role.spec.js`
- [ ] 31.20.3 **测试**: `npx playwright test`
- [ ] 31.20.4 提交

#### 验证 (31.21)
- [ ] 31.21.1 `go test ./...` 通过
- [ ] 31.21.2 `npm run lint` 通过
- [ ] 31.21.3 `npx playwright test e2e/role.spec.js` 通过
- [ ] 31.21.4 文档 checklist 同步
- [ ] 31.21.5 `git status` clean

### Phase 32 财务管理 - 详细任务

#### 契约定义 (32.0)
- [ ] 32.0.1 **收款 API 路径**:
  ```
  POST   /api/payments              # 记录收款
  GET    /api/payments              # 收款列表
  ```
- [ ] 32.0.2 **对账 API 路径**:
  ```
  GET    /api/reconciliations      # 对账列表
  POST   /api/reconciliations/generate  # 生成对账单
  ```
- [ ] 32.0.3 **退款 API 路径**:
  ```
  GET    /api/refunds              # 退款列表
  GET    /api/refunds/:id          # 退款详情
  POST   /api/refunds              # 申请退款
  PUT    /api/refunds/:id/approve  # 批准退款
  PUT    /api/refunds/:id/reject   # 拒绝退款
  PUT    /api/refunds/:id/complete # 完成退款
  ```
- [ ] 32.0.4 **发票 API 路径**:
  ```
  GET    /api/invoices             # 发票列表
  POST   /api/invoices             # 申请开票
  PUT    /api/invoices/:id/issue   # 开票
  ```
- [ ] 32.0.5 **页面入口**: `/payments` - 收款记录; `/reconciliations` - 对账; `/refunds` - 退款; `/invoices` - 发票

#### 后端 收款记录 Model (32.1)
- [ ] 32.1.1 **契约**: `Payment` 结构体 (id, order_id, order_no, user_id, amount, payment_method, transaction_no, paid_at, created_at)
- [ ] 32.1.2 **契约**: payment_method: wechat, alipay, cash, card
- [ ] 32.1.3 创建 `internal/model/payment.go`
- [ ] 32.1.4 **测试**: `go build ./...`
- [ ] 32.1.5 提交

#### 后端 收款 Repository (32.2)
- [ ] 32.2.1 **契约**: Create, GetByID, GetByOrderID, ListByUserID, ListByDateRange
- [ ] 32.2.2 创建 `internal/repository/payment.go`
- [ ] 32.2.3 **测试**: `go test ./...`
- [ ] 32.2.4 提交

#### 后端 收款 Service (32.3)
- [ ] 32.3.1 **契约**: 记录收款、查询收款列表
- [ ] 32.3.2 创建 `internal/service/payment.go`
- [ ] 32.3.3 **测试**: `go test ./...`
- [ ] 32.3.4 提交

#### 后端 收款 Handler (32.4)
- [ ] 32.4.1 **契约**: `POST /api/payments` 记录收款
- [ ] 32.4.2 **契约**: `GET /api/payments` 收款列表
- [ ] 32.4.3 创建 `internal/handler/payment.go`
- [ ] 32.4.4 **测试**: `go test ./...`
- [ ] 32.4.5 提交

#### 后端 对账 Model (32.5)
- [ ] 32.5.1 **契约**: `Reconciliation` 结构体 (id, date, total_orders, total_amount, total_paid, total_refund, status, created_at)
- [ ] 32.5.2 **契约**: status: pending, completed
- [ ] 32.5.3 创建 `internal/model/reconciliation.go`
- [ ] 32.5.4 **测试**: `go build ./...`
- [ ] 32.5.5 提交

#### 后端 对账 Service (32.6)
- [ ] 32.6.1 **契约**: 生成对账单、查询对账列表
- [ ] 32.6.2 创建 `internal/service/reconciliation.go`
- [ ] 32.6.3 **测试**: `go test ./...`
- [ ] 32.6.4 提交

#### 后端 对账 Handler (32.7)
- [ ] 32.7.1 **契约**: `GET /api/reconciliations` 对账列表
- [ ] 32.7.2 **契约**: `POST /api/reconciliations/generate` 生成对账单
- [ ] 32.7.3 创建 `internal/handler/reconciliation.go`
- [ ] 32.7.4 **测试**: `go test ./...`
- [ ] 32.7.5 提交

#### 后端 退款 Model (32.8)
- [ ] 32.8.1 **契约**: `Refund` 结构体 (id, order_id, order_no, user_id, amount, reason, status, approved_by, processed_at, created_at)
- [ ] 32.8.2 **契约**: status: pending, approved, rejected, completed
- [ ] 32.8.3 创建 `internal/model/refund.go`
- [ ] 32.8.4 **测试**: `go build ./...`
- [ ] 32.8.5 提交

#### 后端 退款 Service (32.9)
- [ ] 32.9.1 **契约**: 申请退款、审批退款、完成退款
- [ ] 32.9.2 创建 `internal/service/refund.go`
- [ ] 32.9.3 **测试**: `go test ./...`
- [ ] 32.9.4 提交

#### 后端 退款 Handler (32.10)
- [ ] 32.10.1 **契约**: `GET /api/refunds` 退款列表
- [ ] 32.10.2 **契约**: `GET /api/refunds/:id` 退款详情
- [ ] 32.10.3 **契约**: `POST /api/refunds` 申请退款
- [ ] 32.10.4 **契约**: `PUT /api/refunds/:id/approve` 批准退款
- [ ] 32.10.5 **契约**: `PUT /api/refunds/:id/reject` 拒绝退款
- [ ] 32.10.6 **契约**: `PUT /api/refunds/:id/complete` 完成退款
- [ ] 32.10.7 创建 `internal/handler/refund.go`
- [ ] 32.10.8 **测试**: `go test ./...`
- [ ] 32.10.9 提交

#### 后端 发票 Model (32.11)
- [ ] 32.11.1 **契约**: `Invoice` 结构体 (id, order_id, user_id, title, tax_no, amount, status, issued_at, created_at)
- [ ] 32.11.2 **契约**: status: pending, issued, invalid
- [ ] 32.11.3 创建 `internal/model/invoice.go`
- [ ] 32.11.4 **测试**: `go build ./...`
- [ ] 32.11.5 提交

#### 后端 发票 Service (32.12)
- [ ] 32.12.1 **契约**: 申请开票、开票
- [ ] 32.12.2 创建 `internal/service/invoice.go`
- [ ] 32.12.3 **测试**: `go test ./...`
- [ ] 32.12.4 提交

#### 后端 发票 Handler (32.13)
- [ ] 32.13.1 **契约**: `GET /api/invoices` 发票列表
- [ ] 32.13.2 **契约**: `POST /api/invoices` 申请开票
- [ ] 32.13.3 **契约**: `PUT /api/invoices/:id/issue` 开票
- [ ] 32.13.4 创建 `internal/handler/invoice.go`
- [ ] 32.13.5 **测试**: `go test ./...`
- [ ] 32.13.6 提交

#### 后端 路由 (32.14)
- [ ] 32.14.1 **契约**: 路由组 `/api/payments`, `/api/reconciliations`, `/api/refunds`, `/api/invoices`
- [ ] 32.14.2 在 router 中注册财务相关路由
- [ ] 32.14.3 **测试**: `go build ./...`
- [ ] 32.14.4 提交

#### 前端 Payment API (32.15)
- [ ] 32.15.1 **契约**: 收款记录 CRUD
- [ ] 32.15.2 创建 `src/api/payment.js`
- [ ] 32.15.3 **测试**: `npm test`
- [ ] 32.15.4 提交

#### 前端 收款记录页面 (32.16)
- [ ] 32.16.1 **契约**: 收款列表、按订单/用户/日期筛选
- [ ] 32.16.2 创建 `src/views/payment/PaymentList.vue`
- [ ] 32.16.3 **测试**: `npm test`
- [ ] 32.16.4 提交

#### 前端 退款页面 (32.17)
- [ ] 32.17.1 **契约**: 退款列表、申请退款、审批退款
- [ ] 32.17.2 创建 `src/views/refund/RefundList.vue`
- [ ] 32.17.3 创建 `src/views/refund/RefundForm.vue`
- [ ] 32.17.4 **测试**: `npm test`
- [ ] 32.17.5 提交

#### 前端 对账页面 (32.18)
- [ ] 32.18.1 **契约**: 对账单列表、按日期筛选
- [ ] 32.18.2 创建 `src/views/reconciliation/ReconciliationList.vue`
- [ ] 32.18.3 **测试**: `npm test`
- [ ] 32.18.4 提交

#### 前端 发票页面 (32.19)
- [ ] 32.19.1 **契约**: 发票列表、开票申请
- [ ] 32.19.2 创建 `src/views/invoice/InvoiceList.vue`
- [ ] 32.19.3 创建 `src/views/invoice/InvoiceForm.vue`
- [ ] 32.19.4 **测试**: `npm run build`
- [ ] 32.19.5 提交

#### E2E (32.20)
- [ ] 32.20.1 **契约**: 财务管理 E2E 测试
- [ ] 32.20.2 添加 `e2e/payment.spec.js`, `e2e/refund.spec.js`
- [ ] 32.20.3 **测试**: `npx playwright test`
- [ ] 32.20.4 提交

#### 验证 (32.21)
- [ ] 32.21.1 `go test ./...` 通过
- [ ] 32.21.2 `npm run lint` 通过
- [ ] 32.21.3 `npx playwright test e2e/payment.spec.js e2e/refund.spec.js` 通过
- [ ] 32.21.4 文档 checklist 同步
- [ ] 32.21.5 `git status` clean

### Phase 33 价格策略配置 - 详细任务

#### 契约定义 (33.0)
- [ ] 33.0.1 **定价规则 API 路径**:
  ```
  GET    /api/price-rules              # 规则列表
  GET    /api/price-rules/:id           # 规则详情
  POST   /api/price-rules              # 创建规则
  PUT    /api/price-rules/:id           # 更新规则
  DELETE /api/price-rules/:id           # 删除规则
  GET    /api/price-rules/calculate    # 计算价格
  ```
- [ ] 33.0.2 **折扣活动 API 路径**:
  ```
  GET    /api/discounts                # 活动列表
  POST   /api/discounts                # 创建活动
  PUT    /api/discounts/:id             # 更新活动
  PUT    /api/discounts/:id/enable     # 启用活动
  PUT    /api/discounts/:id/disable    # 禁用活动
  ```
- [ ] 33.0.3 **优惠券 API 路径**:
  ```
  GET    /api/coupons                  # 优惠券列表
  POST   /api/coupons                  # 创建优惠券
  POST   /api/coupons/generate          # 生成优惠券码
  POST   /api/coupons/claim            # 用户领取优惠券
  GET    /api/coupons/my               # 我的优惠券
  POST   /api/coupons/validate         # 校验优惠券
  ```
- [ ] 33.0.4 **页面入口**: `/price-rules` - 定价规则; `/discounts` - 折扣活动; `/coupons` - 优惠券
- [ ] 33.0.5 **价格计算**: OrderForm.vue 显示原价/折扣/优惠券/最终价格
- [ ] 33.1.1 **契约**: `PriceRule` 结构体 (id, name, rule_type, product_id, category_id, base_price, unit_price, min_quantity, created_at, updated_at)
- [ ] 33.1.2 **契约**: rule_type: fixed, time-based, quantity-based
- [ ] 33.1.3 创建 `internal/model/price_rule.go`
- [ ] 33.1.4 **测试**: `go build ./...`
- [ ] 33.1.5 提交

#### 后端 定价规则 Repository (33.2)
- [ ] 33.2.1 **契约**: Create, GetByID, Update, Delete, List, ListByProductID, ListByCategoryID
- [ ] 33.2.2 创建 `internal/repository/price_rule.go`
- [ ] 33.2.3 **测试**: `go test ./...`
- [ ] 33.2.4 提交

#### 后端 定价规则 Service (33.3)
- [ ] 33.3.1 **契约**: 计算价格 (product_id, quantity, duration)
- [ ] 33.3.2 创建 `internal/service/price_rule.go`
- [ ] 33.3.3 **测试**: `go test ./...`
- [ ] 33.3.4 提交

#### 后端 定价规则 Handler (33.4)
- [ ] 33.4.1 **契约**: `GET /api/price-rules` 规则列表
- [ ] 33.4.2 **契约**: `GET /api/price-rules/:id` 规则详情
- [ ] 33.4.3 **契约**: `POST /api/price-rules` 创建规则
- [ ] 33.4.4 **契约**: `PUT /api/price-rules/:id` 更新规则
- [ ] 33.4.5 **契约**: `DELETE /api/price-rules/:id` 删除规则
- [ ] 33.4.6 **契约**: `GET /api/price-rules/calculate` 计算价格
- [ ] 33.4.7 创建 `internal/handler/price_rule.go`
- [ ] 33.4.8 **测试**: `go test ./...`
- [ ] 33.4.9 提交

#### 后端 折扣活动 Model (33.5)
- [ ] 33.5.1 **契约**: `Discount` 结构体 (id, name, discount_type, discount_value, min_order_amount, start_at, end_at, status, created_at, updated_at)
- [ ] 33.5.2 **契约**: discount_type: percentage, fixed, gift
- [ ] 33.5.3 **契约**: status: pending, active, expired, disabled
- [ ] 33.5.4 创建 `internal/model/discount.go`
- [ ] 33.5.5 **测试**: `go build ./...`
- [ ] 33.5.6 提交

#### 后端 折扣活动 Service (33.6)
- [ ] 33.6.1 **契约**: 折扣活动 CRUD、启用/禁用
- [ ] 33.6.2 创建 `internal/service/discount.go`
- [ ] 33.6.3 **测试**: `go test ./...`
- [ ] 33.6.4 提交

#### 后端 折扣活动 Handler (33.7)
- [ ] 33.7.1 **契约**: `GET /api/discounts` 活动列表
- [ ] 33.7.2 **契约**: `GET /api/discounts/:id` 活动详情
- [ ] 33.7.3 **契约**: `POST /api/discounts` 创建活动
- [ ] 33.7.4 **契约**: `PUT /api/discounts/:id` 更新活动
- [ ] 33.7.5 **契约**: `DELETE /api/discounts/:id` 删除活动
- [ ] 33.7.6 **契约**: `PUT /api/discounts/:id/enable` 启用活动
- [ ] 33.7.7 **契约**: `PUT /api/discounts/:id/disable` 禁用活动
- [ ] 33.7.8 创建 `internal/handler/discount.go`
- [ ] 33.7.9 **测试**: `go test ./...`
- [ ] 33.7.10 提交

#### 后端 优惠券 Model (33.8)
- [ ] 33.8.1 **契约**: `Coupon` 结构体 (id, code, name, discount_type, discount_value, min_order_amount, total_count, remain_count, per_user_limit, start_at, end_at, status, created_at, updated_at)
- [ ] 33.8.2 **契约**: status: pending, active, expired, disabled
- [ ] 33.8.3 创建 `internal/model/coupon.go`
- [ ] 33.8.4 **测试**: `go build ./...`
- [ ] 33.8.5 提交

#### 后端 优惠券 Service (33.9)
- [ ] 33.9.1 **契约**: 优惠券 CRUD、生成、领取、校验
- [ ] 33.9.2 创建 `internal/service/coupon.go`
- [ ] 33.9.3 **测试**: `go test ./...`
- [ ] 33.9.4 提交

#### 后端 优惠券 Handler (33.10)
- [ ] 33.10.1 **契约**: `GET /api/coupons` 优惠券列表
- [ ] 33.10.2 **契约**: `GET /api/coupons/:id` 优惠券详情
- [ ] 33.10.3 **契约**: `POST /api/coupons` 创建优惠券
- [ ] 33.10.4 **契约**: `PUT /api/coupons/:id` 更新优惠券
- [ ] 33.10.5 **契约**: `DELETE /api/coupons/:id` 删除优惠券
- [ ] 33.10.6 **契约**: `POST /api/coupons/generate` 生成优惠券码
- [ ] 33.10.7 **契约**: `POST /api/coupons/claim` 用户领取优惠券
- [ ] 33.10.8 **契约**: `GET /api/coupons/my` 我的优惠券
- [ ] 33.10.9 **契约**: `POST /api/coupons/validate` 校验优惠券
- [ ] 33.10.10 创建 `internal/handler/coupon.go`
- [ ] 33.10.11 **测试**: `go test ./...`
- [ ] 33.10.12 提交

#### 后端 价格计算 (33.11)
- [ ] 33.11.1 **契约**: 创建订单时自动计算最优价格(定价规则+折扣+优惠券)
- [ ] 33.11.2 修改 order service 添加价格计算逻辑
- [ ] 33.11.3 **测试**: `go test ./...`
- [ ] 33.11.4 提交

#### 后端 路由 (33.12)
- [ ] 33.12.1 **契约**: 路由组 `/api/price-rules`, `/api/discounts`, `/api/coupons`
- [ ] 33.12.2 在 router 中注册价格策略路由
- [ ] 33.12.3 **测试**: `go build ./...`
- [ ] 33.12.4 提交

#### 前端 PriceRule API (33.13)
- [ ] 33.13.1 **契约**: 定价规则 CRUD
- [ ] 33.13.2 创建 `src/api/price-rule.js`
- [ ] 33.13.3 **测试**: `npm test`
- [ ] 33.13.4 提交

#### 前端 定价规则页面 (33.14)
- [ ] 33.14.1 **契约**: 规则列表、新增、编辑、删除
- [ ] 33.14.2 创建 `src/views/price-rule/PriceRuleList.vue`
- [ ] 33.14.3 创建 `src/views/price-rule/PriceRuleForm.vue`
- [ ] 33.14.4 **测试**: `npm test`
- [ ] 33.14.5 提交

#### 前端 折扣活动页面 (33.15)
- [ ] 33.15.1 **契约**: 活动列表、新增、编辑、启用/禁用
- [ ] 33.15.2 创建 `src/views/discount/DiscountList.vue`
- [ ] 33.15.3 创建 `src/views/discount/DiscountForm.vue`
- [ ] 33.15.4 **测试**: `npm test`
- [ ] 33.15.5 提交

#### 前端 优惠券页面 (33.16)
- [ ] 33.16.1 **契约**: 优惠券列表、生成、领取、校验
- [ ] 33.16.2 创建 `src/views/coupon/CouponList.vue`
- [ ] 33.16.3 创建 `src/views/coupon/CouponForm.vue`
- [ ] 33.16.4 **测试**: `npm test`
- [ ] 33.16.5 提交

#### 前端 下单价格计算 (33.17)
- [ ] 33.17.1 **契约**: OrderForm.vue 添加价格计算展示(原价、折扣、优惠券、红包后价格)
- [ ] 33.17.2 修改 `src/views/order/OrderForm.vue`
- [ ] 33.17.3 **测试**: `npx playwright test`
- [ ] 33.17.4 提交

#### E2E (33.18)
- [ ] 33.18.1 **契约**: 价格策略 E2E 测试
- [ ] 33.18.2 添加 `e2e/price-rule.spec.js`, `e2e/coupon.spec.js`
- [ ] 33.18.3 **测试**: `npx playwright test`
- [ ] 33.18.4 提交

#### 验证 (33.19)
- [ ] 33.19.1 `go test ./...` 通过
- [ ] 33.19.2 `npm run lint` 通过
- [ ] 33.19.3 `npx playwright test e2e/price-rule.spec.js e2e/coupon.spec.js` 通过
- [ ] 33.19.4 文档 checklist 同步
- [ ] 33.19.5 `git status` clean
