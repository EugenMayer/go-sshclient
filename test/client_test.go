package test

import (
	"github.com/eugenmayer/go-exec/runner"
	"github.com/eugenmayer/go-sshclient/sshwrapper"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strconv"
	"testing"
)

// docker-compose up needs to be done before so we have a ssh server up and running
func TestSshKeyAuth(t *testing.T) {
	port, err := strconv.Atoi(os.Getenv("SSHPORT"))
	host := os.Getenv("SSHHOSTNAME")
	sshApi, err := sshwrapper.DefaultSshApiSetup(host, port, "root", "sshkeys/id_rsa")
	if assert.Nil(t, err, "No connection error occurred") {
		stdout, stderr, err := sshApi.Run("echo hello")
		assert.Nil(t, err, "No error during running the command")
		assert.Equal(t, "hello\r\n", stdout, "Command return hello")
		assert.Equal(t, "", stderr, "No output on stderr")
	}
}

// password based auth
func TestSshPasswordAuth(t *testing.T) {
	port, err := strconv.Atoi(os.Getenv("SSHPORT"))
	host := os.Getenv("SSHHOSTNAME")
	sshApi := sshwrapper.NewSshApi(host, port, "root", "")
	sshApi.Password = "test"
	err = sshApi.DefaultSshPasswordSetup()
	if assert.Nil(t, err, "No connection error occurred") {
		stdout, stderr, err := sshApi.Run("echo hello")
		assert.Nil(t, err, "No error during running the command")
		assert.Equal(t, "hello\r\n", stdout, "Command return hello")
		assert.Equal(t, "", stderr, "No output on stderr")
	}
}

// run muliple coammnds to ensure we properly close sessions and connections
func TestMultiplieCommandsRusageWorks(t *testing.T) {
	port, _ := strconv.Atoi(os.Getenv("SSHPORT"))
	host := os.Getenv("SSHHOSTNAME")
	sshApi,err := sshwrapper.DefaultSshApiSetup(host, port, "root", "sshkeys/id_rsa")

	stdout, stderr, err := sshApi.Run("echo hello")
	assert.Nil(t, err, "No error during running the command")
	assert.Equal(t, "hello\r\n", stdout, "Command return hello")
	assert.Equal(t, "", stderr, "No output on stderr")

	stdout, stderr, err = sshApi.Run("echo hoho")
	assert.Nil(t, err, "No error during running the command")
	assert.Equal(t, "hoho\r\n", stdout, "Command return hoho")
	assert.Equal(t, "", stderr, "No output on stderr")

	stdout, stderr, err = sshApi.Run("echo heyhey")
	assert.Nil(t, err, "No error during running the command")
	assert.Equal(t, "heyhey\r\n", stdout, "Command return heyhey")
	assert.Equal(t, "", stderr, "No output on stderr")
}

// test direct scp implementation
func TestScpToRemote(t *testing.T) {
	port, err := strconv.Atoi(os.Getenv("SSHPORT"))
	host := os.Getenv("SSHHOSTNAME")

	sshApi,err := sshwrapper.DefaultSshApiSetup(host, port, "root", "sshkeys/id_rsa")

	// create a 10mb dummy file
	f, err := os.Create("/tmp/dummyfile")
	if err != nil {
		log.Fatal(err)
	}

	if err := f.Truncate(1e7); err != nil {
		log.Fatal(err)
	}
	f.Close()

	// copy to remote
	err = sshApi.CopyToRemote("/tmp/dummyfile", "/tmp/remotefile")
	assert.Nil(t, err, "No error during copying to the remote")
	_, _, err = sshApi.Run("ls /tmp/remotefile")
	assert.Nil(t, err, "File does exist remotely, so local to remote worked")

	// copy from remote
	err = sshApi.CopyFromRemote("/tmp/remotefile", "/tmp/fileisback")
	assert.Nil(t, err, "No error during copying from the remote")
	_, err = os.Stat("/tmp/fileisback")
	assert.Nil(t, err, "File does exist locally too, so remote to local worked")
}

// test scp implementation using the command pattern
func TestScpToRemoteCommandPattern(t *testing.T) {
	port, err := strconv.Atoi(os.Getenv("SSHPORT"))
	host := os.Getenv("SSHHOSTNAME")

	sshApi, err := sshwrapper.DefaultSshApiSetup(host, port, "root", "sshkeys/id_rsa")
	copyToRemote := runner.CopyToRemoteFromLocalRunner{SshApi: sshApi, Verbose: true}
	copyFromRemote := runner.CopyFromRemoteToLocalRunner{SshApi: sshApi, Verbose: true}

	// create a 10mb dummy file
	f, err := os.Create("/tmp/dummyfile")
	if err != nil {
		log.Fatal(err)
	}

	if err := f.Truncate(1e7); err != nil {
		log.Fatal(err)
	}
	f.Close()

	// copy to remote
	_,_, err = copyToRemote.Copy("/tmp/dummyfile", "/tmp/remotefile")
	assert.Nil(t, err, "No error during copying to the remote")
	_, _, err = sshApi.Run("ls /tmp/remotefile")
	assert.Nil(t, err, "File does exist remotely, so local to remote worked")

	// copy from remote
	_,_, err = copyFromRemote.Copy("/tmp/remotefile", "/tmp/fileisback")
	assert.Nil(t, err, "No error during copying from the remote")
	_, err = os.Stat("/tmp/fileisback")
	assert.Nil(t, err, "File does exist locally too, so remote to local worked")
}