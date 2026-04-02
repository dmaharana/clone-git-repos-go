package metrics

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// MetricsExporter handles the export of repository metrics to CSV
type MetricsExporter struct {
	outputPath string
}

// NewMetricsExporter creates a new metrics exporter
func NewMetricsExporter(outputPath string) *MetricsExporter {
	return &MetricsExporter{
		outputPath: outputPath,
	}
}

// ExportMetricsToCSV writes repository metrics to a CSV file
func (e *MetricsExporter) ExportMetricsToCSV(repoPath string, metrics *RepositoryMetrics, velocity VelocityMetrics) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(e.outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(e.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{
		"Repository",
		"Analysis Date",
		"Total Commits",
		"Unique Authors",
		"Average Commit Size",
		"Commits Per Day",
		"Lines Per Day",
		"Velocity Trend",
		"Productivity Score",
		"Average Time Between Commits (hours)",
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	// Calculate additional metrics
	productivityScore := CalculateProductivityScore(metrics)
	avgTimeBetweenCommits := CalculateAverageTimeBetweenCommits(metrics.TimeBetweenCommits)

	// Prepare row data
	row := []string{
		repoPath,
		time.Now().Format("2006-01-02 15:04:05"),
		strconv.Itoa(metrics.TotalCommits),
		strconv.Itoa(len(metrics.UniqueAuthors)),
		fmt.Sprintf("%.2f", metrics.AverageCommitSize),
		fmt.Sprintf("%.2f", velocity.CommitsPerDay),
		fmt.Sprintf("%.2f", velocity.AverageLinesPerDay),
		fmt.Sprintf("%.2f", velocity.TrendSlope),
		fmt.Sprintf("%.2f", productivityScore),
		fmt.Sprintf("%.2f", avgTimeBetweenCommits.Hours()),
	}

	if err := writer.Write(row); err != nil {
		return fmt.Errorf("failed to write metrics row: %w", err)
	}

	// Export detailed author metrics to separate files
	if err := e.exportAuthorMetrics(repoPath, metrics); err != nil {
		return fmt.Errorf("failed to export author metrics: %w", err)
	}

	if err := e.exportTimeBasedMetrics(repoPath, metrics); err != nil {
		return fmt.Errorf("failed to export time-based metrics: %w", err)
	}

	return nil
}

// exportAuthorMetrics exports per-author statistics to a separate CSV file
func (e *MetricsExporter) exportAuthorMetrics(repoPath string, metrics *RepositoryMetrics) error {
	authorFile := strings.TrimSuffix(e.outputPath, ".csv") + "_authors.csv"
	file, err := os.Create(authorFile)
	if err != nil {
		return fmt.Errorf("failed to create author metrics file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{"Author", "Commit Count", "Contribution Percentage"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write author headers: %w", err)
	}

	// Get sorted list of authors
	authors := make([]string, 0, len(metrics.CommitsByAuthor))
	for author := range metrics.CommitsByAuthor {
		authors = append(authors, author)
	}
	sort.Strings(authors)

	// Write author metrics
	for _, author := range authors {
		commitCount := metrics.CommitsByAuthor[author]
		percentage := float64(commitCount) / float64(metrics.TotalCommits) * 100

		row := []string{
			author,
			strconv.Itoa(commitCount),
			fmt.Sprintf("%.2f", percentage),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write author row: %w", err)
		}
	}

	return nil
}

// exportTimeBasedMetrics exports time-based metrics to a separate CSV file
func (e *MetricsExporter) exportTimeBasedMetrics(repoPath string, metrics *RepositoryMetrics) error {
	timeFile := strings.TrimSuffix(e.outputPath, ".csv") + "_timeline.csv"
	file, err := os.Create(timeFile)
	if err != nil {
		return fmt.Errorf("failed to create time-based metrics file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{"Date", "Commit Count", "Code Churn"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write timeline headers: %w", err)
	}

	// Get sorted list of dates
	dates := make([]string, 0, len(metrics.CommitsByDate))
	for date := range metrics.CommitsByDate {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Write time-based metrics
	for _, date := range dates {
		row := []string{
			date,
			strconv.Itoa(metrics.CommitsByDate[date]),
			strconv.Itoa(metrics.CodeChurn[date]),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write timeline row: %w", err)
		}
	}

	return nil
}

// ExportMultiRepoMetrics exports metrics for multiple repositories to CSV
func ExportMultiRepoMetrics(repos []string, outputPath string) error {
	exporter := NewMetricsExporter(outputPath)

	for _, repoPath := range repos {
		analyzer, err := NewAnalyzer(repoPath)
		if err != nil {
			return fmt.Errorf("failed to create analyzer for %s: %w", repoPath, err)
		}

		metrics, err := analyzer.AnalyzeRepository()
		if err != nil {
			return fmt.Errorf("failed to analyze repository %s: %w", repoPath, err)
		}

		// Calculate velocity metrics
		timeRange := TimeRange{
			Start: time.Now().AddDate(0, -1, 0), // Last month
			End:   time.Now(),
		}

		commits := make([]CommitMetrics, 0)
		velocity := CalculateVelocity(commits, timeRange)

		if err := exporter.ExportMetricsToCSV(repoPath, metrics, velocity); err != nil {
			return fmt.Errorf("failed to export metrics for %s: %w", repoPath, err)
		}
	}

	return nil
}
