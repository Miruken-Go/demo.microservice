package cred

// Config is the configuration for the azure credentials.
type Config struct {
	ClientSecret []struct{
		Id     string
		Secret string
	}
}
