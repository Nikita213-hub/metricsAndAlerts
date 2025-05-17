package main

import (
	"fmt"
)

func main() {
	agent := NewAgent()
	err := agent.Configure()
	if err != nil {
		panic(err)
	}
	err = agent.Run()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", *agent.Cfg)
}
