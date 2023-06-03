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

func GetConnection(profile map[string]interface{}) (connection *ssh.Client, err error) {

	config := &ssh.ClientConfig{

		User: fmt.Sprint(profile[CredentialProfile].(map[string]interface{})[Username]),

		Auth: []ssh.AuthMethod{ssh.Password(fmt.Sprint(profile[CredentialProfile].(map[string]interface{})[Password]))},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Timeout: SSHTimeout,
	}

	connection, err = ssh.Dial(TCP, fmt.Sprint(profile[DiscoveryProfile].(map[string]interface{})[IP], Colon, profile[DiscoveryProfile].(map[string]interface{})[Port]), config)

	return
}
