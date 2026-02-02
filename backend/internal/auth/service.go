package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nexus/term/internal/database"
	"github.com/nexus/term/internal/security"
)

var SecretKey = []byte("REPLACE_THIS_WITH_RANDOM_SECRET_IN_PROD")

type Claims struct {
	UserID int    `json:"uid"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// Service handles authentication
type Service struct{}

// In-memory store for encryption keys (UserID -> Key)
// WARNING: In a real multi-instance prod, this needs Redis or similar.
// For this MVP, if server restarts, users need to re-login to populate this.
var userKeys = make(map[int][]byte)

func GetUserKey(userID int) []byte {
	return userKeys[userID]
}

// Login verifies credentials and returns a JWT
// For MVP Phase 3: We implement registration on fly if user count is 0?
// Or we just expose Register endpoint.
func (s *Service) Login(username, password string) (string, error) {
	var id int
	var hash string
	var salt []byte

	err := database.DB.QueryRow("SELECT id, password_hash, encryption_salt FROM users WHERE username = ?", username).Scan(&id, &hash, &salt)
	if err == sql.ErrNoRows {
		return "", errors.New("invalid credentials")
	} else if err != nil {
		return "", err
	}

	// Verify generic hash (simulated here since we didn't implement bcrypt yet, 
	// for Phase 3 we can just use the Derived Key as the 'Hash' for simplicity/consistency 
	// or use proper bcrypt. Let's use PBKDF2 derived key as auth proof for now to reuse code).
	// NOTE: In production, store bcrypt(password) for auth, and keep salt for key derivation separately.
	// Here we check if DerivedKey matches 'hash' (stored as hex/base64).
	
	derived := security.DeriveKey([]byte(password), salt)
	if string(derived) != hash { // Comparing raw bytes string might be unsafe for timing, but ok for MVP
		return "", errors.New("invalid credentials")
	}

	// Cache the key for this session
	userKeys[id] = derived

	// Generate JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: id,
		Name:   username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// Register creates a new user
func (s *Service) Register(username, password string) error {
	salt, err := security.GenerateSalt(16)
	if err != nil {
		return err
	}

	derived := security.DeriveKey([]byte(password), salt)
	
	// Store 'derived' as password_hash suitable for MVP. 
	// Ideally we hash this again.
	
	_, err = database.DB.Exec("INSERT INTO users (username, password_hash, encryption_salt) VALUES (?, ?, ?)", 
		username, string(derived), salt)
	if err != nil {
		// Basic check for SQLite unique constraint error
		// Using string match for MVP; purely modernc.org/sqlite might return specific error types but this is generic.
		if err.Error() == "constraint failed: UNIQUE constraint failed: users.username (2067)" || 
		   (len(err.Error()) > 0 && err.Error()[:10] == "constraint") { 
		   // Simplified check
			return errors.New("username already taken")
		}
		return err
	}
	return nil
}

// EnsureAdmin ensures an admin account exists with the given credentials.
// It creates the user if missing. If the user exists, it verifies the password.
// It ONLY updates the salt/hash if the password has changed, to preserve encryption keys.
func (s *Service) EnsureAdmin(username, password string) error {
    var id int
    var existingHash string
    var existingSalt []byte

    err := database.DB.QueryRow("SELECT id, password_hash, encryption_salt FROM users WHERE username = ?", username).Scan(&id, &existingHash, &existingSalt)
    
    if err == sql.ErrNoRows {
        // Create new user
        salt, _ := security.GenerateSalt(16)
        derived := security.DeriveKey([]byte(password), salt)
        _, err = database.DB.Exec("INSERT INTO users (username, password_hash, encryption_salt) VALUES (?, ?, ?)", 
            username, string(derived), salt)
        return err
    } else if err != nil {
        return err
    }

    // User exists, check if password matches
    derived := security.DeriveKey([]byte(password), existingSalt)
    if string(derived) == existingHash {
        // Password matches, DO NOT update. This preserves the encryption key.
        return nil
    }

    // Password changed, we must update.
    // WARNING: This invalidates all encrypted data for this user!
    // We generate a new salt for security best practice on password rotation.
    newSalt, _ := security.GenerateSalt(16)
    newDerived := security.DeriveKey([]byte(password), newSalt)
    
    _, err = database.DB.Exec("UPDATE users SET password_hash = ?, encryption_salt = ? WHERE id = ?", 
        string(newDerived), newSalt, id)
    return err
}
