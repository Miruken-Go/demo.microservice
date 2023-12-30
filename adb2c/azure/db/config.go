package db

import "time"

type Config struct {
	ConnectionUri string
	Timeout       time.Duration
}

