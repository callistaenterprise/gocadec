package main

import (
	"flag"
	"github.com/callistaenterprise/gocadec/compositeservice/client"
	"github.com/callistaenterprise/gocadec/compositeservice/service"
	ct "github.com/eriklupander/cloudtoolkit"
	"github.com/spf13/viper"
	"log"
	"os"
	"runtime/pprof"
	"time"
        "sync"
        "os/signal"
        "syscall"
)

var appName = "compservice"

// var EnvProfile string = "dev"

var configServerDefaultUrl string // = "http://configserver:8888"
var messageBrokerDefaultUrl string
var profile string

var amqpClient *ct.MessagingClient

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	start := time.Now().UTC()
	ct.Log.Println("Starting " + appName + "...")
	parseFlags()
	initProfilingIfEnabled()
	// Comment in the line below to dump various hostname ips to log to see what mood the DNS resolver is in...
	// ct.DumpDNS()

	ct.LoadSpringCloudConfig(appName, profile, configServerDefaultUrl)
	ct.InitTracingFromConfigProperty(appName)

	// Initialize AMQP connection
	amqpClient = ct.InitMessagingClientFromConnectionString(messageBrokerDefaultUrl)
	defer amqpClient.GetConn().Close()

	ct.ConfigureHystrix([]string{"get_account_secured"}, amqpClient)

	// Configure the HTTP client (disable Keep-alives so Docker Swarm round-robins for us)
	client.ConfigureHttpClient()

	// Starts HTTP service  (async)
	go service.StartWebServer(viper.GetString("server_port"))

	ct.Log.Printf("Started %v in %v", appName, time.Now().UTC().Sub(start))
	// Block...
	wg := sync.WaitGroup{} // Use a WaitGroup to block main() exit
	wg.Add(1)
	wg.Wait()
}
func initProfilingIfEnabled() {
        if *cpuprofile != "" {
                f, err := os.Create(*cpuprofile)
                if err != nil {
                        log.Fatal(err)
                }
                pprof.StartCPUProfile(f)
                handleSigterm()
        }
}

// Handles Ctrl+C gracefully, e.g. dumping CPU profiling to disk before exiting.
func handleSigterm() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt)
        signal.Notify(c, syscall.SIGTERM)
        go func() {
                <-c
                pprof.StopCPUProfile()
                os.Exit(1)
        }()
}

func parseFlags() {
	configServerUrl := flag.String("configserverUrl", "http://configserver:8888", "Address to config server")
	messageBrokerUrl := flag.String("messageBrokerUrl", "amqp://guest:guest@rabbitmq:5672", "Address to config server")
	profilePtr := flag.String("profile", "dev", "Application profile")
	flag.Parse()
	configServerDefaultUrl = *configServerUrl
	messageBrokerDefaultUrl = *messageBrokerUrl
	profile = *profilePtr
}
