package main

import (
	"fmt"
	"os"
)

var appName = "hermes" // name after go install ...

func main() {
	name, ok := os.LookupEnv("APP_NAME")
	if !ok {
		name = appName
	}
	fmt.Printf("name: %v\n", name)
}
