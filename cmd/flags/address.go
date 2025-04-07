package flags

import (
	"errors"
	"strconv"
	"strings"
)

type Address struct {
	host string
	port string
}

func NewAddress(defaultHost string, defaultPort string) *Address {
	return &Address{
		host: defaultHost,
		port: defaultPort,
	}
}

func (a *Address) String() string {
	return a.host + ":" + a.port
}

func (a *Address) Set(s string) error {
	var splt []string
	if splt = strings.Split(s, ":"); len(splt) != 2 {
		return errors.New("input host and port in host:port format")
	}
	_, err := strconv.Atoi(splt[1])
	if err != nil {
		return errors.New("port is incorrect")
	}
	a.host = splt[0]
	a.port = splt[1]
	return nil
}

func (a *Address) GetHost() string {
	return a.host
}

func (a *Address) GetPort() string {
	return a.port
}
