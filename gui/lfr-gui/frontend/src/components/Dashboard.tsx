import React, { useState, useEffect } from 'react'
import {
  Container,
  Header,
  Grid,
  Box,
  SpaceBetween,
  StatusIndicator,
  ProgressBar,
  Button,
  Alert,
  Cards,
  Badge
} from '@cloudscape-design/components'

import { LFRService } from "../../bindings/lfr-gui"

interface UserInfo {
  role: string
  username: string
  project: string
  permissions: string[]
}

interface ProjectInfo {
  name: string
  student_count: number
  running_count: number
  budget_used: number
  budget_total: number
  days_remaining: number
}

interface DashboardProps {
  userInfo: UserInfo
}

const Dashboard: React.FC<DashboardProps> = ({ userInfo }) => {
  const [projectInfo, setProjectInfo] = useState<ProjectInfo | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Load project information
    if (userInfo.project) {
      LFRService.GetProjectInfo(userInfo.project)
        .then((info: ProjectInfo) => {
          setProjectInfo(info)
          setLoading(false)
        })
        .catch((err: any) => {
          console.error('Failed to get project info:', err)
          setLoading(false)
        })
    }
  }, [userInfo.project])

  const renderStudentDashboard = () => (
    <Container header={<Header variant="h1">My Cloud Computer</Header>}>
      <SpaceBetween direction="vertical" size="l">
        <Alert statusIconAriaLabel="Info" type="info">
          Welcome to your cloud computing environment! Your computer will automatically
          start when you connect and sleep when not in use to save costs.
        </Alert>

        <Grid gridDefinition={[{ colspan: 12 }]}>
          <Container header={<Header variant="h2">{userInfo.project}</Header>}>
            <SpaceBetween direction="vertical" size="m">
              <Box>
                <SpaceBetween direction="horizontal" size="m">
                  <StatusIndicator type="success">Ready to use</StatusIndicator>
                  <Badge color="green">Budget: $18.50 / $25.00</Badge>
                  <Badge color="blue">45 days remaining</Badge>
                </SpaceBetween>
              </Box>

              <Box padding={{ vertical: "m" }}>
                <strong>💻 {userInfo.username}-ubuntu-22-04</strong>
                <br />
                Status: Ready to connect
                <br />
                Last used: 2 hours ago
              </Box>

              <SpaceBetween direction="horizontal" size="m">
                <Button variant="primary" iconName="external">
                  Connect Now
                </Button>
                <Button iconName="folder">
                  Shared Files
                </Button>
                <Button iconName="status-info">
                  Get Help
                </Button>
              </SpaceBetween>
            </SpaceBetween>
          </Container>
        </Grid>

        <Container header={<Header variant="h3">💡 Tips</Header>}>
          <ul>
            <li>Your computer sleeps automatically to save money</li>
            <li>Shared files are in /mnt/efs/shared/</li>
            <li>Save your work often</li>
            <li>Ask your teacher or TA if you need help</li>
          </ul>
        </Container>
      </SpaceBetween>
    </Container>
  )

  const renderProfessorDashboard = () => (
    <Container header={<Header variant="h1">{projectInfo?.name || 'Class Overview'}</Header>}>
      <SpaceBetween direction="vertical" size="l">
        {/* Class Summary Cards */}
        <Grid gridDefinition={[{ colspan: 3 }, { colspan: 3 }, { colspan: 3 }, { colspan: 3 }]}>
          <Container>
            <Box textAlign="center">
              <Box fontSize="heading-xl" fontWeight="bold" color="text-status-success">
                {projectInfo?.student_count || 0}
              </Box>
              <Box color="text-body-secondary">Students</Box>
            </Box>
          </Container>

          <Container>
            <Box textAlign="center">
              <Box fontSize="heading-xl" fontWeight="bold" color="text-status-info">
                {projectInfo?.running_count || 0}
              </Box>
              <Box color="text-body-secondary">Online Now</Box>
            </Box>
          </Container>

          <Container>
            <Box textAlign="center">
              <Box fontSize="heading-xl" fontWeight="bold" color="text-status-warning">
                ${projectInfo?.budget_used || 0}
              </Box>
              <Box color="text-body-secondary">Budget Used</Box>
            </Box>
          </Container>

          <Container>
            <Box textAlign="center">
              <Box fontSize="heading-xl" fontWeight="bold">
                {projectInfo?.days_remaining || 0}
              </Box>
              <Box color="text-body-secondary">Days Left</Box>
            </Box>
          </Container>
        </Grid>

        {/* Budget Progress */}
        <Container header={<Header variant="h2">💰 Budget Status</Header>}>
          <SpaceBetween direction="vertical" size="m">
            <ProgressBar
              value={(projectInfo?.budget_used || 0) / (projectInfo?.budget_total || 1) * 100}
              additionalInfo={`$${projectInfo?.budget_used || 0} of $${projectInfo?.budget_total || 0} used`}
              description="Semester budget usage"
              label="Budget Progress"
            />
            {(projectInfo?.budget_used || 0) / (projectInfo?.budget_total || 1) > 0.8 && (
              <Alert
                statusIconAriaLabel="Warning"
                type="warning"
                header="Budget Warning"
              >
                You've used {Math.round((projectInfo?.budget_used || 0) / (projectInfo?.budget_total || 1) * 100)}%
                of your semester budget. Consider applying cost optimization policies.
              </Alert>
            )}
          </SpaceBetween>
        </Container>

        {/* Quick Actions */}
        <Container header={<Header variant="h2">⚡ Quick Actions</Header>}>
          <SpaceBetween direction="horizontal" size="m">
            <Button variant="primary" iconName="upload">
              Start All Instances
            </Button>
            <Button iconName="download">
              Stop All Instances
            </Button>
            <Button iconName="envelope">
              Email Class
            </Button>
            <Button iconName="trending-up">
              Generate Report
            </Button>
          </SpaceBetween>
        </Container>

        {/* Recent Alerts */}
        <Container header={<Header variant="h2">🚨 Recent Alerts</Header>}>
          <SpaceBetween direction="vertical" size="s">
            <Alert statusIconAriaLabel="Info" type="info">
              3 students have requested instance starts in the last 10 minutes
            </Alert>
            <Alert statusIconAriaLabel="Warning" type="warning">
              Student alice is approaching budget limit ($23/$25 used)
            </Alert>
          </SpaceBetween>
        </Container>
      </SpaceBetween>
    </Container>
  )

  if (loading) {
    return (
      <Container>
        <Box textAlign="center" padding="xl">
          <StatusIndicator type="loading">Loading dashboard...</StatusIndicator>
        </Box>
      </Container>
    )
  }

  // Render different dashboards based on user role
  switch (userInfo.role) {
    case 'student':
      return renderStudentDashboard()
    case 'professor':
    case 'admin':
      return renderProfessorDashboard()
    case 'ta':
      return renderProfessorDashboard() // TAs see similar view with limited actions
    default:
      return (
        <Container>
          <Alert statusIconAriaLabel="Error" type="error">
            Unknown user role: {userInfo.role}
          </Alert>
        </Container>
      )
  }
}

export default Dashboard