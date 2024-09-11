package pkg_test

import (
	"encoding/csv"
	"github.com/sebastianreh/compass-code-assessment/pkg"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVConnector_ReadCSV(t *testing.T) {
	connector := pkg.NewCSVConnector()

	t.Run("when reading CSV file successfully, it should return data", func(t *testing.T) {
		content := "header1,header2\nvalue1,value2\n"
		tempFile := createTempCSVFile(t, content)
		defer os.Remove(tempFile.Name())

		data, err := connector.ReadCSV(tempFile.Name())

		assert.Nil(t, err)
		assert.Equal(t, 2, len(data))
		assert.Equal(t, []string{"header1", "header2"}, data[0])
		assert.Equal(t, []string{"value1", "value2"}, data[1])
	})

	t.Run("when the file does not exist, it should return an error", func(t *testing.T) {
		_, err := connector.ReadCSV("non_existent_file.csv")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error opening file")
	})

	t.Run("when reading invalid CSV content, it should return an error", func(t *testing.T) {
		tempFile := createTempCSVFile(t, "header1,header2\nvalue1")
		defer os.Remove(tempFile.Name())

		_, err := connector.ReadCSV(tempFile.Name())
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error reading CSV file")
	})
}

func TestCSVConnector_WriteCSV(t *testing.T) {
	connector := pkg.NewCSVConnector()

	t.Run("when writing CSV file successfully, it should write data", func(t *testing.T) {
		tempFile, err := ioutil.TempFile("", "test_write_csv_*.csv")
		assert.Nil(t, err)
		defer os.Remove(tempFile.Name())

		header := []string{"header1", "header2"}
		data := [][]string{
			{"value1", "value2"},
		}

		err = connector.WriteCSV(tempFile.Name(), header, data)
		assert.Nil(t, err)

		fileContent := readCSVFile(t, tempFile.Name())
		assert.Equal(t, [][]string{{"header1", "header2"}, {"value1", "value2"}}, fileContent)
	})

	t.Run("when the file cannot be created, it should return an error", func(t *testing.T) {
		err := connector.WriteCSV("/invalid_path/output.csv", []string{"header"}, [][]string{{"value1"}})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error creating file")
	})

	t.Run("when writing CSV with empty header, it should write only the data", func(t *testing.T) {
		tempFile, err := ioutil.TempFile("", "test_write_csv_*.csv")
		assert.Nil(t, err)
		defer os.Remove(tempFile.Name())

		data := [][]string{
			{"value1", "value2"},
		}

		err = connector.WriteCSV(tempFile.Name(), nil, data)
		assert.Nil(t, err)

		fileContent := readCSVFile(t, tempFile.Name())
		assert.Equal(t, [][]string{{"value1", "value2"}}, fileContent)
	})
}

// Helper function to create a temporary CSV
func createTempCSVFile(t *testing.T, content string) *os.File {
	tempFile, err := ioutil.TempFile("", "test_csv_*.csv")
	assert.Nil(t, err)

	_, err = tempFile.WriteString(content)
	assert.Nil(t, err)

	err = tempFile.Close()
	assert.Nil(t, err)

	return tempFile
}

// Helper function to read content from a CSV file
func readCSVFile(t *testing.T, filePath string) [][]string {
	file, err := os.Open(filePath)
	assert.Nil(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	assert.Nil(t, err)

	return records
}
