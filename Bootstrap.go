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

	requestString := os.Args[1]

	var request = make(map[string]any)

	err := json.Unmarshal([]byte(requestString), &request)

	if err != nil {
		panic(err.Error())
	}

	/*response := linux.Discover(request)

	jsonStr, err := json.Marshal(response)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(fmt.Sprintf("%v", string(jsonStr)))*/

	response := linux.Collect(request)

	jsonStr, err := json.Marshal(response)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(fmt.Sprintf("%v", string(jsonStr)))
}
