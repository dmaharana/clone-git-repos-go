package metrics

import (
	"math"
	"sort"
	"time"
)

// TimeRange represents a period between two dates
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// VelocityMetrics represents the development velocity metrics
type VelocityMetrics struct {
	CommitsPerDay      float64
	AverageLinesPerDay float64
	TrendSlope         float64 // Positive means increasing velocity
}

// CalculateVelocity computes development velocity metrics over a time range
func CalculateVelocity(commits []CommitMetrics, timeRange TimeRange) VelocityMetrics {
	if len(commits) == 0 {
		return VelocityMetrics{}
	}

	// Sort commits by date
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.Before(commits[j].Date)
	})

	totalDays := timeRange.End.Sub(timeRange.Start).Hours() / 24
	if totalDays < 1 {
		totalDays = 1
	}

	totalLines := 0
	for _, commit := range commits {
		totalLines += commit.LinesAdded + commit.LinesDeleted
	}

	// Calculate trend slope using simple linear regression
	var sumX, sumY, sumXY, sumXX float64
	n := float64(len(commits))

	for i, commit := range commits {
		x := float64(i)
		y := float64(commit.LinesAdded + commit.LinesDeleted)
		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)

	return VelocityMetrics{
		CommitsPerDay:      float64(len(commits)) / totalDays,
		AverageLinesPerDay: float64(totalLines) / totalDays,
		TrendSlope:         slope,
	}
}

// CalculateCommitDistribution returns the distribution of commits across different time periods
func CalculateCommitDistribution(commits []CommitMetrics) map[string]int {
	distribution := make(map[string]int)

	for _, commit := range commits {
		hour := commit.Date.Format("15") // 24-hour format
		distribution[hour]++
	}

	return distribution
}

// CalculateProductivityScore computes a normalized productivity score (0-100)
func CalculateProductivityScore(metrics *RepositoryMetrics) float64 {
	if metrics.TotalCommits == 0 {
		return 0
	}

	// Factors to consider in productivity score
	commitFrequency := float64(metrics.TotalCommits)
	authorCount := float64(len(metrics.UniqueAuthors))
	averageCommitSize := metrics.AverageCommitSize

	// Normalize each factor (these values can be adjusted based on project needs)
	normalizedFrequency := math.Min(commitFrequency/100, 1.0)
	normalizedAuthors := math.Min(authorCount/10, 1.0)
	normalizedSize := math.Min(averageCommitSize/500, 1.0)

	// Calculate weighted score
	score := (normalizedFrequency * 0.4) + (normalizedAuthors * 0.3) + (normalizedSize * 0.3)

	// Convert to 0-100 scale
	return score * 100
}

// CalculateAverageTimeBetweenCommits computes the mean time between commits
func CalculateAverageTimeBetweenCommits(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	var total time.Duration
	for _, d := range durations {
		total += d
	}

	return total / time.Duration(len(durations))
}
