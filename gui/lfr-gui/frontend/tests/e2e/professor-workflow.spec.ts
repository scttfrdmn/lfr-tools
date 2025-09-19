import { test, expect } from '@playwright/test'

test.describe('Professor Workflow', () => {
  test.beforeEach(async ({ page }) => {
    // Mock professor authentication
    await page.route('**/GetUserRole', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          role: 'professor',
          username: 'prof-smith',
          project: 'cs101-fall2024',
          permissions: ['create', 'delete', 'start', 'stop', 'ssh', 'admin']
        })
      })
    })

    // Mock project information
    await page.route('**/GetProjectInfo', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          name: 'cs101-fall2024',
          student_count: 25,
          running_count: 12,
          budget_used: 340.50,
          budget_total: 500.00,
          days_remaining: 45
        })
      })
    })

    // Mock instances list
    await page.route('**/ListInstances', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            name: 'alice-ubuntu',
            state: 'running',
            public_ip: '1.2.3.4',
            blueprint: 'ubuntu_22_04',
            bundle: 'app_standard_xl_1_0',
            username: 'alice'
          },
          {
            name: 'bob-ubuntu',
            state: 'stopped',
            public_ip: '',
            blueprint: 'ubuntu_22_04',
            bundle: 'app_standard_2xl_1_0',
            username: 'bob'
          }
        ])
      })
    })

    await page.goto('http://localhost:3000')
  })

  test('professor sees comprehensive class dashboard', async ({ page }) => {
    // Wait for dashboard to load
    await expect(page.getByText('cs101-fall2024 Overview')).toBeVisible()

    // Check key metrics are displayed
    await expect(page.getByText('25')).toBeVisible() // Student count
    await expect(page.getByText('12')).toBeVisible() // Online count
    await expect(page.getByText('$340')).toBeVisible() // Budget used

    // Check action buttons
    await expect(page.getByText('Start All Instances')).toBeVisible()
    await expect(page.getByText('Stop All Instances')).toBeVisible()
  })

  test('professor can navigate to different sections', async ({ page }) => {
    await expect(page.getByText('Dashboard')).toBeVisible()

    // Navigate to Instances
    await page.getByText('Instances').click()
    await expect(page.getByText('Class Instances')).toBeVisible()

    // Navigate to Analytics
    await page.getByText('Analytics').click()
    await expect(page.getByText('Class Analytics')).toBeVisible()

    // Navigate to Students
    await page.getByText('Students').click()
    await expect(page.getByText('Student Management')).toBeVisible()
  })

  test('professor can manage individual instances', async ({ page }) => {
    // Go to instances page
    await page.getByText('Instances').click()

    // Should see instance table
    await expect(page.getByText('alice-ubuntu')).toBeVisible()
    await expect(page.getByText('bob-ubuntu')).toBeVisible()

    // Check status indicators
    await expect(page.getByText('Running')).toBeVisible()
    await expect(page.getByText('Stopped')).toBeVisible()

    // Should see action buttons
    await expect(page.getByRole('button', { name: 'SSH' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible()
  })

  test('professor can perform bulk operations', async ({ page }) => {
    await page.getByText('Instances').click()

    // Select multiple instances
    await page.getByRole('checkbox').first().click()
    await page.getByRole('checkbox').nth(1).click()

    // Should see bulk action buttons
    await expect(page.getByText('Start Selected')).toBeVisible()
    await expect(page.getByText('Stop Selected')).toBeVisible()
  })

  test('professor sees cost analytics and optimization', async ({ page }) => {
    await page.getByText('Analytics').click()

    // Should see analytics dashboard
    await expect(page.getByText('Weekly Cost Trend')).toBeVisible()
    await expect(page.getByText('Student Usage Patterns')).toBeVisible()
    await expect(page.getByText('Cost Optimization Suggestions')).toBeVisible()

    // Should see live monitoring
    await expect(page.getByText('Live Monitoring Active')).toBeVisible()
    await expect(page.getByText('Students online: 12')).toBeVisible()
  })
})

test.describe('Professor Instance Management', () => {
  test.beforeEach(async ({ page }) => {
    // Mock professor auth and data
    await page.route('**/GetUserRole', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          role: 'professor',
          username: 'prof-smith',
          project: 'cs101',
          permissions: ['create', 'delete', 'start', 'stop', 'ssh', 'admin']
        })
      })
    })

    await page.goto('http://localhost:3000')
  })

  test('professor can start instance with confirmation', async ({ page }) => {
    await page.getByText('Instances').click()

    // Mock stopped instance
    await page.route('**/ListInstances', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([{
          name: 'alice-ubuntu',
          state: 'stopped',
          public_ip: '',
          username: 'alice'
        }])
      })
    })

    await page.reload()

    // Click start button
    await page.getByRole('button', { name: 'Start' }).click()

    // Should see confirmation modal
    await expect(page.getByText('Start Instances')).toBeVisible()
    await expect(page.getByText('Starting instances will begin billing charges')).toBeVisible()

    // Mock start operation
    await page.route('**/StartInstance', route => {
      route.fulfill({ status: 200, body: '{}' })
    })

    // Confirm start
    await page.getByRole('button', { name: /Start.*Instance/ }).click()

    // Should see success feedback
    // Note: In real implementation, this would show progress and status updates
  })

  test('professor can stop instance with warning', async ({ page }) => {
    await page.getByText('Instances').click()

    // Click stop button
    await page.getByRole('button', { name: 'Stop' }).click()

    // Should see warning modal
    await expect(page.getByText('Stop Instances')).toBeVisible()
    await expect(page.getByText('may cause data loss if work is not saved')).toBeVisible()
  })
})