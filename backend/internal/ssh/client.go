package ssh

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

// Connect connects to the remote SSH server
func Connect(host string, port int, user, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
				answers = make([]string, len(questions))
				for i := range questions {
					answers[i] = password
				}
				return answers, nil
			}),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For MVP, ignore host key check
		Timeout:         5 * time.Second,
	}
	
	addr := fmt.Sprintf("%s:%d", host, port)
	return ssh.Dial("tcp", addr, config)
}

// StartShell starts a shell session
func StartShell(client *ssh.Client, cols, rows int) (*ssh.Session, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // Enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	
	if err := session.RequestPty("xterm", rows, cols, modes); err != nil {
		session.Close()
		return nil, err
	}
	
	return session, nil
}
