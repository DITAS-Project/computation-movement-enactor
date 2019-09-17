package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type mainHandler struct {
	port               int
	movementController *movementController
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
	return GetParamsMap(query.Get)
}

func (h mainHandler) MoveVDC(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	vdcID := ps.ByName("vdcID")
	params, err := h.getParams(r.URL.Query(), "sourceInfra", "targetInfra", "blueprintId")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ip, httpErr := h.movementController.MoveVDC(params["blueprintId"], vdcID, params["sourceInfra"], params["targetInfra"])
	if err != nil {
		respondWithError(w, httpErr.code, httpErr.body.Error())
		return
	}

	respond(w, http.StatusOK, []byte(ip), "plain/text")
}

func (h mainHandler) Start() {
	router := httprouter.New()
	router.PUT("/vdc/:vdcID", h.MoveVDC)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", h.port), router))
}
