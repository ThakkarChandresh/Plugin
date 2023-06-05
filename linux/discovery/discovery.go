package discovery

import (
	"Plugin/linux/util"
	"errors"
	"fmt"
	"strings"
)

const (
	Status   string = "status"
	Hostname string = "hostname"
)

func Discover(profile map[string]interface{}) (response map[string]interface{}, err error) {

	response = make(map[string]interface{})

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()

	output, err := util.ExecuteCommand(profile, Hostname)

	if err != nil {
		return
	}

	result := strings.TrimSpace(string(output))

	response[Status] = true

	response[Hostname] = result

	return
}
