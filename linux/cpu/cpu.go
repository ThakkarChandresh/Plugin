package cpu

import (
	"Plugin/linux/util"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

const (
	cpuInfoCommand string = `nproc --all && mpstat -P ALL | awk 'NR>3 {print $4 " " $7 " " $5 " " $14}'`
)

func GetCpuMetrics(connection *ssh.Client, channel chan map[string]interface{}) {
	response := make(map[string]interface{})

	defer func() {
		channel <- response
	}()

	defer func() {
		if r := recover(); r != nil {
			//response[util.SystemCPU] = map[string]interface{}{util.Status: util.Fail, util.Err: r}
		}
	}()

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
		response[util.SystemCPU] = map[string]interface{}{util.Status: util.Fail, util.Err: err.Error()}
		return
	}

	response = make(map[string]interface{})

	allCpuInfo := strings.Split(outputInfo[1], util.SpaceSeparator)

	response[util.CPUCores] = cores

	response[util.CPUPercentage] = allCpuInfo[1]

	response[util.CPUUserPercentage] = allCpuInfo[2]

	response[util.CPUIdlePercentage] = allCpuInfo[3]

	outputInfo = outputInfo[2 : len(outputInfo)-1]

	individualCpuInfo := make([]map[string]interface{}, cores)

	for i := 0; i < len(outputInfo); i++ {
		oneCpuOutput := strings.Split(outputInfo[i], util.SpaceSeparator)

		oneCpuInfo := make(map[string]any)

		oneCpuInfo[util.CPUCore] = oneCpuOutput[0]

		oneCpuInfo[util.CPUPercentage] = oneCpuOutput[1]

		oneCpuInfo[util.CPUUserPercentage] = oneCpuOutput[2]

		oneCpuInfo[util.CPUIdlePercentage] = oneCpuOutput[3]

		individualCpuInfo[i] = oneCpuInfo
	}

	response[util.SystemCPU] = individualCpuInfo
}
