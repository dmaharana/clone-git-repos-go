package main

import (
	"fmt"
	"net/url"
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

	AuthenticationErrorString    = "authentication required"
	DirectoryExistsErrorString  = "repository already exists"
	ResultFileName               = "clone-git-repo-result.csv"
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
		errorCode := cloneRepository(url, repoDir, rs)
		for {
			if errorCode != 0 {
				errCount++
				if errCount > MaxRetries {
					log.Printf("Error cloning repository %s: %v\n", url, err)
					break
				}
			} else {
				break
			}
			switch errorCode {
			case AuthenticationError:
				errorCode = handleAuthenticationError(url, repoDir, cfg, rs)
			case DirectoryExistsError:
				errorCode = handleDirectoryExistsError(url, repoDir, rs)
			default:
				break
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
func cloneRepository(url string, repoDir string, rs *repostatus.RepoStatus) int {
	err := git.CloneRepo(url, repoDir, rs)
	if err != nil {
		return checkError(err)
	}
	return 0
}

func checkError(err error) int {
	switch err.Error() {
	case AuthenticationErrorString:
		return AuthenticationError
	case DirectoryExistsErrorString:
		return DirectoryExistsError
	default:
		return UnknownError
	}
}

// handle authentication error
func handleAuthenticationError(url string, repoDir string, cfg *config.Config, rs *repostatus.RepoStatus) int {
	// Add username and token to the URL
	// escape special characters in username and token
	escapedUsername := url.PathEscape(cfg.Username)
	escapedToken := url.PathEscape(cfg.Token)

	// remove the "https://" prefix
	urlWithCredentials := fmt.Sprintf("%s://%s:%s@%s", "https", escapedUsername, escapedToken, strings.TrimPrefix(url, "https://"))
	if err := git.CloneRepo(urlWithCredentials, repoDir, rs); err != nil {
		log.Printf("Error cloning repository %s: %v\n", url, err)
		return checkError(err)
	}
	return 0
}

// handle if directory already exists, remove it and try again
func handleDirectoryExistsError(url string, repoDir string, rs *repostatus.RepoStatus) int {
	// Remove the partially cloned directory
	os.RemoveAll(repoDir)

	// Clone the repository into the directory
	if err := git.CloneRepo(url, repoDir, rs); err != nil {
		log.Printf("Error cloning repository %s: %v\n", url, err)
		return checkError(err)
	}
	return 0
}