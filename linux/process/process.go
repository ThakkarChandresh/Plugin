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

	allProcessMetrics := strings.Split(strings.Trim(string(output), util.NewLine), util.NewLine)

	result := make([]map[string]interface{}, len(allProcessMetrics))

	for i := 0; i < len(allProcessMetrics); i++ {
		process := strings.Split(allProcessMetrics[i], util.SpaceSeparator)

		processMetrics := make(map[string]interface{})

		processMetrics[processPID] = process[0]

		processMetrics[processCPU] = process[1]

		processMetrics[processMemory] = process[2]

		processMetrics[processUser] = process[3]

		processMetrics[processCommand] = process[4]

		result[i] = processMetrics
	}

	response[systemProcess] = result

}
