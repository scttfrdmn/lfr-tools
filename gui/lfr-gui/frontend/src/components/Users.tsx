import React from 'react'
import { Container, Header, Alert } from '@cloudscape-design/components'

interface UserInfo {
  role: string
  username: string
  project: string
  permissions: string[]
}

interface UsersProps {
  userInfo: UserInfo
}

const Users: React.FC<UsersProps> = ({ userInfo }) => {
  return (
    <Container header={<Header variant="h1">Users & Groups</Header>}>
      <Alert statusIconAriaLabel="Info" type="info">
        User and group management interface coming soon.
        This will provide bulk user creation, group management, and permission controls.
      </Alert>
    </Container>
  )
}

export default Users