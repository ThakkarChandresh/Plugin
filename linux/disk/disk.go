package disk

import (
	"Plugin/linux/util"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

const (
	diskInfoCommand      string = `iostat -dx | awk 'NR>3 {print $1 " " $2 " " $8 " " $3 " " $9}'`
	systemDisk           string = "system.disk"
	systemDiskBytes      string = "system.disk.bytes.per.sec"
	systemDiskWriteBytes string = "system.disk.write.bytes.per.sec"
	systemDiskReadBytes  string = "system.disk.read.bytes.per.sec"
	systemDiskWriteOps   string = "system.disk.write.ops.per.sec"
	systemDiskReadOps    string = "system.disk.read.ops.per.sec"
)

func GetDiskMetrics(profile map[string]interface{}, channel chan map[string]interface{}) {
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

	output, err := session.Output(diskInfoCommand)

	if err != nil {
		return
	}

	allDiskMetrics := strings.Split(strings.TrimSpace(string(output)), util.NewLine)

	result := make([]map[string]interface{}, len(allDiskMetrics))

	for i := 0; i < len(allDiskMetrics); i++ {
		totalBytes := 0.0

		disk := strings.Split(allDiskMetrics[i], util.SpaceSeparator)

		diskMetrics := make(map[string]interface{})

		diskMetrics[systemDisk] = disk[0]

		if readOps, err := strconv.ParseFloat(disk[1], 64); err == nil {

			diskMetrics[systemDiskReadOps] = readOps
		}

		if writeOps, err := strconv.ParseFloat(disk[2], 64); err == nil {

			diskMetrics[systemDiskWriteOps] = writeOps
		}

		if readBytes, err := strconv.ParseFloat(disk[3], 64); err == nil {

			readBytes *= 1024

			diskMetrics[systemDiskReadBytes] = readBytes

			totalBytes += readBytes
		}

		if writeBytes, err := strconv.ParseFloat(disk[4], 64); err == nil {

			writeBytes *= 1024

			diskMetrics[systemDiskWriteBytes] = writeBytes

			totalBytes += writeBytes
		}

		diskMetrics[systemDiskBytes] = totalBytes

		result[i] = diskMetrics
	}

	response[systemDisk] = result
}
