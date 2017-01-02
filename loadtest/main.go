package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"net/url"
)

var Log = logrus.New()

func main() {

	usersPtr := flag.Int("users", 10, "Number of users")
	delayPtr := flag.Int("delay", 1000, "Delay between calls per user")
	flag.Parse()

	users := *usersPtr
	var _ int = *delayPtr

	for i := 0; i < users; i++ {
		go securedTest()
	}

	// Block...
	wg := sync.WaitGroup{} // Use a WaitGroup to block main() exit
	wg.Add(1)
	wg.Wait()

}

func getToken() string {

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Add("client_id", "acme")
	data.Add("scope", "webshop")
	data.Add("username", "user")
	data.Add("password", "password")
	req, err := http.NewRequest("POST", "https://192.168.99.100:9999/uaa/oauth/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		panic(err.Error())
	}
	var DefaultTransport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	headers := make(map[string][]string)
	headers["Authorization"] = []string{"Basic YWNtZTphY21lc2VjcmV0"}
	headers["Content-Type"] = []string{"application/x-www-form-urlencoded"}

	req.Header = headers

	resp, err := DefaultTransport.RoundTrip(req)
	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode > 299 {
		panic("Call to get auth token returned status " + resp.Status)
	}
	respdata, _ := ioutil.ReadAll(resp.Body)
	m := make(map[string]interface{})
	err = json.Unmarshal(respdata, &m)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Got TOKEN: " + string(m["access_token"].(string)))
	return string(m["access_token"].(string))
}

var defaultTransport http.RoundTripper = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func securedTest() {

	var token = getToken()
	for {
		accountId := rand.Intn(99) + 10000
		url := "https://192.168.99.100:8765/api/account/" + strconv.Itoa(accountId)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Connection", "keep-alive")
		req.Header.Add("Keep-Alive", "timeout=10, max=5")
		resp, err := defaultTransport.RoundTrip(req)

		if err != nil {
			panic(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		Log.Println(string(body))
		time.Sleep(time.Second * 1)
	}
}

func standardTest() {

	for {
		accountId := rand.Intn(99) + 10000
		url := "https://192.168.99.100:8765/api/account/" + strconv.Itoa(accountId)

		var DefaultTransport http.RoundTripper = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := DefaultTransport.RoundTrip(req)

		if err != nil {
			panic(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		Log.Println(string(body))
		time.Sleep(time.Second * 1)
	}

}
