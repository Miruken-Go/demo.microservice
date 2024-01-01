package db

import "time"

type Config struct {
	Name          string
	ConnectionUri string
	Timeout       time.Duration
	Provision     bool
}

