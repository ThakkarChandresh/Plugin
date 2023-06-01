package linux

import (
	"Plugin/linux/process"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func getConnection(profile map[string]any) (connection *ssh.Client) {
	config := &ssh.ClientConfig{
		User:            profile["credential_profile"].(map[string]any)["username"].(string),
		Auth:            []ssh.AuthMethod{ssh.Password(profile["credential_profile"].(map[string]any)["password"].(string))},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	connection, err := ssh.Dial("tcp", profile["discovery_profile"].(map[string]any)["ip"].(string)+":"+profile["discovery_profile"].(map[string]any)["port"].(string), config)

	if err != nil {
		panic(err.Error())
	}

	return
}

func executeCommand(command string, connection *ssh.Client) string {
	session, err := connection.NewSession()

	//Session will automatically close

	if err != nil {
		panic(err.Error())
	}

	result, err := session.Output(command)

	if err != nil {
		panic(err.Error())
	}

	return strings.ReplaceAll(string(result), "\n", "")
}

func Discover(profile map[string]any) (response map[string]any) {
	response = make(map[string]any)

	defer func() {
		if r := recover(); r != nil {
			response["status"] = "device not discovered!"
			response["err"] = r
		}
	}()

	connection := getConnection(profile)

	defer func(connection *ssh.Client) {
		err := connection.Close()

		if err != nil {
			panic(err.Error())
		}
	}(connection)

	command := "hostname"

	result := executeCommand(command, connection)

	response["status"] = "device discovered successfully!"
	response["hostname"] = result

	return
}

func Collect(profile map[string]any) (response map[string]any) {
	response = make(map[string]any)

	defer func() {
		if r := recover(); r != nil {
			response["status"] = "failed!"
			response["err"] = r
		}
	}()

	connection := getConnection(profile)

	defer func(connection *ssh.Client) {
		err := connection.Close()

		if err != nil {
			panic(err.Error())
		}
	}(connection)

	channel := make(chan map[string]any)

	go process.GetProcessMetrics(connection, channel)

	response = <-channel

	return
}
