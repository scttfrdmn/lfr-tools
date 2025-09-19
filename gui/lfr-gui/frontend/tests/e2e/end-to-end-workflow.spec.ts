import { test, expect } from '@playwright/test'

test.describe('Complete Educational Workflow', () => {
  test('professor sets up class and student connects', async ({ page, context }) => {
    // Step 1: Professor login
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

    // Professor should see class management interface
    await expect(page.getByText(/Class Overview/)).toBeVisible()
    await expect(page.getByText('Start All Instances')).toBeVisible()

    // Step 2: Professor starts instances for class
    await page.getByText('Instances').click()

    // Mock instances list with stopped instances
    await page.route('**/ListInstances', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            name: 'alice-ubuntu',
            state: 'stopped',
            public_ip: '',
            username: 'alice'
          }
        ])
      })
    })

    await page.reload()

    // Click start instance
    await page.getByRole('button', { name: 'Start' }).click()
    await expect(page.getByText('Start Instances')).toBeVisible()

    // Mock successful start
    await page.route('**/StartInstance', route => {
      route.fulfill({ status: 200, body: '{}' })
    })

    await page.getByRole('button', { name: /Start.*Instance/ }).click()

    // Step 3: Open student view in new page
    const studentPage = await context.newPage()

    // Mock student authentication
    await studentPage.route('**/GetUserRole', route => {
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

    // Mock running instance for student
    await studentPage.route('**/ListInstances', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([{
          name: 'alice-ubuntu',
          state: 'running',
          public_ip: '1.2.3.4',
          username: 'alice'
        }])
      })
    })

    await studentPage.goto('http://localhost:3000')

    // Student should see workspace
    await expect(studentPage.getByText('My Workspace')).toBeVisible()
    await expect(studentPage.getByText('Ready to use')).toBeVisible()

    // Student can access terminal
    await studentPage.getByRole('tab', { name: '💻 Terminal' }).click()
    await expect(studentPage.getByText('SSH terminal')).toBeVisible()

    // Student can access desktop
    await studentPage.getByRole('tab', { name: '🖥️ Desktop' }).click()
    await expect(studentPage.getByText('Remote desktop access')).toBeVisible()

    await studentPage.close()
  })

  test('TA can support students during class', async ({ page }) => {
    // Mock TA authentication
    await page.route('**/GetUserRole', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          role: 'ta',
          username: 'ta-alice',
          project: 'cs101',
          permissions: ['start', 'stop', 'connect', 'support']
        })
      })
    })

    await page.goto('http://localhost:3000')

    // TA should see limited interface
    await expect(page.getByText('Dashboard')).toBeVisible()
    await expect(page.getByText('Instances')).toBeVisible()
    await expect(page.getByText('Student Support')).toBeVisible()

    // Should NOT see admin features
    await expect(page.getByText('Users & Groups')).not.toBeVisible()
    await expect(page.getByText('Settings')).not.toBeVisible()

    // TA can access analytics for class monitoring
    await page.getByText('Analytics').click()
    await expect(page.getByText('Class Analytics')).toBeVisible()
  })

  test('cost monitoring and optimization workflow', async ({ page }) => {
    // Mock professor with cost data
    await page.route('**/GetUserRole', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          role: 'professor',
          username: 'prof-smith',
          project: 'cs101',
          permissions: ['admin']
        })
      })
    })

    // Mock high budget usage
    await page.route('**/GetProjectInfo', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          name: 'cs101',
          student_count: 25,
          running_count: 20,
          budget_used: 450.0,
          budget_total: 500.0,
          days_remaining: 30
        })
      })
    })

    await page.goto('http://localhost:3000')

    // Should see budget warning
    await expect(page.getByText(/Budget Warning/)).toBeVisible()

    // Navigate to analytics
    await page.getByText('Analytics').click()

    // Should see optimization suggestions
    await expect(page.getByText('Cost Optimization Suggestions')).toBeVisible()
    await expect(page.getByText('Apply All Suggestions')).toBeVisible()

    // Professor can apply optimizations
    await page.getByText('Apply All Suggestions').click()

    // Should see confirmation or progress
    // Note: Real implementation would show optimization progress
  })
})

test.describe('Error Handling and Edge Cases', () => {
  test('handles AWS connection failures gracefully', async ({ page }) => {
    // Mock AWS connection failure
    await page.route('**/GetUserRole', route => {
      route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'AWS connection failed' })
      })
    })

    await page.goto('http://localhost:3000')

    // Should show error message
    await expect(page.getByText(/Failed to load user information/)).toBeVisible()
  })

  test('handles no instances gracefully', async ({ page }) => {
    await page.route('**/GetUserRole', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          role: 'student',
          username: 'alice',
          project: 'empty-class',
          permissions: ['connect']
        })
      })
    })

    await page.route('**/ListInstances', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([])
      })
    })

    await page.goto('http://localhost:3000')

    // Should show no workspace message
    await expect(page.getByText(/No Workspace Found/)).toBeVisible()
    await expect(page.getByText(/contact your instructor/)).toBeVisible()
  })

  test('validates user permissions correctly', async ({ page }) => {
    // Mock student trying to access admin features
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

    // Student should not see admin navigation
    await expect(page.getByText('Users & Groups')).not.toBeVisible()
    await expect(page.getByText('Settings')).not.toBeVisible()

    // Attempting to navigate to admin URL should be blocked
    await page.goto('http://localhost:3000/users')
    await expect(page.getByText('My Workspace')).toBeVisible() // Should redirect to dashboard
  })
})

test.describe('Performance and Responsiveness', () => {
  test('app loads within performance budget', async ({ page }) => {
    const startTime = Date.now()

    await page.goto('http://localhost:3000')
    await expect(page.getByText(/LFR Tools|My Workspace|Class Overview/)).toBeVisible()

    const loadTime = Date.now() - startTime

    // Should load within 3 seconds
    expect(loadTime).toBeLessThan(3000)
  })

  test('handles slow network conditions', async ({ page, context }) => {
    // Simulate slow network
    await context.route('**/*', route => {
      setTimeout(() => route.continue(), 500) // 500ms delay
    })

    await page.goto('http://localhost:3000')

    // Should show loading states
    await expect(page.getByText(/Loading/)).toBeVisible()

    // Should eventually load
    await expect(page.getByText(/LFR Tools|My Workspace|Class Overview/)).toBeVisible({ timeout: 10000 })
  })

  test('responsive design works on different screen sizes', async ({ page }) => {
    await page.goto('http://localhost:3000')

    // Test desktop size
    await page.setViewportSize({ width: 1920, height: 1080 })
    await expect(page.getByText(/Dashboard/)).toBeVisible()

    // Test tablet size
    await page.setViewportSize({ width: 768, height: 1024 })
    await expect(page.getByText(/Dashboard/)).toBeVisible()

    // Test mobile size
    await page.setViewportSize({ width: 375, height: 667 })
    await expect(page.getByText(/Dashboard/)).toBeVisible()
  })
})