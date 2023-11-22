FROM golang:1.21.3-alpine3.18 as build

WORKDIR /go/src/app
COPY . .
WORKDIR /go/src/app/adb2c/cmd/auth-srv

RUN env CGO_ENABLED=0 go vet .
RUN env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/auth-srv .

FROM alpine:3.18.0
ARG app_version="0.0.0.0"
ARG app_source_url="https://github.com/Miruken-Go/demo.microservice"
ENV App__Version=$app_version
ENV App__Source__Url=$app_source_url
COPY --from=build /go/bin/auth-srv /app/auth-srv
EXPOSE 8080
CMD ["/app/auth-srv"]