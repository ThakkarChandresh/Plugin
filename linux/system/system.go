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

	systemInfoMetrics := strings.Split(strings.TrimSpace(string(output)), util.NewLine)

	response[upTime] = systemInfoMetrics[1]

	response[osVersion] = strings.TrimSpace(strings.Split(systemInfoMetrics[2], util.Colon)[1])

	systemInfoMetrics = strings.Split(strings.TrimSpace(systemInfoMetrics[0]), util.SpaceSeparator)

	response[systemName] = systemInfoMetrics[0]

	response[osName] = systemInfoMetrics[1]

	response[threads] = systemInfoMetrics[2]

	response[contextSwitches] = systemInfoMetrics[3]

	response[runningProcesses] = systemInfoMetrics[4]

	response[blockProcesses] = systemInfoMetrics[5]
}
