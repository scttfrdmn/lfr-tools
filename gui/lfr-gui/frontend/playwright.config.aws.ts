import { defineConfig, devices } from '@playwright/test'

/**
 * Playwright configuration for AWS integration testing
 * This config runs tests against real AWS services
 */
export default defineConfig({
  testDir: './tests/e2e',
  testMatch: '**/aws-integration.spec.ts',
  fullyParallel: false, // Run sequentially for AWS tests to avoid conflicts
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0, // Less retries for real AWS
  workers: 1, // Single worker for AWS tests
  reporter: 'html',
  timeout: 60000, // Longer timeout for real AWS calls

  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    actionTimeout: 10000, // Longer action timeout for AWS operations
  },

  // Test against fewer browsers for AWS integration (focus on functionality)
  projects: [
    {
      name: 'chromium-aws',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  // Start the Wails GUI application before tests
  webServer: {
    command: 'cd .. && wails3 dev', // Go up one level to gui/lfr-gui and run wails
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
    timeout: 120000, // Longer startup time for Wails app
    env: {
      // Pass through AWS credentials from environment
      AWS_REGION: process.env.AWS_REGION || 'us-east-1',
      AWS_PROFILE: process.env.AWS_PROFILE || 'default',
      AWS_ACCESS_KEY_ID: process.env.AWS_ACCESS_KEY_ID,
      AWS_SECRET_ACCESS_KEY: process.env.AWS_SECRET_ACCESS_KEY,
      AWS_SESSION_TOKEN: process.env.AWS_SESSION_TOKEN,
    }
  },

  // Global test setup for AWS integration
  globalSetup: './tests/setup/aws-setup.ts',
  globalTeardown: './tests/setup/aws-teardown.ts',
})