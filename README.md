## WAT

A wrapper around go's core implementation of `golang.org/x/crypto/ssh` to make ssh and scp usage more convenient in the regular cases you need it

## Library usage

Please see the examples here
 - [ssh](https://github.com/EugenMayer/go-antibash-boilerplate/blob/master/cmd/myssh.go)
 - [scp](https://github.com/EugenMayer/go-antibash-boilerplate/blob/master/cmd/myscp.go)
 


### Detailed
For ssh key based auth use this 
```go
package mystuff

import (
	"github.com/eugenmayer/go-sshclient/sshwrapper"
)

sshApi := sshwrapper.SshApi{Host: "somehost",Port: 22, User: "root"}
err := sshApi.DefaultSshPrivkeySetup(privkey)
if err != nil {
    log.Fatal(err)
}

stdout, stderr, err := sshApi.Run(cmd)
if err != nil {
    log.Print(stdout)
    log.Print(stderr)
    log.Fatal(err)
}
```

If you want to use your local ssh agent

```go
package mystuff

import (
	"github.com/eugenmayer/go-sshclient/sshwrapper"
)

sshApi := sshwrapper.SshApi{Host: "somehost",Port: 22, User: "root"}
err := sshApi.DefaultSshAgentSetup()
if err != nil {
    log.Fatal(err)
}

stdout, stderr, err := sshApi.Run(cmd)
if err != nil {
    log.Print(stdout)
    log.Print(stderr)
    log.Fatal(err)
}
```

And if you want to use (why..please why..) password based connection

```go
package mystuff

import (
	"github.com/eugenmayer/go-sshclient/sshwrapper"
)

sshApi := sshwrapper.SshApi{Host: "somehost",Port: 22, User: "root"}
err := sshApi.DefaultSshPasswordSetup("yoursshpassword")
if err != nil {
    log.Fatal(err)
}

stdout, stderr, err := sshApi.Run(cmd)
if err != nil {
    log.Print(stdout)
    log.Print(stderr)
    log.Fatal(err)
}
```


 