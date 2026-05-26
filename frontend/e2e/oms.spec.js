import { test, expect } from '@playwright/test'

const waitForCategories = (page) =>
  page.waitForResponse((response) =>
    response.url().includes('/api/categories') &&
    response.request().method() === 'GET' &&
    response.ok()
  )

const gotoCategoryList = async (page) => {
  const categoriesResponse = waitForCategories(page)
  await page.locator('.sidebar-menu').getByRole('menuitem', { name: '品类管理' }).click()
  await categoriesResponse
  await expect(page.locator('.el-table')).toBeVisible()
}

const createCategoryViaApi = async (page, prefix = '测试品类') => {
  const name = `${prefix}_${Date.now()}`
  const response = await page.request.post('/api/categories', {
    data: {
      name,
      description: 'E2E 测试品类',
    },
  })
  expect(response.ok()).toBeTruthy()
  return name
}

const selectFirstVisibleOption = async (page, select, selectedText) => {
  const combobox = select.getByRole('combobox').first()
  await select.click()
  const controls = await combobox.getAttribute('aria-controls')
  const options = controls
    ? page.locator(`[id="${controls}"] .el-select-dropdown__item`)
    : page.locator('.el-select-dropdown:visible .el-select-dropdown__item')
  const option = options.first()
  await expect(option).toBeVisible()
  const text = (await option.textContent())?.trim()
  await option.click({ force: true })
  if (selectedText && text) {
    await expect(select).toContainText(text)
  }
}

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
    await gotoCategoryList(page)
    await expect(page.locator('text=添加品类')).toBeVisible()
  })

  test('category - create new category', async ({ page }) => {
    await gotoCategoryList(page)

    const name = '测试品类_' + Date.now()
    await page.getByRole('button', { name: '添加品类' }).click()
    await page.getByLabel('名称').fill(name)
    await page.getByLabel('描述').fill('这是一个测试描述')

    const createResponse = page.waitForResponse((response) =>
      response.url().includes('/api/categories') &&
      response.request().method() === 'POST' &&
      response.ok()
    )
    await page.getByRole('button', { name: '提交' }).click()
    await createResponse

    await expect(page.getByText('保存成功')).toBeVisible()
    await expect(page.locator('.el-table').getByText(name)).toBeVisible()
  })

  test('category - edit existing category', async ({ page }) => {
    const originalName = await createCategoryViaApi(page, '待编辑品类')
    await gotoCategoryList(page)

    const row = page.locator('.el-table__row', { hasText: originalName })
    await expect(row).toBeVisible()
    await row.getByRole('button', { name: '编辑' }).click()

    const updatedName = '修改后品类_' + Date.now()
    await page.getByLabel('名称').fill(updatedName)
    const updateResponse = page.waitForResponse((response) =>
      response.url().includes('/api/categories/') &&
      response.request().method() === 'PUT' &&
      response.ok()
    )
    await page.getByRole('button', { name: '提交' }).click()
    await updateResponse

    await expect(page.getByText('保存成功')).toBeVisible()
    await expect(page.locator('.el-table').getByText(updatedName)).toBeVisible()
  })

  test('category - delete category', async ({ page }) => {
    const categoryName = await createCategoryViaApi(page, '待删除品类')
    await gotoCategoryList(page)

    const row = page.locator('.el-table__row', { hasText: categoryName })
    await expect(row).toBeVisible()

    await row.scrollIntoViewIfNeeded()
    await row.getByRole('button', { name: '删除' }).click()
    const messageBox = page.locator('.el-message-box:visible')
    await expect(messageBox).toBeVisible()

    const deleteResponse = page.waitForResponse((response) =>
      response.url().includes('/api/categories/') &&
      response.request().method() === 'DELETE' &&
      response.ok()
    )
    await messageBox.locator('.el-button--primary').click()
    await deleteResponse

    await expect(page.getByText('删除成功')).toBeVisible()
    await expect(row).toHaveCount(0)
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
    const usersResponse = page.waitForResponse((response) =>
      response.url().includes('/api/users') &&
      response.request().method() === 'GET' &&
      response.ok()
    )
    const productsResponse = page.waitForResponse((response) =>
      response.url().includes('/api/products') &&
      response.request().method() === 'GET' &&
      response.ok()
    )
    await page.getByRole('button', { name: '创建订单' }).click()
    await usersResponse
    await productsResponse

    // Select user
    await selectFirstVisibleOption(page, page.locator('.el-dialog:visible .el-select').first(), true)
    // Select product from the initial product row
    const productSelect = page.locator('.el-dialog:visible .product-item-row .el-select').first()
    await selectFirstVisibleOption(page, productSelect, true)
    await expect(productSelect).not.toContainText('选择产品')

    // Submit
    const createResponse = page.waitForResponse((response) =>
      response.url().includes('/api/orders') &&
      response.request().method() === 'POST' &&
      response.ok()
    )
    await page.locator('.el-dialog:visible').getByRole('button', { name: '提交' }).click()
    await createResponse
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
