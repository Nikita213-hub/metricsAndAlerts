package main

import (
	"fmt"
)

func main() {
	agent := NewAgent()
	agent.Configure()
	fmt.Printf("%v\n", *agent.Cfg)
}
