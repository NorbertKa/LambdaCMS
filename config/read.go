package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadConfig(filepath string) (*Config, error) {
	config := NewConfig()
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
