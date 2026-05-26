import { test, expect } from '@playwright/test'

test.describe('Auth E2E', () => {
  const timestamp = Date.now()
  const testUsername = 'authuser_' + timestamp
  const testEmail = 'auth' + timestamp + '@example.com'
  const testPhone = '13800138000'
  const testPassword = 'password123'

  test('register page loads', async ({ page }) => {
    await page.goto('/register')
    await expect(page.locator('text=用户注册')).toBeVisible()
    await expect(page.getByLabel('用户名')).toBeVisible()
    await expect(page.getByLabel('密码')).toBeVisible()
    await expect(page.getByLabel('邮箱')).toBeVisible()
    await expect(page.getByLabel('手机号')).toBeVisible()
  })

  test('login page loads', async ({ page }) => {
    await page.goto('/login')
    await expect(page.locator('text=用户登录')).toBeVisible()
    await expect(page.getByLabel('用户名')).toBeVisible()
    await expect(page.getByLabel('密码')).toBeVisible()
  })

  test('register and login flow', async ({ page }) => {
    // Register
    await page.goto('/register')
    await page.getByLabel('用户名').fill(testUsername)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByLabel('邮箱').fill(testEmail)
    await page.getByLabel('手机号').fill(testPhone)
    await page.getByRole('button', { name: '注册' }).click()

    // Should redirect to login after success
    await expect(page).toHaveURL(/\/login/)

    // Login
    await page.goto('/login')
    await page.getByLabel('用户名').fill(testUsername)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByRole('button', { name: '登录' }).click()

    // Should redirect to dashboard after success
    await expect(page).toHaveURL(/\/dashboard/)
  })

  test('authenticated pages load protected module data', async ({ page }) => {
    await page.goto('/register')
    const username = 'moduleuser_' + timestamp
    const email = 'module' + timestamp + '@example.com'
    await page.getByLabel('用户名').fill(username)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByLabel('邮箱').fill(email)
    await page.getByLabel('手机号').fill(testPhone)
    await page.getByRole('button', { name: '注册' }).click()
    await expect(page).toHaveURL(/\/login/)

    await page.getByLabel('用户名').fill(username)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByRole('button', { name: '登录' }).click()
    await expect(page).toHaveURL(/\/dashboard/)

    const sidebar = page.locator('.sidebar-menu')
    const loadingErrors = page.locator('text=/加载.*失败/')

    await expect(page.locator('text=用户总数')).toBeVisible()
    await expect(loadingErrors).toHaveCount(0)

    await sidebar.getByRole('menuitem', { name: '用户管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()
    await expect(loadingErrors).toHaveCount(0)

    await sidebar.getByRole('menuitem', { name: '品类管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()
    await expect(loadingErrors).toHaveCount(0)

    await sidebar.getByRole('menuitem', { name: '产品管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()
    await expect(loadingErrors).toHaveCount(0)

    await sidebar.getByRole('menuitem', { name: '订单管理' }).click()
    await expect(page.locator('.el-table')).toBeVisible()
    await expect(loadingErrors).toHaveCount(0)

    await sidebar.getByRole('menuitem', { name: '统计' }).click()
    await expect(page.locator('text=订单统计')).toBeVisible()
    await expect(loadingErrors).toHaveCount(0)
  })

  test('login with wrong password fails', async ({ page }) => {
    // Register first
    await page.goto('/register')
    await page.getByLabel('用户名').fill('wrongtest_' + timestamp)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByLabel('邮箱').fill('wrong' + timestamp + '@example.com')
    await page.getByLabel('手机号').fill(testPhone)
    await page.getByRole('button', { name: '注册' }).click()

    await page.goto('/login')
    await page.getByLabel('用户名').fill('wrongtest_' + timestamp)
    await page.getByLabel('密码').fill('wrongpassword')
    await page.getByRole('button', { name: '登录' }).click()

    // Should show error message
    await expect(page.getByText('用户名或密码错误')).toBeVisible()
  })

  test('unauthenticated user redirected to login', async ({ page }) => {
    // Navigate to login page first, then clear localStorage
    await page.goto('/login')
    await page.evaluate(() => localStorage.clear())
    await page.goto('/dashboard')
    // Should redirect to login
    await expect(page).toHaveURL(/\/login/)
  })

  test('logout clears token and redirects to login', async ({ page }) => {
    // Register and login
    await page.goto('/register')
    const username = 'logoutuser_' + timestamp
    await page.getByLabel('用户名').fill(username)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByLabel('邮箱').fill('logout' + timestamp + '@example.com')
    await page.getByLabel('手机号').fill(testPhone)
    await page.getByRole('button', { name: '注册' }).click()
    await expect(page).toHaveURL(/\/login/)

    await page.getByLabel('用户名').fill(username)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByRole('button', { name: '登录' }).click()
    await expect(page).toHaveURL(/\/dashboard/)

    // Verify logged in
    await expect(page.locator('text=用户总数')).toBeVisible()

    // Click logout button
    await page.getByRole('button', { name: '退出登录' }).click()

    // Should redirect to login
    await expect(page).toHaveURL(/\/login/)

    // Token should be cleared - trying to access protected page should redirect
    await page.goto('/dashboard')
    await expect(page).toHaveURL(/\/login/)
  })
})
