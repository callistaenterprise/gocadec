package client

import (
        "github.com/afex/hystrix-go/hystrix"
        "github.com/opentracing/opentracing-go"
        "net/http"
        "io/ioutil"
        "github.com/callistaenterprise/gocadec/compositeservice/model"
        "encoding/json"
        ct "github.com/eriklupander/cloudtoolkit"

)

var client = &http.Client{}

func ConfigureHttpClient() {
        var transport http.RoundTripper = &http.Transport{
                DisableKeepAlives: true,
        }
        client.Transport = transport
}

func GetAccountImageUrl(accountId string, span opentracing.Span) ([]byte, error) {
        output := make(chan []byte, 1)
        errors := hystrix.Go("get_account_image_url", func() error {
                // talk to other services
                req, _ := http.NewRequest("GET", "http://imageservice:6767/" + accountId, nil)
                ct.AddTracingToReq(req, span)

                resp, err := client.Do(req)
                if err != nil {
                        ct.Log.Errorln("Image service request returned error: " + err.Error())
                        return err
                }
                responseBody, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                        ct.Log.Errorln("Error reading imageservice response: " + err.Error())
                        return err
                }
                output <- responseBody

                // A bit ugly, return nil to indicate nothing bad happened.
                return nil
        }, nil)

        select {
        case out := <-output:
                return out, nil

        case err := <-errors:
                return nil, err
        }
        return nil, nil
}


func GetImageData(accountId string, span opentracing.Span) ([]byte, error) {

        output := make(chan []byte, 1)
        errors := hystrix.Go("get_account_image", func() error {
                // talk to other services
                req, _ := http.NewRequest("GET", "http://imageservice:6767/" + accountId, nil)
                ct.AddTracingToReq(req, span)
                resp, err := client.Do(req)
                if err != nil {
                        ct.Log.Errorln("Image service request returned error: " + err.Error())
                        return err
                }
                responseBody, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                        ct.Log.Errorln("Error reading imageservice response: " + err.Error())
                        return err
                }
                output <- responseBody

                // A bit ugly, return nil to indicate nothing bad happened.
                return nil
        }, nil)

        select {
        case out := <-output:
                return out, nil

        case err := <-errors:
                return nil, err
        }
        return nil, nil
}

func GetAccountData(accountId string, span opentracing.Span) (model.Account, error) {

        output := make(chan []byte, 1)

        errors := hystrix.Go("get_account", func() error {

                req, _ := http.NewRequest("GET", "http://accountservice:7777/accounts/" + accountId, nil)
                ct.AddTracingToReq(req, span)
                resp, err := client.Do(req)
                if err != nil {
                        ct.Log.Errorln("Account service request returned error: " + err.Error())
                        return err
                }
                responseBody, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                        ct.Log.Errorln("Error reading accountservice response: " + err.Error())
                        return err
                }
                output <- responseBody

                // A bit ugly, return nil to indicate nothing bad happened.
                return nil
        }, nil)

        select {
        case out := <-output:
                var account model.Account
                json.Unmarshal(out, &account)
                return account, nil

        case err := <-errors:
                return model.Account{}, err
        }
}


func GetQuotes(span opentracing.Span) (model.Quote, error) {

        output := make(chan []byte, 1)

        errors := hystrix.Go("get_quote", func() error {

                req, _ := http.NewRequest("GET", "http://192.168.99.100:8080/api/quote", nil)
                ct.AddTracingToReq(req, span)
                resp, err := client.Do(req)
                if err != nil {
                        ct.Log.Errorln("Quote service request returned error: " + err.Error())
                        return err
                }
                responseBody, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                        ct.Log.Errorln("Error reading Quote response: " + err.Error())
                        return err
                }
                output <- responseBody

                // A bit ugly, return nil to indicate nothing bad happened.
                return nil
        }, func(err error) error {
                output <- []byte("Every day, a new chance at failure.")
                return nil
        })

        select {
        case out := <-output:
                var quote model.Quote
                json.Unmarshal(out, &quote)
                return quote, nil

        case err := <-errors:
                return model.Quote{}, err
        }
}



func GetData(accountId string) ([]byte, error) {

        output := make(chan []byte, 1)
        errors := hystrix.Go("get_data", func() error {
                output <- getData(accountId)
                return nil
        }, func(err error) error {
                // fallback method here
                return nil
        })

        select {
        case out := <-output:
                return out, nil
        case err := <-errors:
                return nil, err
        }
}

func getData(accountId string) []byte {
     return []byte("")
}