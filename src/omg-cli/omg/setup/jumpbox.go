package setup

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"omg-cli/ssh"
	"os"
	"time"
)

type Jumpbox struct {
	logger              *log.Logger
	session             *ssh.Connection
	terraformConfigPath string
}

func NewJumpbox(logger *log.Logger, ip, username, sshKeyPath, terraformConfigPath string) (*Jumpbox, error) {
	jumpboxLogger := *logger
	jumpboxLogger.SetPrefix(fmt.Sprintf("%s[jumpbox] ", jumpboxLogger.Prefix()))
	key, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		return nil, err
	}

	jb := &Jumpbox{logger: &jumpboxLogger, terraformConfigPath: terraformConfigPath}
	jb.session, err = ssh.NewConnection(&jumpboxLogger, ip, ssh.Port, username, key)
	if err != nil {
		return nil, err
	}

	return jb, nil
}

func (jb *Jumpbox) PoolTillStarted() error {
	timer := time.After(time.Duration(0 * time.Second))
	timeout := time.After(time.Duration(120 * time.Second))
	for {
		select {
		case <-timeout:
			return errors.New("Timeout waiting for Jumpbox to start")
		case <-timer:
			if err := jb.session.EnsureConnected(); err == nil {
				return nil
			}
			jb.logger.Print("waiting for Jumpbox to start")
			timer = time.After(time.Duration(5 * time.Second))
		}
	}
}

// Push the OMG binary, environment config to jumpbox
func (jb *Jumpbox) UploadDependencies() error {
	if err := jb.session.EnsureConnected(); err != nil {
		return err
	}

	me, err := os.Executable()
	if err != nil {
		return err
	}

	for _, f := range []struct {
		local string
		dest  string
	}{
		{me, "omg-cli"},
		{jb.terraformConfigPath, "env.json"},
	} {
		if err := jb.session.UploadFile(f.local, f.dest); err != nil {
			return err
		}
	}

	return nil
}

func (jb *Jumpbox) RunOmg(args string) error {
	if err := jb.session.EnsureConnected(); err != nil {
		return err
	}

	return jb.session.RunCommand(fmt.Sprintf("~/omg-cli %s", args))
}
