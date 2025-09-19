import React, { useState, useEffect } from 'react'
import {
  Container,
  Header,
  Box,
  SpaceBetween,
  Button,
  StatusIndicator,
  Alert,
  ProgressBar,
  Grid,
  Badge,
  Modal
} from '@cloudscape-design/components'

import { LFRService } from "../../bindings/lfr-gui"
import Terminal from './Terminal'

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

interface StudentConnectProps {
  userInfo: UserInfo
}

const StudentConnect: React.FC<StudentConnectProps> = ({ userInfo }) => {
  const [instance, setInstance] = useState<InstanceInfo | null>(null)
  const [loading, setLoading] = useState(true)
  const [showTerminal, setShowTerminal] = useState(false)
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
          <StatusIndicator type="loading">Loading your cloud computer...</StatusIndicator>
        </Box>
      </Container>
    )
  }

  if (!instance) {
    return (
      <Container header={<Header variant="h1">My Cloud Computer</Header>}>
        <Alert statusIconAriaLabel="Error" type="error" header="No Instance Found">
          No cloud computer found for your account. Please contact your instructor.
        </Alert>
      </Container>
    )
  }

  return (
    <Container header={<Header variant="h1">My Cloud Computer</Header>}>
      <SpaceBetween direction="vertical" size="l">
        {getBudgetAlert()}

        {/* Main Instance Card */}
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

              <SpaceBetween direction="horizontal" size="m">
                {instance.state === 'running' ? (
                  <>
                    <Button
                      variant="primary"
                      iconName="call"
                      onClick={() => setShowTerminal(true)}
                    >
                      Open Terminal
                    </Button>
                    <Button
                      iconName="external"
                      href={`ssh://${userInfo.username}@${instance.public_ip}`}
                    >
                      SSH (External)
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

        {/* Tips and Help */}
        <Container header={<Header variant="h3">💡 Tips for Success</Header>}>
          <Box>
            <ul style={{ margin: 0, paddingLeft: '20px' }}>
              <li>Save your work frequently - computers can sleep automatically</li>
              <li>Shared class files are in <code>/mnt/efs/shared/</code></li>
              <li>Submit homework to <code>/mnt/efs/submissions/</code></li>
              <li>Ask your teacher or TA if you need help</li>
              <li>Your computer will sleep when not used to save money</li>
            </ul>
          </Box>
        </Container>

        {/* Recent Activity */}
        <Container header={<Header variant="h3">📈 Recent Activity</Header>}>
          <SpaceBetween direction="vertical" size="xs">
            <Box color="text-body-secondary">Last connected: 2 hours ago</Box>
            <Box color="text-body-secondary">Total usage today: 3.2 hours</Box>
            <Box color="text-body-secondary">Files saved: homework1.py, lab2.py</Box>
          </SpaceBetween>
        </Container>

        {/* Terminal Modal */}
        <Modal
          onDismiss={() => setShowTerminal(false)}
          visible={showTerminal}
          header="SSH Terminal"
          size="max"
          footer={
            <Box float="right">
              <Button variant="link" onClick={() => setShowTerminal(false)}>
                Close Terminal
              </Button>
            </Box>
          }
        >
          <Terminal
            username={userInfo.username}
            project={userInfo.project}
            instanceInfo={instance}
            onConnect={() => setConnected(true)}
            onDisconnect={() => setConnected(false)}
          />
        </Modal>
      </SpaceBetween>
    </Container>
  )
}

export default StudentConnect