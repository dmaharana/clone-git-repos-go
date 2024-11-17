package config

import (
	"flag"
	"log"
)

// Config holds the application configuration
type Config struct {
	RepoCSV   string
	CloneDir  string
	Username  string
	Token     string
	LogDir    string
	LogMaxSize int64
}

const (
	DefaultCSVFile  = "repositories.csv"
	DefaultCloneDir = "clonedir"
)

// ParseFlags parses command line flags and returns a Config struct
func ParseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.RepoCSV, "f", "", "CSV file")
	flag.StringVar(&cfg.CloneDir, "d", "", "Clone directory")
	flag.StringVar(&cfg.Username, "u", "", "Username")
	flag.StringVar(&cfg.Token, "t", "", "Token")
	flag.StringVar(&cfg.LogDir, "logdir", "logs", "Log directory")
	flag.Int64Var(&cfg.LogMaxSize, "logsize", 10*1024*1024, "Maximum log file size in bytes")

	flag.Parse()

	// set default values
	if cfg.RepoCSV == "" {
		cfg.RepoCSV = DefaultCSVFile
	}
	
	if cfg.CloneDir == "" {
		cfg.CloneDir = DefaultCloneDir
	}
	
	if cfg.Username == "" || cfg.Token == "" {
		log.Fatal("Username and Token are required")
	}

	return cfg
}
