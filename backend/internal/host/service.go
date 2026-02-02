package host

import (
	"errors"

	"github.com/nexus/term/internal/auth"
	"github.com/nexus/term/internal/database"
	"github.com/nexus/term/internal/security"
)

type Host struct {
	ID       int    `json:"id"`
	Alias    string `json:"alias"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // Write-only (input), never returned
}

type Service struct{}

func (s *Service) List(userID int) ([]Host, error) {
	rows, err := database.DB.Query("SELECT id, alias, hostname, port, username FROM hosts WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []Host
	for rows.Next() {
		var h Host
		if err := rows.Scan(&h.ID, &h.Alias, &h.Hostname, &h.Port, &h.Username); err != nil {
			continue
		}
		hosts = append(hosts, h)
	}
	return hosts, nil
}

func (s *Service) Add(userID int, h Host) error {
	key := auth.GetUserKey(userID)
	if key == nil {
		return errors.New("key missing: valid session required (try re-login)")
	}

	encryptedPass, err := security.Encrypt([]byte(h.Password), key)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec("INSERT INTO hosts (user_id, alias, hostname, port, username, encrypted_password) VALUES (?, ?, ?, ?, ?, ?)",
		userID, h.Alias, h.Hostname, h.Port, h.Username, encryptedPass)
	return err
}

func (s *Service) Delete(userID, hostID int) error {
	_, err := database.DB.Exec("DELETE FROM hosts WHERE id = ? AND user_id = ?", hostID, userID)
	return err
}

func (s *Service) Update(userID int, h Host) error {
    key := auth.GetUserKey(userID)
    // Only update password if provided
    if h.Password != "" {
        if key == nil {
            return errors.New("key missing")
        }
        encryptedPass, err := security.Encrypt([]byte(h.Password), key)
        if err != nil {
            return err
        }
        _, err = database.DB.Exec("UPDATE hosts SET alias=?, hostname=?, port=?, username=?, encrypted_password=? WHERE id=? AND user_id=?",
            h.Alias, h.Hostname, h.Port, h.Username, encryptedPass, h.ID, userID)
        return err
    } else {
        // Update without password
        _, err := database.DB.Exec("UPDATE hosts SET alias=?, hostname=?, port=?, username=? WHERE id=? AND user_id=?",
            h.Alias, h.Hostname, h.Port, h.Username, h.ID, userID)
        return err
    }
}

func (s *Service) GetDecryptedPassword(userID, hostID int) (string, error) {
	var encrypted []byte
	err := database.DB.QueryRow("SELECT encrypted_password FROM hosts WHERE id = ? AND user_id = ?", hostID, userID).Scan(&encrypted)
	if err != nil {
		return "", err
	}
	
	key := auth.GetUserKey(userID)
	if key == nil {
		return "", errors.New("key missing")
	}
	
	decrypted, err := security.Decrypt(encrypted, key)
	if err != nil {
		return "", err
	}
	
	return string(decrypted), nil
}
