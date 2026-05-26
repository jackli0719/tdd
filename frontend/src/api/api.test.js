import { describe, it, expect } from 'vitest'

// Mock the API responses
const mockUsers = [
  { id: 1, username: 'testuser', email: 'test@example.com', phone: '1234567890' },
  { id: 2, username: 'testuser2', email: 'test2@example.com', phone: '0987654321' },
]

const mockProducts = [
  { id: 1, name: 'Product A', price: 99.99, stock: 100 },
  { id: 2, name: 'Product B', price: 49.99, stock: 50 },
]

const mockOrders = [
  { id: 1, order_no: 'ORD001', user_id: 1, total_amount: 99.99, status: 'pending' },
  { id: 2, order_no: 'ORD002', user_id: 1, total_amount: 199.98, status: 'completed' },
]

// Simulate the API interceptor response format
const createApiResponse = (data) => ({
  code: 0,
  message: 'success',
  data,
})

describe('API Response Parsing', () => {
  it('should extract users from wrapped response format', () => {
    const apiResponse = createApiResponse({
      users: mockUsers,
      total: 2,
      page: 1,
    })
    const users = apiResponse.data?.users || []
    expect(users.length).toBe(2)
    expect(users[0].username).toBe('testuser')
  })

  it('should extract products from wrapped response format', () => {
    const apiResponse = createApiResponse({
      products: mockProducts,
      total: 2,
      page: 1,
    })
    const products = apiResponse.data?.products || []
    expect(products.length).toBe(2)
    expect(products[0].price).toBe(99.99)
  })

  it('should extract orders from wrapped response format', () => {
    const apiResponse = createApiResponse({
      orders: mockOrders,
      total: 2,
      page: 1,
    })
    const orders = apiResponse.data?.orders || []
    expect(orders.length).toBe(2)
    expect(orders[0].total_amount).toBe(99.99)
  })

  it('should handle empty response gracefully', () => {
    const apiResponse = createApiResponse({})
    const users = apiResponse.data?.users || []
    expect(users.length).toBe(0)
  })

  it('should handle null data gracefully', () => {
    const apiResponse = createApiResponse(null)
    const users = (apiResponse.data || {}).users || []
    expect(users.length).toBe(0)
  })
})

describe('Order Status Mapping', () => {
  const getStatusText = (status) => {
    const texts = {
      pending: '待确认',
      confirmed: '已确认',
      in_service: '服务中',
      completed: '已完成',
      cancelled: '已取消',
    }
    return texts[status] || status
  }

  const getStatusType = (status) => {
    const types = {
      pending: 'info',
      confirmed: 'success',
      in_service: 'warning',
      completed: '',
      cancelled: 'danger',
    }
    return status in types ? types[status] : 'info'
  }

  it('should return correct status text for all statuses', () => {
    expect(getStatusText('pending')).toBe('待确认')
    expect(getStatusText('confirmed')).toBe('已确认')
    expect(getStatusText('in_service')).toBe('服务中')
    expect(getStatusText('completed')).toBe('已完成')
    expect(getStatusText('cancelled')).toBe('已取消')
  })

  it('should return correct status type for all statuses', () => {
    expect(getStatusType('pending')).toBe('info')
    expect(getStatusType('confirmed')).toBe('success')
    expect(getStatusType('in_service')).toBe('warning')
    expect(getStatusType('completed')).toBe('')
    expect(getStatusType('cancelled')).toBe('danger')
  })

  it('should return status as-is for unknown status', () => {
    expect(getStatusText('unknown')).toBe('unknown')
    expect(getStatusType('unknown')).toBe('info')
  })
})

