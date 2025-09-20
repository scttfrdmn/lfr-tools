import React, { useState } from 'react'
import {
  Container,
  Header,
  Box,
  SpaceBetween,
  Button,
  StatusIndicator,
  Alert,
  Select,
  Badge,
  Modal,
  ProgressBar,
  Grid
} from '@cloudscape-design/components'

interface DCVViewerProps {
  username: string
  project: string
  instanceInfo?: {
    name: string
    state: string
    public_ip: string
  }
  onConnect?: () => void
  onDisconnect?: () => void
}

interface DCVConnectionInfo {
  id: string
  connected: boolean
  public_ip: string
  username: string
  session_id: string
  dcv_port: number
  viewer_url: string
}

const DCVViewer: React.FC<DCVViewerProps> = ({
  username,
  project,
  instanceInfo,
  onConnect,
  onDisconnect
}) => {
  const [connectionInfo, setConnectionInfo] = useState<DCVConnectionInfo | null>(null)
  const [connecting, setConnecting] = useState(false)
  const [quality, setQuality] = useState('medium')
  const [showAdvanced, setShowAdvanced] = useState(false)
  const [dcvStatus] = useState<'available' | 'installing' | 'error'>('available')

  const qualityOptions = [
    { label: 'Low (15 FPS) - Best for slow connections', value: 'low' },
    { label: 'Medium (30 FPS) - Balanced performance', value: 'medium' },
    { label: 'High (60 FPS) - Best quality', value: 'high' },
    { label: 'Lossless - Perfect quality (requires fast connection)', value: 'lossless' }
  ]

  const handleConnectDCV = async () => {
    if (!instanceInfo) return

    setConnecting(true)

    try {
      // Check instance status
      if (instanceInfo.state !== 'running') {
        throw new Error(`Instance is ${instanceInfo.state}. Start the instance first.`)
      }

      if (!instanceInfo.public_ip) {
        throw new Error('Instance has no public IP address')
      }

      // Simulate DCV connection establishment
      // In production, this would call DCVService.ConnectDCV()
      await new Promise(resolve => setTimeout(resolve, 2000))

      const mockConnectionInfo: DCVConnectionInfo = {
        id: `dcv-${project}-${username}-${Date.now()}`,
        connected: true,
        public_ip: instanceInfo.public_ip,
        username: username,
        session_id: `${username}-session`,
        dcv_port: 8443,
        viewer_url: `https://${instanceInfo.public_ip}:8443/#${username}-session`
      }

      setConnectionInfo(mockConnectionInfo)
      setConnecting(false)
      onConnect?.()

    } catch (error) {
      console.error('DCV connection failed:', error)
      setConnecting(false)
    }
  }

  const handleDisconnectDCV = () => {
    setConnectionInfo(null)
    onDisconnect?.()
  }

  const handleLaunchViewer = () => {
    if (!connectionInfo) return

    // In production, this would call DCVService.LaunchDCVViewer()
    window.open(connectionInfo.viewer_url, '_blank')
  }

  const getConnectionStatus = () => {
    if (connecting) {
      return <StatusIndicator type="in-progress">Connecting to desktop...</StatusIndicator>
    } else if (connectionInfo?.connected) {
      return <StatusIndicator type="success">Desktop ready</StatusIndicator>
    } else {
      return <StatusIndicator type="stopped">Disconnected</StatusIndicator>
    }
  }

  const renderDCVControls = () => (
    <SpaceBetween direction="vertical" size="m">
      {/* Quality Settings */}
      <Container header={<Header variant="h3">🎨 Display Quality</Header>}>
        <SpaceBetween direction="vertical" size="s">
          <Select
            selectedOption={qualityOptions.find(opt => opt.value === quality) || null}
            onChange={({ detail }) => setQuality(detail.selectedOption.value!)}
            options={qualityOptions}
            placeholder="Choose quality setting"
          />

          <Box color="text-body-secondary" fontSize="body-s">
            💡 Higher quality requires more bandwidth. Start with Medium and adjust as needed.
          </Box>
        </SpaceBetween>
      </Container>

      {/* Connection Actions */}
      <SpaceBetween direction="horizontal" size="m">
        {!connectionInfo ? (
          <Button
            variant="primary"
            iconName="contact"
            onClick={handleConnectDCV}
            loading={connecting}
            disabled={!instanceInfo || instanceInfo.state !== 'running'}
          >
            {connecting ? 'Connecting...' : 'Connect to Desktop'}
          </Button>
        ) : (
          <>
            <Button
              variant="primary"
              iconName="external"
              onClick={handleLaunchViewer}
            >
              Open Desktop Viewer
            </Button>
            <Button
              variant="normal"
              iconName="close"
              onClick={handleDisconnectDCV}
            >
              Disconnect
            </Button>
          </>
        )}

        <Button
          variant="link"
          iconName="settings"
          onClick={() => setShowAdvanced(!showAdvanced)}
        >
          Advanced Options
        </Button>
      </SpaceBetween>
    </SpaceBetween>
  )

  const renderDesktopPreview = () => (
    <Container header={<Header variant="h3">🖥️ Desktop Preview</Header>}>
      <Box>
        {connectionInfo ? (
          <SpaceBetween direction="vertical" size="s">
            <Box textAlign="center" padding="m">
              <div
                style={{
                  width: '100%',
                  height: '200px',
                  backgroundColor: '#1e1e1e',
                  borderRadius: '8px',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: 'white',
                  fontSize: '18px'
                }}
              >
                🖥️ Ubuntu Desktop Ready
                <br />
                <small>Click "Open Desktop Viewer" to access</small>
              </div>
            </Box>

            <Grid gridDefinition={[{ colspan: 6 }, { colspan: 6 }]}>
              <Box>
                <strong>Session Info:</strong>
                <br />
                ID: {connectionInfo.session_id}
                <br />
                Port: {connectionInfo.dcv_port}
                <br />
                Quality: {quality}
              </Box>
              <Box>
                <strong>Connection:</strong>
                <br />
                IP: {connectionInfo.public_ip}
                <br />
                Status: <Badge color="green">Active</Badge>
                <br />
                Encryption: <Badge color="blue">TLS</Badge>
              </Box>
            </Grid>
          </SpaceBetween>
        ) : (
          <Box textAlign="center" padding="xl" color="text-body-secondary">
            Desktop will appear here after connecting
          </Box>
        )}
      </Box>
    </Container>
  )

  return (
    <Container
      header={
        <Header
          variant="h2"
          description={`Remote desktop access for ${username}'s cloud computer`}
          actions={getConnectionStatus()}
        >
          🖥️ Remote Desktop (DCV)
        </Header>
      }
    >
      <SpaceBetween direction="vertical" size="l">
        {!instanceInfo && (
          <Alert statusIconAriaLabel="Warning" type="warning">
            No instance information available. Please select an instance first.
          </Alert>
        )}

        {instanceInfo && instanceInfo.state !== 'running' && (
          <Alert statusIconAriaLabel="Info" type="info">
            Instance is {instanceInfo.state}. Start the instance to access the desktop.
          </Alert>
        )}

        {dcvStatus === 'installing' && (
          <Alert statusIconAriaLabel="Info" type="info" header="Setting up Desktop">
            <SpaceBetween direction="vertical" size="s">
              <Box>Installing and configuring NICE DCV on your instance...</Box>
              <ProgressBar value={65} label="Installation progress" />
            </SpaceBetween>
          </Alert>
        )}

        {instanceInfo && instanceInfo.state === 'running' && (
          <Grid gridDefinition={[{ colspan: 8 }, { colspan: 4 }]}>
            <SpaceBetween direction="vertical" size="m">
              {renderDCVControls()}
              {renderDesktopPreview()}
            </SpaceBetween>

            {/* Desktop Info Panel */}
            <Container header={<Header variant="h3">💡 Desktop Features</Header>}>
              <SpaceBetween direction="vertical" size="s">
                <Box>
                  <strong>Available Applications:</strong>
                  <ul style={{ margin: '8px 0', paddingLeft: '20px' }}>
                    <li>Firefox web browser</li>
                    <li>VS Code editor</li>
                    <li>Terminal applications</li>
                    <li>File manager</li>
                    <li>LibreOffice suite</li>
                  </ul>
                </Box>

                <Box>
                  <strong>Desktop Benefits:</strong>
                  <ul style={{ margin: '8px 0', paddingLeft: '20px' }}>
                    <li>Visual programming tools</li>
                    <li>Copy/paste between local and cloud</li>
                    <li>Multiple windows and applications</li>
                    <li>Graphical data visualization</li>
                  </ul>
                </Box>

                <Box fontSize="body-s" color="text-body-secondary">
                  💡 Use Terminal for command-line work, Desktop for visual applications
                </Box>
              </SpaceBetween>
            </Container>
          </Grid>
        )}

        {/* Advanced Options Modal */}
        <Modal
          onDismiss={() => setShowAdvanced(false)}
          visible={showAdvanced}
          header="Advanced DCV Options"
          footer={
            <Box float="right">
              <Button variant="link" onClick={() => setShowAdvanced(false)}>
                Close
              </Button>
            </Box>
          }
        >
          <SpaceBetween direction="vertical" size="m">
            <Container header={<Header variant="h3">Connection Settings</Header>}>
              <SpaceBetween direction="vertical" size="s">
                <Box>
                  <strong>Keyboard Layout:</strong> US English (auto-detected)
                </Box>
                <Box>
                  <strong>Screen Resolution:</strong> Auto-fit to viewer window
                </Box>
                <Box>
                  <strong>Color Depth:</strong> 24-bit (16.7 million colors)
                </Box>
                <Box>
                  <strong>Compression:</strong> Automatic based on connection speed
                </Box>
              </SpaceBetween>
            </Container>

            <Container header={<Header variant="h3">Performance Tips</Header>}>
              <ul style={{ margin: 0, paddingLeft: '20px' }}>
                <li>Close unused applications on your local computer for better performance</li>
                <li>Use "Medium" quality for most work, "High" for detailed graphics</li>
                <li>Fullscreen mode provides the best desktop experience</li>
                <li>Right-click for context menus just like on your local computer</li>
              </ul>
            </Container>
          </SpaceBetween>
        </Modal>
      </SpaceBetween>
    </Container>
  )
}

export default DCVViewer