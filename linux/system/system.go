package system

import (
	"Plugin/linux/util"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

const (
	SystemInfoMetricsCommand string = `hostname |tr '\n' " " && uname |tr '\n' " " && ps -eo nlwp | awk '{ num_threads += $1 } END { print num_threads }' | tr '\n' " " && vmstat | tail -n 1 | awk '{print $12}' | tr '\n' " " && ps axo state | grep "R" | wc -l | tr '\n' " " && ps axo stat | grep "D" | wc -l && uptime -p | awk 'gsub("up ","")' && hostnamectl | grep "Operating System"`
	SystemOSVersion          string = "system.os.version"
	SystemName               string = "system.name"
	SystemOSName             string = "system.os.name"
	SystemUpTime             string = "system.uptime"
	SystemThreads            string = "system.threads"
	SystemContextSwitches    string = "system.context.switches"
	SystemRunningProcesses   string = "system.running.processes"
	SystemBlockProcesse      string = "system.block.processes"
	SystemInfoError          string = "system.info.error"
)

func GetSystemInformationMetrics(profile map[string]interface{}, channel chan map[string]interface{}) {
	response := make(map[string]interface{})

	defer func() {
		channel <- response
	}()

	defer func() {
		if r := recover(); r != nil {
			response[SystemInfoError] = fmt.Sprint(r)
		}
	}()

	connection, err := util.GetConnection(profile)

	if err != nil {
		response[SystemInfoError] = fmt.Sprint(err)

		return
	}

	defer func(connection *ssh.Client) {
		if err = connection.Close(); err != nil {

			response[SystemInfoError] = fmt.Sprint(err)
		}
	}(connection)

	session, err := connection.NewSession()

	//Session will automatically close

	if err != nil {
		response[SystemInfoError] = fmt.Sprint(err)

		return
	}

	output, err := session.Output(SystemInfoMetricsCommand)

	if err != nil {
		response[SystemInfoError] = fmt.Sprint(err)

		return
	}

	systemInfoMetrics := strings.Split(strings.TrimSpace(string(output)), util.NewLine)

	response[SystemUpTime] = systemInfoMetrics[1]

	response[SystemOSVersion] = strings.TrimSpace(strings.Split(systemInfoMetrics[2], util.Colon)[1])

	systemInfoMetrics = strings.Split(strings.TrimSpace(systemInfoMetrics[0]), util.SpaceSeparator)

	response[SystemName] = systemInfoMetrics[0]

	response[SystemOSName] = systemInfoMetrics[1]

	if threads, err := strconv.Atoi(systemInfoMetrics[2]); err == nil {

		response[SystemThreads] = threads
	}

	if contextSwitches, err := strconv.Atoi(systemInfoMetrics[3]); err == nil {

		response[SystemContextSwitches] = contextSwitches
	}

	if runningProcesses, err := strconv.Atoi(systemInfoMetrics[4]); err == nil {

		response[SystemRunningProcesses] = runningProcesses
	}

	if blockProcesses, err := strconv.Atoi(systemInfoMetrics[5]); err == nil {

		response[SystemBlockProcesse] = blockProcesses
	}

}
