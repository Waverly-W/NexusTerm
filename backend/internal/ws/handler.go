package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/nexus/term/internal/auth"
	"github.com/nexus/term/internal/host"
	"github.com/nexus/term/internal/ssh"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for MVP
	},
}

type Handler struct {
	Manager *ssh.SessionManager
	HostSvc *host.Service
}

func NewHandler(manager *ssh.SessionManager, hostSvc *host.Service) *Handler {
	return &Handler{
		Manager: manager,
		HostSvc: hostSvc,
	}
}

// ServeHTTP handles the websocket connection
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Printf("[WS] Connection attempt from %s, Origin: %s", r.RemoteAddr, r.Header.Get("Origin"))

	// Parse Token to get UserID (if available)
	tokenStr := r.URL.Query().Get("token")
	var userID int
	if tokenStr != "" {
		token, err := jwt.ParseWithClaims(tokenStr, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return auth.SecretKey, nil
		})
        
        if err != nil {
            log.Printf("[WS] Token parse error: %v", err)
        }
		
		if token != nil {
            if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
                userID = claims.UserID
                log.Printf("[WS] Authenticated UserID: %d", userID)
            } else {
                log.Printf("[WS] Token invalid or claims mismatch")
            }
        }
	} else {
        log.Printf("[WS] No token provided")
    }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[WS] Upgrade failed:", err)
		return
	}
    log.Println("[WS] Upgrade successful, connection established")
	defer ws.Close()

	// ...

	var currentSession *ssh.Session

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		
		if mt == websocket.TextMessage {
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("json error:", err)
				continue
			}
			
			msgType, _ := msg["t"].(string)
			
			switch msgType {
			case "connect": // Resume or Start
				// { "t": "connect", "host": "...", "user": "...", "password": "..." }
				// OR { "t": "connect", "host_id": 123 }
				
				hostIDStr, _ := msg["host_id"].(string) 
				var hostID float64 
				if v, ok := msg["host_id"].(float64); ok {
					hostID = v
				} else if hostIDStr != "" {
					if v, err := strconv.ParseFloat(hostIDStr, 64); err == nil {
						hostID = v
					}
				}

				// Check for Session Resumption
				sessionID, _ := msg["session_id"].(string)
				keepAliveMin, _ := msg["keep_alive"].(float64)

				if sessionID != "" && h.Manager != nil {
					existingSession := h.Manager.Get(sessionID)
					if existingSession != nil {
						// Resume
						log.Printf("Resuming session %s", sessionID)
						existingSession.AttachWebSocket(ws)
						currentSession = existingSession
						
						// Update KeepAlive if provided (optional, or stick to initial)
						if keepAliveMin > 0 {
							existingSession.MaxKeepAlive = time.Duration(keepAliveMin) * time.Minute
						}

						ws.WriteJSON(map[string]string{
							"status": "connected",
							"id": currentSession.ID,
						})

						// Send recent history
						// Need thread safety for reading history? RingBuffer is sync-safe-ish but verify
						hist := existingSession.HistoryBuf.Bytes()
						if len(hist) > 0 {
							ws.WriteMessage(websocket.BinaryMessage, hist)
						}

						// Refresh size?
						existingSession.SSHSession.WindowChange(existingSession.Dimensions.Rows, existingSession.Dimensions.Cols)
						
						continue
					} else {
						// Session not found (expired?)
						ws.WriteJSON(map[string]string{"error": "session_expired_or_not_found"})
						// Client should try fresh connect
						continue
					}
				}

				hostAddr, _ := msg["host"].(string)
				port := 22
				user, _ := msg["user"].(string)
				password, _ := msg["password"].(string)

				if hostID > 0 && h.HostSvc != nil && userID > 0 {
					decryptedPass, err := h.HostSvc.GetDecryptedPassword(userID, int(hostID))
					if err == nil {
						hosts, _ := h.HostSvc.List(userID)
						for _, hst := range hosts {
							if hst.ID == int(hostID) {
								hostAddr = hst.Hostname
								user = hst.Username
								port = hst.Port
								password = decryptedPass
								break
							}
						}
					} else {
						// Log this error
						log.Printf("Failed to decrypt password for hostID %0.f user %d: %v", hostID, userID, err)
					}
				} else {
					if hostID > 0 {
						log.Printf("Lookup skipped: HostSvc=%v UserID=%d", h.HostSvc != nil, userID)
					}
				}
				
				if hostAddr == "" || user == "" || password == "" {
					errMsg := "missing_credentials"
					if hostID > 0 && password == "" { errMsg = "host_lookup_failed_or_key_missing" }
					log.Println("Connection failed:", errMsg)
					ws.WriteJSON(map[string]string{"error": errMsg})
					continue
				}

				client, err := ssh.Connect(hostAddr, port, user, password)
				if err != nil {
					ws.WriteJSON(map[string]string{"error": err.Error()})
					continue
				}
				
				// Request Pty
				rows := 24
				cols := 80
				sshSess, err := ssh.StartShell(client, cols, rows)
				if err != nil {
					ws.WriteJSON(map[string]string{"error": err.Error()})
					client.Close()
					continue
				}
				
				// Get Pipes
				stdin, err := sshSess.StdinPipe()
				if err != nil {
					ws.WriteJSON(map[string]string{"error": "stdin pipe failed"})
					sshSess.Close()
					client.Close()
					continue
				}
				stdout, err := sshSess.StdoutPipe()
				if err != nil {
					ws.WriteJSON(map[string]string{"error": "stdout pipe failed"})
					sshSess.Close()
					client.Close()
					continue
				}
				// We also need Stderr, usually merge with Stdout for MVP
				// stderr, _ := sshSess.StderrPipe()
				
				// Combine pipes
				// Note: PtyWrapper only takes one reader. 
				// For now, let's just use stdout. Merging stderr is better but complicated without multi-reader.
				// Or we can use io.MultiReader if we want.
				// StartShell uses RequestPty which usually merges stderr to stdout on server side?
				// Yes, pty merges them.
				
				pty := &ssh.PtyWrapper{
					Stdin:  stdin,
					Stdout: stdout,
				}
				
				// Start remote shell
				if err := sshSess.Shell(); err != nil {
					ws.WriteJSON(map[string]string{"error": "shell failed"})
					sshSess.Close()
					client.Close()
					continue
				}
				
				// Create Session
				currentSession = ssh.NewSession(client, sshSess, pty, rows, cols)
				currentSession.WsConn = ws
				
				// Add to manager (optional for simple loopback, but good practice)
				if h.Manager != nil {
					h.Manager.Add(currentSession)
				}
				
				// Set KeepAlive
				if keepAliveMin > 0 {
					currentSession.MaxKeepAlive = time.Duration(keepAliveMin) * time.Minute
				}
				
				// Start Pumper in background
				go currentSession.StartPumper()
				// Also read stderr if separated? Assume merged by PTY.

				ws.WriteJSON(map[string]string{
					"status": "connected",
					"id": currentSession.ID,
				})
				
			case "r": // Resize
				rows, _ := msg["rows"].(float64)
				cols, _ := msg["cols"].(float64)
				if currentSession != nil {
					currentSession.Resize(int(rows), int(cols))
				}
			case "d": // Data
				dataStr, _ := msg["d"].(string)
				if currentSession != nil && currentSession.Pty != nil {
					currentSession.Pty.Write([]byte(dataStr))
				}
			}
		} else if mt == websocket.BinaryMessage {
			// Binary input from xterm
			if currentSession != nil && currentSession.Pty != nil {
				currentSession.Pty.Write(message)
			}
		}
	}
	
	if currentSession != nil {
		if h.Manager != nil {
			if currentSession.MaxKeepAlive > 0 {
				log.Printf("Detaching session %s (KeepAlive: %v)", currentSession.ID, currentSession.MaxKeepAlive)
				currentSession.DetachWebSocket()
				currentSession.KeepAliveUntil = time.Now().Add(currentSession.MaxKeepAlive)
			} else {
				h.Manager.Remove(currentSession.ID)
			}
		} else {
			currentSession.Close()
		}
	}
}
