package cpu

import (
	"Plugin/linux/util"
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
	SystemCPUError          string = "system.cpu.error"
)

func GetCpuMetrics(profile map[string]interface{}, channel chan map[string]interface{}) {
	response := make(map[string]interface{})

	defer func() {
		channel <- response
	}()

	defer func() {
		if r := recover(); r != nil {
			response[SystemCPUError] = fmt.Sprint(r)
		}
	}()

	connection, err := util.GetConnection(profile)

	if err != nil {
		response[SystemCPUError] = fmt.Sprint(err)

		return
	}

	defer func(connection *ssh.Client) {
		if err = connection.Close(); err != nil {

			response[SystemCPUError] = fmt.Sprint(err)
		}
	}(connection)

	session, err := connection.NewSession()

	//Session will automatically close

	if err != nil {
		response[SystemCPUError] = fmt.Sprint(err)

		return
	}

	output, err := session.Output(SystemCPUMetricsCommand)

	if err != nil {

		response[SystemCPUError] = fmt.Sprint(err)

		return
	}

	allCPUMetrics := strings.Split(strings.TrimSpace(string(output)), util.NewLine)

	response = make(map[string]interface{})

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
}
