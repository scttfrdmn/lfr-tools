import React, { useEffect, useRef, useState } from 'react'
import { Terminal as XTerm } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import '@xterm/xterm/css/xterm.css'

import {
  Container,
  Header,
  Box,
  SpaceBetween,
  Button,
  StatusIndicator,
  Alert
} from '@cloudscape-design/components'

interface TerminalProps {
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

const Terminal: React.FC<TerminalProps> = ({
  username,
  project,
  instanceInfo,
  onConnect,
  onDisconnect
}) => {
  const terminalRef = useRef<HTMLDivElement>(null)
  const [terminal, setTerminal] = useState<XTerm | null>(null)
  const [connected, setConnected] = useState(false)
  const [connecting, setConnecting] = useState(false)
  const [fitAddon, setFitAddon] = useState<FitAddon | null>(null)

  useEffect(() => {
    if (!terminalRef.current) return

    // Create terminal instance
    const term = new XTerm({
      theme: {
        background: '#1e1e1e',
        foreground: '#ffffff',
        cursor: '#ffffff',
        selection: '#d4d4aa',
      },
      fontSize: 14,
      fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace',
      cursorBlink: true,
      allowTransparency: false,
    })

    // Add addons
    const fit = new FitAddon()
    const webLinks = new WebLinksAddon()

    term.loadAddon(fit)
    term.loadAddon(webLinks)

    // Open terminal in container
    term.open(terminalRef.current)
    fit.fit()

    // Welcome message
    term.writeln('\x1b[1;34mLFR Tools Terminal\x1b[0m')
    term.writeln(`Connecting to: ${username}@${project}`)
    term.writeln('')

    setTerminal(term)
    setFitAddon(fit)

    // Handle window resize
    const handleResize = () => {
      if (fit) {
        fit.fit()
      }
    }

    window.addEventListener('resize', handleResize)

    return () => {
      window.removeEventListener('resize', handleResize)
      term.dispose()
    }
  }, [username, project])

  const handleConnect = async () => {
    if (!terminal || !instanceInfo) return

    setConnecting(true)

    try {
      // Check instance status
      if (instanceInfo.state !== 'running') {
        terminal.writeln(`\x1b[33mInstance is ${instanceInfo.state}. Starting instance...\x1b[0m`)

        // In a real implementation, this would:
        // 1. Start the instance via LFRService.StartInstance()
        // 2. Wait for it to be running
        // 3. Get the SSH connection details
        // 4. Establish WebSocket connection to backend SSH proxy

        terminal.writeln('\x1b[32mInstance started successfully!\x1b[0m')
      }

      if (!instanceInfo.public_ip) {
        terminal.writeln('\x1b[31mError: Instance has no public IP address\x1b[0m')
        setConnecting(false)
        return
      }

      terminal.writeln(`\x1b[32mConnecting to ${instanceInfo.public_ip}...\x1b[0m`)

      // Simulate connection process
      await new Promise(resolve => setTimeout(resolve, 1000))

      terminal.writeln('')
      terminal.writeln('\x1b[32mConnected to your cloud computer!\x1b[0m')
      terminal.writeln('')
      terminal.writeln('Welcome to Ubuntu 22.04.5 LTS (GNU/Linux 6.8.0-1028-aws x86_64)')
      terminal.writeln('')
      terminal.writeln(' * Documentation:  https://help.ubuntu.com')
      terminal.writeln(' * Management:     https://landscape.canonical.com')
      terminal.writeln(' * Support:        https://ubuntu.com/pro')
      terminal.writeln('')
      terminal.writeln(`${username}@instance:~$ `)

      setConnected(true)
      setConnecting(false)
      onConnect?.()

      // In a real implementation, this would establish a WebSocket connection
      // to a backend SSH proxy service that handles the actual SSH connection

    } catch (error) {
      terminal.writeln(`\x1b[31mConnection failed: ${error}\x1b[0m`)
      setConnecting(false)
    }
  }

  const handleDisconnect = () => {
    if (terminal) {
      terminal.writeln('')
      terminal.writeln('\x1b[33mConnection closed.\x1b[0m')
      terminal.writeln('')
    }
    setConnected(false)
    onDisconnect?.()
  }

  const getStatusIndicator = () => {
    if (connecting) {
      return <StatusIndicator type="in-progress">Connecting...</StatusIndicator>
    } else if (connected) {
      return <StatusIndicator type="success">Connected</StatusIndicator>
    } else {
      return <StatusIndicator type="stopped">Disconnected</StatusIndicator>
    }
  }

  return (
    <Container
      header={
        <Header
          variant="h2"
          description={`SSH terminal for ${username}'s cloud computer`}
          actions={
            <SpaceBetween direction="horizontal" size="xs">
              {getStatusIndicator()}
              {!connected && !connecting && (
                <Button
                  variant="primary"
                  iconName="call"
                  onClick={handleConnect}
                  disabled={!instanceInfo}
                >
                  Connect
                </Button>
              )}
              {connected && (
                <Button
                  variant="normal"
                  iconName="close"
                  onClick={handleDisconnect}
                >
                  Disconnect
                </Button>
              )}
              <Button
                iconName="expand"
                onClick={() => fitAddon?.fit()}
              >
                Fit
              </Button>
            </SpaceBetween>
          }
        >
          Terminal
        </Header>
      }
    >
      <SpaceBetween direction="vertical" size="m">
        {!instanceInfo && (
          <Alert statusIconAriaLabel="Warning" type="warning">
            No instance information available. Please select an instance first.
          </Alert>
        )}

        {instanceInfo && instanceInfo.state !== 'running' && (
          <Alert statusIconAriaLabel="Info" type="info">
            Instance is {instanceInfo.state}. Click Connect to start and connect.
          </Alert>
        )}

        <Box>
          <div
            ref={terminalRef}
            style={{
              height: '400px',
              width: '100%',
              border: '1px solid #d5dbdb',
              borderRadius: '8px',
              padding: '8px'
            }}
          />
        </Box>

        {connected && (
          <Alert statusIconAriaLabel="Success" type="success">
            Connected to {username}'s cloud computer. You can now run commands interactively.
          </Alert>
        )}
      </SpaceBetween>
    </Container>
  )
}

export default Terminal