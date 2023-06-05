package main

import (
	"Plugin/linux/cpu"
	"Plugin/linux/discovery"
	"Plugin/linux/disk"
	"Plugin/linux/memory"
	"Plugin/linux/process"
	"Plugin/linux/system"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	Error              string = "error"
	DeviceType         string = "device_type"
	Linux              string = "linux"
	RequestType        string = "request_type"
	InvalidRequestJson string = "invalid request json"
	Discovery          string = "discovery"
	Polling            string = "polling"
	MetricGroup        string = "metric_group"
	SystemCPU          string = "system.cpu"
	SystemDisk         string = "system.disk"
	SystemMemory       string = "system.memory"
	SystemProcess      string = "system.process"
	SystemInfo         string = "system.info"
	Result             string = "result"
)

func main() {
	var request = make(map[string]interface{})

	defer func() {

		jsonStr, _ := json.Marshal(request)

		fmt.Println(fmt.Sprint(string(jsonStr)))
	}()

	defer func() {

		if r := recover(); r != nil {

			request[Error] = fmt.Sprint(r)
		}
	}()

	input := os.Args[1]

	err := json.Unmarshal([]byte(input), &request)

	if err != nil {

		request[Error] = fmt.Sprint(err)

		return
	}

	var response = make(map[string]interface{})

	switch {

	case strings.EqualFold(fmt.Sprint(request[DeviceType]), Linux) && strings.EqualFold(fmt.Sprint(request[RequestType]), Discovery):

		response, err = discovery.Discover(request)

	case strings.EqualFold(fmt.Sprint(request[DeviceType]), Linux) && strings.EqualFold(fmt.Sprint(request[RequestType]), Polling) && strings.EqualFold(fmt.Sprint(request[MetricGroup]), SystemCPU):

		response, err = cpu.Collect(request)

	case strings.EqualFold(fmt.Sprint(request[DeviceType]), Linux) && strings.EqualFold(fmt.Sprint(request[RequestType]), Polling) && strings.EqualFold(fmt.Sprint(request[MetricGroup]), SystemDisk):

		response, err = disk.Collect(request)

	case strings.EqualFold(fmt.Sprint(request[DeviceType]), Linux) && strings.EqualFold(fmt.Sprint(request[RequestType]), Polling) && strings.EqualFold(fmt.Sprint(request[MetricGroup]), SystemMemory):

		response, err = memory.Collect(request)

	case strings.EqualFold(fmt.Sprint(request[DeviceType]), Linux) && strings.EqualFold(fmt.Sprint(request[RequestType]), Polling) && strings.EqualFold(fmt.Sprint(request[MetricGroup]), SystemProcess):

		response, err = process.Collect(request)

	case strings.EqualFold(fmt.Sprint(request[DeviceType]), Linux) && strings.EqualFold(fmt.Sprint(request[RequestType]), Polling) && strings.EqualFold(fmt.Sprint(request[MetricGroup]), SystemInfo):

		response, err = system.Collect(request)

	default:

		err = errors.New(InvalidRequestJson)
	}

	if err != nil {

		request[Error] = fmt.Sprint(err)

		return
	}

	if _, err = json.Marshal(response); err != nil {

		request[Error] = fmt.Sprint(err)

		return
	}

	request[Result] = response
}
