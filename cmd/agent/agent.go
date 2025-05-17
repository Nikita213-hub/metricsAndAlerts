package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

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

type Flags struct {
	CfgFileName    *string
	PollInterval   *int
	ReportInterval *int
	RateLimit      *int
	PrivateKey     *string
}

func (f *Flags) Parse() {
	pollInterval := flag.Int("poll", -1, "Poll interval")
	reportInterval := flag.Int("report", -1, "Report interval")
	rateLimit := flag.Int("rate", -1, "Rate Limit")
	privateKey := flag.String("pk", "", "Private key value")
	fName := flag.String("cfg", "config.yaml", "Config file name")
	flag.Parse()
	f.CfgFileName = fName
	f.PollInterval = pollInterval
	f.RateLimit = rateLimit
	f.ReportInterval = reportInterval
	f.PrivateKey = privateKey
}

func (a *Agent) Configure() error {
	flags := &Flags{}
	flags.Parse()
	a.Cfg = newAgentConfig()
	return a.Cfg.RunConfiguration(flags)
}

func (a *Agent) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		mwp := NewMetricsWP(a.Cfg.Agent.RateLimit)
		td := time.Duration(a.Cfg.Agent.PollInterval) * time.Second
		collectMetrics(mwp.jobs, ctx, &td)
		mwp.Run(a.Cfg.Host.PrivateKey, a.Cfg.Host.Address, ctx)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
	// metrics generator call
	// invoke worker pool with RATE_LIMIT workers
}

type AgentConfig struct {
	Agent struct {
		PollInterval   int `yaml:"poll_interval" env:"POLL_INTERVAL"`
		ReportInterval int `yaml:"report_interval" env:"REPORT_INTERVAL"`
		RateLimit      int `yaml:"rate_limit" env:"RATE_LIMIT"`
	} `yaml:"agent"`
	Host struct {
		Address    string `yaml:"address" env:"HOST_ADDRESS"`
		PrivateKey string `env:"PRIVATE_KEY"`
	} `yaml:"host"`
}

func newAgentConfig() *AgentConfig {
	return &AgentConfig{}
}

func (ac *AgentConfig) RunConfiguration(flags *Flags) error {
	err := getYamlCfg(ac, *flags.CfgFileName)
	if err != nil {
		return err
	}
	err = getEnvCfg(ac)
	if err != nil {
		return err
	}

	if *flags.PollInterval != -1 {
		ac.Agent.PollInterval = *flags.PollInterval
	}
	if *flags.ReportInterval != -1 {
		ac.Agent.ReportInterval = *flags.ReportInterval
	}
	if *flags.RateLimit != -1 {
		ac.Agent.RateLimit = *flags.RateLimit
	}
	if *flags.PrivateKey != "" {
		ac.Host.PrivateKey = *flags.PrivateKey
	}
	return nil
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
