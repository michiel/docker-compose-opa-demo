package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/open-policy-agent/opa/plugins"
	"github.com/open-policy-agent/opa/plugins/logs"
	"github.com/open-policy-agent/opa/runtime"
	"github.com/open-policy-agent/opa/util"
)

type Config struct {
	Stderr         bool   `json:"stderr"`
	SplunkHECURI   string `json:"splunk_hec_uri"`
	SplunkHECToken string `json:"splunk_hec_token"`
}

type Factory struct{}

func (Factory) New(_ *plugins.Manager, config interface{}) plugins.Plugin {
	return &PrintlnLogger{
		config: config.(Config),
	}
}

func (Factory) Validate(_ *plugins.Manager, config []byte) (interface{}, error) {
	parsedConfig := Config{}
	return parsedConfig, util.Unmarshal(config, &parsedConfig)
}

type PrintlnLogger struct {
	config Config
	mtx    sync.Mutex
}

func (p *PrintlnLogger) Start(ctx context.Context) error {
	return nil
}

func (p *PrintlnLogger) Stop(ctx context.Context) {
}

func (p *PrintlnLogger) Reconfigure(ctx context.Context, config interface{}) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.config = config.(Config)
}

func (p *PrintlnLogger) Log(ctx context.Context, event logs.EventV1) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	w := os.Stdout
	if p.config.Stderr {
		w = os.Stderr
	}
	fmt.Fprintln(w, event) // ignoring errors!

	/* hackity hack */

	url := p.config.SplunkHECURI
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return nil
}

func Init() error {
	runtime.RegisterPlugin("println_decision_logger", Factory{})
	return nil
}
