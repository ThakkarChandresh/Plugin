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

	outputInfo := strings.Split(string(output), util.NewLine)

	cores, err := strconv.Atoi(outputInfo[0])

	if err != nil {
		return
	}

	response = make(map[string]interface{})

	allCpuInfo := strings.Split(outputInfo[1], util.SpaceSeparator)

	response[cpuCores] = cores

	response[cpuPercentage] = allCpuInfo[1]

	response[cpuUserPercentage] = allCpuInfo[2]

	response[cpuIdlePercentage] = allCpuInfo[3]

	outputInfo = outputInfo[2 : len(outputInfo)-1]

	result := make([]map[string]interface{}, cores)

	for i := 0; i < len(outputInfo); i++ {
		oneCpuOutput := strings.Split(outputInfo[i], util.SpaceSeparator)

		oneCpuInfo := make(map[string]any)

		oneCpuInfo[cpuCore] = oneCpuOutput[0]

		oneCpuInfo[cpuPercentage] = oneCpuOutput[1]

		oneCpuInfo[cpuUserPercentage] = oneCpuOutput[2]

		oneCpuInfo[cpuIdlePercentage] = oneCpuOutput[3]

		result[i] = oneCpuInfo
	}

	response[SystemCPU] = result
}
