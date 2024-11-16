package config

import (
	"flag"
	"log"
)

type Config struct {
	Username string
	Token    string
	RepoCSV  string
	CloneDir string
}

const (
	DefaultCSVFile  = "repositories.csv"
	DefaultCloneDir = "clonedir"
)

// ParseFlags parses command line flags and returns a Config struct
func ParseFlags() *Config {
	config := &Config{}
	
	flag.StringVar(&config.Username, "u", "", "Username")
	flag.StringVar(&config.Token, "t", "", "Token")
	flag.StringVar(&config.RepoCSV, "f", "", "CSV file")
	flag.StringVar(&config.CloneDir, "d", "", "Clone directory")
	flag.Parse()

	if config.Username == "" || config.Token == "" {
		// print usage and exit
		flag.Usage()
		log.Fatal("Username and Token are required")
	}

	// set default values
	if config.RepoCSV == "" {
		config.RepoCSV = DefaultCSVFile
	}

	if config.CloneDir == "" {
		config.CloneDir = DefaultCloneDir
	}

	return config
}
