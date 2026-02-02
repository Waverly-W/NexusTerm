package main

"time"
import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nexus/term/internal/auth"
	"github.com/nexus/term/internal/database"
	"github.com/nexus/term/internal/host"
	"github.com/nexus/term/internal/ssh"
	"github.com/nexus/term/internal/ws"
)

func main() {
	// Initialize Database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data.db"
	}
	if err := database.Init(dbPath); err != nil {
		log.Fatal("Database init:", err)
	}

	// Initialize SessionManager
	sessionMgr := ssh.NewSessionManager()

	// Services
	authSvc := &auth.Service{}
	hostSvc := &host.Service{}

	// Admin Bootstrap
	adminUser := os.Getenv("ADMIN_USER")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminUser != "" && adminPass != "" {
		if err := authSvc.EnsureAdmin(adminUser, adminPass); err != nil {
			log.Printf("Failed to ensure admin user: %v", err)
		} else {
			log.Printf("Admin user '%s' ensured.", adminUser)
		}
	}

	// Initialize Handler (Dependency Injection)
	handler := ws.NewHandler(sessionMgr, hostSvc)

	// Helper for CORS
	enableCors := func(w *http.ResponseWriter) {
		(*w).Header().Set("Access-Control-Allow-Origin", "*")
		(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	// Setup Routes
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" {
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// ... existing logic
		username := r.FormValue("username")
		password := r.FormValue("password")

		token, err := authSvc.Login(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		w.Write([]byte(token))
	})

	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" {
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if os.Getenv("DISABLE_REGISTRATION") == "true" {
			http.Error(w, "Registration is disabled by administrator.", http.StatusForbidden)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if err := authSvc.Register(username, password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write([]byte("registered"))
	})

	http.HandleFunc("/api/hosts", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" {
			return
		}

		// Auth Check (header or query)
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		// Strip "Bearer "
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}

		token, _ := jwt.ParseWithClaims(tokenStr, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return auth.SecretKey, nil
		})

		// Get Claims
		var userID int
		if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
			userID = claims.UserID
		} else {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if r.Method == "GET" {
			hosts, err := hostSvc.List(userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(hosts)
		} else if r.Method == "POST" {
			var h host.Host
			if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := hostSvc.Add(userID, h); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("ok"))
		} else if r.Method == "PUT" {
			var h host.Host
			if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := hostSvc.Update(userID, h); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("updated"))
		}
	})

	http.HandleFunc("/api/test-connection", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" {
			return
		}
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			User     string `json:"user"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		client, err := ssh.Connect(req.Host, req.Port, req.User, req.Password)
		if err != nil {
			http.Error(w, "Connection failed: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer client.Close()

		w.Write([]byte("success"))
	})

	http.HandleFunc("/api/ws", handler.ServeHTTP)

	// Serve Frontend Static Files
	// In production (Docker), frontend is built to /static
	// For dev, we might use Vite proxy, but here we serve /static as well
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Logging Middleware
    logRequest := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
            next.ServeHTTP(lrw, r)
            log.Printf("%s %s %d %s", r.Method, r.URL.Path, lrw.statusCode, time.Since(start))
        })
    }
    
    // Wrap handler
    loggedMux := logRequest(http.DefaultServeMux)

	port := "8080"
	log.Printf("Starting NexusTerm Server on :%s", port)
	if err := http.ListenAndServe(":"+port, loggedMux); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
    lrw.statusCode = code
    lrw.ResponseWriter.WriteHeader(code)
}
