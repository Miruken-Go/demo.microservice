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
    http://localhost:8080

## Build and run the container localy

### Build Docker Image

    cd teamsrv
    docker build -t -e APPLICATION_VERSION=local teamsrv:local .
    
### Run Docker container locally

    docker run -it -p 8080:8080 teamsrv:local

## Build and Push Docker Image to Azure Container Repo

    cd teamsrv
    TAG=$(date +%s); echo $TAG
    IMAGE_NAME="teamsrvdevmichael.azurecr.io/teamsrv:$TAG"; echo $IMAGE_NAME
    docker build --build-arg application_version=$TAG -t $IMAGE_NAME .
    az login
    az acr login -n teamsrvdevmichael   
    docker push $IMAGE_NAME

---

### Running a named image detached

    cd teamsrv
    docker run -itd -p 8080:8080 -v $(pwd):/go/src/app -w /go/src/app golang:1.20-alpine3.18
    docker exec -it go_server sh
    docker rm -f go_server
