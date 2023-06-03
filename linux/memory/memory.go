package memory

import (
	"Plugin/linux/util"
	"golang.org/x/crypto/ssh"
	"strings"
)

const (
	memoryInfoCommand    string = `free -b | awk 'NR>1 {print $2" " $3" " ((($2 - $7) * 100) / $2) " " $4 " " (($4 * 100) / $2) " " $7}'| head -n 1|tr '\n' " " && free -b | awk 'NR>2 {print $2}'`
	installedMemory      string = "system.memory.installed.bytes"
	usedMemory           string = "system.memory.used.bytes"
	usedMemoryPercentage string = "system.memory.used.percentage"
	freeMemory           string = "system.memory.free.bytes"
	freeMemoryPercentage string = "system.memory.free.percentage"
	availableMemory      string = "system.memory.available.bytes"
	swapMemory           string = "system.memory.swap.bytes"
)

func GetMemoryMetrics(profile map[string]interface{}, channel chan map[string]interface{}) {
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

	output, err := session.Output(memoryInfoCommand)

	if err != nil {
		return
	}

	memoryMetrics := strings.Split(strings.TrimSpace(string(output)), util.SpaceSeparator)

	response[installedMemory] = memoryMetrics[0]

	response[usedMemory] = memoryMetrics[1]

	response[usedMemoryPercentage] = memoryMetrics[2]

	response[freeMemory] = memoryMetrics[3]

	response[freeMemoryPercentage] = memoryMetrics[4]

	response[availableMemory] = memoryMetrics[5]

	response[swapMemory] = memoryMetrics[6]
}
