import React from 'react'
import { Container, Header, Alert } from '@cloudscape-design/components'

interface UserInfo {
  role: string
  username: string
  project: string
  permissions: string[]
}

interface StudentsProps {
  userInfo: UserInfo
}

const Students: React.FC<StudentsProps> = ({ userInfo }) => {
  return (
    <Container header={<Header variant="h1">Student Management</Header>}>
      <Alert statusIconAriaLabel="Info" type="info">
        Student management interface coming soon.
        This will provide class setup, token generation, and student support features.
      </Alert>
    </Container>
  )
}

export default Students