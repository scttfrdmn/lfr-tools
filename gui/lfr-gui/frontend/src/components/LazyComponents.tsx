import React, { lazy, Suspense } from 'react'
import { Container, StatusIndicator, Box } from '@cloudscape-design/components'

// Lazy load heavy components for better performance
const LazyAnalytics = lazy(() => import('./Analytics'))
const LazyTerminal = lazy(() => import('./Terminal'))
const LazyInstances = lazy(() => import('./Instances'))

// Loading component
const ComponentLoader: React.FC<{ componentName: string }> = ({ componentName }) => (
  <Container>
    <Box textAlign="center" padding="xl">
      <StatusIndicator type="loading">Loading {componentName}...</StatusIndicator>
    </Box>
  </Container>
)

// Wrapped components with Suspense
export const Analytics: React.FC<any> = (props) => (
  <Suspense fallback={<ComponentLoader componentName="Analytics" />}>
    <LazyAnalytics {...props} />
  </Suspense>
)

export const Terminal: React.FC<any> = (props) => (
  <Suspense fallback={<ComponentLoader componentName="Terminal" />}>
    <LazyTerminal {...props} />
  </Suspense>
)

export const Instances: React.FC<any> = (props) => (
  <Suspense fallback={<ComponentLoader componentName="Instances" />}>
    <LazyInstances {...props} />
  </Suspense>
)

// Error boundary for lazy components
export class LazyComponentErrorBoundary extends React.Component<
  { children: React.ReactNode; fallback?: React.ReactNode },
  { hasError: boolean }
> {
  constructor(props: any) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(): { hasError: boolean } {
    return { hasError: true }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Lazy component error:', error, errorInfo)
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <Container>
          <Box textAlign="center" padding="xl">
            <StatusIndicator type="error">
              Failed to load component. Please refresh the page.
            </StatusIndicator>
          </Box>
        </Container>
      )
    }

    return this.props.children
  }
}