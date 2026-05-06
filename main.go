package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type configJSON struct {
	CALENDAR_SERVICE_ACCOUNT_JSON string `json:"CALENDAR_SERVICE_ACCOUNT_JSON"`
	CALENDAR_SERVICE_ACCOUNT_B64 string `json:"CALENDAR_SERVICE_ACCOUNT_B64"`
	CALENDAR_ID                  string `json:"CALENDAR_ID"`
	CALENDAR_MAX_RESULTS         string `json:"CALENDAR_MAX_RESULTS"`
	CALENDAR_TIMEZONE            string `json:"CALENDAR_TIMEZONE"`
}

func main() {
	configStr := flag.String("config", "", "JSON string with Calendar settings (overrides .env, overridden by env vars)")
	flag.Parse()

	if *configStr != "" {
		var cfg configJSON
		if err := json.Unmarshal([]byte(*configStr), &cfg); err != nil {
			log.Fatalf("Failed to parse --config JSON: %v", err)
		}
		setIfNotEnv := func(key, val string) {
			if val != "" && os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
		setIfNotEnv("CALENDAR_SERVICE_ACCOUNT_JSON", cfg.CALENDAR_SERVICE_ACCOUNT_JSON)
		setIfNotEnv("CALENDAR_SERVICE_ACCOUNT_B64", cfg.CALENDAR_SERVICE_ACCOUNT_B64)
		setIfNotEnv("CALENDAR_ID", cfg.CALENDAR_ID)
		setIfNotEnv("CALENDAR_MAX_RESULTS", cfg.CALENDAR_MAX_RESULTS)
		setIfNotEnv("CALENDAR_TIMEZONE", cfg.CALENDAR_TIMEZONE)
	}

	log.Println("Starting Oido Google Calendar MCP Server v1.0.0...")
	RunMCPServer()
}
