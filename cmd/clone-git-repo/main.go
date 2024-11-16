package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/your-username/clone-git-repo/internal/pkg/config"
	"github.com/your-username/clone-git-repo/internal/pkg/csv"
	"github.com/your-username/clone-git-repo/internal/pkg/git"
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

	AuthenticationErrorString = "authentication required"
	DirectoryExistsErrorString = "repository already exists"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Print version information
	if Version != "" {
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Git Branch: %s\n", GitBranch)
	}

	// Parse command-line flags
	cfg := config.ParseFlags()

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
		

		errCount := 0
		maxRetries := 3
		errorCode := cloneRepository(url, repoDir)
		for {
			if errorCode != 0 {
				errCount++
				if errCount > maxRetries {
					log.Printf("Error cloning repository %s: %v\n", url, err)
					break
				}
			} else {
				break
			}
			switch errorCode {
			case AuthenticationError:
				errorCode = handleAuthenticationError(url, repoDir, cfg)
			case DirectoryExistsError:
				errorCode = handleDirectoryExistsError(url, repoDir)
			default:
				break
			}
		}

	}
}

// perform git clone and return error
func cloneRepository(url string, repoDir string) int {
	err := git.CloneRepo(url, repoDir)

	if err != nil {
		return checkError(err)
	}

	return 0
}

func checkError (err error) int {
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
func handleAuthenticationError(url string, repoDir string, cfg *config.Config) int {
	// Add username and token to the URL
	// remove the "https://" prefix
	urlWithCredentials := fmt.Sprintf("%s://%s:%s@%s", "https", cfg.Username, cfg.Token, strings.TrimPrefix(url, "https://"))
	if err := git.CloneRepo(urlWithCredentials, repoDir); err != nil {
		log.Printf("Error cloning repository %s: %v\n", url, err)
		return checkError(err)
	}

	return 0
}

// handle if directory already exists, remove it and try again
func handleDirectoryExistsError(url string, repoDir string) int {
	// Remove the partially cloned directory
	os.RemoveAll(repoDir)

	// Clone the repository into the directory
	if err := git.CloneRepo(url, repoDir); err != nil {
		log.Printf("Error cloning repository %s: %v\n", url, err)
		return checkError(err)
	}

	return 0
}