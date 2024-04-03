package smdsconfig

import "os"

const (
	configFileFullPath string = "/opt/com/zmanhobbies/smds/smds-config.json"
)

type jsonSMDSConfig struct {
	ServerIDJSON    string      `json:"server-id"`
	EnvironmentJSON Environment `json:"environment"`
}

func loadConfigFromDisk() (SMDSConfig, error) {
	configFile, err := os.Open(configFileFullPath)

	if os.IsNotExist(err) {
		newConfigFile, err := os.Create(configFileFullPath)
		
		if err != nil {

		}
	}
}
