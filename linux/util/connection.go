package util

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

const (
	credentialProfile string        = "credential_profile"
	username          string        = "username"
	password          string        = "password"
	sshTimout         time.Duration = 30
	tcp               string        = "tcp"
	discoveryProfile  string        = "discovery_profile"
	ip                string        = "ip"
	port              string        = "port"
)

func GetConnection(profile map[string]interface{}) (connection *ssh.Client, err error) {

	config := &ssh.ClientConfig{
		User: fmt.Sprint(profile[credentialProfile].(map[string]interface{})[username]),

		Auth: []ssh.AuthMethod{ssh.Password(fmt.Sprint(profile[credentialProfile].(map[string]interface{})[password]))},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Timeout: sshTimout * time.Second,
	}

	connection, err = ssh.Dial(tcp, fmt.Sprint(profile[discoveryProfile].(map[string]interface{})[ip], Colon, profile[discoveryProfile].(map[string]interface{})[port]), config)

	return
}
