# Demo Microservice

## Local Development

### Updating golang version

* Search and replace the golang docker image `golang:1.21.3-alpine3.18` with the new version
* Update the `demo.microservice.build` container golang version 
    * deployment/Dockerfile
            RUN wget https://go.dev/dl/go1.21.3.linux-amd64.tar.gz
            RUN tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
    * check in dockerfile and the new image is built automatically
    * Seach and replace the container version `ghcr.io/miruken-go/demo.microservice.build:1699298856` in the .github folder

### Docker image

    docker pull golang:1.21.3-alpine3.18

## Local Development Workflow

### Run the docker container

    cd team-srv
    docker run -it --rm -p 8080:8080 -v $(pwd):/go/src/app -w /go/src/app golang:1.21.3-alpine3.18

### Build the application

    env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/team-srv ./cmd

### Execute tests

    go test ./...

### Run the team-srv web app

    /go/bin/team-srv
    http://localhost:8080

## Build and run the container localy

### Build Docker Image

    cd team-srv
    docker build -t -e APPLICATION_VERSION=local team-srv:local .
    
### Run Docker container locally

    docker run -it --rm -p 8080:8080 team-srv:local

## Build and Push Docker Image to Azure Container Repo

    cd team-srv
    TAG=$(date +%s); echo $TAG
    IMAGE_NAME="teamsrvdevmichael.azurecr.io/team-srv:$TAG"; echo $IMAGE_NAME
    docker build --build-arg application_version=$TAG -t $IMAGE_NAME .
    az login
    az acr login -n teamsrvdevmichael   
    docker push $IMAGE_NAME

---

### Running a named image detached

    cd team-srv
    docker run -itd -p 8080:8080 -v $(pwd):/go/src/app -w /go/src/app golang:1.21.3-alpine3.18
    docker exec -it go_server sh
    docker rm -f go_server
