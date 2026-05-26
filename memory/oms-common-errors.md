# OMS 常见错误记录

记录开发中犯过的错误，避免重复犯错。

---

## 1. 地址管理模块 (Phase 26)

### 1.1 前端路由未注册
**问题**: 创建了 `AddressList.vue` 但忘记在 `router/index.js` 注册 `/addresses` 路由，导致 `page.goto('/addresses')` 404。

**教训**: 新建页面组件后，必须同时注册路由和菜单。

**预防**: 每次新增 view，都同步检查 router/index.js 和 Layout.vue menu。

### 1.2 组件复用时 props 变化未监听
**问题**: `AddressList.vue` 用 `props.userId` 加载地址，但 watch 缺失导致切换用户时数据不刷新。

**教训**: 通过 prop 传递并在内部使用的值，如果会影响 API 请求，必须 watch。

**预防**:
```javascript
watch(() => props.userId, (newId) => {
  effectiveUserId.value = newId
  loadAddresses()
})
```

### 1.3 Playwright toast 冲突 (strict mode)
**问题**: 连续创建两个地址，页面同时显示两个"创建成功" toast，`page.getByText('创建成功')` 触发 strict mode 错误。

**教训**: 当操作结果可能重复时，用 `.first()` 或更精确的定位。

**解决**:
```javascript
// 不好：strict mode 会报错
await expect(page.getByText('创建成功')).toBeVisible()

// 好：限定到第一个 toast
await expect(page.locator('.el-message').first().getByText('创建成功')).toBeVisible()

// 或者等待 toast 消失后再进行下一步
await page.waitForTimeout(2000)
```

### 1.4 Playwright message-box 按钮定位失败
**问题**: `page.locator('.el-message-box__btns button').filter({ hasText: '确定' })` 超时，Element Plus 的 MessageBox 有特殊 DOM 结构。

**教训**: message-box 是通过 Portal 渲染到 body 外层，定位器需要考虑实际 DOM 结构。

**解决**:
```javascript
// 尝试多种定位方式
await page.locator('.el-message-box__btns button').nth(1).click()  // 第二个按钮是确认
await page.locator('.el-button').filter({ hasText: '确定', hasNotText: '取消' })
await page.getByRole('button', { name: '确定', exact: true }).click({ force: true })
```

### 1.5 E2E 测试数据隔离 (user_id)
**问题**: E2E 测试用 `?user_id=1` 硬编码访问地址，但用户1未必是当前登录用户，导致看到别人的数据。

**教训**: 认证后的页面访问应该使用当前登录用户的数据，通过 API 获取 user_id。

**解决**:
```javascript
// 从 /api/auth/me 获取当前用户 ID
const userId = await getUserId()
await getAddresses(userId)
```

---

## 2. 认证模块 (Phase 24)

### 2.1 setupAuthGuard 重复调用
**问题**: `main.js` 和 `router/index.js` 都调用了 `setupAuthGuard(router)`，虽然测试没失败，但是冗余。

**教训**: guard 只需要在 router 配置处调用一次。

**预防**: 项目规范中明确 guard 只在一处初始化。

---

## 3. 服务人员模块 (Phase 27)

### 3.1 E2E 并发测试手机号冲突
**问题**: beforeEach 用固定手机号 `13800138000`，多个 worker 并发注册冲突。

**教训**: 并行测试中测试数据必须唯一，用 workerIndex + timestamp 组合。

**解决**:
```javascript
const phone = `138${String(testInfo.workerIndex).padStart(2,'0')}${String(Date.now()).slice(-6)}`
```

### 3.2 连续创建后 toast 残留导致 DOM 污染
**问题**: 快速连续创建两个 staff，第二个 staff 的 toast 还在显示时，删除第一个就会有两个 toast 残留。

**教训**: 在需要断言结果的操作之间加足够的等待，让 toast 自然消失。

---

## 4. 通用前端问题

### 4.1 v-model 绑定在 prop 上无效
**问题**: `<DatePicker v-model="modelValue">` 绑定到 props.modelValue，Vue 报 warn 但不报错。

**教训**: prop 是只读的，v-model 需要绑定到本地的 ref。

**解决**:
```vue
<!-- 不好 -->
<DatePicker v-model="modelValue">

<!-- 好 -->
<DatePicker :model-value="modelValue" @update:model-value="val => emit('update:modelValue', val)">
```

### 4.2 组件内使用 useRoute() 但没有引入
**问题**: `AddressList.vue` 引入了 `useRoute` 但只用了 `route.query`，导致 `props.userId` 被忽略。

**教训**: 组件优先使用明确的 prop，避免隐式依赖 URL query。

---

## 5. 常见错误模式总结

| 错误类型 | 预防措施 |
|---------|---------|
| 路由/菜单漏注册 | 新增 view 时同步检查 router 和 Layout |
| prop 变化不监听 | 凡是影响 API 请求的 prop 都要 watch |
| 并发测试数据冲突 | 用 workerIndex + timestamp 保证唯一性 |
| toast 冲突 | 操作间加 wait，或用 `.first()` 限定 |
| message-box 定位失败 | 用 nth(1) 选确认按钮，或 force click |