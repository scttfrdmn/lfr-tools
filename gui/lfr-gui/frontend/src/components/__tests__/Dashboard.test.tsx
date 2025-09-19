import React from 'react'
import { render, screen } from '@testing-library/react'
import '@testing-library/jest-dom'
import Dashboard from '../Dashboard'

// Mock the Wails bindings
jest.mock('../../bindings/lfr-gui', () => ({
  LFRService: {
    GetProjectInfo: jest.fn().mockResolvedValue({
      name: 'test-project',
      student_count: 5,
      running_count: 3,
      budget_used: 150.0,
      budget_total: 500.0,
      days_remaining: 45
    })
  }
}))

const mockUserInfo = {
  role: 'professor',
  username: 'test-professor',
  project: 'test-project',
  permissions: ['create', 'delete', 'start', 'stop']
}

const mockStudentInfo = {
  role: 'student',
  username: 'alice',
  project: 'cs101',
  permissions: ['connect']
}

describe('Dashboard Component', () => {
  test('renders professor dashboard correctly', async () => {
    render(<Dashboard userInfo={mockUserInfo} />)

    // Check for professor-specific elements
    expect(screen.getByText('Class Overview')).toBeInTheDocument()
    expect(screen.getByText('Budget Status')).toBeInTheDocument()
    expect(screen.getByText('Quick Actions')).toBeInTheDocument()

    // Wait for async data loading
    await screen.findByText('Students')
    expect(screen.getByText('Online Now')).toBeInTheDocument()
  })

  test('renders student dashboard correctly', () => {
    render(<Dashboard userInfo={mockStudentInfo} />)

    // Check for student-specific elements
    expect(screen.getByText('My Cloud Computer')).toBeInTheDocument()
    expect(screen.getByText('Connect Now')).toBeInTheDocument()
    expect(screen.getByText('Tips')).toBeInTheDocument()

    // Should not show professor features
    expect(screen.queryByText('Quick Actions')).not.toBeInTheDocument()
    expect(screen.queryByText('Budget Status')).not.toBeInTheDocument()
  })

  test('shows loading state initially', () => {
    // Mock delayed response
    const mockLFRService = require('../../bindings/lfr-gui').LFRService
    mockLFRService.GetProjectInfo.mockImplementation(
      () => new Promise(resolve => setTimeout(() => resolve({}), 1000))
    )

    render(<Dashboard userInfo={mockUserInfo} />)

    expect(screen.getByText('Loading dashboard...')).toBeInTheDocument()
  })

  test('handles different user roles appropriately', () => {
    const taUserInfo = {
      role: 'ta',
      username: 'ta-alice',
      project: 'cs101',
      permissions: ['start', 'stop', 'connect']
    }

    render(<Dashboard userInfo={taUserInfo} />)

    // TA should see professor-style dashboard but with limited permissions
    expect(screen.getByText('Class Overview')).toBeInTheDocument()
  })

  test('handles unknown role gracefully', () => {
    const unknownUserInfo = {
      role: 'unknown',
      username: 'unknown-user',
      project: 'test',
      permissions: []
    }

    render(<Dashboard userInfo={unknownUserInfo} />)

    expect(screen.getByText(/Unknown user role/)).toBeInTheDocument()
  })
})