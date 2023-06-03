package main

import (
	"Plugin/linux"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	discovery   string = "Discovery"
	polling     string = "Polling"
	requestType string = "type"
)

func main() {
	var request = make(map[string]interface{})

	defer func() {

		jsonStr, _ := json.Marshal(request)

		fmt.Println(fmt.Sprintf("%v", string(jsonStr)))
	}()

	defer func() {

		if r := recover(); r != nil {

			request["err"] = fmt.Sprintf("%v", r)
		}
	}()

	input := os.Args[1]

	err := json.Unmarshal([]byte(input), &request)

	if err != nil {

		request["err"] = fmt.Sprintf("%v", err)

		return
	}

	var response = make(map[string]interface{})

	if strings.EqualFold(request[requestType].(string), discovery) {

		response, err = linux.Discover(request)

	} else if strings.EqualFold(request[requestType].(string), polling) {

		response, err = linux.Collect(request)
	}

	if err != nil {

		request["err"] = fmt.Sprintf("%v", err)

		return
	}

	if _, err = json.Marshal(response); err != nil {

		request["err"] = fmt.Sprintf("%v", err)

		return
	}

	request["result"] = response
}
