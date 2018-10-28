package test

import (
	"github.com/eugenmayer/go-sshclient/sshwrapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

// docker-compose up needs to be done before so we have a ssh server up and running
func TestSshKeyAuth(t *testing.T) {
	sshApi := sshwrapper.SshApi{Host: "sshserver",Port: 22, User: "root"}
	err := sshApi.DefaultSshPrivkeySetup("sshkeys/id_rsa")
	if assert.Nil(t, err, "No connection error occurred") {
		stdout, stderr, err := sshApi.Run("echo hello")
		assert.Nil(t, err, "No error during running the command")
		assert.Equal(t, stdout, "hello\r\n", "Command return hello")
		assert.Equal(t, stderr, "", "No output on stderr")
	}
}

// password based auth
func TestSshPasswordAuth(t *testing.T) {
	sshApi := sshwrapper.SshApi{Host: "sshserver",Port: 22, User: "root", Password: "test"}
	err := sshApi.DefaultSshPasswordSetup()
	if assert.Nil(t, err, "No connection error occurred") {
		stdout, stderr, err := sshApi.Run("echo hello")
		assert.Nil(t, err, "No error during running the command")
		assert.Equal(t, stdout, "hello\r\n", "Command return hello")
		assert.Equal(t, stderr, "", "No output on stderr")
	}
}