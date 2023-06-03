package disk

import (
	"Plugin/linux/util"
	"fmt"
	"golang.org/x/crypto/ssh"
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
	SystemDiskError            string = "system.disk.error"
)

func GetDiskMetrics(profile map[string]interface{}, channel chan map[string]interface{}) {
	response := make(map[string]interface{})

	defer func() {
		channel <- response
	}()

	defer func() {
		if r := recover(); r != nil {
			response[SystemDiskError] = fmt.Sprint(r)
		}
	}()

	connection, err := util.GetConnection(profile)

	if err != nil {
		response[SystemDiskError] = fmt.Sprint(err)

		return
	}

	defer func(connection *ssh.Client) {
		if err = connection.Close(); err != nil {

			response[SystemDiskError] = fmt.Sprint(err)
		}
	}(connection)

	session, err := connection.NewSession()

	//Session will automatically close

	if err != nil {
		response[SystemDiskError] = fmt.Sprint(err)

		return
	}

	output, err := session.Output(SystemDiskMetricsCommand)

	if err != nil {
		response[SystemDiskError] = fmt.Sprint(err)

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
}
