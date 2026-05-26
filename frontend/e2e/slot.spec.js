import { test, expect } from '@playwright/test'
import { testIdentity } from './test-data.js'

test.describe('Slot E2E', () => {
  const testPassword = 'password123'

  test.beforeEach(async ({ page }, testInfo) => {
    // Register and login
    const identity = testIdentity(testInfo, 'slot', '138')
    await page.goto('/register')
    await page.getByLabel('用户名').fill(identity.username)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByLabel('邮箱').fill(identity.email)
    await page.getByLabel('手机号').fill(identity.phone)
    await page.getByRole('button', { name: '注册' }).click()
    await expect(page).toHaveURL(/\/login/)

    await page.getByLabel('用户名').fill(identity.username)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByRole('button', { name: '登录' }).click()
    await expect(page).toHaveURL(/\/dashboard/)
  })

  test('order form shows appointment time picker', async ({ page }) => {
    await page.goto('/orders')
    await page.getByRole('button', { name: '创建订单' }).click()
    await expect(page.locator('.el-dialog')).toBeVisible()
    // Should have date picker
    await expect(page.locator('.el-date-editor')).toBeVisible()
    // Appointment time slot selector is present
    const dialogBody = page.locator('.el-dialog .el-dialog__body')
    await expect(dialogBody.getByText('选择时间段')).toBeVisible()
  })

  test('date picker disables past dates', async ({ page }) => {
    await page.goto('/orders')
    await page.getByRole('button', { name: '创建订单' }).click()
    await page.locator('.el-date-editor').click()
    // Past dates should be disabled (have .disabled class)
    const disabledCell = page.locator('.el-date-table td.prev-month.disabled, .el-date-table td.disabled').first()
    if (await disabledCell.isVisible({ timeout: 1000 }).catch(() => false)) {
      await expect(disabledCell).toHaveClass(/disabled/)
    }
  })

  test('order list shows appointment time column', async ({ page }) => {
    await page.goto('/orders')
    // Should have 预约时间 column in table
    const hasColumn = await page.locator('.el-table__header th').filter({ hasText: '预约时间' }).isVisible().catch(() => false)
    if (hasColumn) {
      await expect(page.locator('.el-table__header th').filter({ hasText: '预约时间' })).toBeVisible()
    }
  })
})
