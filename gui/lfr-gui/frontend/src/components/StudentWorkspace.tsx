import React, { useState, useEffect } from 'react'
import {
  Container,
  Header,
  Box,
  SpaceBetween,
  Button,
  StatusIndicator,
  Alert,
  Tabs,
  Badge,
  Grid,
  ProgressBar
} from '@cloudscape-design/components'

import { LFRService } from "../../bindings/lfr-gui"
import Terminal from './Terminal'
import DCVViewer from './DCVViewer'

interface UserInfo {
  role: string
  username: string
  project: string
  permissions: string[]
}

interface InstanceInfo {
  name: string
  state: string
  public_ip: string
  blueprint: string
  bundle: string
  username: string
}

interface StudentWorkspaceProps {
  userInfo: UserInfo
}

const StudentWorkspace: React.FC<StudentWorkspaceProps> = ({ userInfo }) => {
  const [instance, setInstance] = useState<InstanceInfo | null>(null)
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState('overview')
  const [budgetInfo, setBudgetInfo] = useState({
    used: 18.50,
    total: 25.00,
    percentage: 74
  })

  useEffect(() => {
    loadStudentInstance()
  }, [userInfo.username, userInfo.project])

  const loadStudentInstance = async () => {
    setLoading(true)
    try {
      const instances = await LFRService.ListInstances(userInfo.project)
      const studentInstance = instances?.find(inst => inst.username === userInfo.username)
      setInstance(studentInstance || null)
    } catch (err) {
      console.error('Failed to load student instance:', err)
    }
    setLoading(false)
  }

  const handleStartInstance = async () => {
    if (!instance) return

    try {
      await LFRService.StartInstance(instance.name)
      // Refresh instance status
      setTimeout(loadStudentInstance, 2000)
    } catch (err) {
      console.error('Failed to start instance:', err)
    }
  }

  const getInstanceStatusIndicator = () => {
    if (!instance) return <StatusIndicator type="error">No instance</StatusIndicator>

    switch (instance.state.toLowerCase()) {
      case 'running':
        return <StatusIndicator type="success">Ready to use</StatusIndicator>
      case 'stopped':
        return <StatusIndicator type="stopped">Sleeping</StatusIndicator>
      case 'starting':
      case 'pending':
        return <StatusIndicator type="in-progress">Starting up</StatusIndicator>
      case 'stopping':
        return <StatusIndicator type="in-progress">Going to sleep</StatusIndicator>
      default:
        return <StatusIndicator type="info">{instance.state}</StatusIndicator>
    }
  }

  const getBudgetAlert = () => {
    if (budgetInfo.percentage >= 95) {
      return (
        <Alert statusIconAriaLabel="Error" type="error" header="Budget Almost Exhausted">
          You've used {budgetInfo.percentage}% of your budget. Your computer will be suspended soon.
          Contact your instructor if you need more budget.
        </Alert>
      )
    } else if (budgetInfo.percentage >= 80) {
      return (
        <Alert statusIconAriaLabel="Warning" type="warning" header="Budget Warning">
          You've used {budgetInfo.percentage}% of your budget. Consider stopping your computer
          when not in use to save money.
        </Alert>
      )
    }
    return null
  }

  if (loading) {
    return (
      <Container>
        <Box textAlign="center" padding="xl">
          <StatusIndicator type="loading">Loading your workspace...</StatusIndicator>
        </Box>
      </Container>
    )
  }

  if (!instance) {
    return (
      <Container header={<Header variant="h1">My Workspace</Header>}>
        <Alert statusIconAriaLabel="Error" type="error" header="No Workspace Found">
          No cloud computer found for your account. Please contact your instructor.
        </Alert>
      </Container>
    )
  }

  const tabsContent = [
    {
      id: 'overview',
      label: 'Overview',
      content: (
        <SpaceBetween direction="vertical" size="l">
          {getBudgetAlert()}

          {/* Main Workspace Card */}
          <Container>
            <Grid gridDefinition={[{ colspan: 8 }, { colspan: 4 }]}>
              <SpaceBetween direction="vertical" size="m">
                <Box>
                  <SpaceBetween direction="horizontal" size="m" alignItems="center">
                    <Box fontSize="heading-m" fontWeight="bold">
                      📚 {userInfo.project}
                    </Box>
                    {getInstanceStatusIndicator()}
                  </SpaceBetween>
                </Box>

                <Box>
                  <SpaceBetween direction="vertical" size="xs">
                    <Box>
                      <strong>💻 {instance.name}</strong>
                    </Box>
                    <Box color="text-body-secondary">
                      {instance.blueprint} • {instance.bundle.replace('app_standard_', '').replace('_1_0', '').toUpperCase()}
                    </Box>
                    {instance.public_ip && (
                      <Box color="text-body-secondary">
                        IP: {instance.public_ip}
                      </Box>
                    )}
                  </SpaceBetween>
                </Box>

                <Box>
                  <strong>🎯 Choose your work environment:</strong>
                </Box>

                <SpaceBetween direction="horizontal" size="m">
                  {instance.state === 'running' ? (
                    <>
                      <Button
                        variant="primary"
                        iconName="call"
                        onClick={() => setActiveTab('terminal')}
                      >
                        💻 Terminal
                      </Button>
                      <Button
                        variant="primary"
                        iconName="contact"
                        onClick={() => setActiveTab('desktop')}
                      >
                        🖥️ Desktop
                      </Button>
                    </>
                  ) : (
                    <Button
                      variant="primary"
                      iconName="upload"
                      onClick={handleStartInstance}
                      loading={instance.state === 'starting' || instance.state === 'pending'}
                    >
                      {instance.state === 'starting' || instance.state === 'pending'
                        ? 'Starting...'
                        : 'Start Computer'
                      }
                    </Button>
                  )}

                  <Button iconName="folder">
                    Shared Files
                  </Button>

                  <Button iconName="status-info">
                    Get Help
                  </Button>
                </SpaceBetween>
              </SpaceBetween>

              {/* Budget Info */}
              <Container>
                <SpaceBetween direction="vertical" size="s">
                  <Box textAlign="center">
                    <Box fontSize="heading-l" fontWeight="bold">
                      ${budgetInfo.used}
                    </Box>
                    <Box color="text-body-secondary">of ${budgetInfo.total}</Box>
                  </SpaceBetween>

                  <ProgressBar
                    value={budgetInfo.percentage}
                    description="Budget used this semester"
                    size="small"
                  />

                  <Box textAlign="center" color="text-body-secondary" fontSize="body-s">
                    💡 Computer sleeps automatically to save money
                  </Box>
                </SpaceBetween>
              </Container>
            </Grid>
          </Container>

          {/* Work Environment Comparison */}
          <Container header={<Header variant="h3">🤔 Which Environment Should I Use?</Header>}>
            <Grid gridDefinition={[{ colspan: 6 }, { colspan: 6 }]}>
              <Box>
                <SpaceBetween direction="vertical" size="s">
                  <Box fontWeight="bold">💻 Terminal (Command Line)</Box>
                  <Box color="text-body-secondary">Best for:</Box>
                  <ul style={{ margin: 0, paddingLeft: '20px' }}>
                    <li>Programming and coding</li>
                    <li>Running scripts and commands</li>
                    <li>File management</li>
                    <li>Git version control</li>
                  </ul>
                </SpaceBetween>
              </Box>

              <Box>
                <SpaceBetween direction="vertical" size="s">
                  <Box fontWeight="bold">🖥️ Desktop (Visual Interface)</Box>
                  <Box color="text-body-secondary">Best for:</Box>
                  <ul style={{ margin: 0, paddingLeft: '20px' }}>
                    <li>Visual programming (Jupyter, RStudio)</li>
                    <li>Data visualization and charts</li>
                    <li>Document editing and presentations</li>
                    <li>Web browsing and research</li>
                  </ul>
                </SpaceBetween>
              </Box>
            </Grid>
          </Container>
        </SpaceBetween>
      )
    },
    {
      id: 'terminal',
      label: '💻 Terminal',
      content: (
        <Terminal
          username={userInfo.username}
          project={userInfo.project}
          instanceInfo={instance}
        />
      )
    },
    {
      id: 'desktop',
      label: '🖥️ Desktop',
      content: (
        <DCVViewer
          username={userInfo.username}
          project={userInfo.project}
          instanceInfo={instance}
        />
      )
    }
  ]

  return (
    <Container header={<Header variant="h1">My Workspace</Header>}>
      <Tabs
        tabs={tabsContent}
        activeTabId={activeTab}
        onChange={({ detail }) => setActiveTab(detail.activeTabId)}
      />
    </Container>
  )
}

export default StudentWorkspace