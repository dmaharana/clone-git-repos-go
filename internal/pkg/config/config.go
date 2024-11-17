package config

import (
	"flag"
	"log"

	"gopkg.in/ini.v1"
)

// Config holds the application configuration
type Config struct {
	RepoCSV    string
	CloneDir   string
	Username   string
	Token      string
	LogDir     string
	LogMaxSize int64
}

const (
	DefaultCSVFile    = "repositories.csv"
	DefaultCloneDir   = "clonedir"
	DefaultConfigFile = "config.ini"
)

// ParseFlags parses command line flags and config file, returns a Config struct
func ParseFlags() *Config {
	cfg := &Config{}
	var configFile string

	// Only parse the config file path from command line
	flag.StringVar(&configFile, "c", DefaultConfigFile, "Path to config file")
	flag.Parse()

	// Load config file
	iniFile, err := ini.Load(configFile)
	if err != nil {
		log.Printf("Warning: Could not load config file: %v\n", err)
		log.Printf("Using command line arguments instead\n")
		return parseCommandLineArgs()
	}

	// Read from config file
	credentials := iniFile.Section("credentials")
	cfg.Username = credentials.Key("username").String()
	cfg.Token = credentials.Key("token").String()

	paths := iniFile.Section("paths")
	cfg.RepoCSV = paths.Key("csv_file").MustString(DefaultCSVFile)
	cfg.CloneDir = paths.Key("clone_dir").MustString(DefaultCloneDir)

	logging := iniFile.Section("logging")
	cfg.LogDir = logging.Key("log_dir").MustString("logs")
	cfg.LogMaxSize = logging.Key("log_max_size").MustInt64(10 * 1024 * 1024)

	// Validate required fields
	if cfg.Username == "" || cfg.Token == "" {
		log.Fatal("Username and Token are required in config file")
	}

	return cfg
}

// parseCommandLineArgs parses command line arguments as fallback
func parseCommandLineArgs() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.RepoCSV, "f", DefaultCSVFile, "CSV file")
	flag.StringVar(&cfg.CloneDir, "d", DefaultCloneDir, "Clone directory")
	flag.StringVar(&cfg.Username, "u", "", "Username")
	flag.StringVar(&cfg.Token, "t", "", "Token")
	flag.StringVar(&cfg.LogDir, "logdir", "logs", "Log directory")
	flag.Int64Var(&cfg.LogMaxSize, "logsize", 10*1024*1024, "Maximum log file size in bytes")

	flag.Parse()

	if cfg.Username == "" || cfg.Token == "" {
		log.Fatal("Username and Token are required")
	}

	return cfg
}
