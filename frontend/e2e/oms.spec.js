import { test, expect } from '@playwright/test'

test.describe('OMS Frontend E2E', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/dashboard')
  })

  test('health check - backend API', async ({ page }) => {
    const response = await page.request.get('http://localhost:8080/health')
    expect(response.ok()).toBeTruthy()
    const body = await response.json()
    expect(body.status).toBe('ok')
  })

  test('dashboard loads with stats', async ({ page }) => {
    await expect(page.locator('text=用户总数')).toBeVisible()
    await expect(page.locator('text=产品总数')).toBeVisible()
    await expect(page.locator('text=订单总数')).toBeVisible()
  })

  test('dashboard shows current time', async ({ page }) => {
    // Check time card is visible
    const timeCard = page.locator('.current-time')
    await expect(timeCard).toBeVisible()
    // Time format: YYYY/MM/DD HH:mm:ss
    const timeRegex = /^\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2}$/
    await expect(timeCard).toHaveText(timeRegex)
  })

  test('navigation menu works', async ({ page }) => {
    // Use sidebar menu for clicks to avoid matching page title
    const sidebar = page.locator('.sidebar-menu')
    await sidebar.getByRole('menuitem', { name: '用户管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()

    await sidebar.getByRole('menuitem', { name: '产品管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()

    await sidebar.getByRole('menuitem', { name: '订单管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()

    await sidebar.getByRole('menuitem', { name: '统计' }).click()
    await expect(page.locator('text=订单统计')).toBeVisible()
  })

  test('user list page loads', async ({ page }) => {
    await page.locator('.sidebar-menu').getByRole('menuitem', { name: '用户管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()
    await expect(page.locator('text=添加用户')).toBeVisible()
  })

  test('product list page loads', async ({ page }) => {
    await page.locator('.sidebar-menu').getByRole('menuitem', { name: '产品管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()
    await expect(page.locator('text=添加产品')).toBeVisible()
  })

  test('category list page loads', async ({ page }) => {
    await page.locator('.sidebar-menu').getByRole('menuitem', { name: '品类管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()
    await expect(page.locator('text=添加品类')).toBeVisible()
  })

  test('category - create new category', async ({ page }) => {
    await page.locator('.sidebar-menu').getByRole('menuitem', { name: '品类管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()

    // Click add button
    await page.getByRole('button', { name: '添加品类' }).click()
    // Fill form
    await page.getByLabel('名称').fill('测试品类_' + Date.now())
    await page.getByLabel('描述').fill('这是一个测试描述')
    // Submit
    await page.getByRole('button', { name: '提交' }).click()
    // Verify success
    await expect(page.getByText('保存成功')).toBeVisible()
  })

  test('category - edit existing category', async ({ page }) => {
    await page.locator('.sidebar-menu').getByRole('menuitem', { name: '品类管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()

    // Click first edit button
    await page.getByRole('button', { name: '编辑' }).first().click()
    // Modify name
    await page.getByLabel('名称').fill('修改后品类_' + Date.now())
    // Submit
    await page.getByRole('button', { name: '提交' }).click()
    // Verify success
    await expect(page.getByText('保存成功')).toBeVisible()
  })

  test('category - delete category', async ({ page }) => {
    await page.locator('.sidebar-menu').getByRole('menuitem', { name: '品类管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()

    // Get initial row count
    const initialRows = await page.locator('.el-table tbody tr').count()

    // Click first delete button
    await page.getByRole('button', { name: '删除' }).first().click()
    // Confirm dialog
    await page.getByRole('button', { name: '确定' }).click()
    // Verify success
    await expect(page.getByText('删除成功')).toBeVisible()
  })

  test('order list page loads', async ({ page }) => {
    await page.locator('.sidebar-menu').getByRole('menuitem', { name: '订单管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()
    await expect(page.locator('text=创建订单')).toBeVisible()
  })

  test('order - create and confirm order', async ({ page }) => {
    await page.locator('.sidebar-menu').getByRole('menuitem', { name: '订单管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()

    // Create order
    await page.getByRole('button', { name: '创建订单' }).click()
    // Select user
    await page.locator('.el-select').first().click()
    await page.locator('.el-select-dropdown__item').first().click()
    // Add product
    await page.locator('.el-button', { hasText: '添加产品' }).click()
    await page.locator('.el-select').nth(1).click()
    await page.locator('.el-select-dropdown__item').first().click()
    // Submit
    await page.getByRole('button', { hasText: '提交' }).click()
    // Verify success
    await expect(page.getByText('创建成功')).toBeVisible()

    // Confirm order (pending → confirmed)
    const firstPending = page.locator('.el-tag', { hasText: '待确认' }).first()
    if (await firstPending.isVisible()) {
      await page.getByRole('button', { name: '确认' }).first().click()
      await expect(page.getByText('状态更新成功')).toBeVisible()
    }
  })

  test('order - state transitions', async ({ page }) => {
    await page.locator('.sidebar-menu').getByRole('menuitem', { name: '订单管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()

    // Look for confirmed order to start service
    const confirmedTag = page.locator('.el-tag', { hasText: '已确认' }).first()
    if (await confirmedTag.isVisible()) {
      await page.getByRole('button', { name: '开始服务' }).first().click()
      await expect(page.getByText('状态更新成功')).toBeVisible()
    }

    // Look for in_service order to complete
    const inServiceTag = page.locator('.el-tag', { hasText: '服务中' }).first()
    if (await inServiceTag.isVisible()) {
      await page.getByRole('button', { name: '完成' }).first().click()
      await expect(page.getByText('状态更新成功')).toBeVisible()
    }
  })

  test('stats page loads', async ({ page }) => {
    await page.locator('.sidebar-menu').getByRole('menuitem', { name: '统计' }).click()
    await expect(page.locator('text=订单统计')).toBeVisible()
    await expect(page.locator('text=营收统计')).toBeVisible()
  })
})