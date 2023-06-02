package process

import (
	"Plugin/linux/util"
	"golang.org/x/crypto/ssh"
	"strings"
)

const (
	processInfoCommand string = `ps aux | awk 'NR> 1 {print $2 " " $3 "% " $4 "% " $1" "$11}'`
	systemProcess      string = "system.process"
	processPID         string = "system.process.pid"
	processCPU         string = "system.process.cpu"
	processMemory      string = "system.process.memory"
	processUser        string = "system.process.user"
	processCommand     string = "system.process.command"
)

func GetProcessMetrics(profile map[string]interface{}, channel chan map[string]interface{}) {
	response := make(map[string]interface{})

	defer func() {
		channel <- response
	}()

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	connection, err := util.GetConnection(profile)

	if err != nil {
		return
	}

	defer func(connection *ssh.Client) {
		if e := connection.Close(); e != nil {

			err = e

		}
	}(connection)

	session, err := connection.NewSession()

	//Session will automatically close

	if err != nil {
		return
	}

	output, err := session.Output(processInfoCommand)

	if err != nil {
		return
	}

	outputInfo := strings.Split(strings.TrimSpace(strings.ReplaceAll(string(output), util.NewLine, util.SpaceSeparator)), util.SpaceSeparator)

	resultLength := strings.Count(string(output), "\n")

	result := make([]map[string]interface{}, resultLength)

	for i, j := 0, 0; i < len(outputInfo); i++ {

		processInfo := make(map[string]interface{})

		processInfo[processPID] = outputInfo[i]

		i++

		processInfo[processCPU] = outputInfo[i]

		i++

		processInfo[processMemory] = outputInfo[i]

		i++

		processInfo[processUser] = outputInfo[i]

		i++

		processInfo[processCommand] = outputInfo[i]

		result[j] = processInfo

		j++
	}

	response[systemProcess] = result

}
