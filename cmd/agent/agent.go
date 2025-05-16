package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v3"
)

type Agent struct {
	Cfg    *AgentConfig
	Logger *slog.Logger
}

func NewAgent() *Agent {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	return &Agent{
		Logger: logger,
		Cfg:    nil,
	}
}

func (a *Agent) Configure() {
	a.Cfg = newAgentConfig()
	// add get file name logic
	a.Cfg.RunConfiguration("config.yaml")
}

type AgentConfig struct {
	Agent struct {
		PollInterval   int `yaml:"poll_interval" env:"POLL_INTERVAL"`
		ReportInterval int `yaml:"report_interval" env:"REPORT_INTERVAL"`
		RateLimit      int `yaml:"rate_limit" env:"RATE_LIMIT"`
	} `yaml:"agent"`
	Host struct {
		Address string `yaml:"address" env:"HOST_ADDRESS"`
	} `yaml:"host"`
}

func newAgentConfig() *AgentConfig {
	return &AgentConfig{}
}

func (ac *AgentConfig) RunConfiguration(cfgFileName string) {
	err := getYamlCfg(ac, cfgFileName)
	if err != nil {
		fmt.Println(err)
	}
	err = getEnvCfg(ac)
	if err != nil {
		fmt.Println(err)
	}
	pollInterval := flag.Int("poll", -1, "Poll interval")
	reportInterval := flag.Int("report", -1, "Report interval")
	rateLimit := flag.Int("rate", -1, "Rate Limit")
	flag.Parse()
	if *pollInterval != -1 {
		ac.Agent.PollInterval = *pollInterval
	}
	if *reportInterval != -1 {
		ac.Agent.ReportInterval = *reportInterval
	}
	if *rateLimit != -1 {
		ac.Agent.RateLimit = *rateLimit
	}

}

func getYamlCfg(src *AgentConfig, cfgFileName string) error {
	data, err := os.ReadFile(cfgFileName)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &src)
	if err != nil {
		return err
	}
	return nil
}

func getEnvCfg(src *AgentConfig) error {
	return env.Parse(src)
}
