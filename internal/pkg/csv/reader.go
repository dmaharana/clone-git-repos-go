package csv

import (
	"encoding/csv"
	"os"
)

// ReadRepositoryURLs reads repository URLs from a CSV file
func ReadRepositoryURLs(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var repositoryURLs []string
	for _, record := range records[1:] { // Skip the header
		repositoryURLs = append(repositoryURLs, record[0])
	}

	return repositoryURLs, nil
}
