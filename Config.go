package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

type RadarConfig struct {
	IP        string `json:"Radar_IP"`
	Port      string `json:"Port"`
	OutputDir string `json:"Output_Directory"`
}

func LoadConfig() (*RadarConfig, error) {

	configFile, err := os.OpenFile("./config.json", os.O_RDONLY, 0666)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	configBytes, _ := io.ReadAll(configFile)

	var config RadarConfig = RadarConfig{}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {

		return nil, err
	}

	log.Infof("\nLoaded config: %+v", config)

	return &config, nil

}
