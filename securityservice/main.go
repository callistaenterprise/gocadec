/**
The MIT License (MIT)

Copyright (c) 2016 Callista Enterprise

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
	"github.com/callistaenterprise/gocadec/securityservice/service"
	ct "github.com/eriklupander/cloudtoolkit"
	"github.com/spf13/viper"
	"sync"
)

var appName = "securityservice"
var configServerDefaultUrl = "http://configserver:8888"

var amqpClient *ct.MessagingClient

func main() {
	ct.Log.Println("Starting " + appName + "...")
	ct.LoadSpringCloudConfig(appName, ct.ResolveProfile(), configServerDefaultUrl)
	ct.InitTracingFromConfigProperty(appName)

	amqpClient = ct.InitMessagingClientFromConfigProperty()
	defer amqpClient.GetConn().Close()

	ct.ConfigureHystrix([]string{"get_account_secured"}, amqpClient)
	service.ConfigureClient() // Disable keep-alives so Docker Swarm can round-robin for us.

	go service.StartWebServer(viper.GetString("server_port")) // Starts HTTP service  (async)

	// Block...
	wg := sync.WaitGroup{} // Use a WaitGroup to block main() exit
	wg.Add(1)
	wg.Wait()
}
