FROM golang:1.20-alpine3.18 as build

WORKDIR /go/src/app
COPY . .

RUN env CGO_ENABLED=0 go vet ./cmd
RUN env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/teamsrv ./cmd

FROM alpine:3.18.0
ARG app_version="0.0.0.0"
ARG app_source_url="https://github.com/Miruken-Go/demo.microservice"
ENV App__Version=$app_version
ENV App__Source__Url=$app_source_url
ENV Login__OAuth__0__Module=login.jwt
ENV Login__OAuth__0__Options__JWKS__Uri=https://teamsrvdevcraig.b2clogin.com/teamsrvdevcraig.onmicrosoft.com/b2c_1_signin/discovery/v2.0/keys
COPY --from=build /go/bin/teamsrv /app/teamsrv
EXPOSE 8080
CMD ["/app/teamsrv"]