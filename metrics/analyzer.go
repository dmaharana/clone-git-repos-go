package metrics

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CommitMetrics represents metrics for a single commit
type CommitMetrics struct {
	Hash         string
	Author       string
	Date         time.Time
	LinesAdded   int
	LinesDeleted int
	FilesChanged int
}

// RepositoryMetrics contains aggregated metrics for the entire repository
type RepositoryMetrics struct {
	TotalCommits       int
	UniqueAuthors      map[string]bool
	CommitsByAuthor    map[string]int
	CommitsByDate      map[string]int
	AverageCommitSize  float64
	CodeChurn          map[string]int // Lines added/deleted per day
	CommitFrequency    map[string]int // Commits per day
	TimeBetweenCommits []time.Duration
}

// Analyzer handles repository metric calculations
type Analyzer struct {
	repo *git.Repository
}

// NewAnalyzer creates a new repository analyzer
func NewAnalyzer(repoPath string) (*Analyzer, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}
	return &Analyzer{repo: repo}, nil
}

// AnalyzeRepository performs a complete analysis of the repository
func (a *Analyzer) AnalyzeRepository() (*RepositoryMetrics, error) {
	metrics := &RepositoryMetrics{
		UniqueAuthors:   make(map[string]bool),
		CommitsByAuthor: make(map[string]int),
		CommitsByDate:   make(map[string]int),
		CodeChurn:       make(map[string]int),
		CommitFrequency: make(map[string]int),
	}

	ref, err := a.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository head: %w", err)
	}

	commits, err := a.repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	var lastCommitTime *time.Time
	totalLines := 0

	err = commits.ForEach(func(c *object.Commit) error {
		// Track unique authors
		metrics.UniqueAuthors[c.Author.Email] = true

		// Count commits by author
		metrics.CommitsByAuthor[c.Author.Email]++

		// Count commits by date
		dateStr := c.Author.When.Format("2006-01-02")
		metrics.CommitsByDate[dateStr]++
		metrics.CommitFrequency[dateStr]++

		// Calculate time between commits
		if lastCommitTime != nil {
			timeBetween := lastCommitTime.Sub(c.Author.When)
			metrics.TimeBetweenCommits = append(metrics.TimeBetweenCommits, timeBetween)
		}
		commitTime := c.Author.When
		lastCommitTime = &commitTime

		// Get commit stats
		stats, err := c.Stats()
		if err != nil {
			return fmt.Errorf("failed to get commit stats: %w", err)
		}

		for _, stat := range stats {
			totalLines += stat.Addition + stat.Deletion
			metrics.CodeChurn[dateStr] += stat.Addition + stat.Deletion
		}

		metrics.TotalCommits++
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to analyze commits: %w", err)
	}

	if metrics.TotalCommits > 0 {
		metrics.AverageCommitSize = float64(totalLines) / float64(metrics.TotalCommits)
	}

	return metrics, nil
}

// GetCommitFrequencyByAuthor returns the number of commits per author
func (a *Analyzer) GetCommitFrequencyByAuthor() (map[string]int, error) {
	commitsByAuthor := make(map[string]int)

	ref, err := a.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository head: %w", err)
	}

	commits, err := a.repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	err = commits.ForEach(func(c *object.Commit) error {
		commitsByAuthor[c.Author.Email]++
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to analyze commits: %w", err)
	}

	return commitsByAuthor, nil
}

// GetCodeChurnByAuthor returns the total lines changed (added + deleted) per author
func (a *Analyzer) GetCodeChurnByAuthor() (map[string]int, error) {
	churnByAuthor := make(map[string]int)

	ref, err := a.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository head: %w", err)
	}

	commits, err := a.repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	err = commits.ForEach(func(c *object.Commit) error {
		stats, err := c.Stats()
		if err != nil {
			return fmt.Errorf("failed to get commit stats: %w", err)
		}

		for _, stat := range stats {
			churnByAuthor[c.Author.Email] += stat.Addition + stat.Deletion
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to analyze commits: %w", err)
	}

	return churnByAuthor, nil
}
