package ssh

import (
	"io"
)

// PtyWrapper combines stdin and stdout pipes into a ReadWriteCloser
type PtyWrapper struct {
	Stdin  io.WriteCloser
	Stdout io.Reader
}

func (p *PtyWrapper) Read(b []byte) (n int, err error) {
	return p.Stdout.Read(b)
}

func (p *PtyWrapper) Write(b []byte) (n int, err error) {
	return p.Stdin.Write(b)
}

func (p *PtyWrapper) Close() error {
	return p.Stdin.Close()
}
