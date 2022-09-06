package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/go-playground/validator"
)

type Environment struct {
	IsProduction bool   `json:"is_production" validate:"required"`
	LogLevel     int    `json:"log_level"`
	Host         string `json:"host" validate:"required"`
	Port         int    `json:"port" validate:"required"`
	SyncCycle    int    `json:"sync_cycle_in_day" validate:"required"`
}

type Locate2UConfig struct {
	BaseUrl      string `json:"base_url"`
	ClientId     string `json:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope"`

	// This field is required for creating `stop` in Locate2U
	AssignedTeamMemberID string `json:"assignedTeamMemberId"`

	// This field is used for creating a link
	LinkMessage string `json:"link_message"`
}

type ApiConfig struct {
	// ...
	Token   string `json:"token"`
	BaseUrl string `json:"base_url"`
}

type Config struct {
	Environment Environment    `json:"environment" validate:"required"`
	Locate2U    Locate2UConfig `json:"locate2u_config" validate:"required"`
	ApiConfig   ApiConfig      `json:"api_config" validate:"required"`
}

func GetConfig() *Config {
	configBytes, errConfFile := os.ReadFile("config.json")
	if errConfFile != nil {
		log.Fatalf("failed to open config.json file: %s", errConfFile.Error())
	}

	config := Config{}
	confDecErr := json.NewDecoder(bytes.NewReader(configBytes)).Decode(&config)
	if confDecErr != nil {
		log.Fatalf("failed to decode config: %s", confDecErr.Error())
	}

	validate := validator.New()
	confValErr := validate.Struct(config)
	if confValErr != nil {
		log.Fatalf("failed to validate config: %s", confValErr)
	}

	return &config
}
