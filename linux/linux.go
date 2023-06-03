package linux

import (
	"Plugin/linux/cpu"
	"Plugin/linux/disk"
	"Plugin/linux/memory"
	"Plugin/linux/process"
	"Plugin/linux/system"
	"Plugin/linux/util"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

const (
	status   string = "status"
	hostname string = "hostname"
	success  string = "success"
)

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
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()

	connection, err := util.GetConnection(profile)

	if err != nil {
		return
	}

	defer func(connection *ssh.Client) {
		if closeErr := connection.Close(); closeErr != nil {
			err = closeErr
		}
	}(connection)

	result, err := executeCommand(hostname, connection)

	if err != nil {
		return
	}

	response[status] = success

	response[hostname] = result

	return
}

func Collect(profile map[string]interface{}) (response map[string]interface{}, err error) {
	response = make(map[string]interface{})

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()

	channel := make(chan map[string]interface{}, 5)

	defer func() {
		close(channel)
	}()

	go process.GetProcessMetrics(profile, channel)

	go memory.GetMemoryMetrics(profile, channel)

	go system.GetSystemInformationMetrics(profile, channel)

	go cpu.GetCpuMetrics(profile, channel)

	go disk.GetDiskMetrics(profile, channel)

	for i := 0; i < 5; i++ {
		output := <-channel

		for key, value := range output {
			response[key] = value
		}
	}

	return
}
