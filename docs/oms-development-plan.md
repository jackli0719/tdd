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
