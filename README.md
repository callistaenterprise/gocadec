# gocadec
Various Go-based microservices for the Cadec 2017 presentation

### The services

- Accountservice: Serves account objects as JSON, based on data pre-seeded on startup into a local BoltDB instance.
- Compositeservice: Contains no data of its own, orchestrates calls onward to the Accountservice and Imageservice to build the full Account object.
- Eventservice: Listens to for messages on an AMQP queue named "hello". Just logs.
- Imageservice: Returns a dummy URL _or_ actually performs a Sepia filtering on a static resource.
- Securityservice: Intermediary between the Zuul Edge server and the Compositeservice. Extracts the OAuth token from incoming requests and validates it against the Spring Security authserver.

### Dependencies
All of these Microservices are built to run _without_ a Discovery Service present, e.g. with Docker Swarm or possibly Kubernetes (not tested). They are also built as a showcase on how to nicely integrate with a Spring Cloud / Netflix OSS landscape.
For the sake of code reuse, most of the integration with the Spring/Netflix services (config, auth, turbine, hystrix dashboard etc.) are handled by a "cloud toolkit" I've written:

[Cloud Toolkit](https://github.com/eriklupander/cloudtoolkit)

### Supporting subprojects

- Loadtest: Simple Go program that hammers the EDGE server with requests over HTTPs.

     go run *.go -users=100
     

### Dockerization

Each microservice comes with a Dockerfile. Typically, all microservices require a working Netflix/Spring Cloud environment with at least the following Docker Swarm services deployed:

- configserver: E.g. Spring Cloud config that serves some basic conf, see NNNN
- rabbitmq: Used for client discovery for Turbine, e.g. turbine listens to a specific queue where the Go services can announce their hystrix streams.
- authserver: Only required if calling the secured APIs. For demo purposes, most of the microservices expose themselves publically.
- zipkin: Opentracing server.
- zuul: Acts as EDGE server, terminating HTTPs among things. Just there for demo purposes, configuration needs to be added here.

## Building, running....
Mostly TODO. You can try to build using the gradle plugins, otherwise you can try creating a shell script in the root folder (gocadec/) with the following content:
 
     #!/bin/bash
     export GOOS=linux
     export CGO_ENABLED=0
     
     cd imageservice;go get;go build -o imageservice-linux-amd64;echo built `pwd`;cd ..
     cd accountservice;go get;go build -o accountservice-linux-amd64;echo built `pwd`;cd ..
     cd compositeservice;go get;go build -o compositeservice-linux-amd64;echo built `pwd`;cd ..
     cd eventservice;go get;go build -o eventservice-linux-amd64;echo built `pwd`;cd ..
     cd securityservice;go get;go build -o securityservice-linux-amd64;echo built `pwd`;cd ..
     
     export GOOS=darwin

I.e. it will cd into each microservice directory, fetch deps and then build a linux/amd64 binary. Change the GOOS according to your OS (darwin,windows,linux are good bets).

You may also want to write some kind of shell script that builds all the docker images and creates Docker Swarm services for them. An example could be something like:

    #!/bin/bash
    
    # Remove existing services
    docker service rm accountservice
    docker service rm compservice
    docker service rm imageservice
    docker service rm eventservice
    docker service rm securityservice
    # docker service rm dvizz
    
    # Build and tag
    docker build go-accountservice -t yourself/accountservice
    docker build go-composite-service -t yourself/compservice
    docker build go-imageservice -t yourself/imageservice
    docker build go-eventservice -t yourself/eventservice
    docker build go-securityservice -t yourself/securityservice
    
    # Push images if you're running more than one Swarm Node
    # docker push yourself/accountservice
    # docker push yourself/compservice
    # docker push yourself/imageservice
    # docker push yourself/eventservice
    # docker push yourself/securityservice
     
    # Create new services
    docker service create --replicas 1 --name accountservice -p 7777:7777 --network my-network --update-delay 10s --with-registry-auth  --update-parallelism 1 yourself/accountservice
    docker service create --replicas 1 --name imageservice -p 6767:6767 --network my-network --update-delay 10s --with-registry-auth  --update-parallelism 1 yourself/imageservice
    docker service create --replicas 1 --name compservice -p 6565:6565 -p 8181:8181 --network my-network --update-delay 10s --with-registry-auth  --update-parallelism 1 yourself/compservice
    docker service create --replicas 1 --name eventservice -p 6868:6868 --network my-network --update-delay 10s --with-registry-auth  --update-parallelism 1 yourself/eventservice
    docker service create --replicas 1 --name securityservice -p 6666:6666 --network my-network --update-delay 10s --with-registry-auth  --update-parallelism 1 yourself/securityservice

Please note that the microservices won't work properly since they still require at least a running configserver swarm service (TODO add this to these docs), rabbitmq etc.

### Docker Swarm deployment

This is a somewhat advanced topic, which is TODO to document at the time of writing.

### 