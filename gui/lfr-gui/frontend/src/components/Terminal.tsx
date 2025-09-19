import React, { useEffect, useRef, useState } from 'react'
import { Terminal as XTerm } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import '@xterm/xterm/css/xterm.css'

import { useWebSocketSSH } from '../hooks/useWebSocketSSH'

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
  const [fitAddon, setFitAddon] = useState<FitAddon | null>(null)

  // Use WebSocket SSH hook for real communication
  const { status, connect, sendInput, resize, disconnect, setOnOutput } = useWebSocketSSH(username, project)

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

    try {
      // Check instance status
      if (instanceInfo.state !== 'running') {
        terminal.writeln(`\x1b[33mInstance is ${instanceInfo.state}. Starting instance...\x1b[0m`)
        // TODO: Integrate with LFRService.StartInstance() for real instance starting
        return
      }

      if (!instanceInfo.public_ip) {
        terminal.writeln('\x1b[31mError: Instance has no public IP address\x1b[0m')
        return
      }

      terminal.writeln(`\x1b[32mConnecting to ${instanceInfo.public_ip}...\x1b[0m`)

      // Connect via WebSocket SSH proxy
      await connect()

    } catch (error) {
      terminal.writeln(`\x1b[31mConnection failed: ${error}\x1b[0m`)
    }
  }

  // Set up terminal input handling
  useEffect(() => {
    if (!terminal) return

    const handleTerminalInput = (input: string) => {
      sendInput(input)
    }

    // Handle terminal data input (when user types)
    const disposable = terminal.onData(handleTerminalInput)

    return () => {
      disposable.dispose()
    }
  }, [terminal, sendInput])

  // Set up output handling from WebSocket
  useEffect(() => {
    if (!terminal) return

    setOnOutput((output: string) => {
      terminal.write(output)
    })
  }, [terminal, setOnOutput])

  // Handle connection status changes
  useEffect(() => {
    if (!terminal) return

    if (status.connected && !status.connecting) {
      terminal.writeln('\x1b[32m✅ Connected to your cloud computer!\x1b[0m')
      onConnect?.()
    }

    if (status.error) {
      terminal.writeln(`\x1b[31m❌ ${status.error}\x1b[0m`)
    }
  }, [status, terminal, onConnect])

  // Handle terminal resize
  const handleResize = useCallback(() => {
    if (fitAddon && terminal) {
      fitAddon.fit()
      resize(terminal.cols, terminal.rows)
    }
  }, [fitAddon, terminal, resize])

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