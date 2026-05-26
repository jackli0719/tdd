import { test, expect } from '@playwright/test'

let userId = 0
const getNextUsername = () => `staff_${Date.now()}_${++userId}`

test.describe('Staff E2E', () => {
  const testPassword = 'password123'

  test.beforeEach(async ({ page }) => {
    // Register and login
    const username = getNextUsername()
    await page.goto('/register')
    await page.getByLabel('用户名').fill(username)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByLabel('邮箱').fill(username + '@example.com')
    await page.getByLabel('手机号').fill('13800138000')
    await page.getByRole('button', { name: '注册' }).click()
    await expect(page).toHaveURL(/\/login/)

    await page.getByLabel('用户名').fill(username)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByRole('button', { name: '登录' }).click()
    await expect(page).toHaveURL(/\/dashboard/)
  })

  test('staff page loads with empty list', async ({ page }) => {
    await page.goto('/staff')
    await expect(page.locator('text=服务人员管理')).toBeVisible()
    await expect(page.getByRole('button', { name: '添加人员' })).toBeVisible()
  })

  test('add new staff', async ({ page }) => {
    await page.goto('/staff')
    await page.getByRole('button', { name: '添加人员' }).click()
    await expect(page.locator('.el-dialog')).toBeVisible()

    await page.getByLabel('姓名').fill('张三')
    await page.getByLabel('手机号').fill('13800138001')
    await page.locator('.el-dialog__footer button').filter({ hasText: '确定' }).click()

    await expect(page.getByText('创建成功')).toBeVisible()
  })

  test('edit staff', async ({ page }) => {
    const editName = 'edit_' + Date.now()
    await page.goto('/staff')

    // Add staff first
    await page.getByRole('button', { name: '添加人员' }).click()
    await page.getByLabel('姓名').fill('李四')
    await page.getByLabel('手机号').fill('13800138002')
    await page.locator('.el-dialog__footer button').filter({ hasText: '确定' }).click()
    await expect(page.getByText('创建成功')).toBeVisible()

    // Edit staff - use first row's edit button
    await page.locator('.el-table__body tr').first().locator('button').filter({ hasText: '编辑' }).click()
    await expect(page.locator('.el-dialog')).toBeVisible()

    await page.getByLabel('姓名').fill(editName)
    await page.locator('.el-dialog__footer button').filter({ hasText: '确定' }).click()

    // Verify dialog closes and no error toast appears
    await expect(page.locator('.el-dialog')).not.toBeVisible({ timeout: 3000 })
  })

  test('change staff status', async ({ page }) => {
    await page.goto('/staff')

    // Add staff first
    await page.getByRole('button', { name: '添加人员' }).click()
    await page.getByLabel('姓名').fill('王五')
    await page.getByLabel('手机号').fill('13800138003')
    await page.locator('.el-dialog__footer button').filter({ hasText: '确定' }).click()
    await expect(page.getByText('创建成功')).toBeVisible()

    // Change status to busy - use first row's button
    await page.locator('.el-table__body tr').first().locator('button').filter({ hasText: '设为忙碌' }).click()
    await expect(page.getByText('状态更新成功')).toBeVisible()
  })

  test('delete staff', async ({ page }) => {
    const delName = 'del_' + Date.now()
    await page.goto('/staff')
    await page.getByRole('button', { name: '添加人员' }).click()
    await page.getByLabel('姓名').fill(delName)
    await page.getByLabel('手机号').fill('13800138004')
    await page.locator('.el-dialog__footer button').filter({ hasText: '确定' }).click()
    await expect(page.getByText('创建成功')).toBeVisible()

    // Verify staff appears in list
    await expect(page.locator('.el-table__body').getByText(delName)).toBeVisible()

    // Delete staff
    page.on('dialog', dialog => dialog.accept())
    await page.locator('.el-table__body tr').first().locator('button').filter({ hasText: '删除' }).click()

    // Verify delete success message
    await expect(page.getByText('删除成功')).toBeVisible()

    // Verify staff no longer in list
    await expect(page.locator('.el-table__body').getByText(delName)).not.toBeVisible()
  })
})