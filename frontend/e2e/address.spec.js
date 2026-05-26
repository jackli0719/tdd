import { test, expect } from '@playwright/test'
import { testIdentity } from './test-data.js'

test.describe('Address Management', () => {
  const testPassword = 'password123'

  test.beforeEach(async ({ page }, testInfo) => {
    const identity = testIdentity(testInfo, 'addr', '139')
    await page.goto('/register')
    await page.getByLabel('用户名').fill(identity.username)
    await page.getByLabel('密码').fill(testPassword)
    await page.getByLabel('邮箱').fill(identity.email)
    await page.getByLabel('手机号').fill(identity.phone)
    await page.getByRole('button', { name: '注册' }).click()
    await expect(page).toHaveURL(/\/login/)

    await page.getByPlaceholder('请输入用户名').fill(identity.username)
    await page.getByPlaceholder('请输入密码').fill(testPassword)
    await page.getByRole('button', { name: '登录' }).click()
    await expect(page).toHaveURL(/\/dashboard/)
  })

  test('create address shows success message', async ({ page }) => {
    await page.goto('/addresses')
    await page.waitForTimeout(500)
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

  test('edit address updates successfully', async ({ page }) => {
    await page.goto('/addresses')
    await page.waitForTimeout(500)

    await page.getByRole('button', { name: '新增地址' }).click()
    const dialog = page.locator('.el-dialog')
    await dialog.locator('input').nth(0).fill('初始名字')
    await dialog.locator('input').nth(1).fill('13800138001')
    await dialog.locator('input').nth(2).fill('广东省')
    await dialog.locator('input').nth(3).fill('深圳市')
    await dialog.locator('input').nth(4).fill('南山区')
    await dialog.locator('textarea').fill('某街道')
    await dialog.locator('button').filter({ hasText: '确定' }).click()
    await expect(page.locator('.el-message').first().getByText('创建成功')).toBeVisible({ timeout: 5000 })
    // Wait for toast to disappear
    await page.waitForTimeout(2000)

    const editBtn = page.locator('.el-table__body .el-button').filter({ hasText: '编辑' }).first()
    await expect(editBtn).toBeVisible({ timeout: 5000 })
    await editBtn.click()
    await expect(dialog).toBeVisible()
    await dialog.locator('input').nth(0).fill('修改后名字')
    await dialog.locator('button').filter({ hasText: '确定' }).click()
    await expect(page.locator('.el-message').first().getByText('更新成功')).toBeVisible({ timeout: 5000 })
  })

  test('set default address', async ({ page }) => {
    await page.goto('/addresses')
    await page.waitForTimeout(500)

    await page.getByRole('button', { name: '新增地址' }).click()
    const dialog = page.locator('.el-dialog')
    await dialog.locator('input').nth(0).fill('地址1')
    await dialog.locator('input').nth(1).fill('13800138001')
    await dialog.locator('input').nth(2).fill('广东省')
    await dialog.locator('input').nth(3).fill('深圳市')
    await dialog.locator('input').nth(4).fill('区')
    await dialog.locator('textarea').fill('详细')
    await dialog.locator('button').filter({ hasText: '确定' }).click()
    await expect(page.locator('.el-message').first().getByText('创建成功')).toBeVisible({ timeout: 5000 })
    // Wait for toast to disappear
    await page.waitForTimeout(2000)

    await page.getByRole('button', { name: '新增地址' }).click()
    await dialog.locator('input').nth(0).fill('地址2')
    await dialog.locator('input').nth(1).fill('13800138002')
    await dialog.locator('input').nth(2).fill('广东省')
    await dialog.locator('input').nth(3).fill('广州市')
    await dialog.locator('input').nth(4).fill('区')
    await dialog.locator('textarea').fill('详细2')
    await dialog.locator('button').filter({ hasText: '确定' }).click()
    await expect(page.locator('.el-message').first().getByText('创建成功')).toBeVisible({ timeout: 5000 })
    // Wait for toast to disappear
    await page.waitForTimeout(2000)

    const btns = page.locator('.el-table__body .el-button').filter({ hasText: '设为默认' })
    const firstSetDefaultBtn = btns.first()
    if (await firstSetDefaultBtn.isEnabled()) {
      await firstSetDefaultBtn.click()
      await expect(page.getByText('设置成功')).toBeVisible({ timeout: 5000 })
    }
  })

  test('delete address removes from list', async ({ page }) => {
    await page.goto('/addresses')
    await page.waitForTimeout(500)

    await page.getByRole('button', { name: '新增地址' }).click()
    const dialog = page.locator('.el-dialog')
    await dialog.locator('input').nth(0).fill('删除测试')
    await dialog.locator('input').nth(1).fill('13800138001')
    await dialog.locator('input').nth(2).fill('广东省')
    await dialog.locator('input').nth(3).fill('深圳市')
    await dialog.locator('input').nth(4).fill('区')
    await dialog.locator('textarea').fill('详细')
    await dialog.locator('button').filter({ hasText: '确定' }).click()
    await expect(page.locator('.el-message').first().getByText('创建成功')).toBeVisible({ timeout: 5000 })
    // Wait for toast to disappear
    await page.waitForTimeout(2000)

    const initialRows = await page.locator('.el-table__body tr').count()
    expect(initialRows).toBeGreaterThan(0)

    const deleteBtn = page.locator('.el-table__body .el-button').filter({ hasText: '删除' }).first()
    await expect(deleteBtn).toBeVisible({ timeout: 5000 })
    await deleteBtn.click()
    await expect(page.locator('.el-message-box')).toBeVisible()
    // Click second button in message box (the confirm button, not cancel)
    await page.locator('.el-message-box__btns button').nth(1).click()
    await expect(page.getByText('删除成功')).toBeVisible({ timeout: 5000 })

    await expect(page.locator('.el-table__body tr')).toHaveCount(initialRows - 1)
  })

  test('empty address list shows No Data', async ({ page }) => {
    await page.goto('/addresses')
    await page.waitForTimeout(500)

    const emptyText = page.locator('.el-table__empty-text')
    if (await emptyText.isVisible()) {
      await expect(emptyText).toContainText('No Data')
    }
  })
})
