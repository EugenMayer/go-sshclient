package scpwrapper

import (
	"golang.org/x/crypto/ssh"
	"os"
)

// scp a file from a remote server using ssh / scp
func CopyFromRemote(source string, dest string, session *ssh.Session) error {
	r, err := session.Output("dd if=" + source)
	if err != nil {
		return err
	}
	defer session.Close()
	file, err := os.Create(dest)
	if err != nil {
		return err
	}

	_, err = file.Write(r)
	return err
}