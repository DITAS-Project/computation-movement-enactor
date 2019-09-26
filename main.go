package main

import (
	"github.com/spf13/viper"
)

const (
	cMEPortProperty             = "port"
	deploymentEngineURLProperty = "deployment_engine.url"
	preSharedKeyProperty        = "sign.key"
	tombstoneSecureProperty     = "tombstone.use_ssl"

	cMEPortDefaultValue         = 8080
	tombstoneSecureDefaultValue = false
)

func getProperties(params ...string) (map[string]string, error) {
	return GetParamsMap(viper.GetString, params...)
}

func main() {

	viper.SetConfigName("cme")
	viper.AddConfigPath("/etc/ditas/")

	viper.SetDefault(cMEPortProperty, cMEPortDefaultValue)
	viper.SetDefault(tombstoneSecureProperty, tombstoneSecureDefaultValue)

	viper.ReadInConfig()

	properties, err := getProperties(cMEPortProperty, deploymentEngineURLProperty, preSharedKeyProperty)
	if err != nil {
		panic(err.Error())
	}

	port := viper.GetInt(cMEPortProperty)
	useSSL := viper.GetBool(tombstoneSecureProperty)

	deploymentEngineURL, _ := properties[deploymentEngineURLProperty]

	preSharedKey, _ := properties[preSharedKeyProperty]

	movementController, err := NewMovementController(deploymentEngineURL, preSharedKey, useSSL)
	if err != nil {
		panic(err)
	}

	handler := mainHandler{
		port:               port,
		movementController: movementController,
	}

	handler.Start()
}
