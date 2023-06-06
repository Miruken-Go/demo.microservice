# Demo Microservice

## Local Development

### Docker image

    docker pull golang:1.20-alpine3.18

## Local Development Workflow

### Run the docker container

    cd teamsrv
    docker run -it -p 8080:8080 -v $(pwd):/go/src/app -w /go/src/app golang:1.20-alpine3.18

### Build the application

    env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/teamsrv ./cmd

### Execute tests

    go test ./...

### Run the teamsrv web app

    /go/bin/teamsrv
    http://localhost:8080/docs
    http://localhost:8080/openapi

## Build and run the container localy

### Build Docker Image

    cd teamsrv
    docker build -t teamsrv:local .
    
### Run Docker container locally

    docker run -it -p 8080:8080 teamsrv:local

## Push Docker Image to Azure Container Repo

    az login
    az acr login -n mirukengo   
    docker tag teamsrv:local mirukengo.azurecr.io/teamsrv:20230207
    docker push mirukengo.azurecr.io/teamsrv:20230207

---

### Running a named image detached

    cd teamsrv
    docker run -itd -p 8080:8080 -v $(pwd):/go/src/app -w /go/src/app golang:1.20-alpine3.18
    docker exec -it go_server sh
    docker rm -f go_server
