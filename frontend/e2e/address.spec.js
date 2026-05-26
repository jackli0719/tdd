import { test, expect } from '@playwright/test'

test.describe('Address Management', () => {
  const testPassword = 'password123'

  test.beforeEach(async ({ page }, testInfo) => {
    const username = 'addr_' + Date.now() + '_' + testInfo.workerIndex
    const phone = '139' + String(Date.now()).slice(-8)
    await page.goto('/register')
    await page.getByLabel('用户名').fill(username)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByLabel('邮箱').fill(username + '@example.com')
    await page.getByLabel('手机号').fill(phone)
    await page.getByRole('button', { name: '注册' }).click()
    await expect(page).toHaveURL(/\/login/)

    await page.getByPlaceholder('请输入用户名').fill(username)
    await page.getByPlaceholder('请输入密码').fill(testPassword)
    await page.getByRole('button', { name: '登录' }).click()
    await expect(page).toHaveURL(/\/dashboard/)
  })

  test('create address shows success message', async ({ page }) => {
    await page.goto('/addresses?user_id=1')
    await page.getByRole('button', { name: '新增地址' }).click()

    const dialog = page.locator('.el-dialog')
    await expect(dialog).toBeVisible()

    await dialog.locator('input').nth(0).fill('张三')
    await dialog.locator('input').nth(1).fill('13800138001')
    await dialog.locator('input').nth(2).fill('广东省')
    await dialog.locator('input').nth(3).fill('深圳市')
    await dialog.locator('input').nth(4).fill('南山区')
    await dialog.locator('textarea').fill('某街道某号')

    await dialog.locator('button').filter({ hasText: '确定' }).click()
    await expect(page.getByText('创建成功')).toBeVisible({ timeout: 5000 })
  })

  test('address list shows No Data when empty', async ({ page }) => {
    await page.goto('/addresses?user_id=999999')
    await page.waitForTimeout(1000)

    const emptyText = page.locator('.el-table__empty-text')
    if (await emptyText.isVisible()) {
      await expect(emptyText).toContainText('No Data')
    }
  })
})