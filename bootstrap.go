package main

import (
	"Plugin/linux"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	input := os.Args[1]

	var request = make(map[string]interface{})

	err := json.Unmarshal([]byte(input), &request)

	if err != nil {
		fmt.Println(err)
		return
	}

	/*response, err := linux.Discover(request)

	if err != nil {
		fmt.Println(err)
		return
	}

	jsonStr, err := json.Marshal(response)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("%v", string(jsonStr)))*/

	response, err := linux.Collect(request)

	if err != nil {
		fmt.Println(err)
		return
	}

	jsonStr, err := json.Marshal(response)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("%v", string(jsonStr)))
}
