# defaultContainerImage

## Local Development

### Docker image

    docker pull golang:1.21.3-alpine3.18

### Run the docker container

    cd defaultContainerImage
    docker run -it --rm -p 8080:8080 -v $(pwd):/go/src/app -w /go/src/app golang:1.21.3-alpine3.18

### Build the application

    env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/teamsrv ./cmd

### Execute tests

    go test ./...

### Run the app

    /go/bin/teamsrv
    http://localhost:8080

## Build and run the container localy

### Build Docker Image

    cd defaultContainerImage/app
    docker build --build-arg ENV=local -t teamsrv:default .
    
### Run Docker container locally

    docker run -it -p 8080:8080 teamsrv:default
