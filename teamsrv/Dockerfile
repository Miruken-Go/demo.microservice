FROM golang:1.20-alpine3.18 as build

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 go vet ./cmd
RUN env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/teamsrv ./cmd

FROM alpine:3.18.0
ARG application_version="0.0.0.0"
ENV APPLICATION_VERSION=$application_version
COPY --from=build /go/bin/teamsrv /app/teamsrv
EXPOSE 8080
CMD ["/app/teamsrv"]