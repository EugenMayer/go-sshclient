package sshwrapper

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
	"os"
	"time"
)

type SshApi struct {
	SshConfig *ssh.ClientConfig
	Client    *ssh.Client
	Session   *ssh.Session
	User      string
	Password  string
	Host      string
	Port      int

	StdOut bytes.Buffer
	StdErr bytes.Buffer
}

func DefaultApiSetup(user string, host string, port int, key string) (sshApi *SshApi, err error) {
	sshApi = &SshApi{
		User: user,
		Host: host,
		Port: port,
	}

	if key == "" {
		err = sshApi.DefaultSshAgentSetup()
	} else {
		err = sshApi.DefaultSshPrivkeySetup(key)
	}
	return sshApi, err
}

func (sshApi *SshApi) DefaultSshAgentSetup() (error) {
	sshAgent, err := SSHAgent()
	if err != nil {
		return err
	}

	sshApi.SshConfig = &ssh.ClientConfig{
		User:            sshApi.User,
		Auth:            []ssh.AuthMethod{sshAgent},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return nil
}

func (sshApi *SshApi) DefaultSshPasswordSetup() (error) {
	sshApi.SshConfig = &ssh.ClientConfig{
		User:            sshApi.User,
		Auth:            []ssh.AuthMethod{ssh.Password(sshApi.Password)},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return nil
}

func (sshApi *SshApi) DefaultSshPrivkeySetup(keyPath string) (error) {
	privateKey, err := PrivateKeyFile(keyPath)
	if err != nil {
		return err
	}
	sshApi.SshConfig = &ssh.ClientConfig{
		User:            sshApi.User,
		Auth:            []ssh.AuthMethod{privateKey},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return nil
}

func SSHAgent() (ssh.AuthMethod, error) {
	socket := os.Getenv("SSH_AUTH_SOCK")
	if sshAgent, err := net.Dial("unix", socket); err != nil {
		return nil, err
	} else {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers), nil
	}
}

func PrivateKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}
