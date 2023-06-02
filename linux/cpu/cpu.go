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

	allCPUMetrics := strings.Split(strings.Trim(string(output), util.NewLine), util.NewLine)

	cores, err := strconv.Atoi(allCPUMetrics[0])

	if err != nil {
		return
	}

	response = make(map[string]interface{})

	avgCPUMetrics := strings.Split(allCPUMetrics[1], util.SpaceSeparator)

	response[cpuCores] = cores

	response[cpuPercentage] = avgCPUMetrics[1]

	response[cpuUserPercentage] = avgCPUMetrics[2]

	response[cpuIdlePercentage] = avgCPUMetrics[3]

	allCPUMetrics = allCPUMetrics[2:]

	result := make([]map[string]interface{}, cores)

	for i := 0; i < len(allCPUMetrics); i++ {
		cpu := strings.Split(allCPUMetrics[i], util.SpaceSeparator)

		cpuMetrics := make(map[string]any)

		cpuMetrics[cpuCore] = cpu[0]

		cpuMetrics[cpuPercentage] = cpu[1]

		cpuMetrics[cpuUserPercentage] = cpu[2]

		cpuMetrics[cpuIdlePercentage] = cpu[3]

		result[i] = cpuMetrics
	}

	response[SystemCPU] = result
}
