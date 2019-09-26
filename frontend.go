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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type mainHandler struct {
	port               int
	movementController *MovementController
}

func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
	return
}

func respond(w http.ResponseWriter, code int, payload []byte, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	w.Write(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respond(w, code, []byte(message), "plain/text")
}

func (h mainHandler) getParams(query url.Values, parameters ...string) (map[string]string, error) {
	return GetParamsMap(query.Get, parameters...)
}

// MoveVDC moves a VDC instance from one cluster to another
func (h mainHandler) MoveVDC(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// swagger:operation PUT /vdc/{vdcId} vdc MoveVDC
	//
	//
	// Moves a VDC from one infrastructure to another by:
	// - Creating a copy of the VDC in the target infrastructure if it doesn't exist
	// - Setting it to "serve mode" if it already exists
	// - Setting the VDC in the source infrastructure to "redirect mode" to the one in the target infrastructure
	// - Returns the IP and port of the VDC copy serving requests
	//
	// ---
	// consumes:
	// - text/plain
	//
	// produces:
	// - application/json
	// - text/plain
	//
	// parameters:
	// - name: vdcId
	//   in: path
	//   type: string
	//   required: true
	//   description: The indentifier of the VDC to move
	// - name: sourceInfra
	//   in: query
	//   type: string
	//   required: true
	//   description: The identifier of the infrastructure in which the VDC is actually serving requests.
	// - name: targetInfra
	//   in: query
	//   type: string
	//   required: true
	//   description: The identifier of the infrastructure that must serve requests from now on.
	//
	// responses:
	//   200:
	//     description: The IP and port of the VDC instance serving requests
	//     schema:
	//       $ref: "#/definitions/infrastructureInformation"
	//   400:
	//     description: Bad request
	//   500:
	//     description: Internal error
	vdcID := ps.ByName("vdcID")
	params, err := h.getParams(r.URL.Query(), "sourceInfra", "targetInfra", "blueprintId")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	conf, httpErr := h.movementController.MoveVDC(vdcID, params["sourceInfra"], params["targetInfra"])
	if httpErr != nil {
		respondWithError(w, httpErr.code, httpErr.body.Error())
		return
	}

	respondJSON(w, http.StatusOK, conf)
}

func (h mainHandler) Start() {
	router := httprouter.New()
	router.PUT("/vdc/:vdcID", h.MoveVDC)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", h.port), router))
}
