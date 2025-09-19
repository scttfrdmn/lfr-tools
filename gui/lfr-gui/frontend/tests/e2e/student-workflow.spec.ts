import { test, expect } from '@playwright/test'

test.describe('Student Workflow', () => {
  test.beforeEach(async ({ page }) => {
    // Mock AWS services for testing
    await page.route('**/api/**', route => {
      const url = route.request().url()

      if (url.includes('GetUserRole')) {
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            role: 'student',
            username: 'test-student',
            project: 'test-class',
            permissions: ['connect']
          })
        })
      } else if (url.includes('ListInstances')) {
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify([{
            name: 'test-student-ubuntu',
            state: 'running',
            public_ip: '1.2.3.4',
            blueprint: 'ubuntu_22_04',
            bundle: 'app_standard_xl_1_0',
            username: 'test-student'
          }])
        })
      } else {
        route.continue()
      }
    })

    await page.goto('http://localhost:3000')
  })

  test('student sees workspace dashboard', async ({ page }) => {
    // Wait for app to load
    await expect(page.getByText('My Workspace')).toBeVisible()

    // Check for student-specific elements
    await expect(page.getByText('test-class')).toBeVisible()
    await expect(page.getByText('Ready to use')).toBeVisible()

    // Should see budget information
    await expect(page.getByText(/Budget:/)).toBeVisible()

    // Should see connection options
    await expect(page.getByText('Terminal')).toBeVisible()
    await expect(page.getByText('Desktop')).toBeVisible()
  })

  test('student can switch between terminal and desktop tabs', async ({ page }) => {
    await expect(page.getByText('My Workspace')).toBeVisible()

    // Click Terminal tab
    await page.getByRole('tab', { name: '💻 Terminal' }).click()
    await expect(page.getByText('SSH terminal')).toBeVisible()

    // Click Desktop tab
    await page.getByRole('tab', { name: '🖥️ Desktop' }).click()
    await expect(page.getByText('Remote desktop access')).toBeVisible()
  })

  test('student sees appropriate budget warnings', async ({ page }) => {
    // Mock high budget usage
    await page.route('**/GetProjectInfo', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          budget_used: 23.5,
          budget_total: 25.0
        })
      })
    })

    await page.reload()

    // Should see budget warning
    await expect(page.getByText(/Budget Warning/)).toBeVisible()
    await expect(page.getByText(/94%/)).toBeVisible()
  })

  test('student can access help and guidance', async ({ page }) => {
    await expect(page.getByText('My Workspace')).toBeVisible()

    // Check for help elements
    await expect(page.getByText('Get Help')).toBeVisible()
    await expect(page.getByText('Which Environment Should I Use?')).toBeVisible()

    // Check for educational guidance
    await expect(page.getByText('Best for:')).toBeVisible()
    await expect(page.getByText('Programming and coding')).toBeVisible()
    await expect(page.getByText('Visual programming')).toBeVisible()
  })
})

test.describe('Student Authentication', () => {
  test('handles missing authentication gracefully', async ({ page }) => {
    // Mock authentication failure
    await page.route('**/GetUserRole', route => {
      route.fulfill({
        status: 401,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'No valid authentication found' })
      })
    })

    await page.goto('http://localhost:3000')

    // Should show authentication error
    await expect(page.getByText(/Failed to load user information/)).toBeVisible()
  })

  test('student role shows limited interface', async ({ page }) => {
    await page.route('**/GetUserRole', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          role: 'student',
          username: 'alice',
          project: 'cs101',
          permissions: ['connect']
        })
      })
    })

    await page.goto('http://localhost:3000')

    // Should not see admin features
    await expect(page.getByText('Users & Groups')).not.toBeVisible()
    await expect(page.getByText('Settings')).not.toBeVisible()

    // Should see limited navigation
    await expect(page.getByText('Dashboard')).toBeVisible()
    await expect(page.getByText('Instances')).toBeVisible()
  })
})