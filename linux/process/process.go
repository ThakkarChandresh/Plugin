package process

import (
	"Plugin/linux/util"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	SystemProcessMetricsCommand string = `ps aux | awk 'NR> 1 {print $2 " " $3 " " $4 " " $1" "$11}'`
	SystemProcess               string = "system.process"
	SystemProcessPID            string = "system.process.pid"
	SystemProcessCPU            string = "system.process.cpu"
	SystemProcessMemory         string = "system.process.memory"
	SystemProcessUser           string = "system.process.user"
	SystemProcessCommand        string = "system.process.command"
)

func Collect(profile map[string]interface{}) (response map[string]interface{}, err error) {

	response = make(map[string]interface{})

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()

	output, err := util.ExecuteCommand(profile, SystemProcessMetricsCommand)

	if err != nil {
		return
	}

	allProcessMetrics := strings.Split(strings.TrimSpace(string(output)), util.NewLine)

	result := make([]map[string]interface{}, len(allProcessMetrics))

	for i := 0; i < len(allProcessMetrics); i++ {
		process := strings.Split(allProcessMetrics[i], util.SpaceSeparator)

		processMetrics := make(map[string]interface{})

		if processPID, err := strconv.Atoi(process[0]); err == nil {

			processMetrics[SystemProcessPID] = processPID
		}

		if processCPU, err := strconv.ParseFloat(process[1], 64); err == nil {

			processMetrics[SystemProcessCPU] = processCPU
		}

		if processMemory, err := strconv.ParseFloat(process[2], 64); err == nil {

			processMetrics[SystemProcessMemory] = processMemory
		}

		processMetrics[SystemProcessUser] = process[3]

		processMetrics[SystemProcessCommand] = process[4]

		result[i] = processMetrics
	}

	response[SystemProcess] = result

	return
}
