package cpu

import (
	"Plugin/linux/util"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

const (
	cpuInfoCommand    string = `nproc --all && mpstat -P ALL | awk 'NR>3 {print $4 " " $7 " " $5 " " $14}'`
	SystemCPU         string = "system.cpu"
	cpuCores          string = "system.cpu.cores"
	cpuCore           string = "system.cpu.core"
	cpuPercentage     string = "system.cpu.percentage"
	cpuUserPercentage string = "system.cpu.user.percentage"
	cpuIdlePercentage string = "system.cpu.idle.percentage"
)

func GetCpuMetrics(profile map[string]interface{}, channel chan map[string]interface{}) {
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

	output, err := session.Output(cpuInfoCommand)

	if err != nil {
		return
	}

	allCPUMetrics := strings.Split(strings.TrimSpace(string(output)), util.NewLine)

	response = make(map[string]interface{})

	avgCPUMetrics := strings.Split(allCPUMetrics[1], util.SpaceSeparator)

	if cores, err := strconv.Atoi(allCPUMetrics[0]); err == nil {

		response[cpuCores] = cores
	}

	if percentage, err := strconv.ParseFloat(avgCPUMetrics[1], 64); err == nil {

		response[cpuPercentage] = percentage
	}

	if userPercentage, err := strconv.ParseFloat(avgCPUMetrics[2], 64); err == nil {

		response[cpuUserPercentage] = userPercentage
	}

	if idlePercentage, err := strconv.ParseFloat(avgCPUMetrics[3], 64); err == nil {

		response[cpuIdlePercentage] = idlePercentage
	}

	allCPUMetrics = allCPUMetrics[2:]

	result := make([]map[string]interface{}, len(allCPUMetrics))

	for i := 0; i < len(allCPUMetrics); i++ {
		cpu := strings.Split(allCPUMetrics[i], util.SpaceSeparator)

		cpuMetrics := make(map[string]any)

		if core, err := strconv.Atoi(cpu[0]); err == nil {

			cpuMetrics[cpuCore] = core
		}

		if percentage, err := strconv.ParseFloat(cpu[1], 64); err == nil {

			cpuMetrics[cpuPercentage] = percentage
		}

		if userPercentage, err := strconv.ParseFloat(cpu[2], 64); err == nil {

			cpuMetrics[cpuUserPercentage] = userPercentage
		}

		if idlePercentage, err := strconv.ParseFloat(cpu[3], 64); err == nil {

			cpuMetrics[cpuIdlePercentage] = idlePercentage
		}

		result[i] = cpuMetrics
	}

	response[SystemCPU] = result
}
