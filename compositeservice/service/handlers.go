package service

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/callistaenterprise/gocadec/compositeservice/client"
	"encoding/json"
	"github.com/callistaenterprise/gocadec/compositeservice/model"
	"github.com/opentracing/opentracing-go"
	ct "github.com/eriklupander/cloudtoolkit"
)

func GetAccount(w http.ResponseWriter, r *http.Request) {

	span := ct.StartHTTPTrace(r, "GetAccount")
	defer span.Finish()
	vars := mux.Vars(r)
	var accountId = vars["accountId"]

	child := ct.Tracer.StartSpan("GetAccountImageUrl", opentracing.ChildOf(span.Context()))
	accountImageUrl, _ := client.GetAccountImageUrl(accountId, child)
	child.Finish()

	child = ct.Tracer.StartSpan("GetAccountData", opentracing.ChildOf(span.Context()))
	account, _ := client.GetAccountData(accountId, child)
	child.Finish()

	child = ct.Tracer.StartSpan("GetImageData", opentracing.ChildOf(span.Context()))
	imageBytes, _ := client.GetImageData(accountId, child)
	child.Finish()

	userAccount := model.UserAccount{
		Id: account.Id,
		Name: account.Name,
		ImageData: imageBytes,
		ImageUrl: string(accountImageUrl),
		AccountServedBy: account.ServedBy,
		ServedBy: ct.GetLocalIP(),
	}

	json, _ := json.Marshal(userAccount)

	writeAndReturn(w, json)
}

func HealtCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte("OK"))
}

func writeAndReturn(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func writeServerError(w http.ResponseWriter, msg string) {
	fmt.Println(msg)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}