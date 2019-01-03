package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

// Config - main configuration struct.
type Config struct {
	dir  string
	file string

	CoreAddress string `yaml:"core_address"`
}

// ConfigDirName - ConfigDirName
var ConfigDirName = "pasd-cli"

// InitConfig creates config.
func InitConfig() Config {
	conf := Config{}

	// Parse user's configs
	conf.loadUserConfig()

	return conf
}

// loadUserConfig loads config file and merge it with default
// one. If there's no such file, it creates new.
func (conf *Config) loadUserConfig() {
	// Get dirname
	confHome := os.Getenv("XDG_CONFIG_HOME")
	conf.dir = path.Join(confHome, ConfigDirName)
	if confHome == "" {
		homeDir := os.Getenv("HOME")
		if homeDir == "/root" {
			conf.dir = "/etc/pasd-cli"
		} else {
			conf.dir = path.Join(os.Getenv("HOME"), ".config", ConfigDirName)
		}
	}

	// Create dir if nothing
	if _, err := os.Stat(conf.dir); os.IsNotExist(err) {
		os.MkdirAll(conf.dir, os.ModePerm)
	}

	// Get filename
	conf.file = path.Join(conf.dir, "config.yml")

	// Open file
	file, err := os.OpenFile(conf.file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("Cannot read/create config file.")
	}
	defer file.Close()

	// Read file
	content, err := ioutil.ReadAll(file)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatalln("Cannot read config file.", err)
	}

	// Write default config
	if len(content) == 0 {
		content = []byte(DefaultConfigContent)
		file.WriteString(DefaultConfigContent)
	}

	// Parse
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		log.Fatalln("Cannot parse config file.", err)
	}
}
