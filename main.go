/**
 * Copyright 2018 Atos
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy of
 * the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 *
 * This is being developed for the DITAS Project: https://www.ditas-project.eu/
 */
package main

import (
	"github.com/spf13/viper"
)

const (
	cMEPortProperty             = "port"
	deploymentEngineURLProperty = "deployment_engine.url"
	blueprintIDProperty         = "blueprint.id"
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

	properties, err := getProperties(cMEPortProperty, deploymentEngineURLProperty, preSharedKeyProperty, blueprintIDProperty)
	if err != nil {
		panic(err.Error())
	}

	port := viper.GetInt(cMEPortProperty)
	useSSL := viper.GetBool(tombstoneSecureProperty)

	deploymentEngineURL, _ := properties[deploymentEngineURLProperty]

	preSharedKey, _ := properties[preSharedKeyProperty]

	blueprintID, _ := properties[blueprintIDProperty]

	movementController, err := NewMovementController(blueprintID, deploymentEngineURL, preSharedKey, useSSL)
	if err != nil {
		panic(err)
	}

	handler := mainHandler{
		port:               port,
		movementController: movementController,
	}

	handler.Start()
}
