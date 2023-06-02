package linux

import (
	"Plugin/linux/cpu"
	"Plugin/linux/information"
	"Plugin/linux/memory"
	"Plugin/linux/process"
	"Plugin/linux/util"
	"errors"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func getConnection(profile map[string]interface{}) (connection *ssh.Client, err error) {
	config := &ssh.ClientConfig{
		User:            profile[util.CredentialProfile].(map[string]interface{})[util.Username].(string),
		Auth:            []ssh.AuthMethod{ssh.Password(profile[util.CredentialProfile].(map[string]interface{})[util.Password].(string))},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         util.SSHTimout * time.Second,
	}

	connection, err = ssh.Dial(util.TCP, profile[util.DiscoveryProfile].(map[string]interface{})[util.IP].(string)+util.Colon+profile[util.DiscoveryProfile].(map[string]interface{})[util.Port].(string), config)

	return
}

func executeCommand(command string, connection *ssh.Client) (result string, err error) {
	session, err := connection.NewSession()

	//Session will automatically close

	if err != nil {
		return
	}

	output, err := session.Output(command)

	if err != nil {
		return
	}

	result = strings.ReplaceAll(string(output), util.NewLine, util.Empty)
	return
}

func Discover(profile map[string]interface{}) (response map[string]interface{}, err error) {
	response = make(map[string]interface{})

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()

	connection, err := getConnection(profile)

	if err != nil {
		return
	}

	defer func(connection *ssh.Client) {
		if e := connection.Close(); e != nil {
			err = e
		}
	}(connection)

	command := util.Hostname

	result, err := executeCommand(command, connection)

	if err != nil {
		return
	}

	response[util.Status] = util.Discovered

	response[util.Hostname] = result

	return
}

func Collect(profile map[string]interface{}) (response map[string]interface{}, err error) {
	response = make(map[string]interface{})

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()

	connection, err := getConnection(profile)

	if err != nil {
		return
	}

	defer func(connection *ssh.Client) {
		if e := connection.Close(); e != nil {

			err = e

		}
	}(connection)

	channel := make(chan map[string]interface{}, 4)

	go process.GetProcessMetrics(connection, channel)

	go memory.GetMemoryMetrics(connection, channel)

	go information.GetSystemInformationMetrics(connection, channel)

	go cpu.GetCpuMetrics(connection, channel)

	for i := 0; i < 4; i++ {
		output := <-channel

		for key, value := range output {
			response[key] = value
		}
	}

	return
}
