import { defineConfig, devices } from '@playwright/test'

/**
 * Playwright configuration for LocalStack testing
 * This config runs tests against LocalStack for safe, isolated testing
 */
export default defineConfig({
  testDir: './tests/e2e',
  testMatch: '**/aws-integration.spec.ts',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 1,
  workers: process.env.CI ? 1 : 2,
  reporter: 'html',
  timeout: 30000,

  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },

  projects: [
    {
      name: 'chromium-localstack',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  // Start both LocalStack and the GUI application
  webServer: [
    {
      command: 'docker-compose -f ../../docker-compose.test.yml up',
      url: 'http://localhost:4566',
      timeout: 30000,
      reuseExistingServer: false,
    },
    {
      command: 'cd .. && wails3 dev',
      url: 'http://localhost:3000',
      timeout: 60000,
      reuseExistingServer: !process.env.CI,
      env: {
        AWS_ENDPOINT_URL: 'http://localhost:4566',
        AWS_ACCESS_KEY_ID: 'test',
        AWS_SECRET_ACCESS_KEY: 'test',
        AWS_REGION: 'us-east-1',
      }
    }
  ],

  // Global setup for LocalStack
  globalSetup: './tests/setup/localstack-setup.ts',
  globalTeardown: './tests/setup/localstack-teardown.ts',
})