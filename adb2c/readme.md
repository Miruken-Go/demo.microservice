# adb2c

## Docker image

    docker pull golang:1.21.3-alpine3.18

## api-connector-srv

### Run the docker container

    cd adb2c
    docker run -it --rm -p 8080:8080 -v $(pwd):/go/src/app -w /go/src/app golang:1.21.3-alpine3.18

### Build the application

    env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/api-connector-srv ./cmd/api-connector-srv

### Execute tests

    go test ./...

### Run api-connector-srv

    /go/bin/api-connector-srv
    http://localhost:8080/enrich

### Setting env variables
For convenience, once the container is up and running this script will 
set the require env variables, build, and run the app.

    ./cmd/api-connector-srv/build_and_run.dev.sh 

### Build the container

    docker build -f cmd/api-connector-srv/Dockerfile -t adb2c-connector-srv:local .

## auth-srv

### Run the docker container

    cd adb2c
    docker run -it --rm -p 8080:8080 -v $(pwd):/go/src/app -w /go/src/app golang:1.21.3-alpine3.18

### Build the application

    env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/auth-srv ./cmd/auth-srv

### Execute tests

    go test ./...

### Run auth-srv

    /go/bin/auth-srv
    http://localhost:8080

### Setting env variables
For convenience, once the container is up and running this script will 
set the require env variables, build, and run the app.

    ./cmd/auth-srv/build_and_run.dev.sh 

### Build the container

    docker build -f cmd/auth-srv/Dockerfile -t adb2c-connector-srv:local .

