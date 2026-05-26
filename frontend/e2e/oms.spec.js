import { test, expect } from '@playwright/test'

test.describe('OMS Frontend E2E', () => {
  // Setup: register and login before tests
  test.beforeEach(async ({ page }) => {
    // Navigate to login page first to set up the context
    await page.goto('/login')

    // Register a test user
    const timestamp = Date.now()
    const registerRes = await page.request.post('http://localhost:8080/api/auth/register', {
      data: {
        username: 'testuser_' + timestamp,
        password: 'password123',
        email: 'test' + timestamp + '@example.com',
        phone: '13800138000',
      },
    })
    // Login
    const loginRes = await page.request.post('http://localhost:8080/api/auth/login', {
      data: {
        username: 'testuser_' + timestamp,
        password: 'password123',
      },
    })
    // Set token in localStorage and navigate to dashboard
    if (loginRes.ok()) {
      const loginData = await loginRes.json()
      const token = loginData?.data?.token
      if (token) {
        await page.evaluate((t) => {
          localStorage.setItem('auth_token', t)
        }, token)
      }
    }
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

  test('order list page loads', async ({ page }) => {
    await page.locator('.sidebar-menu').getByRole('menuitem', { name: '订单管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()
    await expect(page.locator('text=创建订单')).toBeVisible()
  })

  test('stats page loads', async ({ page }) => {
    await page.locator('.sidebar-menu').getByRole('menuitem', { name: '统计' }).click()
    await expect(page.locator('text=订单统计')).toBeVisible()
    await expect(page.locator('text=营收统计')).toBeVisible()
  })
})