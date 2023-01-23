FROM golang:1.19-alpine3.17 as build

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 go vet ./cmd
RUN env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/teamsrv ./cmd

FROM alpine:3.17.1
COPY --from=build /go/bin /app
EXPOSE 8080
CMD ["/app/teamsrv"]