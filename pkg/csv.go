package pkg

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVConnector interface {
	ReadCSV(filePath string) ([][]string, error)
	WriteCSV(filePath string, header []string, data [][]string) error
}

type csvConnector struct{}

func NewCSVConnector() CSVConnector {
	return &csvConnector{}
}

func (csvConnector) ReadCSV(filePath string) ([][]string, error) {
	output := make([][]string, 0)
	file, err := os.Open(filePath)
	if err != nil {
		return output, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return output, fmt.Errorf("error reading CSV file: %w", err)
	}

	return records, nil
}

func (csvConnector) WriteCSV(filePath string, header []string, data [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if len(header) > 0 {
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("error writing header: %w", err)
		}
	}

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	return nil
}
