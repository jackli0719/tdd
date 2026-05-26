import { test, expect } from '@playwright/test'

let userId = 0
const apiBase = 'http://localhost:8080/api'
const getNextUsername = () => `review_${Date.now()}_${++userId}`

const registerAndLogin = async (page, testInfo) => {
  const username = getNextUsername()
  const phone = `139${String(testInfo.workerIndex).padStart(2, '0')}${String(Date.now()).slice(-6)}`

  await page.goto('/login')
  await page.request.post(`${apiBase}/auth/register`, {
    data: {
      username,
      password: 'password123',
      email: `${username}@example.com`,
      phone,
    },
  })
  const loginRes = await page.request.post(`${apiBase}/auth/login`, {
    data: { username, password: 'password123' },
  })
  const loginData = await loginRes.json()
  const token = loginData.data.token
  const currentUserId = loginData.data.user.id
  await page.evaluate((authToken) => {
    localStorage.setItem('auth_token', authToken)
  }, token)

  return { token, userId: currentUserId }
}

const apiPost = async (page, token, path, data) => {
  const res = await page.request.post(`${apiBase}${path}`, {
    headers: { Authorization: `Bearer ${token}` },
    data,
  })
  expect(res.ok()).toBeTruthy()
  return res.json()
}

const createCompletedOrder = async (page, token, currentUserId) => {
  const timestamp = Date.now()
  const category = await apiPost(page, token, '/categories', {
    name: `评价品类${timestamp}`,
    description: 'review e2e',
  })
  const product = await apiPost(page, token, '/products', {
    category_id: category.data.id,
    name: `评价产品${timestamp}`,
    price: 99,
    stock: 10,
  })
  const staff = await apiPost(page, token, '/staff', {
    name: `评价人员${timestamp}`,
    phone: `137${String(timestamp).slice(-8)}`,
  })
  const order = await apiPost(page, token, '/orders', {
    user_id: currentUserId,
    staff_id: staff.data.id,
    address: '测试服务地址',
    items: [{ product_id: product.data.id, quantity: 1 }],
  })

  await apiPost(page, token, `/orders/${order.data.id}/paid`, {})
  await apiPost(page, token, `/orders/${order.data.id}/ship`, {})
  await apiPost(page, token, `/orders/${order.data.id}/complete`, {})

  return { orderId: order.data.id, staffId: staff.data.id }
}

test.describe('Review E2E', () => {
  test('completed order can be reviewed and listed', async ({ page }, testInfo) => {
    const session = await registerAndLogin(page, testInfo)
    const data = await createCompletedOrder(page, session.token, session.userId)

    await page.goto('/orders')
    const targetRow = page.locator('.el-table__body tr').first()
    await expect(targetRow).toBeVisible()
    await expect(targetRow.locator('td').first()).toContainText(String(data.orderId))
    await targetRow.getByRole('button', { name: '评价' }).click()
    const reviewDialog = page.getByRole('dialog', { name: '订单评价' })
    await expect(reviewDialog).toBeVisible()

    const reviewComment = `服务准时，体验很好 ${data.orderId}`
    await page.locator('.el-rate__item').nth(4).click()
    await reviewDialog.getByPlaceholder('请输入评价内容').fill(reviewComment)
    await reviewDialog.locator('.el-dialog__footer .el-button--primary').click({ force: true })
    await expect(page.locator('.el-message').first().getByText('评价成功')).toBeVisible()

    await page.goto('/reviews')
    await expect(page.locator('.main-content').getByText('评价管理')).toBeVisible()
    await expect(page.locator('.el-table__body').getByText(reviewComment)).toBeVisible()

    await page.goto('/staff')
    const staffRow = page.locator('.el-table__body tr').filter({ hasText: String(data.staffId) }).first()
    await expect(staffRow).toContainText('5.0 (1)')
  })
})
