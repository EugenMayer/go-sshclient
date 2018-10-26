package scpwrapper

import (
	"golang.org/x/crypto/ssh"
	"os"
)

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
	//write to local file
	file.Write(r)
	return nil
}