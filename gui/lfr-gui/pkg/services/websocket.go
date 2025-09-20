// Package services provides WebSocket transport for SSH and real-time communication
package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

// WebSocketService provides WebSocket transport for SSH and monitoring
type WebSocketService struct {
	connections map[string]*WebSocketConnection
	mu          sync.RWMutex
	upgrader    websocket.Upgrader
	sshService  *SSHProxyService
}

// WebSocketConnection represents a WebSocket connection with SSH session
type WebSocketConnection struct {
	ID         string
	WebSocket  *websocket.Conn
	SSHSession *ssh.Session
	SSHClient  *ssh.Client
	Username   string
	Project    string
	Active     bool
}

// WebSocketMessage represents messages sent over WebSocket
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	ID      string      `json:"id,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// TerminalMessage represents terminal input/output
type TerminalMessage struct {
	Input  string `json:"input,omitempty"`
	Output string `json:"output,omitempty"`
	Resize *struct {
		Cols int `json:"cols"`
		Rows int `json:"rows"`
	} `json:"resize,omitempty"`
}

// NewWebSocketService creates a new WebSocket service
func NewWebSocketService(sshService *SSHProxyService) *WebSocketService {
	return &WebSocketService{
		connections: make(map[string]*WebSocketConnection),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow connections from Wails app
				origin := r.Header.Get("Origin")
				return origin == "wails://localhost" || origin == "http://localhost:3000"
			},
		},
		sshService: sshService,
	}
}

// HandleWebSocketConnection handles WebSocket connection upgrades
func (ws *WebSocketService) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Generate connection ID
	connectionID := fmt.Sprintf("ws-%d", len(ws.connections))

	// Create WebSocket connection
	wsConn := &WebSocketConnection{
		ID:        connectionID,
		WebSocket: conn,
		Active:    true,
	}

	ws.mu.Lock()
	ws.connections[connectionID] = wsConn
	ws.mu.Unlock()

	// Handle messages
	ws.handleMessages(wsConn)

	// Cleanup
	ws.mu.Lock()
	delete(ws.connections, connectionID)
	ws.mu.Unlock()
}

// handleMessages processes WebSocket messages
func (ws *WebSocketService) handleMessages(wsConn *WebSocketConnection) {
	for {
		var msg WebSocketMessage
		err := wsConn.WebSocket.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		switch msg.Type {
		case "ssh_connect":
			ws.handleSSHConnect(wsConn, msg)
		case "ssh_input":
			ws.handleSSHInput(wsConn, msg)
		case "ssh_resize":
			ws.handleSSHResize(wsConn, msg)
		case "ssh_disconnect":
			ws.handleSSHDisconnect(wsConn)
		case "dcv_connect":
			ws.handleDCVConnect(wsConn, msg)
		case "status_subscribe":
			ws.handleStatusSubscribe(wsConn, msg)
		default:
			ws.sendError(wsConn, fmt.Sprintf("Unknown message type: %s", msg.Type))
		}
	}
}

// handleSSHConnect establishes SSH connection
func (ws *WebSocketService) handleSSHConnect(wsConn *WebSocketConnection, msg WebSocketMessage) {
	connectData, ok := msg.Data.(map[string]interface{})
	if !ok {
		ws.sendError(wsConn, "Invalid SSH connect data")
		return
	}

	username := connectData["username"].(string)
	project := connectData["project"].(string)

	// Get SSH connection info
	sshInfo, err := ws.sshService.InitiateSSHConnection(username, project)
	if err != nil {
		ws.sendError(wsConn, fmt.Sprintf("SSH connection failed: %v", err))
		return
	}

	// Get the actual SSH connection
	ws.mu.RLock()
	sshConnections := ws.sshService.connections
	ws.mu.RUnlock()

	var sshConn *SSHConnection
	for _, conn := range sshConnections {
		if conn.Username == username && conn.Project == project {
			sshConn = conn
			break
		}
	}

	if sshConn == nil {
		ws.sendError(wsConn, "SSH connection not found")
		return
	}

	// Associate SSH session with WebSocket
	wsConn.SSHSession = sshConn.Session
	wsConn.SSHClient = sshConn.Client
	wsConn.Username = username
	wsConn.Project = project

	// Start interactive session
	err = ws.sshService.StartInteractiveSession(sshInfo.ID)
	if err != nil {
		ws.sendError(wsConn, fmt.Sprintf("Failed to start interactive session: %v", err))
		return
	}

	// Start reading SSH output
	go ws.readSSHOutput(wsConn)

	// Send success response
	ws.sendMessage(wsConn, WebSocketMessage{
		Type: "ssh_connected",
		Data: map[string]interface{}{
			"username":  username,
			"public_ip": sshInfo.PublicIP,
			"connected": true,
		},
	})
}

// handleSSHInput sends input to SSH session
func (ws *WebSocketService) handleSSHInput(wsConn *WebSocketConnection, msg WebSocketMessage) {
	if wsConn.SSHSession == nil {
		ws.sendError(wsConn, "No active SSH session")
		return
	}

	terminalMsg, ok := msg.Data.(map[string]interface{})
	if !ok {
		ws.sendError(wsConn, "Invalid terminal message")
		return
	}

	input := terminalMsg["input"].(string)

	// Send input to SSH session
	stdin, err := wsConn.SSHSession.StdinPipe()
	if err != nil {
		ws.sendError(wsConn, fmt.Sprintf("Failed to get stdin: %v", err))
		return
	}

	_, err = stdin.Write([]byte(input))
	if err != nil {
		ws.sendError(wsConn, fmt.Sprintf("Failed to send input: %v", err))
		return
	}
}

// handleSSHResize handles terminal resize events
func (ws *WebSocketService) handleSSHResize(wsConn *WebSocketConnection, msg WebSocketMessage) {
	if wsConn.SSHSession == nil {
		ws.sendError(wsConn, "No active SSH session")
		return
	}

	resizeData, ok := msg.Data.(map[string]interface{})
	if !ok {
		ws.sendError(wsConn, "Invalid resize data")
		return
	}

	cols := int(resizeData["cols"].(float64))
	rows := int(resizeData["rows"].(float64))

	// Request window change
	err := wsConn.SSHSession.WindowChange(rows, cols)
	if err != nil {
		log.Printf("Failed to resize terminal: %v", err)
	}
}

// handleSSHDisconnect closes SSH connection
func (ws *WebSocketService) handleSSHDisconnect(wsConn *WebSocketConnection) {
	if wsConn.SSHSession != nil {
		wsConn.SSHSession.Close()
		wsConn.SSHSession = nil
	}
	if wsConn.SSHClient != nil {
		wsConn.SSHClient.Close()
		wsConn.SSHClient = nil
	}

	ws.sendMessage(wsConn, WebSocketMessage{
		Type: "ssh_disconnected",
		Data: map[string]interface{}{
			"username": wsConn.Username,
		},
	})
}

// handleDCVConnect establishes DCV connection
func (ws *WebSocketService) handleDCVConnect(wsConn *WebSocketConnection, msg WebSocketMessage) {
	connectData, ok := msg.Data.(map[string]interface{})
	if !ok {
		ws.sendError(wsConn, "Invalid DCV connect data")
		return
	}

	username := connectData["username"].(string)
	project := connectData["project"].(string)
	quality := connectData["quality"].(string)

	// Create DCV service
	dcvService := NewDCVService()

	// Connect DCV
	dcvInfo, err := dcvService.ConnectDCV(username, project, DCVQuality(quality))
	if err != nil {
		ws.sendError(wsConn, fmt.Sprintf("DCV connection failed: %v", err))
		return
	}

	// Send DCV connection info
	ws.sendMessage(wsConn, WebSocketMessage{
		Type: "dcv_connected",
		Data: dcvInfo,
	})
}

// handleStatusSubscribe subscribes to real-time status updates
func (ws *WebSocketService) handleStatusSubscribe(wsConn *WebSocketConnection, msg WebSocketMessage) {
	subscribeData, ok := msg.Data.(map[string]interface{})
	if !ok {
		ws.sendError(wsConn, "Invalid subscribe data")
		return
	}

	project := subscribeData["project"].(string)

	// TODO: Implement real-time status subscription
	// This would integrate with the monitoring service

	ws.sendMessage(wsConn, WebSocketMessage{
		Type: "status_subscribed",
		Data: map[string]interface{}{
			"project":    project,
			"subscribed": true,
		},
	})
}

// readSSHOutput reads output from SSH session and sends to WebSocket
func (ws *WebSocketService) readSSHOutput(wsConn *WebSocketConnection) {
	if wsConn.SSHSession == nil {
		return
	}

	stdout, err := wsConn.SSHSession.StdoutPipe()
	if err != nil {
		log.Printf("Failed to get stdout: %v", err)
		return
	}

	stderr, err := wsConn.SSHSession.StderrPipe()
	if err != nil {
		log.Printf("Failed to get stderr: %v", err)
		return
	}

	// Read stdout
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := stdout.Read(buffer)
			if err != nil {
				if err != io.EOF {
					log.Printf("SSH stdout read error: %v", err)
				}
				break
			}

			output := string(buffer[:n])
			ws.sendMessage(wsConn, WebSocketMessage{
				Type: "ssh_output",
				Data: map[string]interface{}{
					"output": output,
					"stream": "stdout",
				},
			})
		}
	}()

	// Read stderr
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := stderr.Read(buffer)
			if err != nil {
				if err != io.EOF {
					log.Printf("SSH stderr read error: %v", err)
				}
				break
			}

			output := string(buffer[:n])
			ws.sendMessage(wsConn, WebSocketMessage{
				Type: "ssh_output",
				Data: map[string]interface{}{
					"output": output,
					"stream": "stderr",
				},
			})
		}
	}()
}

// sendMessage sends a message over WebSocket
func (ws *WebSocketService) sendMessage(wsConn *WebSocketConnection, msg WebSocketMessage) {
	if !wsConn.Active {
		return
	}

	err := wsConn.WebSocket.WriteJSON(msg)
	if err != nil {
		log.Printf("Failed to send WebSocket message: %v", err)
		wsConn.Active = false
	}
}

// sendError sends an error message over WebSocket
func (ws *WebSocketService) sendError(wsConn *WebSocketConnection, errorMsg string) {
	ws.sendMessage(wsConn, WebSocketMessage{
		Type:  "error",
		Error: errorMsg,
	})
}

// BroadcastStatusUpdate broadcasts status updates to all connected clients
func (ws *WebSocketService) BroadcastStatusUpdate(project string, update interface{}) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	message := WebSocketMessage{
		Type: "status_update",
		Data: update,
	}

	for _, conn := range ws.connections {
		if conn.Project == project && conn.Active {
			ws.sendMessage(conn, message)
		}
	}
}

// StartWebSocketServer starts the WebSocket server
func (ws *WebSocketService) StartWebSocketServer(port int) error {
	http.HandleFunc("/ws", ws.HandleWebSocketConnection)

	log.Printf("Starting WebSocket server on port %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}