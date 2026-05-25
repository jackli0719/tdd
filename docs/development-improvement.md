# OMS 开发改进建议

## 背景

本次 OMS 项目开发中，虽然整体流程规范（任务驱动 → 契约先行 → 测试验证 → 实现提交），但仍出现了大量错误。根本原因是缺乏足够的集成验证和跨 Agent 协调机制。

---

## 发现的根本问题

### 1. API 契约缺失

**问题**：后端和前端并行开发，但没有先定义好 API 的响应结构和字段名。

**表现**：
- 后端返回 `{ code, message, data: { orders: [...] } }`
- 前端期望 `res.data` 或 `res.orders` 不确定
- 字段名不一致：`total_price` vs `total_amount`，`today/week/month` vs `total_revenue/pending_revenue/...`

**改进**：
- 每个模块开发前，先编写 API 契约文档，包含：
  - 请求路径和方法
  - 响应结构（含 JSON 示例）
  - 字段类型和命名规范
  - 错误码定义
- 契约需前后端双方签字确认后再开发

---

### 2. 缺乏集成验证

**问题**：每个 Phase 完成后只做后端单元测试，没有前后端联调测试。

**表现**：
- 后端 API 返回的字段名与前端期望不一致，上线前才发现
- 列表数据解析错误长期未被发现
- 统计页字段名错了多个版本

**改进**：
- 每个 Phase 完成后，执行简单的联调测试（curl 或 Postman）
- 前后端各建一个 `api-contract.md`，定时同步字段变更
- 集成测试应纳入 CI/CD 流水线

---

### 3. 多 Agent 开发协调不足

**问题**：Phase 8-14 是多个 Agent 并行开发后端和前端，但没有定期同步机制。

**表现**：
- 前端 Agent 不知道后端实际返回的字段名
- 后端改了字段名，前端没有及时更新
- 没有统一的字段命名标准检查工具

**改进**：
- Agent 分配时明确主责模块和依赖模块
- 每天定时同步接口变更（简短报告）
- 重要变更（如字段改名）需通知相关 Agent
- 使用共享的接口文档（如 API.md）

---

### 4. 命名规范未严格执�行

**问题**：`docs/naming-convention.md` 定义了 snake_case，但实际开发中仍有遗漏。

**表现**：
- 后端返回 `total_revenue`（正确）
- 前端用了 `total_price`（错误）
- Dashboard 用 `today/week/month`（错误）

**改进**：
- 后端：使用 `golint` 或 `gofmt` 检查代码风格
- 前端：在 ESLint 中配置 `camelCase` vs `snake_case` 转换规则
- CI 中加入命名规范检查步骤

---

### 5. 缺少前端 API Mock

**问题**：前端开发时没有 mock 后端响应，导致字段名错误无法提前发现。

**表现**：
- 前端组件开发时无法验证响应解析是否正确
- 必须在后端运行后才能发现问题

**改进**：
- 前端每个 API 应该有 mock 数据文件（如 `__mocks__/api.js`）
- 使用 MSW (Mock Service Worker) 拦截 API 请求
- Vitest 测试中应该包含 API 解析逻辑的单元测试

---

## 推荐的开发流程

```
1. 需求分析 → 输出任务清单
           ↓
2. API 契约设计 → 输出 api-contract.md（前后端签字确认）
           ↓
3. 并行开发（Agent A 后端，Agent B 前端 mock 数据）
           ↓
4. 后端单元测试 + 前端 mock 测试
           ↓
5. 集成验证（前后端联调）
           ↓
6. E2E 测试（Playwright）
           ↓
7. 提交代码
```

---

## API 契约模板

```markdown
# [模块名] API 契约

## 用户列表
### GET /api/users

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 10 |

**成功响应 (200)**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "test",
        "email": "test@example.com",
        "phone": "1234567890",
        "created_at": "2026-05-25T10:00:00+08:00"
      }
    ],
    "total": 100,
    "page": 1
  }
}
```

**字段说明**：
| 字段 | 类型 | 说明 |
|------|------|------|
| users | array | 用户列表 |
| total | int | 总记录数 |
| page | int | 当前页码 |

**错误响应**：
- 500: 服务器内部错误
```

---

## 命名规范检查工具

### 后端（Go）
- 使用 `golint` 检查命名规范
- 在 `go.mod tidy` 后运行检查

### 前端（Vue/JS）
- ESLint 配置 `snake_case` 命名规则
- Vitest 测试中包含字段名验证

---

## 记忆库更新

已将本次错误记录到 `memory/oms-common-errors.md`，下次开发新项目时应：
1. 先阅读 `naming-convention.md` 和 `oms-common-errors.md`
2. 开发新功能前检查是否涉及已知的错误模式
3. 遵守 API 契约设计流程