import { test, expect } from '@playwright/test'

test.describe('Address Management', () => {
  test.beforeEach(async ({ page }) => {
    // Login first
    await page.goto('/login')
    await page.getByPlaceholder('用户名').fill('admin')
    await page.getByPlaceholder('密码').fill('admin123')
    await page.getByRole('button', { name: '登录' }).click()
    await page.waitForURL('/')
  })

  test('create address', async ({ page }) => {
    // Navigate to address page via user list - first create a user
    await page.goto('/users')
    await page.getByRole('button', { name: '新增用户' }).click()
    const phone = '139' + String(Date.now()).slice(-8)
    await page.getByPlaceholder('请输入用户名').fill('地址测试用户' + Date.now())
    await page.getByPlaceholder('请输入手机号').fill(phone)
    await page.getByPlaceholder('请输入邮箱').fill(`test${Date.now()}@example.com`)
    await page.getByRole('button', { name: '确定' }).click()
    await expect(page.getByText('创建成功')).toBeVisible({ timeout: 5000 })

    // Navigate to address management for that user
    await page.goto('/addresses?user_id=1')

    // Click add address button
    await page.getByRole('button', { name: '新增地址' }).click()

    // Wait for dialog
    const dialog = page.locator('.el-dialog')
    await expect(dialog).toBeVisible()

    // Fill address form
    await dialog.locator('input').first().fill('张三')
    await dialog.locator('input').nth(1).fill('13800138001')
    await dialog.locator('input').nth(2).fill('广东省')
    await dialog.locator('input').nth(3).fill('深圳市')
    await dialog.locator('input').nth(4).fill('南山区')
    await dialog.locator('textarea').fill('某街道某号')

    // Submit
    await dialog.locator('button').filter({ hasText: '确定' }).click()

    // Should show success toast
    await expect(page.getByText('创建成功')).toBeVisible({ timeout: 5000 })
  })

  test('edit address', async ({ page }) => {
    // Navigate to address page
    await page.goto('/addresses?user_id=1')

    // Wait for table to load
    await page.waitForSelector('.el-table')

    // Check if there's an address to edit
    const editBtn = page.locator('.el-table .el-button').filter({ hasText: '编辑' }).first()
    if (await editBtn.isVisible({ timeout: 2000 })) {
      await editBtn.click()

      // Wait for dialog
      const dialog = page.locator('.el-dialog')
      await expect(dialog).toBeVisible()

      // Modify name
      await dialog.locator('input').first().fill('李四修改')

      // Submit
      await dialog.locator('button').filter({ hasText: '确定' }).click()

      // Should show success toast
      await expect(page.getByText('更新成功')).toBeVisible({ timeout: 5000 })
    }
  })

  test('set default address', async ({ page }) => {
    // Navigate to address page
    await page.goto('/addresses?user_id=1')

    // Wait for table to load
    await page.waitForSelector('.el-table')

    // Check if there's a non-default address to set
    const setDefaultBtn = page.locator('.el-table .el-button').filter({ hasText: '设为默认' }).first()
    if (await setDefaultBtn.isVisible({ timeout: 2000 })) {
      await setDefaultBtn.click()

      // Should show success toast
      await expect(page.getByText('设置成功')).toBeVisible({ timeout: 5000 })

      // The button should now be disabled
      await expect(setDefaultBtn).toBeDisabled()
    }
  })

  test('delete address', async ({ page }) => {
    // Navigate to address page
    await page.goto('/addresses?user_id=1')

    // Wait for table to load
    await page.waitForSelector('.el-table')

    // Get initial row count
    const initialRows = await page.locator('.el-table__body tr').count()

    // Check if there's an address to delete
    const deleteBtn = page.locator('.el-table .el-button').filter({ hasText: '删除' }).first()
    if (await deleteBtn.isVisible({ timeout: 2000 })) {
      await deleteBtn.click()

      // Confirm in dialog
      await expect(page.locator('.el-message-box')).toBeVisible()
      await page.locator('.el-message-box__btns button').filter({ hasText: '确定' }).click()

      // Should show success toast
      await expect(page.getByText('删除成功')).toBeVisible({ timeout: 5000 })

      // Row count should decrease
      await expect(page.locator('.el-table__body tr')).toHaveCount(initialRows - 1)
    }
  })
})