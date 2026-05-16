package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	DirName    = ".sentinel"
	TokenFile  = "token"
	ConfigFile = "config.json"
)

type Config struct {
	ActiveProfile string             `json:"active_profile"`
	Profiles      map[string]Profile `json:"profiles"`
}

type Profile struct {
	Name    string `json:"name"`
	APIURL  string `json:"api_url"`
	Email   string `json:"email"`
}

func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, DirName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return dir, nil
}

func SaveToken(token string) error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, TokenFile)
	return os.WriteFile(path, []byte(token), 0600)
}

func LoadToken() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, TokenFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func RemoveToken() error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, TokenFile)
	return os.Remove(path)
}

func LoadConfig() (*Config, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, ConfigFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				ActiveProfile: "default",
				Profiles: map[string]Profile{
					"default": {
						Name:   "default",
						APIURL: "https://api.hackura.app",
					},
				},
			}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, ConfigFile)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
