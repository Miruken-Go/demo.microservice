# Demo Microservice

## Docker

    # Naming 
    docker pull golang:1.19-alpine3.17
    docker run -itd -p 80:80 -v $(pwd):/app --name go_server golang:1.19-alpine3.17
    docker exec -it go_server sh
    docker rm -f go_server
    
    docker run -it -p 8080:8080 -v $(pwd):/app golang:1.19-alpine3.17
    env GOOS=linux CGO_ENABLED=0 go build -o teamsrv ./cmd
    http://localhost:8080/process




cd teamsrv

docker build -t teamsrv:latest .
docker run -it -p 8080:8080 teamsrv:latest

## Local Dev
    docker run -it -p 80:80 -v $(pwd):/go/src golang:1.19-alpine3.17

## Build and Deploy

    cd teamsrv
    docker build -t teamsrv:latest .

    az login
    az acr login -n mirukengo   
    docker tag teamsrv:latest mirukengo.azurecr.io/teamsrv:20230207
    docker push mirukengo.azurecr.io/teamsrv:20230207