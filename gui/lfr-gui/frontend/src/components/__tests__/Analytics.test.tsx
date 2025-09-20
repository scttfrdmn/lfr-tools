// import React from 'react'
import { render, screen } from '@testing-library/react'
import '@testing-library/jest-dom'
import Analytics from '../Analytics'

// Mock the real-time status hook
jest.mock('../hooks/useRealTimeStatus', () => ({
  useClassMonitoring: () => ({
    projectStatus: {
      name: 'test-project',
      students_online: 5,
      budget_used: 340.50,
      alert_count: 2
    },
    studentsOnline: 5,
    budgetUsed: 340.50,
    alerts: ['alice connected', 'bob disconnected'],
    lastUpdate: '2024-09-19T10:30:00Z'
  })
}))

const mockUserInfo = {
  role: 'professor',
  username: 'test-professor',
  project: 'test-project',
  permissions: ['create', 'delete', 'start', 'stop']
}

describe('Analytics Component', () => {
  test('renders analytics dashboard correctly', () => {
    render(<Analytics userInfo={mockUserInfo} />)

    // Check for main sections
    expect(screen.getByText(/Class Analytics/)).toBeInTheDocument()
    expect(screen.getByText('Weekly Cost Trend')).toBeInTheDocument()
    expect(screen.getByText('Student Usage Patterns')).toBeInTheDocument()
    expect(screen.getByText('Cost Optimization Suggestions')).toBeInTheDocument()
  })

  test('shows real-time monitoring status', () => {
    render(<Analytics userInfo={mockUserInfo} />)

    expect(screen.getByText('Live Monitoring Active')).toBeInTheDocument()
    expect(screen.getByText(/Students online: 5/)).toBeInTheDocument()
  })

  test('displays budget overview correctly', () => {
    render(<Analytics userInfo={mockUserInfo} />)

    expect(screen.getByText('Budget Overview')).toBeInTheDocument()
    expect(screen.getByText('Remaining')).toBeInTheDocument()
    expect(screen.getByText('Days Left')).toBeInTheDocument()
  })

  test('shows optimization suggestions', () => {
    render(<Analytics userInfo={mockUserInfo} />)

    expect(screen.getByText('Cost Optimization Suggestions')).toBeInTheDocument()
    expect(screen.getByText('Apply All Suggestions')).toBeInTheDocument()
  })

  test('displays live student status', () => {
    render(<Analytics userInfo={mockUserInfo} />)

    expect(screen.getByText('Live Student Status')).toBeInTheDocument()
    expect(screen.getByText('Online Now')).toBeInTheDocument()
    expect(screen.getByText('Working')).toBeInTheDocument()
    expect(screen.getByText('Idle')).toBeInTheDocument()
    expect(screen.getByText('Offline')).toBeInTheDocument()
  })
})