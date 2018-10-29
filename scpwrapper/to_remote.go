package scpwrapper

import (
	"github.com/kballard/go-shellquote"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)
// scp a file to a remote server using ssh / scp
func CopyToRemote(source string, dest string, session *ssh.Session) error {
	session.Stdout = nil
	r, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(source, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	cmd := shellquote.Join("scp", "-qt", dest)
	if err := session.Start(cmd); err != nil {
		return err
	}

	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}

	if err := session.Wait(); err != nil {
		return err
	}
	session.Close()
	return nil
}