package information

import (
	"Plugin/linux/util"
	"golang.org/x/crypto/ssh"
	"strings"
)

const (
	systemInfoCommand string = `hostname |tr '\n' " " && uname |tr '\n' " " && ps -eo nlwp | awk '{ num_threads += $1 } END { print num_threads }' | tr '\n' " " && vmstat | tail -n 1 | awk '{print $12}' | tr '\n' " " && ps axo state | grep "R" | wc -l | tr '\n' " " && ps axo stat | grep "D" | wc -l && uptime -p | awk 'gsub("up ","")' && hostnamectl | grep "Operating System"`
)

func GetSystemInformationMetrics(connection *ssh.Client, channel chan map[string]interface{}) {
	response := make(map[string]interface{})

	defer func() {
		channel <- response
	}()

	defer func() {
		if r := recover(); r != nil {
			//response[util.SystemInfo] = map[string]interface{}{util.Status: util.Fail, util.Err: r}
		}
	}()

	session, err := connection.NewSession()

	//Session will automatically close

	if err != nil {
		return
	}

	output, err := session.Output(systemInfoCommand)

	if err != nil {
		return
	}

	outputInfo := strings.Split(string(output), util.NewLine)

	uptime := outputInfo[1]

	operatingSystemOutput := outputInfo[2]

	outputInfo = strings.Split(strings.TrimSpace(outputInfo[0]), util.SpaceSeparator)

	operatingSystemInfo := strings.Split(strings.TrimSpace(operatingSystemOutput), util.Colon)

	response[util.OsVersion] = operatingSystemInfo[1]

	response[util.SystemName] = outputInfo[0]

	response[util.OsName] = outputInfo[1]

	response[util.Uptime] = uptime

	response[util.Threads] = outputInfo[2]

	response[util.ContextSwitches] = outputInfo[3]

	response[util.RunningProcesses] = outputInfo[4]

	response[util.BlockProcesses] = outputInfo[5]
}
