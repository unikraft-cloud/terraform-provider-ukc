package config

// Version contains the current software version from git.
var Version = "dev"

type Config struct {
	AccessToken string
	Metro       string
}
