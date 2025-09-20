import { useState } from 'react'
import '@cloudscape-design/global-styles/index.css'
import {
  AppLayout,
  TopNavigation,
  SideNavigation,
  Container,
  Header,
  Button,
  SpaceBetween,
  Box,
  Grid,
  StatusIndicator
} from '@cloudscape-design/components'

interface UserInfo {
  role: string
  username: string
  project: string
  permissions: string[]
}

function App() {
  const [userInfo] = useState<UserInfo>({
    role: 'professor',
    username: 'demo-professor',
    project: 'CS101-Fall2024',
    permissions: ['create', 'delete', 'start', 'stop', 'ssh', 'admin']
  })

  const navigationItems = [
    { type: 'link', text: 'Dashboard', href: '#dashboard' },
    { type: 'link', text: 'Instances', href: '#instances' },
    { type: 'link', text: 'Students', href: '#students' },
    { type: 'link', text: 'Analytics', href: '#analytics' }
  ]

  const topNavigation = (
    <TopNavigation
      identity={{ title: "LFR Tools", href: "/" }}
      utilities={[{
        type: "menu-dropdown",
        text: userInfo.username,
        description: userInfo.role,
        items: [
          { id: "settings", text: "Settings" },
          { id: "signout", text: "Sign Out" }
        ]
      }]}
    />
  )

  const sideNavigation = (
    <SideNavigation
      header={{ text: userInfo.project, href: "#" }}
      items={navigationItems}
    />
  )

  const content = (
    <Container header={<Header variant="h1">Class Dashboard</Header>}>
      <SpaceBetween direction="vertical" size="l">
        <Grid gridDefinition={[{ colspan: 6 }, { colspan: 6 }]}>
          <Box>
            <Container header={<Header variant="h2">Instance Status</Header>}>
              <SpaceBetween direction="vertical" size="s">
                <Box>
                  <StatusIndicator type="success">Running: 1</StatusIndicator>
                </Box>
                <Box>
                  <StatusIndicator type="stopped">Stopped: 5</StatusIndicator>
                </Box>
                <Box>
                  Total Instances: 6
                </Box>
              </SpaceBetween>
            </Container>
          </Box>

          <Box>
            <Container header={<Header variant="h2">Quick Actions</Header>}>
              <SpaceBetween direction="vertical" size="s">
                <Button variant="primary">Start All Instances</Button>
                <Button>Stop All Instances</Button>
                <Button>View Student Status</Button>
              </SpaceBetween>
            </Container>
          </Box>
        </Grid>

        <Container header={<Header variant="h2">Recent Activity</Header>}>
          <Box>
            <p>• duke-gpu instance running (GPU enabled)</p>
            <p>• 5 instances available for class use</p>
            <p>• All systems operational</p>
          </Box>
        </Container>
      </SpaceBetween>
    </Container>
  )

  return (
    <>
      {topNavigation}
      <AppLayout
        navigation={sideNavigation}
        content={content}
        disableContentHeaderOverlap={true}
      />
    </>
  )
}

export default App