package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	// Parse Token to get UserID (if available)
	tokenStr := r.URL.Query().Get("token")
	var userID int
	if tokenStr != "" {
		token, _ := jwt.ParseWithClaims(tokenStr, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return auth.SecretKey, nil
		})
		
		if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
			userID = claims.UserID
		}
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
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
						log.Printf("Failed to decrypt password for hostID %d user %d: %v", hostID, userID, err)
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
			h.Manager.Remove(currentSession.ID)
		} else {
			currentSession.Close()
		}
	}
}
