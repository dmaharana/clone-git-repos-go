package repostatus

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// RepoStatus represents the status information of a Git repository
type RepoStatus struct {
	RepoPath     string
	IsCloned     bool
	BranchCount  int
	TagCount     int
}

// GetRepoStatus retrieves the status information for a given repository path
func GetRepoStatus(repoPath string) (*RepoStatus, error) {
	status := &RepoStatus{
		RepoPath: repoPath,
	}

	// Check if repository is cloned by verifying .git directory
	if _, err := os.Stat(repoPath + "/.git"); err == nil {
		status.IsCloned = true
	}

	if !status.IsCloned {
		return status, nil
	}

	// Get branch count
	branchCmd := exec.Command("git", "branch", "-a")
	branchCmd.Dir = repoPath
	branchOutput, err := branchCmd.Output()
	if err == nil {
		branches := strings.Split(string(branchOutput), "\n")
		// Remove empty lines
		count := 0
		for _, branch := range branches {
			if strings.TrimSpace(branch) != "" {
				count++
			}
		}
		status.BranchCount = count
	}

	// Get tag count
	tagCmd := exec.Command("git", "tag")
	tagCmd.Dir = repoPath
	tagOutput, err := tagCmd.Output()
	if err == nil {
		tags := strings.Split(string(tagOutput), "\n")
		// Remove empty lines
		count := 0
		for _, tag := range tags {
			if strings.TrimSpace(tag) != "" {
				count++
			}
		}
		status.TagCount = count
	}

	return status, nil
}

// PrintStatusTable prints a table with repository status information
func PrintStatusTable(statuses []*RepoStatus) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Repository", "Cloned", "Branches", "Tags"})

	for _, status := range statuses {
		clonedStatus := "No"
		if status.IsCloned {
			clonedStatus = "Yes"
		}
		
		table.Append([]string{
			status.RepoPath,
			clonedStatus,
			fmt.Sprintf("%d", status.BranchCount),
			fmt.Sprintf("%d", status.TagCount),
		})
	}

	table.Render()
}

// write slice of RepoStatus to a CSV file
func WriteStatusToCSV(statuses []*RepoStatus, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header
	header := []string{"Repository", "Cloned", "Branches", "Tags"}
	writer := csv.NewWriter(file)
	err = writer.Write(header)
	if err != nil {
		return err
	}
	// writer.Flush()

	// writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, status := range statuses {
		err := writer.Write([]string{
			status.RepoPath,
			fmt.Sprintf("%t", status.IsCloned),
			fmt.Sprintf("%d", status.BranchCount),
			fmt.Sprintf("%d", status.TagCount),
		})
		if err != nil {
			return err
		}
	}
	return nil
}