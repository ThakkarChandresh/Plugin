package util

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

const (
	CredentialProfile string        = "credential_profile"
	Username          string        = "username"
	Password          string        = "password"
	SSHTimeout        time.Duration = 30 * time.Second
	TCP               string        = "tcp"
	DiscoveryProfile  string        = "discovery_profile"
	IP                string        = "ip"
	Port              string        = "port"
)

func getConnection(profile map[string]interface{}) (connection *ssh.Client, err error) {

	config := &ssh.ClientConfig{

		User: fmt.Sprint(profile[CredentialProfile].(map[string]interface{})[Username]),

		Auth: []ssh.AuthMethod{ssh.Password(fmt.Sprint(profile[CredentialProfile].(map[string]interface{})[Password]))},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Timeout: SSHTimeout,
	}

	connection, err = ssh.Dial(TCP, fmt.Sprint(profile[DiscoveryProfile].(map[string]interface{})[IP], Colon, profile[DiscoveryProfile].(map[string]interface{})[Port]), config)

	return
}

func ExecuteCommand(profile map[string]interface{}, command string) (output []byte, err error) {

	connection, err := getConnection(profile)

	if err != nil {
		return
	}

	defer func(connection *ssh.Client) {
		if closeErr := connection.Close(); closeErr != nil {
			err = closeErr
		}
	}(connection)

	session, err := connection.NewSession()

	//Session will automatically close

	if err != nil {
		return
	}

	output, err = session.Output(command)

	if err != nil {
		return
	}

	return
}
