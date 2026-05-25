# OMS 命名规范

## 一、总体原则

**统一命名风格：snake_case（蛇形命名法）**

适用于：数据库、API、Go 后端、TypeScript 前端

---

## 二、数据库命名

| 类型 | 规范 | 示例 |
|------|------|------|
| 表名 | snake_case, 复数名词 | `users`, `order_items` |
| 字段名 | snake_case | `user_id`, `order_no` |
| 主键 | `id` | 自增主键 |
| 外键 | `{table}_id` | `user_id`, `product_id` |
| 索引 | `idx_{table}_{field}` | `idx_orders_user_id` |
| 唯一索引 | `uk_{table}_{field}` | `uk_users_username` |
| 外键约束 | `fk_{table}_{ref}` | `fk_orders_user` |

---

## 三、API 命名

| 类型 | 规范 | 示例 |
|------|------|------|
| URL 路径 | `/api/{resource}` | `/api/users`, `/api/orders` |
| REST 方法 | GET/POST/PUT/DELETE | GET /api/users |
| 请求参数 | snake_case | `user_id`, `order_no` |
| 响应字段 | snake_case | `{ "user_id": 1 }` |

### 3.1 路由示例

```
GET    /api/users           # 列表
GET    /api/users/:id        # 详情
POST   /api/users            # 创建
PUT    /api/users/:id        # 更新
DELETE /api/users/:id        # 删除
```

---

## 四、Go 后端命名

| 类型 | 规范 | 示例 |
|------|------|------|
| 包名 | 小写，单词 | `handler`, `service` |
| 结构体 | PascalCase | `User`, `OrderItem` |
| 变量/函数 | camelCase | `userID`, `CreateUser` |
| 数据库字段 | snake_case | `user_id`, `order_no` |
| JSON 字段 | snake_case | `json:"user_id"` |
| 表名 | 复数 snake_case | `users`, `order_items` |

### 4.1 命名映射示例

```go
// 数据库 -> Go 结构体
user_id      -> UserID
order_no     -> OrderNo
created_at   -> CreatedAt

// Go -> JSON
UserID       -> "user_id"
OrderNo      -> "order_no"
```

---

## 五、TypeScript 前端命名

| 类型 | 规范 | 示例 |
|------|------|------|
| 接口名 | PascalCase | `User`, `OrderItem` |
| 变量/函数 | camelCase | `userId`, `orderNo` |
| API 路径 | snake_case | `/api/users` |
| 请求参数 | snake_case | `user_id`, `order_no` |
| 响应数据 | snake_case | `{ user_id: 1 }` |

---

## 六、命名检查清单

- [ ] 数据库表名使用复数：`users` 不是 `user`
- [ ] 字段名全小写：`user_id` 不是 `userId`
- [ ] API 路径用 `-` 分隔：`/api/order-items`
- [ ] Go 结构体用 PascalCase
- [ ] JSON 字段用 snake_case

---

## 七、跨语言转换对照

| 数据库 | Go | TypeScript | JSON |
|--------|-----|------------|------|
| user_id | UserID | userId | user_id |
| order_no | OrderNo | orderNo | order_no |
| created_at | CreatedAt | createdAt | created_at |
