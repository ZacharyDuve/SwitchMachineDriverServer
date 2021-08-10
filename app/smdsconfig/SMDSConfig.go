package smdsconfig

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"os"
)

const (
	configFilePath string = "config.json"
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
	return errors.New("UNIMPLEMENTED")
}

func writeSMDSConfigToJSONFile(config *smdsConfig) error {
	configFile, err := os.Open(configFilePath)
	defer configFile.Close()

	if err == nil {
		err = writeSMDSConfigAsJSONToWriter(config, configFile)
	}

	return err
}

func newDefaultConfig() *smdsConfig {
	config := &smdsConfig{}
	config.NumControllerBoards = 0
	randBytes := make([]byte, 8)
	rand.Read(randBytes)

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
			configFile, err := createSMDSConfigFile()
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

func getSMDSConfigFile() (io.ReadWriteCloser, error) {
	return os.Open(configFilePath)
}

func createSMDSConfigFile() (io.ReadWriteCloser, error) {
	return os.Create(configFilePath)
}

func readSMDSConfigFromFile() *smdsConfig {
	config := &smdsConfig{}

	configFile, err := getSMDSConfigFile()
	defer configFile.Close()

	if err != nil {
		panic(err)
	}

	configFileData, err := io.ReadAll(configFile)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configFileData, config)

	if err != nil {
		panic(err)
	}

	return config
}
