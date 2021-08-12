package smdsconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	configFilePath            string = "config.json"
	maxNumberControllerBoards uint   = 20
)

type SMDSConfig interface {
	NumberControllerBoards() uint
}

type smdsConfig struct {
	NumControllerBoards uint `json:"number-controller-boards"`
}

func (this *smdsConfig) NumberControllerBoards() uint {
	return this.NumControllerBoards
}

var curSMDSConfig *smdsConfig

func UpdateAndSaveSMDSConfig(newSMDSConfig SMDSConfig) error {
	var err error

	if newSMDSConfig.NumberControllerBoards() > maxNumberControllerBoards {
		err = errors.New(fmt.Sprintf("Requested Number of Controller Boards %d exceeds max of %d", newSMDSConfig.NumberControllerBoards(), maxNumberControllerBoards))
	} else {
		writeableConfig := &smdsConfig{NumControllerBoards: newSMDSConfig.NumberControllerBoards()}
		err = writeSMDSConfigToJSONFile(writeableConfig)
		if err == nil {
			curSMDSConfig.NumControllerBoards = writeableConfig.NumControllerBoards
		}
	}

	return err
}

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
	config.NumControllerBoards = 0
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

// func getSMDSConfigFile() (io.ReadWriteCloser, error) {
// 	return os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0755)
// }

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
