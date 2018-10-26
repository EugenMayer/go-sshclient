package sshwrapper

import (
	"github.com/EugenMayer/go-sshclient"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func (sshApi *SshApi) ConnectAndSession() (err error) {
	if client, err := sshApi.Connect(); err != nil {
		return err
	} else {
		sshApi.Client = client
	}

	return sshApi.SessionDefault()
}

func (sshApi *SshApi) SessionDefault() (err error) {
	if session, err := sshApi.Client.NewSession(); err != nil {
		return err
	} else {
		sshApi.Session = session
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := sshApi.Session.RequestPty("xterm", 80, 40, modes); err != nil {
		sshApi.Session.Close()
		return err
	}

	sshApi.Session.Stdout = &sshApi.StdOut
	sshApi.Session.Stderr = &sshApi.StdErr

	//if sshApi.StdIn, err = sshApi.Session.StdinPipe(); err != nil {
	//	go io.Copy(sshApi.StdIn, os.Stdin)
	//} else {
	//	sshApi.Session.Close()
	//	return err
	//}
	//
	//if sshApi.StdOut, err = sshApi.Session.StdoutPipe(); err != nil {
	//	go io.Copy(os.Stdout, sshApi.StdOut)
	//} else {
	//	sshApi.Session.Close()
	//	return err
	//}
	//
	//if sshApi.StdErr, err = sshApi.Session.StderrPipe(); err != nil {
	//	go io.Copy(os.Stderr, sshApi.StdErr)
	//} else {
	//	sshApi.Session.Close()
	//	return err
	//}
	return nil
}

func (sshApi *SshApi) Connect() (*ssh.Client, error) {
	var addr = fmt.Sprintf("%s:%d", sshApi.Host, sshApi.Port)
	conn, err := net.DialTimeout("tcp", addr, sshApi.SshConfig.Timeout)
	if err != nil {
		return nil, err
	}

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	c, chans, reqs, err := ssh.NewClientConn(conn, addr, sshApi.SshConfig)
	if err != nil {
		return nil, err
	}

	err = conn.SetReadDeadline(time.Time{})
	return ssh.NewClient(c, chans, reqs), err
}

func (sshApi *SshApi) GetStdOut() string {
	return sshApi.StdOut.String()
}

func (sshApi *SshApi) GetStdErr() string {
	return sshApi.StdErr.String()
}

func (sshApi *SshApi) CopyToRemote(source string, dest string) (err error) {
	sshApi.ConnectAndSession()
	err = scpwrapper.CopyToRemote(source, dest, sshApi.Session)
	sshApi.Session.Close()
	return err
}

func (sshApi *SshApi) CopyFromRemote(source string, dest string) (err error) {
	sshApi.ConnectAndSession()
	err = scpwrapper.CopyFromRemote(source, dest, sshApi.Session)
	sshApi.Session.Close()
	return err
}