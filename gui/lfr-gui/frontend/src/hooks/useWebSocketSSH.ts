import { useState, useEffect, useRef, useCallback } from 'react'

export interface WebSocketMessage {
  type: string
  data: any
  id?: string
  error?: string
}

export interface SSHConnectionStatus {
  connected: boolean
  connecting: boolean
  error: string | null
  publicIP: string | null
}

// Custom hook for WebSocket SSH communication
export const useWebSocketSSH = (username: string, project: string) => {
  const [status, setStatus] = useState<SSHConnectionStatus>({
    connected: false,
    connecting: false,
    error: null,
    publicIP: null
  })

  const wsRef = useRef<WebSocket | null>(null)
  const onOutputRef = useRef<((output: string) => void) | null>(null)

  const connect = useCallback(async () => {
    if (status.connecting || status.connected) return

    setStatus(prev => ({ ...prev, connecting: true, error: null }))

    try {
      // Connect to WebSocket server
      const ws = new WebSocket('ws://localhost:8080/ws')
      wsRef.current = ws

      ws.onopen = () => {
        console.log('WebSocket connected')

        // Send SSH connection request
        const connectMessage: WebSocketMessage = {
          type: 'ssh_connect',
          data: {
            username,
            project
          }
        }

        ws.send(JSON.stringify(connectMessage))
      }

      ws.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)

          switch (message.type) {
            case 'ssh_connected':
              setStatus({
                connected: true,
                connecting: false,
                error: null,
                publicIP: message.data.public_ip
              })
              break

            case 'ssh_output':
              if (onOutputRef.current) {
                onOutputRef.current(message.data.output)
              }
              break

            case 'error':
              setStatus(prev => ({
                ...prev,
                connecting: false,
                error: message.error || 'Unknown error'
              }))
              break

            case 'ssh_disconnected':
              setStatus({
                connected: false,
                connecting: false,
                error: null,
                publicIP: null
              })
              break
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      ws.onerror = (error) => {
        console.error('WebSocket error:', error)
        setStatus(prev => ({
          ...prev,
          connecting: false,
          error: 'WebSocket connection failed'
        }))
      }

      ws.onclose = () => {
        console.log('WebSocket disconnected')
        setStatus(prev => ({
          ...prev,
          connected: false,
          connecting: false
        }))
      }

    } catch (error) {
      setStatus(prev => ({
        ...prev,
        connecting: false,
        error: `Connection failed: ${error}`
      }))
    }
  }, [username, project, status.connecting, status.connected])

  const sendInput = useCallback((input: string) => {
    if (!wsRef.current || !status.connected) return

    const message: WebSocketMessage = {
      type: 'ssh_input',
      data: { input }
    }

    wsRef.current.send(JSON.stringify(message))
  }, [status.connected])

  const resize = useCallback((cols: number, rows: number) => {
    if (!wsRef.current || !status.connected) return

    const message: WebSocketMessage = {
      type: 'ssh_resize',
      data: { resize: { cols, rows } }
    }

    wsRef.current.send(JSON.stringify(message))
  }, [status.connected])

  const disconnect = useCallback(() => {
    if (!wsRef.current) return

    const message: WebSocketMessage = {
      type: 'ssh_disconnect',
      data: {}
    }

    wsRef.current.send(JSON.stringify(message))
    wsRef.current.close()
    wsRef.current = null
  }, [])

  const setOnOutput = useCallback((callback: (output: string) => void) => {
    onOutputRef.current = callback
  }, [])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [])

  return {
    status,
    connect,
    sendInput,
    resize,
    disconnect,
    setOnOutput
  }
}

// Hook for DCV WebSocket communication
export const useWebSocketDCV = (username: string, project: string) => {
  const [dcvStatus, setDcvStatus] = useState({
    connected: false,
    connecting: false,
    error: null as string | null,
    viewerURL: null as string | null
  })

  const wsRef = useRef<WebSocket | null>(null)

  const connectDCV = useCallback(async (quality: string = 'medium') => {
    if (dcvStatus.connecting || dcvStatus.connected) return

    setDcvStatus(prev => ({ ...prev, connecting: true, error: null }))

    try {
      const ws = new WebSocket('ws://localhost:8080/ws')
      wsRef.current = ws

      ws.onopen = () => {
        const connectMessage: WebSocketMessage = {
          type: 'dcv_connect',
          data: { username, project, quality }
        }
        ws.send(JSON.stringify(connectMessage))
      }

      ws.onmessage = (event) => {
        const message: WebSocketMessage = JSON.parse(event.data)

        switch (message.type) {
          case 'dcv_connected':
            setDcvStatus({
              connected: true,
              connecting: false,
              error: null,
              viewerURL: message.data.viewer_url
            })
            break

          case 'error':
            setDcvStatus(prev => ({
              ...prev,
              connecting: false,
              error: message.error || 'DCV connection failed'
            }))
            break
        }
      }

    } catch (error) {
      setDcvStatus(prev => ({
        ...prev,
        connecting: false,
        error: `DCV connection failed: ${error}`
      }))
    }
  }, [username, project, dcvStatus.connecting, dcvStatus.connected])

  return {
    dcvStatus,
    connectDCV
  }
}