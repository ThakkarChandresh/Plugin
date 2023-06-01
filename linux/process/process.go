package process

import (
	"golang.org/x/crypto/ssh"
	"strings"
)

const (
	processInfoCommand string = "ps aux | awk 'NR>1 {print $2 \" \" $3 \"% \" $4 \"% \" $1\" \"$11}'"
)

func GetProcessMetrics(connection *ssh.Client, channel chan map[string]any) {
	response := make(map[string]any)

	defer func() {
		channel <- response
	}()

	defer func() {
		if r := recover(); r != nil {
			response["system.process"] = map[string]any{"status": "failed!", "err": r}
		}
	}()

	session, err := connection.NewSession()

	//Session will automatically close

	if err != nil {
		panic(err.Error())
	}

	output, err := session.Output(processInfoCommand)

	if err != nil {
		panic(err.Error())
	}

	outputSlice := strings.Split(strings.TrimSpace(strings.ReplaceAll(string(output), "\n", " ")), " ")

	resultLength := strings.Count(string(output), "\n")

	resultSlice := make([]map[string]any, resultLength)

	for i, j := 0, 0; i < len(outputSlice); i++ {
		processInfo := make(map[string]any)
		processInfo["system.process.pid"] = outputSlice[i]
		i++
		processInfo["system.process.cpu"] = outputSlice[i]
		i++
		processInfo["system.process.memory"] = outputSlice[i]
		i++
		processInfo["system.process.user"] = outputSlice[i]
		i++
		processInfo["system.process.command"] = outputSlice[i]
		resultSlice[j] = processInfo
		j++
	}

	response["system.process"] = map[string]any{"status": "success", "data": resultSlice}

}