describe('Revenue Stats Field Mapping', () => {
  it('should map backend fields correctly', () => {
    const backendResponse = {
      total_revenue: 1000.00,
      pending_revenue: 100.00,
      confirmed_revenue: 200.00,
      in_service_revenue: 300.00,
      completed_revenue: 400.00,
    }

    const revenue = {
      total_revenue: backendResponse.total_revenue || 0,
      pending_revenue: backendResponse.pending_revenue || 0,
      confirmed_revenue: backendResponse.confirmed_revenue || 0,
      in_service_revenue: backendResponse.in_service_revenue || 0,
      completed_revenue: backendResponse.completed_revenue || 0,
    }

    expect(revenue.total_revenue).toBe(1000.00)
    expect(revenue.completed_revenue).toBe(400.00)
  })

  it('should handle missing fields with defaults', () => {
    const backendResponse = {}
    const revenue = {
      total_revenue: backendResponse.total_revenue || 0,
      pending_revenue: backendResponse.pending_revenue || 0,
    }
    expect(revenue.total_revenue).toBe(0)
    expect(revenue.pending_revenue).toBe(0)
  })
})

describe('Order Form Validation', () => {
  it('should validate user selection is required', () => {
    const form = { user_id: null, items: [{ product_id: 1, quantity: 1 }] }
    const isValid = form.user_id !== null
    expect(isValid).toBe(false)
  })

  it('should validate at least one product is required', () => {
    const form = { user_id: 1, items: [{ product_id: null, quantity: 1 }] }
    const validItems = form.items.filter(item => item.product_id && item.quantity > 0)
    expect(validItems.length).toBe(0)
  })

  it('should filter out invalid items', () => {
    const form = {
      user_id: 1,
      items: [
        { product_id: 1, quantity: 1 },
        { product_id: null, quantity: 0 },
        { product_id: 2, quantity: 2 },
      ],
    }
    const validItems = form.items.filter(item => item.product_id && item.quantity > 0)
    expect(validItems.length).toBe(2)
  })
})

describe('Time Formatting', () => {
  const formatTime = (date) => {
    const now = date || new Date()
    const year = now.getFullYear()
    const month = String(now.getMonth() + 1).padStart(2, '0')
    const day = String(now.getDate()).padStart(2, '0')
    const hours = String(now.getHours()).padStart(2, '0')
    const minutes = String(now.getMinutes()).padStart(2, '0')
    const seconds = String(now.getSeconds()).padStart(2, '0')
    return `${year}/${month}/${day} ${hours}:${minutes}:${seconds}`
  }

  it('should format time with zero padding', () => {
    const date = new Date(2026, 0, 5, 9, 5, 7) // 2026/01/05 09:05:07
    const result = formatTime(date)
    expect(result).toBe('2026/01/05 09:05:07')
  })

  it('should format time with double digit values', () => {
    const date = new Date(2026, 11, 25, 14, 30, 45) // 2026/12/25 14:30:45
    const result = formatTime(date)
    expect(result).toBe('2026/12/25 14:30:45')
  })

  it('should format current time when no date provided', () => {
    const result = formatTime()
    expect(result).toMatch(/^\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2}$/)
  })
})

describe('Date Formatting', () => {
  const formatDate = (date) => {
    if (!date) return '-'
    return new Date(date).toLocaleString()
  }

  it('should format valid date', () => {
    const result = formatDate('2026-05-25T18:00:00+08:00')
    expect(result).not.toBe('-')
  })

  it('should return dash for null date', () => {
    expect(formatDate(null)).toBe('-')
  })

  it('should return dash for undefined date', () => {
    expect(formatDate(undefined)).toBe('-')
  })
})

describe('UserForm Validation', () => {
  const validateUserForm = (form) => {
    const errors = []
    if (!form.username || form.username.length < 3) {
      errors.push('用户名至少3个字符')
    }
    if (!form.email || !form.email.includes('@')) {
      errors.push('邮箱格式不正确')
    }
    if (!form.phone || form.phone.length < 1) {
      errors.push('电话不能为空')
    }
    return errors
  }

  it('should pass with valid form data', () => {
    const form = { username: 'testuser', email: 'test@example.com', phone: '1234567890' }
    const errors = validateUserForm(form)
    expect(errors.length).toBe(0)
  })

  it('should fail with missing username', () => {
    const form = { username: '', email: 'test@example.com', phone: '1234567890' }
    const errors = validateUserForm(form)
    expect(errors).toContain('用户名至少3个字符')
  })

  it('should fail with invalid email', () => {
    const form = { username: 'testuser', email: 'invalid-email', phone: '1234567890' }
    const errors = validateUserForm(form)
    expect(errors).toContain('邮箱格式不正确')
  })

  it('should fail with missing phone', () => {
    const form = { username: 'testuser', email: 'test@example.com', phone: '' }
    const errors = validateUserForm(form)
    expect(errors).toContain('电话不能为空')
  })

  it('should return multiple errors for multiple issues', () => {
    const form = { username: '', email: 'invalid', phone: '' }
    const errors = validateUserForm(form)
    expect(errors.length).toBe(3)
  })
})

