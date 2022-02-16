package smdsconfig

import (
	"encoding/json"
	"io"
	"os"

	"github.com/google/uuid"
)

const (
	configFilePath string = "server-config.json"
)

type SMDSConfig interface {
	SMDSId() string
}

type smdsConfig struct {
	id string
}

func (this *smdsConfig) SMDSId() string {
	return this.id
}

var curSMDSConfig *smdsConfig

func writeSMDSConfigToJSONFile(config *smdsConfig) error {
	configFile, err := os.Create(configFilePath)
	defer configFile.Close()

	if err == nil {
		err = writeSMDSConfigAsJSONToWriter(config, configFile)
	}

	return err
}

func newDefaultConfig() *smdsConfig {
	config := &smdsConfig{}
	config.id = uuid.New().String()
	return config
}

func writeSMDSConfigAsJSONToWriter(config *smdsConfig, w io.Writer) error {
	configAsJSONBytes, err := json.Marshal(config)

	if err == nil {
		_, err = w.Write(configAsJSONBytes)
	}

	return err
}

func GetSMDSConfig() SMDSConfig {
	if curSMDSConfig == nil {
		if configFileExists() {
			curSMDSConfig = readSMDSConfigFromFile()
		} else {
			//Need to generate configuration
			curSMDSConfig = newDefaultConfig()
			configFile, err := os.Create(configFilePath)
			defer configFile.Close()
			if err != nil {
				panic(err)
			}
			writeSMDSConfigAsJSONToWriter(curSMDSConfig, configFile)
		}

	}

	return curSMDSConfig
}

func configFileExists() bool {
	_, err := os.Stat(configFilePath)

	return !os.IsNotExist(err)
}

func readSMDSConfigFromFile() *smdsConfig {
	config := &smdsConfig{}

	configFile, err := os.Open(configFilePath)
	defer configFile.Close()

	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(configFile).Decode(config)

	if err != nil {
		panic(err)
	}

	return config
}
