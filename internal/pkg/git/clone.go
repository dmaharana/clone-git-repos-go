package git

import (
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	gitOrigin = "origin"
)

// CloneRepo clones a Git repository and checks out all its branches
func CloneRepo(url string, dir string) error {
	// clone repo
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Println("Error cloning repo:", err)
		return err
	}

	// checkout one branch at a time and search for terms
	bList, err := findAllBranches(r)
	if err != nil {
		log.Println("Error getting branches:", err)
		return err
	}

	log.Println("Branches: ", bList)

	log.Printf("Repository cloned to %s\n", dir)

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

// CreateDirectoryIfNotExist creates a directory if it doesn't exist
func CreateDirectoryIfNotExist(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
