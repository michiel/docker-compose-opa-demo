package main

import (
	"bytes"
	"context"
	"encoding/json"
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

type HECPayload struct {
	Host       string   `json: host`
	Source     string   `json: source`
	SourceType string   `json: sourcetype`
	Index      string   `json: index`
	Event      int      `json: event`
	Fields     []string `json: temp_forecast`
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

	/*
		b, err := json.Marshal(M{"query": M{"query_string": M{"query": "query goes here"}}})
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("   As Map:", string(b))
	*/
	/* hackity hack */

	url := p.config.SplunkHECURI
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	hecPayload := HECPayload{
		Host:       "host",
		Source:     "SomeSource",
		SourceType: "json",
		Index:      "main",
	}

	hecPayloadJson, err := json.Marshal(hecPayload)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}

	req1, err1 := http.NewRequest("POST", url, bytes.NewBuffer(hecPayloadJson))
	req1.Header.Set("Content-Type", "application/json")

	client1 := &http.Client{}
	resp1, err1 := client1.Do(req)
	if err1 != nil {
		panic(err)
	}
	defer resp1.Body.Close()

	return nil
}

func Init() error {
	runtime.RegisterPlugin("println_decision_logger", Factory{})
	return nil
}
