import '@testing-library/jest-dom'

// Mock Wails runtime
jest.mock('@wailsio/runtime', () => ({
  Events: {
    On: jest.fn(),
    Emit: jest.fn(),
  },
  WML: {
    Reload: jest.fn(),
  },
}))

// Mock LFR Service bindings
jest.mock('./bindings/lfr-gui', () => ({
  LFRService: {
    GetUserRole: jest.fn(),
    ListInstances: jest.fn(),
    StartInstance: jest.fn(),
    StopInstance: jest.fn(),
    GetProjectInfo: jest.fn(),
    ConnectToInstance: jest.fn(),
  },
  GreetService: {
    Greet: jest.fn(),
  },
}))

// Mock AWS Cloudscape Design components for testing
jest.mock('@cloudscape-design/components', () => {
  const actual = jest.requireActual('@cloudscape-design/components')
  return {
    ...actual,
    // Add any specific mocks if needed
  }
})