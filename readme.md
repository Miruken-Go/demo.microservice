# Demo Microservice

## Docker

    docker pull golang:1.19-alpine3.17
    docker run -itd -p 80:80 -v $(pwd):/app --name go_server golang:1.19-alpine3.17
    docker exec -it go_server sh
    docker rm -f go_server

