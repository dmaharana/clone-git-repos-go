package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmaharana/clone-git-repo/internal/pkg/config"
	"github.com/dmaharana/clone-git-repo/internal/pkg/csv"
	"github.com/dmaharana/clone-git-repo/internal/pkg/git"
	"github.com/dmaharana/clone-git-repo/internal/pkg/logger"
	"github.com/dmaharana/clone-git-repo/internal/repostatus"
)

// Build information. Populated at build-time.
var (
	Version   string
	BuildTime string
	GitCommit string
	GitBranch string
)

// create enum for exit codes
const (
	CreateDirectoryError = iota + 1
	AuthenticationError
	DirectoryExistsError
	UnknownError

	AuthenticationErrorString   = "authentication required"
	DirectoryExistsErrorString  = "repository already exists"
	ResultFileName              = "clone-git-repo-result.csv"
	DirectoryExistsErrorMessage = "Directory already exists, remove it and try again"

	MaxRetries = 3
)

var log *logger.Logger

func main() {
	// Parse command-line flags
	cfg := config.ParseFlags()

	// Initialize logger
	logCfg := &logger.Config{
		LogDir:     cfg.LogDir,
		MaxSize:    cfg.LogMaxSize,
		TimeFormat: "2006-01-02",
	}

	var err error
	log, err = logger.New(logCfg)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	// Set log flags
	log.SetFlags(logger.LstdFlags | logger.Lshortfile)

	// Print version information
	if Version != "" {
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Git Branch: %s\n", GitBranch)
	}

	// Read Git URLs from CSV file
	repositoryURLs, err := csv.ReadRepositoryURLs(cfg.RepoCSV)
	if err != nil {
		log.Fatal(err)
	}

	// create array to hold clone status
	cloneStatus := []*repostatus.RepoStatus{}

	// Clone each repository into a separate directory
	for _, url := range repositoryURLs {
		rs := &repostatus.RepoStatus{
			RepoPath: url,
		}
		// Generate a unique directory name based on the repository URL
		repoName := filepath.Base(url)
		repoDir := filepath.Join(cfg.CloneDir, repoName)

		errCount := 0
		errorCode := cloneRepository(url, repoDir, rs, cfg)
		for {
			if errorCode != 0 {
				errCount++
				if errCount > MaxRetries {
					log.Printf("Error cloning repository %s: code %d\n", url, errorCode)
					break
				}
			} else {
				break
			}
			switch errorCode {
			case AuthenticationError:
				errorCode = handleAuthenticationError(url, repoDir, cfg, rs)
			case DirectoryExistsError:
				errorCode = handleDirectoryExistsError(url, repoDir, cfg, rs)
			default:
				// Exit the loop for unknown errors
				errCount = MaxRetries + 1
			}
		}

		rs.IsCloned = errorCode == 0
		if !rs.IsCloned {
			rs.BranchCount = 0
			rs.TagCount = 0
		}
		cloneStatus = append(cloneStatus, rs)
	}

	// Print status table
	repostatus.PrintStatusTable(cloneStatus)

	// Write status to CSV file
	if err := repostatus.WriteStatusToCSV(cloneStatus, ResultFileName); err != nil {
		log.Printf("Error writing status to CSV file: %v\n", err)
	}
}

// perform git clone and return error
func cloneRepository(url string, repoDir string, rs *repostatus.RepoStatus, cfg *config.Config) int {
	err := git.CloneRepo(url, repoDir, rs, cfg.Username, cfg.Token)
	if err != nil {
		return checkError(err)
	}
	return 0
}

func checkError(err error) int {
	if err == nil {
		return 0
	}
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "auth") || strings.Contains(errStr, "authentication") || strings.Contains(errStr, "credentials") {
		return AuthenticationError
	}
	if strings.Contains(errStr, "exists") || strings.Contains(errStr, "already exists") {
		return DirectoryExistsError
	}
	return UnknownError
}

// handle authentication error
func handleAuthenticationError(rurl string, repoDir string, cfg *config.Config, rs *repostatus.RepoStatus) int {
	// If authentication failed even with a token, we stop retrying
	if cfg.Token != "" {
		return AuthenticationError
	}

	return AuthenticationError
}

// handle if directory already exists, remove it and try again
func handleDirectoryExistsError(url string, repoDir string, cfg *config.Config, rs *repostatus.RepoStatus) int {
	// Remove the partially cloned directory
	os.RemoveAll(repoDir)

	// Clone the repository into the directory
	return cloneRepository(url, repoDir, rs, cfg)
}

// go through all the repositories in the clonedir and save the developer productivity metrics in a csv file
func calculateMetrics(cfg *config.Config) {
	// Read Git URLs from CSV file
	repositoryURLs, err := csv.ReadRepositoryURLs(cfg.RepoCSV)
	if err != nil {
		log.Fatal(err)
	}

	// Clone each repository into a separate directory
	for _, url := range repositoryURLs {
		// Generate a unique directory name based on the repository URL
		repoName := filepath.Base(url)
		repoDir := filepath.Join(cfg.CloneDir, repoName)
		// Clone the repository into the directory
		errorCode := cloneRepository(url, repoDir, nil, cfg)
		if errorCode != 0 {
			log.Printf("Error cloning repository %s: code %d\n", url, errorCode)
			os.Exit(1)
		}
	}
}
