# teamsrv

## Local Development Workflow

### Docker image

    docker pull golang:1.21.3-alpine3.18

### Run the docker container

    cd team-srv
    docker run -it --rm -p 8080:8080 -v $(pwd):/go/src/app -w /go/src/app golang:1.21.3-alpine3.18

### Build the application

    env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/team-srv ./cmd

### Execute tests

    go test ./...

### Run the teamsrv web app

    /go/bin/team-srv
    http://localhost:8080

### Setting env variables
For convenience, once the container is up and running `build_and_run.sh` will 
set the require env variables, build, and run the app.
