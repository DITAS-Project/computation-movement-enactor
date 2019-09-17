package main

import (
	"github.com/spf13/viper"
)

const (
	DitasCMEPortProperty             = "port"
	DitasDeploymentEngineURLProperty = "deployment_engine.url"
	DitasPreSharedKeyProperty        = "sign.key"
	DitasTombstoneSecureProperty     = "tombstone.use_ssl"

	DitasCMEPortDefaultValue         = 8080
	DitasTombstoneSecureDefaultValue = false
)

func getProperties(pamams ...string) (map[string]string, error) {
	return GetParamsMap(viper.GetString)
}

func main() {

	viper.SetConfigName("cme")
	viper.AddConfigPath("/etc/ditas/")

	viper.SetDefault(DitasCMEPortProperty, DitasCMEPortDefaultValue)
	viper.SetDefault(DitasTombstoneSecureProperty, DitasTombstoneSecureDefaultValue)

	viper.ReadInConfig()

	properties, err := getProperties(DitasCMEPortProperty, DitasDeploymentEngineURLProperty, DitasPreSharedKeyProperty)
	if err != nil {
		panic(err.Error())
	}

	port := viper.GetInt(DitasCMEPortProperty)
	useSSL := viper.GetBool(DitasTombstoneSecureProperty)

	deploymentEngineURL, _ := properties[DitasDeploymentEngineURLProperty]

	preSharedKey, _ := properties[DitasPreSharedKeyProperty]

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
