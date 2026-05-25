import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: {
    command: 'cd /Users/liwei/res/minimaxi/tdd/oms && go build -o server ./cmd/server && ./server & sleep 2 && cd /Users/liwei/res/minimaxi/tdd/frontend && npm run dev',
    port: 5173,
    reuseExistingServer: false, // Always start fresh, don't depend on residual services
    timeout: 180000,
    stdout: 'pipe',
    stderr: 'pipe',
  },
})