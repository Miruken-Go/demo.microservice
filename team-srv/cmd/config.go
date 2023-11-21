package main

import "github.com/miruken-go/miruken/api/http/httpsrv/openapi"

type Config struct {
	App struct {
		Version string
		Source  struct {
			Url string
		}
		Port    string
	}
	OpenApi openapi.Config
}
