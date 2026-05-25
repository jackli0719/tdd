# OMS API 契约

本文档定义 OMS 系统所有 API 的请求和响应结构。

**Base URL**: `http://localhost:8080`
**API 前缀**: `/api`

**统一响应包装**:
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

**错误响应**:
```json
{
  "code": 400,
  "message": "error description"
}
```

---

## 用户管理

### 用户列表
```
GET /api/users?page=1&page_size=10
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "testuser",
        "email": "test@example.com",
        "phone": "1234567890",
        "created_at": "2026-05-25T10:00:00+08:00",
        "updated_at": "2026-05-25T10:00:00+08:00"
      }
    ],
    "total": 100,
    "page": 1
  }
}
```

### 获取单个用户
```
GET /api/users/:id
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "phone": "1234567890",
    "created_at": "2026-05-25T10:00:00+08:00",
    "updated_at": "2026-05-25T10:00:00+08:00"
  }
}
```

### 创建用户
```
POST /api/users
```

**请求体**:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "phone": "1234567890"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "phone": "1234567890",
    "created_at": "2026-05-25T10:00:00+08:00",
    "updated_at": "2026-05-25T10:00:00+08:00"
  }
}
```

### 更新用户
```
PUT /api/users/:id
```

**请求体**:
```json
{
  "username": "updateduser",
  "email": "updated@example.com",
  "phone": "9876543210"
}
```

**响应**: 同获取单个用户

### 删除用户
```
DELETE /api/users/:id
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

---

## 产品管理

### 产品列表
```
GET /api/products?page=1&page_size=10
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "products": [
      {
        "id": 1,
        "name": "空调清洗",
        "price": 98,
        "stock": 95,
        "created_at": "2026-05-25T10:00:00+08:00",
        "updated_at": "2026-05-25T10:00:00+08:00"
      }
    ],
    "total": 50,
    "page": 1
  }
}
```

### 获取单个产品
```
GET /api/products/:id
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "空调清洗",
    "price": 98,
    "stock": 95,
    "created_at": "2026-05-25T10:00:00+08:00",
    "updated_at": "2026-05-25T10:00:00+08:00"
  }
}
```

### 创建产品
```
POST /api/products
```

**请求体**:
```json
{
  "name": "空调清洗",
  "price": 98,
  "stock": 100
}
```

**响应**: 同获取单个产品

### 更新产品
```
PUT /api/products/:id
```

**请求体**:
```json
{
  "name": "空调深度清洗",
  "price": 128,
  "stock": 80
}
```

**响应**: 同获取单个产品

### 删除产品
```
DELETE /api/products/:id
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

---

## 订单管理

### 订单列表
```
GET /api/orders?page=1&page_size=10
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "orders": [
      {
        "id": 1,
        "order_no": "ORD1779708285563400000",
        "user_id": 1,
        "total_amount": 490,
        "status": "paid",
        "created_at": "2026-05-25T19:24:45.563439+08:00",
        "updated_at": "2026-05-25T19:24:47.217833+08:00",
        "items": [
          {
            "id": 1,
            "order_id": 1,
            "product_id": 1,
            "price": 98,
            "quantity": 5,
            "subtotal": 490,
            "created_at": "2026-05-25T19:24:45.563981+08:00"
          }
        ]
      }
    ],
    "total": 10,
    "page": 1
  }
}
```

### 订单状态说明
| 状态 | 说明 |
|------|------|
| pending | 待支付 |
| paid | 已支付 |
| shipped | 已发货 |
| completed | 已完成 |
| cancelled | 已取消 |

### 状态转换规则
- pending → paid / cancelled
- paid → shipped / cancelled
- shipped → completed
- completed → (终态)
- cancelled → (终态)

### 获取单个订单
```
GET /api/orders/:id
```

**响应**: 同订单列表中的单个订单对象

### 创建订单
```
POST /api/orders
```

**请求体**:
```json
{
  "user_id": 1,
  "items": [
    {
      "product_id": 1,
      "quantity": 5
    }
  ]
}
```

**响应**: 同获取单个订单

### 删除订单
```
DELETE /api/orders/:id
```

**约束**: 只能删除 pending 状态的订单

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

### 订单状态操作
```
POST /api/orders/:id/paid    # 支付
POST /api/orders/:id/ship    # 发货
POST /api/orders/:id/complete # 完成
POST /api/orders/:id/cancel  # 取消
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

---

## 统计

### 订单统计
```
GET /api/stats/orders
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,
    "pending": 10,
    "paid": 20,
    "shipped": 30,
    "completed": 35,
    "cancelled": 5
  }
}
```

### 营收统计
```
GET /api/stats/revenue
```

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_revenue": 50000,
    "pending_revenue": 5000,
    "paid_revenue": 10000,
    "shipped_revenue": 15000,
    "completed_revenue": 20000
  }
}
```

**注意**: cancelled 订单不计入营收

---

## 错误码

| HTTP Status | Code | 说明 |
|-------------|------|------|
| 400 | 400 | 请求参数错误 |
| 404 | 404 | 资源不存在 |
| 409 | 409 | 资源冲突（如用户名已存在） |
| 500 | 500 | 服务器内部错误 |

---

## 字段命名规范

**统一使用 snake_case**:
- JSON 字段: `user_id`, `order_no`, `total_amount`
- API 路径: `/api/users`
- 时间字段: `created_at`, `updated_at`