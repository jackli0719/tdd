import { describe, it, expect } from 'vitest'

describe('Dashboard Time Display', () => {
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

  it('should format time correctly for single digit values', () => {
    const date = new Date(2026, 0, 5, 9, 5, 7)
    expect(formatTime(date)).toBe('2026/01/05 09:05:07')
  })

  it('should format time correctly for double digit values', () => {
    const date = new Date(2026, 11, 25, 14, 30, 45)
    expect(formatTime(date)).toBe('2026/12/25 14:30:45')
  })

  it('should format time with year 2026', () => {
    const date = new Date(2026, 4, 26, 9, 30, 0)
    expect(formatTime(date)).toBe('2026/05/26 09:30:00')
  })

  it('should return valid time format when no date provided', () => {
    const result = formatTime()
    expect(result).toMatch(/^\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2}$/)
  })

  it('should handle first day of month', () => {
    const date = new Date(2026, 0, 1, 0, 0, 0)
    expect(formatTime(date)).toBe('2026/01/01 00:00:00')
  })

  it('should handle last day of month', () => {
    const date = new Date(2026, 11, 31, 23, 59, 59)
    expect(formatTime(date)).toBe('2026/12/31 23:59:59')
  })
})