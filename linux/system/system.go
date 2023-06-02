package system

import (
	"Plugin/linux/util"
	"golang.org/x/crypto/ssh"
	"strings"
)

const (
	systemInfoCommand string = `hostname |tr '\n' " " && uname |tr '\n' " " && ps -eo nlwp | awk '{ num_threads += $1 } END { print num_threads }' | tr '\n' " " && vmstat | tail -n 1 | awk '{print $12}' | tr '\n' " " && ps axo state | grep "R" | wc -l | tr '\n' " " && ps axo stat | grep "D" | wc -l && uptime -p | awk 'gsub("up ","")' && hostnamectl | grep "Operating System"`
	osVersion         string = "system.os.version"
	systemName        string = "system.name"
	osName            string = "system.os.name"
	upTime            string = "system.uptime"
	threads           string = "system.threads"
	contextSwitches   string = "system.context.switches"
	runningProcesses  string = "system.running.processes"
	blockProcesses    string = "system.block.processes"
)

func GetSystemInformationMetrics(profile map[string]interface{}, channel chan map[string]interface{}) {
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

	output, err := session.Output(systemInfoCommand)

	if err != nil {
		return
	}

	outputInfo := strings.Split(string(output), util.NewLine)

	upTiming := outputInfo[1]

	operatingSystemOutput := outputInfo[2]

	outputInfo = strings.Split(strings.TrimSpace(outputInfo[0]), util.SpaceSeparator)

	operatingSystemInfo := strings.Split(strings.TrimSpace(operatingSystemOutput), util.Colon)

	response[osVersion] = operatingSystemInfo[1]

	response[systemName] = outputInfo[0]

	response[osName] = outputInfo[1]

	response[upTime] = upTiming

	response[threads] = outputInfo[2]

	response[contextSwitches] = outputInfo[3]

	response[runningProcesses] = outputInfo[4]

	response[blockProcesses] = outputInfo[5]
}
