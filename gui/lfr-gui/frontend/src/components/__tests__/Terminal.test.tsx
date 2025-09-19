import React from 'react'
import { render, screen, fireEvent } from '@testing-library/react'
import '@testing-library/jest-dom'
import Terminal from '../Terminal'

// Mock xterm
jest.mock('@xterm/xterm', () => ({
  Terminal: jest.fn().mockImplementation(() => ({
    open: jest.fn(),
    dispose: jest.fn(),
    writeln: jest.fn(),
    loadAddon: jest.fn(),
  }))
}))

jest.mock('@xterm/addon-fit', () => ({
  FitAddon: jest.fn().mockImplementation(() => ({
    fit: jest.fn(),
  }))
}))

jest.mock('@xterm/addon-web-links', () => ({
  WebLinksAddon: jest.fn()
}))

const mockInstanceInfo = {
  name: 'alice-ubuntu',
  state: 'running',
  public_ip: '34.221.15.92'
}

describe('Terminal Component', () => {
  test('renders terminal container correctly', () => {
    render(
      <Terminal
        username="alice"
        project="cs101"
        instanceInfo={mockInstanceInfo}
      />
    )

    expect(screen.getByText('Terminal')).toBeInTheDocument()
    expect(screen.getByText(/SSH terminal for alice's cloud computer/)).toBeInTheDocument()
  })

  test('shows connect button when disconnected', () => {
    render(
      <Terminal
        username="alice"
        project="cs101"
        instanceInfo={mockInstanceInfo}
      />
    )

    expect(screen.getByText('Connect')).toBeInTheDocument()
    expect(screen.getByText('Disconnected')).toBeInTheDocument()
  })

  test('shows warning when no instance info provided', () => {
    render(
      <Terminal
        username="alice"
        project="cs101"
      />
    )

    expect(screen.getByText(/No instance information available/)).toBeInTheDocument()
  })

  test('shows instance status alert for stopped instance', () => {
    const stoppedInstance = {
      ...mockInstanceInfo,
      state: 'stopped'
    }

    render(
      <Terminal
        username="alice"
        project="cs101"
        instanceInfo={stoppedInstance}
      />
    )

    expect(screen.getByText(/Instance is stopped/)).toBeInTheDocument()
  })

  test('handles connect button click', () => {
    const onConnect = jest.fn()

    render(
      <Terminal
        username="alice"
        project="cs101"
        instanceInfo={mockInstanceInfo}
        onConnect={onConnect}
      />
    )

    const connectButton = screen.getByText('Connect')
    fireEvent.click(connectButton)

    // Terminal should start connection process
    // onConnect will be called when connection is established
  })

  test('provides fit button for terminal resizing', () => {
    render(
      <Terminal
        username="alice"
        project="cs101"
        instanceInfo={mockInstanceInfo}
      />
    )

    expect(screen.getByText('Fit')).toBeInTheDocument()
  })
})