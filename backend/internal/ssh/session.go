package ssh

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nexus/term/internal/utils"
	"golang.org/x/crypto/ssh"
)

type WindowSize struct {
	Rows int `json:"rows"`
	Cols int `json:"cols"`
}

// Session represents an active SSH connection
type Session struct {
	ID         string
	Client     *ssh.Client
	SSHSession *ssh.Session
	Pty        io.ReadWriteCloser
	
	HistoryBuf *utils.RingBuffer
	
	WsConn *websocket.Conn
	WsLock sync.Mutex
	
	LastActive time.Time
	Dimensions WindowSize
}

// NewSession creates a new session
func NewSession(client *ssh.Client, sshSession *ssh.Session, pty io.ReadWriteCloser, rows, cols int) *Session {
	id := generateID()
	return &Session{
		ID:         id,
		Client:     client,
		SSHSession: sshSession,
		Pty:        pty,
		HistoryBuf: utils.NewRingBuffer(50 * 1024), // 50KB buffer
		LastActive: time.Now(),
		Dimensions: WindowSize{Rows: rows, Cols: cols},
	}
}

func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// WriteToWebSocket writes data to the websocket if connected
func (s *Session) WriteToWebSocket(data []byte) error {
	s.WsLock.Lock()
	defer s.WsLock.Unlock()
	
	if s.WsConn == nil {
		return nil
	}
	
	// Binary message for raw terminal data
	return s.WsConn.WriteMessage(websocket.BinaryMessage, data)
}

// Close closes the SSH session and connection
func (s *Session) Close() {
	if s.Pty != nil {
		s.Pty.Close()
	}
	if s.SSHSession != nil {
		s.SSHSession.Close()
	}
	if s.Client != nil {
		s.Client.Close()
	}
	
	s.WsLock.Lock()
	defer s.WsLock.Unlock()
	if s.WsConn != nil {
		s.WsConn.Close()
		s.WsConn = nil
	}
}

// StartPumper starts the copy loop from Pty to RingBuffer and WebSocket
func (s *Session) StartPumper() {
	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		// 1. Read from Pty/SSH
		n, err := s.Pty.Read(buf)
		if err != nil {
			if err != io.EOF {
				// Log error?
			}
			s.Close()
			return
		}
		
		data := buf[:n]
		
		// 2. Write to RingBuffer
		s.HistoryBuf.Write(data)
		
		// 3. Write to WebSocket if connected
		// We ignore error here, if WS fails, Receiver will handle cleanup, 
		// but Pumper should continue (unless we want to detect Zombi session here too?)
		// Actually, if WS is closed, we just don't write.
		// Use a local lock to check check WS status safely.
		s.WsLock.Lock()
		ws := s.WsConn
		s.WsLock.Unlock()
		
		if ws != nil {
			err := ws.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				// Verification: If write fails, maybe WS is closed. 
				// We don't close SSH here, just detach WS.
				s.DetachWebSocket()
			}
		}
		
		s.LastActive = time.Now()
	}
}

// DetachWebSocket removes the websocket connection
func (s *Session) DetachWebSocket() {
	s.WsLock.Lock()
	defer s.WsLock.Unlock()
	if s.WsConn != nil {
		s.WsConn.Close()
		s.WsConn = nil
	}
}

// Resize resizes the terminal
func (s *Session) Resize(rows, cols int) error {
	s.Dimensions = WindowSize{Rows: rows, Cols: cols}
	return s.SSHSession.WindowChange(rows, cols)
}

