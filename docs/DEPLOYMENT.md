# OMS 部署文档

## 环境要求

### 必要环境
- Docker 20.10+
- Docker Compose 2.0+
- Git

### 可选环境（本地开发）
- Go 1.25+
- Node.js 22+
- Rust 1.75+ (用于 Rust 服务器)

## 本地部署

### 前置准备

1. 克隆代码仓库
```bash
git clone <repository-url>
cd tdd
```

2. 配置环境变量

**Backend (oms/.env)**
```bash
cp oms/.env.example oms/.env
# 编辑 oms/.env 设置合适的值
```

**Frontend (frontend/.env)**
```bash
cp frontend/.env.example frontend/.env
# 编辑 frontend/.env 设置合适的值
```

### 使用 Docker 部署

1. 构建并启动所有服务
```bash
docker-compose up -d
```

2. 验证服务状态
```bash
docker-compose ps
```

3. 查看日志
```bash
docker-compose logs -f
```

4. 停止服务
```bash
docker-compose down
```

### 本地开发模式

#### 启动 Backend
```bash
cd oms
go run ./cmd/server
```

#### 启动 Frontend
```bash
cd frontend
npm install
npm run dev
```

## Docker 部署

### 快速开始

```bash
# 构建并启动
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f oms
docker-compose logs -f frontend

# 停止服务
docker-compose down
```

### 验证配置

```bash
# 验证 docker-compose 配置
docker-compose config

# 检查端口占用
lsof -i :80
lsof -i :8080
```

### 生产环境注意事项

1. **数据持久化**
   - 默认情况下 SQLite 数据库存储在容器内
   - 生产环境建议使用 MySQL 等外部数据库

2. **环境变量**
   - 生产环境请修改默认密码和端口
   - 参考 `oms/.env.example` 配置

3. **健康检查**
   - Backend 提供 `/health` 端点
   - 使用 `docker-compose ps` 查看健康状态

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| frontend | 80 | Nginx 前端服务 |
| oms | 8080 | Go Backend API |

## 目录结构

```
tdd/
├── docker-compose.yml    # Docker Compose 配置
├── oms/                  # Go Backend
│   ├── Dockerfile
│   ├── .env.example      # 环境变量模板
│   └── cmd/server/       # 入口文件
└── frontend/            # Vue Frontend
    ├── Dockerfile
    ├── nginx.conf       # Nginx 配置
    └── .env.example     # 环境变量模板
```

## 故障排除

### Backend 无法启动
```bash
# 检查日志
docker-compose logs oms

# 检查端口占用
lsof -i :8080
```

### Frontend 无法访问
```bash
# 检查日志
docker-compose logs frontend

# 检查 nginx 配置
docker-compose exec frontend cat /etc/nginx/conf.d/default.conf
```

### 数据库连接问题
```bash
# 检查 DSN 配置
docker-compose exec oms env | grep DSN

# 测试数据库连接
docker-compose exec oms sh -c 'wget -q --spider http://localhost:8080/health'
```