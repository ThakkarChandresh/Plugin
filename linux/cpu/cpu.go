package cpu

import (
	"Plugin/linux/util"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

const (
	SystemCPUMetricsCommand string = `nproc --all && mpstat -P ALL | awk 'NR>3 {print $4 " " $7 " " $5 " " $14}'`
	SystemCPU               string = "system.cpu"
	SystemCPUCores          string = "system.cpu.cores"
	SystemCPUCore           string = "system.cpu.core"
	SystemCPUPercentage     string = "system.cpu.percentage"
	SystemCPUUserPercentage string = "system.cpu.user.percentage"
	SystemCPUIdlePercentage string = "system.cpu.idle.percentage"
)

func Collect(profile map[string]interface{}) (response map[string]interface{}, err error) {

	response = make(map[string]interface{})

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()

	connection, err := util.GetConnection(profile)

	if err != nil {
		return
	}

	defer func(connection *ssh.Client) {
		if closeErr := connection.Close(); closeErr != nil {
			err = closeErr
		}
	}(connection)

	session, err := connection.NewSession()

	//Session will automatically close

	if err != nil {
		return
	}

	output, err := session.Output(SystemCPUMetricsCommand)

	if err != nil {
		return
	}

	allCPUMetrics := strings.Split(strings.TrimSpace(string(output)), util.NewLine)

	avgCPUMetrics := strings.Split(allCPUMetrics[1], util.SpaceSeparator)

	if cpuCores, err := strconv.Atoi(allCPUMetrics[0]); err == nil {

		response[SystemCPUCores] = cpuCores
	}

	if cpuPercentage, err := strconv.ParseFloat(avgCPUMetrics[1], 64); err == nil {

		response[SystemCPUPercentage] = cpuPercentage
	}

	if cpuUserPercentage, err := strconv.ParseFloat(avgCPUMetrics[2], 64); err == nil {

		response[SystemCPUUserPercentage] = cpuUserPercentage
	}

	if cpuIdlePercentage, err := strconv.ParseFloat(avgCPUMetrics[3], 64); err == nil {

		response[SystemCPUIdlePercentage] = cpuIdlePercentage
	}

	allCPUMetrics = allCPUMetrics[2:]

	result := make([]map[string]interface{}, len(allCPUMetrics))

	for i := 0; i < len(allCPUMetrics); i++ {
		cpu := strings.Split(allCPUMetrics[i], util.SpaceSeparator)

		cpuMetrics := make(map[string]any)

		if cpuCore, err := strconv.Atoi(cpu[0]); err == nil {

			cpuMetrics[SystemCPUCore] = cpuCore
		}

		if cpuPercentage, err := strconv.ParseFloat(cpu[1], 64); err == nil {

			cpuMetrics[SystemCPUPercentage] = cpuPercentage
		}

		if cpuUserPercentage, err := strconv.ParseFloat(cpu[2], 64); err == nil {

			cpuMetrics[SystemCPUUserPercentage] = cpuUserPercentage
		}

		if cpuIdlePercentage, err := strconv.ParseFloat(cpu[3], 64); err == nil {

			cpuMetrics[SystemCPUIdlePercentage] = cpuIdlePercentage
		}

		result[i] = cpuMetrics
	}

	response[SystemCPU] = result

	return
}
