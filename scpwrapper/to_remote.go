package scpwrapper

import (
	"golang.org/x/crypto/ssh"
	"github.com/tmc/scp"
)

func CopyToRemote(source string, dest string, session *ssh.Session) error {
	err := scp.CopyPath(source, dest, session)
	session.Close()
	return err
}