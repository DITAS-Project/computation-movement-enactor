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
	"errors"
	"fmt"
	"net/http"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/go-resty/resty/v2"
)

// HTTPError is the representation of an error that arises from a HTTP call.
// It contains the error code and the error message
type HTTPError struct {
	code int
	body error
}

// infrastructureInformation contains information of a running VDC
// swagger:model
type infrastructureInformation struct {
	// IP of the infrastructure
	IP string
	// Port of the tombstone
	TombstonePort int
	// Port in which the CAF is serving
	CAFPort int
}

type vdcConfiguration struct {
	Infrastructures map[string]infrastructureInformation
}

// MovementController is the structure to control movement of VDCs
type MovementController struct {
	deploymentEngineURL    string
	tombstonePrefix        string
	blueprintID            string
	deploymentEngineClient *resty.Client
	tombstoneClient        *resty.Client
}

// NewMovementController creates a new Computation Movement Controller
// with the deployment engine location, a pre shared key to sign the tombstone requests and a
// boolean to indicate if it should communicate with the tombstone service in a secure (https) or insecure manner (http)
func NewMovementController(blueprintID, deploymentEngineURL, preSharedKey string, tombstoneSecure bool) (*MovementController, error) {
	result := MovementController{
		deploymentEngineURL:    deploymentEngineURL,
		deploymentEngineClient: resty.New(),
		tombstonePrefix:        "http",
		blueprintID:            blueprintID,
	}

	if tombstoneSecure {
		result.tombstonePrefix = "https"
	}

	secret := jwt.NewHS256([]byte(preSharedKey))
	token, err := jwt.Sign(jwt.Payload{}, secret)
	if err != nil {
		return &result, err
	}

	tombstoneClient := resty.New().OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		req.SetAuthToken(string(token))
		return nil
	})
	result.tombstoneClient = tombstoneClient

	return &result, nil
}

func (c *MovementController) decodeError(resp *resty.Response, err error) *HTTPError {
	if err != nil {
		return &HTTPError{
			code: http.StatusInternalServerError,
			body: err,
		}
	}
	if resp.IsError() {
		return &HTTPError{
			code: resp.StatusCode(),
			body: errors.New(string(resp.Body())),
		}
	}
	return nil
}

func (c *MovementController) getVDCInfo(vdcID string) (vdcConfiguration, *HTTPError) {
	url := fmt.Sprintf("%s/blueprint/%s/vdc/%s", c.deploymentEngineURL, c.blueprintID, vdcID)
	var config vdcConfiguration
	err := c.decodeError(c.deploymentEngineClient.R().SetResult(&config).Get(url))
	return config, err
}

func (c *MovementController) moveVDC(vdcID, targetInfra string) (vdcConfiguration, *HTTPError) {
	url := fmt.Sprintf("%s/blueprint/%s/vdc/%s?targetInfra=%s", c.deploymentEngineURL, c.blueprintID, vdcID, targetInfra)
	var config vdcConfiguration
	err := c.decodeError(c.deploymentEngineClient.R().SetResult(&config).Put(url))
	return config, err
}

func (c *MovementController) getTombstoneURL(ip string, port int, path string) string {
	return fmt.Sprintf("%s://%s:%d/%s", c.tombstonePrefix, ip, port, path)
}

func (c *MovementController) setReviveMode(ip string, port int) *HTTPError {
	return c.decodeError(c.tombstoneClient.R().Post(c.getTombstoneURL(ip, port, "revive")))
}

func (c *MovementController) setRedirectMode(ip string, port int, targetIP string, targetPort int) *HTTPError {
	request := c.tombstoneClient.R().SetBody(fmt.Sprintf("%s:%d", targetIP, targetPort))
	return c.decodeError(request.Post(c.getTombstoneURL(ip, port, "tombstone")))
}

// MoveVDC moves a VDC from one infrastructure to another by:
// - Creating a copy of the VDC in the target infrastructure if it doesn't exist
// - Setting it to "serve mode" if it already exists
// - Setting the VDC in the source infrastructure to "redirect mode" to the one in the target infrastructure
// - Returns the IP of the VDC copy serving requests
func (c MovementController) MoveVDC(vdcID, sourceInfraID, targetInfraID string) (infrastructureInformation, *HTTPError) {
	var targetInfraConfig infrastructureInformation
	config, err := c.getVDCInfo(vdcID)
	if err != nil {
		return targetInfraConfig, err
	}
	sourceInfraConfig, ok := config.Infrastructures[sourceInfraID]
	if !ok {
		return targetInfraConfig, &HTTPError{
			code: 500,
			body: fmt.Errorf("Can't find VDC %s configuration in infrastructure %s", vdcID, sourceInfraID),
		}
	}

	targetInfraConfig, ok = config.Infrastructures[targetInfraID]
	if ok {
		err = c.setReviveMode(targetInfraConfig.IP, targetInfraConfig.TombstonePort)
		if err != nil {
			return targetInfraConfig, err
		}
	} else {
		config, err = c.moveVDC(vdcID, targetInfraID)
		if err != nil {
			return targetInfraConfig, err
		}
	}

	targetInfraConfig, ok = config.Infrastructures[targetInfraID]
	if !ok {
		return targetInfraConfig, &HTTPError{
			code: 500,
			body: fmt.Errorf("Can't find configuration of VDC %s in target infrastructure %s", vdcID, targetInfraID),
		}
	}

	err = c.setRedirectMode(sourceInfraConfig.IP, sourceInfraConfig.TombstonePort, targetInfraConfig.IP, targetInfraConfig.CAFPort)

	return targetInfraConfig, err
}
