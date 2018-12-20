/*
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ssh

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

// Connection is the details necessary for creating SSH connections to a target host.
type Connection struct {
	logger       *log.Logger
	output       io.Writer
	client       *ssh.Client
	hostname     string
	port         int
	clientConfig *ssh.ClientConfig
}

// Port is the default SSH port.
const Port = 22

// NewConnection creates a new Connection.
// It does not actually SSH to anything.
func NewConnection(logger *log.Logger, output io.Writer, hostname string, port int, username string, key []byte) (*Connection, error) {
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	cfg := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			logger.Printf("WARNING: ignoring SSH certificate issue for %s", hostname)
			return nil
		},
	}

	c := &Connection{logger: logger, output: output, hostname: hostname, port: port, clientConfig: cfg}

	return c, nil
}

// EnsureConnected performs a test SSH session to ensure the connection is valid.
func (c *Connection) EnsureConnected() error {
	if c.client == nil {
		var err error
		c.client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.hostname, c.port), c.clientConfig)
		if err != nil {
			return err
		}
	}

	ses, err := c.client.NewSession()
	if err != nil {
		return err
	}
	return ses.Run("exit 0")
}

// UploadFile uploads a file to the ssh target host.
func (c *Connection) UploadFile(path, destName string) error {
	ses, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("uploading file: %v", err)
	}

	ses.Stdout = os.Stdout
	ses.Stderr = os.Stderr

	defer ses.Close()

	c.logger.Printf("uploading file %s as %s", path, destName)
	return scp.CopyPath(path, destName, ses)
}

// Mkdir makes a directory on the ssh target host.
func (c *Connection) Mkdir(path string) error {
	ses, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("mkdir: %v", err)
	}
	defer ses.Close()

	c.logger.Printf("mkdir -p ~/%q", path)
	return ses.Run(fmt.Sprintf("mkdir -p ~/%q", path))
}

// RunCommand executes a command on the ssh target host.
func (c *Connection) RunCommand(cmd string) error {
	ses, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("uploading file: %v", err)
	}
	defer ses.Close()

	c.logger.Printf("running command %s", cmd)

	ses.Stdout = c.output
	ses.Stderr = c.output
	if err := ses.Run(cmd); err != nil {
		return fmt.Errorf("running command: %v", err)
	}

	return nil
}

// Close closes the ssh connection.
func (c *Connection) Close() {
	c.client.Close()
}