describe('ProductForm Validation', () => {
  const validateProductForm = (form) => {
    const errors = []
    if (!form.name || form.name.length < 1) {
      errors.push('产品名称不能为空')
    }
    if (!form.price || form.price <= 0) {
      errors.push('价格必须大于0')
    }
    if (form.stock < 0) {
      errors.push('库存不能为负数')
    }
    return errors
  }

  it('should pass with valid product data', () => {
    const form = { name: 'Test Product', price: 99.99, stock: 100 }
    const errors = validateProductForm(form)
    expect(errors.length).toBe(0)
  })

  it('should fail with missing name', () => {
    const form = { name: '', price: 99.99, stock: 100 }
    const errors = validateProductForm(form)
    expect(errors).toContain('产品名称不能为空')
  })

  it('should fail with zero price', () => {
    const form = { name: 'Test Product', price: 0, stock: 100 }
    const errors = validateProductForm(form)
    expect(errors).toContain('价格必须大于0')
  })

  it('should fail with negative stock', () => {
    const form = { name: 'Test Product', price: 99.99, stock: -5 }
    const errors = validateProductForm(form)
    expect(errors).toContain('库存不能为负数')
  })
})

describe('OrderStateTransition', () => {
  const getNextActions = (status) => {
    const transitions = {
      pending: ['confirm', 'cancel'],
      confirmed: ['start', 'cancel'],
      in_service: ['complete'],
      completed: [],
      cancelled: [],
    }
    return transitions[status] || []
  }

  it('should allow confirm and cancel from pending', () => {
    const actions = getNextActions('pending')
    expect(actions).toContain('confirm')
    expect(actions).toContain('cancel')
  })

  it('should allow start and cancel from confirmed', () => {
    const actions = getNextActions('confirmed')
    expect(actions).toContain('start')
    expect(actions).toContain('cancel')
  })

  it('should allow complete from in_service', () => {
    const actions = getNextActions('in_service')
    expect(actions).toContain('complete')
  })

  it('should allow no actions from completed', () => {
    const actions = getNextActions('completed')
    expect(actions.length).toBe(0)
  })

  it('should allow no actions from cancelled', () => {
    const actions = getNextActions('cancelled')
    expect(actions.length).toBe(0)
  })
})

describe('API Error Handling', () => {
  const parseError = (errorResponse) => {
    if (!errorResponse) return 'Unknown error'
    if (errorResponse.code === 404) return 'Resource not found'
    if (errorResponse.code === 400) return errorResponse.message || 'Bad request'
    if (errorResponse.code === 409) return 'Conflict'
    if (errorResponse.code === 500) return 'Server error'
    return errorResponse.message || 'Unknown error'
  }

  it('should parse 404 error as resource not found', () => {
    const error = { code: 404, message: 'user not found' }
    expect(parseError(error)).toBe('Resource not found')
  })

  it('should parse 400 error with message', () => {
    const error = { code: 400, message: 'invalid request' }
    expect(parseError(error)).toBe('invalid request')
  })

  it('should parse 409 error as conflict', () => {
    const error = { code: 409, message: 'user already exists' }
    expect(parseError(error)).toBe('Conflict')
  })

  it('should parse 500 error as server error', () => {
    const error = { code: 500, message: 'internal server error' }
    expect(parseError(error)).toBe('Server error')
  })

  it('should handle null error gracefully', () => {
    expect(parseError(null)).toBe('Unknown error')
  })
})