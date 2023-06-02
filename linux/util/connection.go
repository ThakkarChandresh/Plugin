package util

import (
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
		User:            profile[credentialProfile].(map[string]interface{})[username].(string),
		Auth:            []ssh.AuthMethod{ssh.Password(profile[credentialProfile].(map[string]interface{})[password].(string))},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         sshTimout * time.Second,
	}

	connection, err = ssh.Dial(tcp, profile[discoveryProfile].(map[string]interface{})[ip].(string)+Colon+profile[discoveryProfile].(map[string]interface{})[port].(string), config)

	return
}
