export const uniqueSuffix = (testInfo, label = 'e2e') => {
  const safeLabel = label.replace(/[^a-zA-Z0-9]/g, '').slice(0, 12) || 'e2e'
  const worker = String(testInfo.workerIndex ?? 0).padStart(2, '0')
  const retry = String(testInfo.retry ?? 0)
  const stamp = Date.now().toString(36)
  const random = Math.random().toString(36).slice(2, 8)
  return `${safeLabel}_${worker}_${retry}_${stamp}_${random}`
}

export const uniquePhone = (testInfo, prefix = '13') => {
  const worker = String(testInfo.workerIndex ?? 0).padStart(2, '0')
  const retry = String(testInfo.retry ?? 0).slice(-1)
  const stamp = String(Date.now()).slice(-2)
  const random = String(Math.floor(Math.random() * 10000)).padStart(4, '0')
  return `${prefix}${worker}${retry}${stamp}${random}`
}

export const testIdentity = (testInfo, label = 'user', prefix = '13') => {
  const suffix = uniqueSuffix(testInfo, label)
  return {
    username: `${label}_${suffix}`,
    email: `${label}_${suffix}@example.com`,
    phone: uniquePhone(testInfo, prefix),
  }
}
