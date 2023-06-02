package process

import (
	"Plugin/linux/util"
	"golang.org/x/crypto/ssh"
	"strings"
)

const (
	processInfoCommand string = `ps aux | awk 'NR> 1 {print $2 " " $3 "% " $4 "% " $1" "$11}'`
)

func GetProcessMetrics(connection *ssh.Client, channel chan map[string]interface{}) {
	response := make(map[string]interface{})

	defer func() {
		channel <- response
	}()

	defer func() {
		if r := recover(); r != nil {
			//response[util.SystemProcess] = map[string]interface{}{util.Status: util.Fail, util.Err: r}
		}
	}()

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

		processInfo[util.ProcessPID] = outputInfo[i]

		i++

		processInfo[util.ProcessCPU] = outputInfo[i]

		i++

		processInfo[util.ProcessMemory] = outputInfo[i]

		i++

		processInfo[util.ProcessUser] = outputInfo[i]

		i++

		processInfo[util.ProcessCommand] = outputInfo[i]

		result[j] = processInfo

		j++
	}

	response[util.SystemProcess] = result

}
