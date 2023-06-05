package memory

import (
	"Plugin/linux/util"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	SystemMemoryMetricsCommand string = `free -b | awk 'NR>1 {print $2" " $3" " ((($2 - $7) * 100) / $2) " " $4 " " (($4 * 100) / $2) " " $7}'| head -n 1|tr '\n' " " && free -b | awk 'NR>2 {print $2}'`
	SystemMemoryInstalledBytes string = "system.memory.installed.bytes"
	SystemMemoryUsedBytes      string = "system.memory.used.bytes"
	SystemMemoryUsedPercentage string = "system.memory.used.percentage"
	SystemMemoryFreeBytes      string = "system.memory.free.bytes"
	SystemMemoryFreePercentage string = "system.memory.free.percentage"
	SystemMemoryAvailableBytes string = "system.memory.available.bytes"
	SystemMemorySwapBytes      string = "system.memory.swap.bytes"
)

func Collect(profile map[string]interface{}) (response map[string]interface{}, err error) {

	response = make(map[string]interface{})

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()

	output, err := util.ExecuteCommand(profile, SystemMemoryMetricsCommand)

	if err != nil {
		return
	}

	memoryMetrics := strings.Split(strings.TrimSpace(string(output)), util.SpaceSeparator)

	if installedMemoryBytes, err := strconv.Atoi(memoryMetrics[0]); err == nil {

		response[SystemMemoryInstalledBytes] = installedMemoryBytes
	}

	if usedMemoryBytes, err := strconv.Atoi(memoryMetrics[1]); err == nil {

		response[SystemMemoryUsedBytes] = usedMemoryBytes
	}

	if usedMemoryPercentage, err := strconv.ParseFloat(memoryMetrics[2], 64); err == nil {

		response[SystemMemoryUsedPercentage] = usedMemoryPercentage
	}

	if freeMemoryBytes, err := strconv.Atoi(memoryMetrics[3]); err == nil {

		response[SystemMemoryFreeBytes] = freeMemoryBytes
	}

	if freeMemoryPercentagee, err := strconv.ParseFloat(memoryMetrics[4], 64); err == nil {

		response[SystemMemoryFreePercentage] = freeMemoryPercentagee
	}

	if availableMemoryBytes, err := strconv.Atoi(memoryMetrics[5]); err == nil {

		response[SystemMemoryAvailableBytes] = availableMemoryBytes
	}

	if swapMemoryBytes, err := strconv.Atoi(memoryMetrics[6]); err == nil {

		response[SystemMemorySwapBytes] = swapMemoryBytes
	}
	return
}
