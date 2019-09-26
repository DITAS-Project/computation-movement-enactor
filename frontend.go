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

func (h mainHandler) MoveVDC(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	vdcID := ps.ByName("vdcID")
	params, err := h.getParams(r.URL.Query(), "sourceInfra", "targetInfra", "blueprintId")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	conf, httpErr := h.movementController.MoveVDC(params["blueprintId"], vdcID, params["sourceInfra"], params["targetInfra"])
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
