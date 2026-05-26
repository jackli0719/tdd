# 错误记录

## 2026-05-26: GitHub Actions CI 失败与 E2E 时序问题

### 现象

远程 CI 在 PR `feature/ci-infrastructure` 上失败：

- `build`: `Process completed with exit code 127`
- `test-unit`: `Process completed with exit code 127`
- `test-e2e`: `Process completed with exit code 127`
- `lint`: `Process completed with exit code 3`
- E2E 中品类管理、订单创建出现 UI timing 类失败

本地复现时还观察到：

- `make pre-commit` 在 lint 阶段失败，提示 `internal/handler/stats.go`、`internal/model/order.go` 未 `gofmt`
- Playwright 品类删除偶发找不到确认按钮或等不到 `DELETE /api/categories/:id`
- 订单创建偶发未真正选中产品，下拉选项串到了用户下拉，导致没有发起 `POST /api/orders`
- Playwright 配置里存在本机绝对路径，CI 环境不可用

### 根因

1. CI 依赖安装方式不稳定

   原 workflow 通过 `deps` job 上传 `frontend/node_modules` artifact，再由其他 job 下载复用。`build` 和 `test-unit` job 还缺少独立的 `setup-node`。在 GitHub Actions 的隔离 job 环境里，这容易导致 `npm`、`vite`、`vitest`、`playwright` 等命令不可用，表现为 exit code `127`。

2. lint 失败来自 Go 格式问题

   `oms/internal/handler/stats.go` 和 `oms/internal/model/order.go` 未按 `gofmt` 格式化。本地没有 `golangci-lint` 时会走 `gofmt` fallback；远程 CI 安装了 `golangci-lint`，因此 lint job 直接失败。

3. E2E 用例依赖脆弱 UI 状态

   品类编辑/删除原先点击第一行按钮，依赖已有数据顺序和列表加载时机。订单创建原先使用全局 `.el-select-dropdown__item`，Element Plus 多个下拉层同时存在时可能选到错误下拉。

4. Playwright webServer 配置不可移植

   配置中写死 `/Users/liwei/res/minimaxi/tdd/...`，只能在本机工作，远程 CI 无法使用。

### 修复

- CI workflow 改为每个需要前端的 job 都独立执行：
  - `actions/setup-node@v4`
  - `cache: npm`
  - `cd frontend && npm ci`
- 不再跨 job 上传/下载 `node_modules` 作为依赖恢复方式。
- 对 Go 文件执行 `gofmt`，清除 lint 格式问题。
- Playwright 配置改用相对路径，并允许本地复用已有前端服务。
- 品类 E2E 改为先通过 API 创建独立测试数据，再按具体名称定位对应行。
- 删除确认框定位改为当前可见 `.el-message-box` 的主按钮。
- 订单创建下拉选择改为基于当前 combobox 的 `aria-controls` 精确定位对应 dropdown。
- 提交操作改为等待真实 API response，而不是只等待 toast 文案。

### 验证

本地验证通过：

```bash
cd frontend && npm run test:e2e
```

结果：

```text
14 passed
```

完整门禁通过：

```bash
make pre-commit
```

结果：

```text
All checks passed!
```

### 后续注意

- CI job 之间不要传递 `node_modules`，应使用包管理器缓存和 `npm ci` 保证可复现安装。
- E2E 不要依赖“第一行”“第一个下拉项”等不稳定目标，优先创建独立测试数据并按业务唯一值定位。
- Element Plus 的浮层组件要限定当前可见弹层或使用 `aria-controls` 绑定关系定位。
- Playwright 配置不要写本机绝对路径。
- 提交前必须跑 `make pre-commit`；涉及前端流程时额外跑 `cd frontend && npm run test:e2e`。
