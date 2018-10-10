package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"../Helper"
	"../Helper/Http"
)

//********** Data **********

const (
	ConfigFilename = "Configuration.json"
)

var defaultConfiguration = Data{
	API: httpServer{
		Host: "localhost",
		Port: 8080,
		TimeoutSec: helperhttp.Timeout{
			Read:  15,
			Write: 15,
			Idle:  60,
		},
	},
	Monitoring: httpServer{
		Host: "localhost",
		Port: 8081,
		TimeoutSec: helperhttp.Timeout{
			Read:  15,
			Write: 15,
			Idle:  60,
		},
	},
}

//********** Methods **********

type httpServer = helperhttp.ServerConfig

type Data struct {
	API        httpServer `json:"api"`
	Monitoring httpServer `json:"monitoring"`
}

func (obj *Data) read(filename string) (err error) {
	content, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(content, obj)
	}
	return
}

func (obj Data) write(filename string) (err error) {
	content, err := json.Marshal(obj)
	if err == nil {
		err = ioutil.WriteFile(filename, content, 0644)
	}
	return
}

func ReadAndCreate(filename string) error {
	config := defaultConfiguration
	err := config.read(filename)

	fullpath, _ := filepath.Abs(filename)
	setConfig(fullpath, config)

	if os.IsNotExist(err) || err == nil {
		err = config.write(filename)
	}
	return err
}

func ReadSpringCloudConfig() error {
	springURI, found := os.LookupEnv("APP_SPRING_CONFIG_URI")
	if found == false {
		return fmt.Errorf("Spring cloud config uri not found")
	}
	profile, found := os.LookupEnv("APP_PROFILES")
	if found == false {
		return fmt.Errorf("Spring cloud config uri not found")
	}

	configPath := fmt.Sprintf("%v/%v-%v.json", springURI, Get().Name, profile)

	var restClient helper.RestClient
	content, statusCode, err := restClient.Request(configPath)
	if statusCode != 200 {
		return fmt.Errorf("Spring cloud config %v status code: %v", configPath, statusCode)
	} else if err != nil {
		return err
	}

	config := defaultConfiguration
	err = json.Unmarshal(content, &config)

	setConfig(configPath, config)

	return err
}

/*
// Simplifiée sans méthode de class read et write
func ReadAndCreate(filename string) (config Data, err error) {
	config = defaultConfiguration

	// Lecture du contenu du fichier et remplacement/complétion des valeurs par défaut
	if content, err := ioutil.ReadFile(filename); err == nil {
		err = json.Unmarshal(content, config)
	}

	// Ecriture du fichier s'il n'existe pas
	if os.IsNotExist(err) {
		if content, err := json.Marshal(config); err == nil {
			err = ioutil.WriteFile(filename, content, 0644)
		}
	}

	return
}
*/
