package disk

import (
	"Plugin/linux/util"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	SystemDiskMetricsCommand   string = `iostat -dx | awk 'NR>3 {print $1 " " $2 " " $8 " " $3 " " $9}'`
	SystemDisk                 string = "system.disk"
	SystemDiskBytesPerSec      string = "system.disk.bytes.per.sec"
	SystemDiskWriteBytesPerSec string = "system.disk.write.bytes.per.sec"
	SystemDiskReadBytesPerSec  string = "system.disk.read.bytes.per.sec"
	SystemDiskWriteOpsPerSec   string = "system.disk.write.ops.per.sec"
	SystemDiskReadOpsPerSec    string = "system.disk.read.ops.per.sec"
)

func Collect(profile map[string]interface{}) (response map[string]interface{}, err error) {

	response = make(map[string]interface{})

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()

	output, err := util.ExecuteCommand(profile, SystemDiskMetricsCommand)

	if err != nil {
		return
	}

	allDiskMetrics := strings.Split(strings.TrimSpace(string(output)), util.NewLine)

	result := make([]map[string]interface{}, len(allDiskMetrics))

	for i := 0; i < len(allDiskMetrics); i++ {
		totalBytes := 0.0

		disk := strings.Split(allDiskMetrics[i], util.SpaceSeparator)

		diskMetrics := make(map[string]interface{})

		diskMetrics[SystemDisk] = disk[0]

		if readOps, err := strconv.ParseFloat(disk[1], 64); err == nil {

			diskMetrics[SystemDiskReadOpsPerSec] = readOps
		}

		if writeOps, err := strconv.ParseFloat(disk[2], 64); err == nil {

			diskMetrics[SystemDiskWriteOpsPerSec] = writeOps
		}

		if readBytes, err := strconv.ParseFloat(disk[3], 64); err == nil {

			readBytes *= 1024

			diskMetrics[SystemDiskReadBytesPerSec] = readBytes

			totalBytes += readBytes
		}

		if writeBytes, err := strconv.ParseFloat(disk[4], 64); err == nil {

			writeBytes *= 1024

			diskMetrics[SystemDiskWriteBytesPerSec] = writeBytes

			totalBytes += writeBytes
		}

		diskMetrics[SystemDiskBytesPerSec] = totalBytes

		result[i] = diskMetrics
	}

	response[SystemDisk] = result

	return
}
