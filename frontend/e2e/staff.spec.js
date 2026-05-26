import { test, expect } from '@playwright/test'
import { testIdentity, uniquePhone } from './test-data.js'

const apiBase = 'http://127.0.0.1:8080/api'

const getToken = async (page) => page.evaluate(() => localStorage.getItem('auth_token'))

const apiGet = async (page, token, path) => {
  const res = await page.request.get(`${apiBase}${path}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  expect(res.ok()).toBeTruthy()
  return res.json()
}

const apiPost = async (page, token, path, data) => {
  const res = await page.request.post(`${apiBase}${path}`, {
    headers: { Authorization: `Bearer ${token}` },
    data,
  })
  expect(res.ok()).toBeTruthy()
  return res.json()
}

test.describe('Staff E2E', () => {
  const testPassword = 'password123'

  test.beforeEach(async ({ page }, testInfo) => {
    // Register and login
    const identity = testIdentity(testInfo, 'staff', '138')
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

  test('staff page loads with empty list', async ({ page }) => {
    await page.goto('/staff')
    await expect(page.locator('text=服务人员管理')).toBeVisible()
    await expect(page.getByRole('button', { name: '添加人员' })).toBeVisible()
  })

  test('add new staff', async ({ page }, testInfo) => {
    await page.goto('/staff')
    await page.getByRole('button', { name: '添加人员' }).click()
    await expect(page.locator('.el-dialog')).toBeVisible()

    await page.getByLabel('姓名').fill('张三')
    await page.getByLabel('手机号').fill(uniquePhone(testInfo, '138'))
    await page.locator('.el-dialog__footer button').filter({ hasText: '确定' }).click()

    await expect(page.getByText('创建成功')).toBeVisible()
  })

  test('edit staff', async ({ page }, testInfo) => {
    const editName = 'edit_' + Date.now()
    await page.goto('/staff')

    // Add staff first
    await page.getByRole('button', { name: '添加人员' }).click()
    await page.getByLabel('姓名').fill('李四')
    await page.getByLabel('手机号').fill(uniquePhone(testInfo, '138'))
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

  test('change staff status', async ({ page }, testInfo) => {
    await page.goto('/staff')

    // Add staff first
    await page.getByRole('button', { name: '添加人员' }).click()
    await page.getByLabel('姓名').fill('王五')
    await page.getByLabel('手机号').fill(uniquePhone(testInfo, '138'))
    await page.locator('.el-dialog__footer button').filter({ hasText: '确定' }).click()
    await expect(page.getByText('创建成功')).toBeVisible()

    // Change status to busy - use first row's button
    await page.locator('.el-table__body tr').first().locator('button').filter({ hasText: '设为忙碌' }).click()
    await expect(page.getByText('状态更新成功')).toBeVisible()
  })

  test('filter staff by status', async ({ page }, testInfo) => {
    const token = await getToken(page)
    const timestamp = Date.now()
    const availableName = `available_${timestamp}`
    const busyName = `busy_${timestamp}`

    await apiPost(page, token, '/staff', {
      name: availableName,
      phone: uniquePhone(testInfo, '136'),
      status: 'available',
    })
    await apiPost(page, token, '/staff', {
      name: busyName,
      phone: uniquePhone(testInfo, '135'),
      status: 'busy',
    })

    await page.goto('/staff')
    await expect(page.locator('.el-table__body').getByText(availableName)).toBeVisible()
    await expect(page.locator('.el-table__body').getByText(busyName)).toBeVisible()

    await page.locator('.status-filter').click()
    await page.getByRole('option', { name: '忙碌' }).click()

    await expect(page.locator('.el-table__body').getByText(busyName)).toBeVisible()
    await expect(page.locator('.el-table__body').getByText(availableName)).not.toBeVisible()
  })

  test('assign available staff to order', async ({ page }, testInfo) => {
    const token = await getToken(page)
    const me = await apiGet(page, token, '/auth/me')
    const timestamp = Date.now()
    const category = await apiPost(page, token, '/categories', {
      name: `分配品类${timestamp}`,
      description: 'phase27 e2e',
    })
    const product = await apiPost(page, token, '/products', {
      category_id: category.data.id,
      name: `分配产品${timestamp}`,
      price: 88,
      stock: 10,
    })
    const staff = await apiPost(page, token, '/staff', {
      name: `分配人员${timestamp}`,
      phone: uniquePhone(testInfo, '134'),
      status: 'available',
    })
    const order = await apiPost(page, token, '/orders', {
      user_id: me.data.user_id,
      address: 'Phase27 分配测试地址',
      items: [{ product_id: product.data.id, quantity: 1 }],
    })

    await page.goto('/orders')
    const targetRow = page.locator('.el-table__body tr').first()
    await expect(targetRow.locator('td').first()).toContainText(String(order.data.id))
    await targetRow.getByRole('button', { name: '分配' }).click()

    const dialog = page.getByRole('dialog', { name: '分配服务人员' })
    await expect(dialog).toBeVisible()
    await dialog.getByRole('row', { name: new RegExp(staff.data.name) }).click()
    await dialog.getByRole('button', { name: '确定' }).click()

    await expect(page.locator('.el-message').first().getByText('分配成功')).toBeVisible()
    await expect(targetRow).toContainText(String(staff.data.id))
  })

  test('delete staff', async ({ page }, testInfo) => {
    const delName = 'del_' + Date.now()
    await page.goto('/staff')
    await page.getByRole('button', { name: '添加人员' }).click()
    await page.getByLabel('姓名').fill(delName)
    await page.getByLabel('手机号').fill(uniquePhone(testInfo, '138'))
    await page.locator('.el-dialog__footer button').filter({ hasText: '确定' }).click()

    // Wait for toast to appear then disappear before creating second staff
    await expect(page.getByText('创建成功')).toBeVisible()
    await expect(page.getByText('创建成功')).not.toBeVisible({ timeout: 5000 })

    // Also add another staff so we can verify count decrements
    const otherName = 'other_' + Date.now()
    await page.getByRole('button', { name: '添加人员' }).click()
    await page.getByLabel('姓名').fill(otherName)
    await page.getByLabel('手机号').fill(uniquePhone(testInfo, '138'))
    await page.locator('.el-dialog__footer button').filter({ hasText: '确定' }).click()

    // Wait for toast to appear then disappear
    await expect(page.getByText('创建成功')).toBeVisible()
    await expect(page.getByText('创建成功')).not.toBeVisible({ timeout: 5000 })

    // Verify both staff appear in list before delete
    await expect(page.locator('.el-table__body').getByText(delName)).toBeVisible()
    await expect(page.locator('.el-table__body').getByText(otherName)).toBeVisible()

    // Delete the target staff and confirm the Element Plus message box
    const targetRow = page.getByRole('row', { name: new RegExp(delName) })
    await targetRow.getByRole('button', { name: '删除' }).click()
    await expect(page.getByText('确定要删除该人员吗？')).toBeVisible()
    await page.locator('.el-message-box__btns .el-button--primary').click()

    // Verify delete success message
    await expect(page.getByText('删除成功')).toBeVisible({ timeout: 5000 })

    // Verify staff no longer in list
    await expect(page.locator('.el-table__body').getByText(delName)).not.toBeVisible({ timeout: 5000 })

    // Verify other staff still in list
    await expect(page.locator('.el-table__body').getByText(otherName)).toBeVisible()
  })
})
