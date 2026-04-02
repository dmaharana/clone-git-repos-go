package git

import (
	"log"
	"os"
	"strings"

	"github.com/dmaharana/clone-git-repo/internal/repostatus"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

const (
	gitOrigin = "origin"
)

// CloneRepo clones a Git repository and checks out all its branches
func CloneRepo(url string, dir string, rs *repostatus.RepoStatus, username, token string) error {
	var auth transport.AuthMethod
	cloneURL := url

	// If token is provided, default to HTTPS and use authentication
	if token != "" {
		if strings.HasPrefix(url, "git@") {
			// Convert SSH to HTTPS
			// git@github.com:user/repo.git -> https://github.com/user/repo.git
			cloneURL = strings.Replace(url, ":", "/", 1)
			cloneURL = strings.Replace(cloneURL, "git@", "https://", 1)
		} else if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
			// If it doesn't have a protocol, assume it's a path and try to make it HTTPS if it looks like one
			// This is a bit speculative, but common for some inputs
			if strings.Contains(url, "github.com") || strings.Contains(url, "gitlab.com") || strings.Contains(url, "bitbucket.org") {
				cloneURL = "https://" + url
			}
		}

		if strings.HasPrefix(cloneURL, "https://") || strings.HasPrefix(cloneURL, "http://") {
			auth = &http.BasicAuth{
				Username: username,
				Password: token,
			}
		}
	}

	// clone repo
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      cloneURL,
		Auth:     auth,
		Progress: os.Stdout,
	})

	if err != nil {
		// Do not log the error here if it's authentication related,
		// let the caller handle it to avoid logging sensitive URLs if any
		return err
	}

	// checkout one branch at a time and search for terms
	bList, err := findAllBranches(r)
	if err != nil {
		log.Println("Error getting branches:", err)
		return err
	}

	log.Println("Branches: ", bList)

	// update repo status
	rs.IsCloned = true
	rs.BranchCount = len(bList)

	log.Printf("Repository cloned to %s\n", dir)

	// list all tags
	tList, err := findAllTags(r)
	if err != nil {
		log.Println("Error getting tags:", err)
	}

	log.Println("Tags: ", tList)

	// update repo status
	rs.TagCount = len(tList)

	// set up worktree
	w, err := r.Worktree()
	if err != nil {
		log.Println("Error getting worktree:", err)
		return err
	}

	// checkout all branches
	for _, branch := range bList {
		log.Println("Checking out branch: ", branch)
		// replace "refs/remotes/origin/" at the beginning of the remote branch name with blank
		localBranch := strings.Replace(branch, "refs/remotes/origin/", "", 1)

		w.Pull(&git.PullOptions{RemoteName: gitOrigin})

		// checkout the branch
		log.Println("Checking out branch: ", localBranch)
		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(localBranch),
			Create: true, // Create the branch if it doesn't exist locally
			Force:  true, // Force checkout
		})
		if err != nil {
			log.Println("Error checking out branch: ", err)
			continue
		}
	}

	return nil
}

func findAllBranches(r *git.Repository) ([]string, error) {
	log.Println("Branches: ")
	branches, err := r.References()
	if err != nil {
		return nil, nil
	}

	branchList := make([]string, 0)
	count := 0

	branches.ForEach(func(b *plumbing.Reference) error {
		if b.Type() != plumbing.HashReference {
			return nil
		}

		bname := b.Name().String()
		if strings.Contains(bname, "origin") {
			count++
			branchList = append(branchList, bname)
		}
		return nil
	})

	log.Printf("Total branch(es): %d\n", count)

	return branchList, nil
}

// find all tags
func findAllTags(r *git.Repository) ([]string, error) {
	log.Println("Tags: ")
	tags, err := r.Tags()
	if err != nil {
		return nil, nil
	}

	tagList := make([]string, 0)
	count := 0

	tags.ForEach(func(t *plumbing.Reference) error {
		if t.Type() != plumbing.HashReference {
			return nil
		}

		tname := t.Name().String()
		count++
		tagList = append(tagList, tname)

		return nil
	})

	log.Printf("Total tag(s): %d\n", count)

	return tagList, nil
}

// CreateDirectoryIfNotExist creates a directory if it doesn't exist
func CreateDirectoryIfNotExist(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
