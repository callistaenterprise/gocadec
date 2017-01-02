# Image processing microservice

### Building
    
    gradle build
    
or

    go build *.go

### Docker packaging

Some stuff I always end up forgetting about on OS X:

    # Start docker machine
    docker-machine start
    
    # Run this command to configure your shell: 
    eval "$(docker-machine env default)"

Build the docker container

    docker build -t imageservice .

