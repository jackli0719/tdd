# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目: OMS 订单管理系统

Go + Gin 全功能订单管理系统，支持用户、产品、订单 CRUD 及统计。

## 技术栈

- Go 1.21+
- Gin (HTTP 框架)
- GORM (ORM)
- MySQL 8.0+

## 环境要求

```bash
go version  # >= 1.21
mysql --version
```

## 开发命令

```bash
# 安装依赖
go mod tidy

# 编译
go build ./...

# 运行
go run cmd/server/main.go

# 测试
curl http://localhost:8080/api/users
```

## 项目结构

```
oms/
├── cmd/server/main.go     # 入口
├── internal/
│   ├── handler/           # HTTP 处理 (order/user/product)
│   ├── service/           # 业务逻辑
│   ├── repository/        # 数据访问
│   ├── model/             # 数据模型
│   ├── router/            # 路由
│   └── config/            # 配置
└── docs/
    └── oms-development-plan.md  # 开发文档
```

## 开发文档

- `docs/oms-development-plan.md` - 任务清单、开发计划
- `docs/naming-convention.md` - 命名规范（必须遵循）
- `docs/development-improvement.md` - 开发改进建议（含 API 契约模板）

## 命名规范

**统一使用 snake_case**：
- 数据库表/字段：`user_id`, `order_no`
- API 路径：`/api/users`
- JSON 字段：`{"user_id": 1}`
- Go 结构体字段：`UserID`, `OrderNo`（JSON tag 用 snake_case）

## 多 Agent 开发

使用分层代理模式，参考 docs/oms-development-plan.md 的 Agent 分工表。

## 注意事项

**开发前必读**：
1. 先阅读 `docs/naming-convention.md` 了解命名规范
2. 前后端接口字段名必须保持一致（参考后端 handler 返回的 JSON 结构）
3. 使用 `go build ./...` 和 `go test ./...` 验证后再提交
4. 重要：axios interceptor 返回 `response.data`，后端响应是 `{ code, message, data: {...} }`，取列表数据用 `res.data.orders` 等

## 常见错误

参考 `memory/oms-common-errors.md` 避免重复犯错。
