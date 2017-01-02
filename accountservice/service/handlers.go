package service

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
	"github.com/opentracing/opentracing-go"
	"github.com/eriklupander/cloudtoolkit"
)

/**
 * Loads an account object from the underlying
 */
func GetAccount(w http.ResponseWriter, r *http.Request) {

	// Extract tracing context, if possible
	span := cloudtoolkit.StartHTTPTrace(r, "GetAccount")
	defer endSpan(span)

	vars := mux.Vars(r)
	var accountId = vars["accountId"]
	cloudtoolkit.Log.Println("Getting account " + accountId)
	account, err := QueryAccount(accountId, span)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	cloudtoolkit.Log.Println("Done getting account " + account.Name)
	// Enrich with IP of serving node
	account.ServedBy = cloudtoolkit.GetLocalIP()
	json, _ := json.Marshal(account)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(json)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
}

func endSpan(span opentracing.Span) {
	span.Finish()
}
