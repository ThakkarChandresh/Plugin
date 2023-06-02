package memory

import (
	"Plugin/linux/util"
	"golang.org/x/crypto/ssh"
	"strings"
)

const (
	memoryInfoCommand string = `free -b | awk 'NR>1 {print $2" " $3" " ((($2 - $7) * 100) / $2) " " $4 " " (($4 * 100) / $2) " " $7}'| head -n 1|tr '\n' " " && free -b | awk 'NR>2 {print $2}'`
)

func GetMemoryMetrics(connection *ssh.Client, channel chan map[string]interface{}) {
	response := make(map[string]interface{})

	defer func() {
		channel <- response
	}()

	defer func() {
		if r := recover(); r != nil {
			//response[util.SystemMemory] = map[string]interface{}{util.Status: util.Fail, util.Err: r}
		}
	}()

	session, err := connection.NewSession()

	//Session will automatically close

	if err != nil {
		//response[util.SystemMemory] = map[string]interface{}{util.Status: util.Fail, util.Err: err.Error()}
		return
	}

	output, err := session.Output(memoryInfoCommand)

	if err != nil {
		//response[util.SystemMemory] = map[string]interface{}{util.Status: util.Fail, util.Err: err.Error()}
		return
	}

	outputInfo := strings.Split(strings.TrimSpace(strings.ReplaceAll(string(output), util.NewLine, util.SpaceSeparator)), util.SpaceSeparator)

	response[util.InstalledMemory] = outputInfo[0]

	response[util.UsedMemory] = outputInfo[1]

	response[util.UsedMemoryPercentage] = outputInfo[2]

	response[util.FreeMemory] = outputInfo[3]

	response[util.FreeMemoryPercentage] = outputInfo[4]

	response[util.AvailableMemory] = outputInfo[5]

	response[util.SwapMemory] = outputInfo[6]
}
