import React from 'react'
import { Container, Header, Alert } from '@cloudscape-design/components'

interface UserInfo {
  role: string
  username: string
  project: string
  permissions: string[]
}

interface SettingsProps {
  userInfo: UserInfo
}

const Settings: React.FC<SettingsProps> = ({ userInfo }) => {
  return (
    <Container header={<Header variant="h1">Settings</Header>}>
      <Alert statusIconAriaLabel="Info" type="info">
        Settings interface coming soon.
        This will provide configuration management, preferences, and system settings.
      </Alert>
    </Container>
  )
}

export default Settings