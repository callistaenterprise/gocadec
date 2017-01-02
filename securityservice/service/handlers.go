package service

import (
	"net/http"
	ct "github.com/eriklupander/cloudtoolkit"
	"io/ioutil"
	"github.com/afex/hystrix-go/hystrix"
)

var client = &http.Client{}

func ConfigureClient() {
	var transport http.RoundTripper = &http.Transport{
		DisableKeepAlives: true,
	}
	client.Transport = transport
}
/**
 * Takes the POST body, decodes, processes and finally writes the result to the response.
 */
func SecuredGetAccount(w http.ResponseWriter, r *http.Request) {

	span := ct.StartHTTPTrace(r, "SecuredGetAccount")
	defer span.Finish()
	output := make(chan []byte, 1)
	errors := hystrix.Go("get_account_secured", func() error {

		req, _ := http.NewRequest("GET", "http://compservice:6565" + r.URL.Path, nil)
		ct.AddTracingToReq(req, span)
		resp, err := client.Do(req)

		if err != nil {
			return err
		}
		data, _ := ioutil.ReadAll(resp.Body)
		output <- data
		return nil
	}, nil)

	select {
	case out := <-output:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(out)

	case err := <-errors:
		writeServerError(w, err.Error())
	}
}


func writeServerError(w http.ResponseWriter, msg string) {
	ct.Log.Println(msg)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}
