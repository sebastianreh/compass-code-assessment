package contact_test

import (
	"errors"
	"github.com/sebastianreh/internal/contact"
	"github.com/sebastianreh/mocks"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestContactRepository_GetContactData(t *testing.T) {
	logger := logrus.New()

	t.Run("when reading contact data successfully, it should return contacts", func(t *testing.T) {
		mockCsv := new(mocks.CsvMock)
		repo := contact.NewContactRepository(logger, mockCsv)
		mockedCSVData := [][]string{
			{"ContactID", "FirstName", "LastName", "Email", "ZipCode", "Address"},
			{"1", "John", "Doe", "john@example.com", "12345", "123 Main St"},
			{"2", "Jane", "Smith", "jane@example.com", "54321", "456 Oak St"},
		}

		mockCsv.On("ReadCSV", "files/input.csv").Return(mockedCSVData, nil)

		contacts, err := repo.GetContactData()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(contacts))
		assert.Equal(t, "John", contacts[0].FirstName)
		assert.Equal(t, "Doe", contacts[0].LastName)

		mockCsv.AssertExpectations(t)
	})

	t.Run("when CSV read fails, it should return an error", func(t *testing.T) {
		mockCsv := new(mocks.CsvMock)
		repo := contact.NewContactRepository(logger, mockCsv)
		mockCsv.On("ReadCSV", "files/input.csv").Return([][]string{}, errors.New("failed to read CSV"))

		_, err := repo.GetContactData()

		assert.NotNil(t, err)
		assert.Equal(t, "failed to read CSV", err.Error())

		mockCsv.AssertExpectations(t)
	})
}

func TestContactRepository_WriteContactData(t *testing.T) {
	logger := logrus.New()

	t.Run("when writing contact data successfully, it should not return an error", func(t *testing.T) {
		mockCsv := new(mocks.CsvMock)
		repo := contact.NewContactRepository(logger, mockCsv)
		mockedHeader := []string{"ContactIDSource", "ContactIDMatch", "Accuracy"}
		mockedData := [][]string{
			{"1", "2", "High"},
		}
		mockCsv.On("WriteCSV", "files/output.csv", mockedHeader, mockedData).Return(nil)

		output := []contact.ProcessOutput{
			{ContactIDSource: "1", ContactIDMatch: "2", AccuracyLevel: 4},
		}

		err := repo.WriteContactData(output)
		assert.Nil(t, err)

		mockCsv.AssertExpectations(t)
	})

	t.Run("when writing contact data fails, it should return an error", func(t *testing.T) {
		mockCsv := new(mocks.CsvMock)
		repo := contact.NewContactRepository(logger, mockCsv)
		mockedHeader := []string{"ContactIDSource", "ContactIDMatch", "Accuracy"}
		mockedData := [][]string{
			{"1", "2", "High"},
		}
		mockCsv.On("WriteCSV", "files/output.csv", mockedHeader, mockedData).Return(errors.New("failed to write CSV"))

		output := []contact.ProcessOutput{
			{ContactIDSource: "1", ContactIDMatch: "2", AccuracyLevel: 4},
		}

		err := repo.WriteContactData(output)
		assert.NotNil(t, err)
		assert.Equal(t, "failed to write CSV", err.Error())

		mockCsv.AssertExpectations(t)
	})
}

func TestContactRepository_WriteDuplicateContacts(t *testing.T) {
	logger := logrus.New()

	t.Run("when writing duplicate contacts successfully, it should not return an error", func(t *testing.T) {
		mockCsv := new(mocks.CsvMock)
		repo := contact.NewContactRepository(logger, mockCsv)
		mockedHeader := []string{"ContactID", "FirstName", "LastName", "Email", "ZipCode", "Address"}
		mockedData := [][]string{
			{"1", "John", "Doe", "john@example.com", "12345", "123 Main St"},
		}
		mockCsv.On("WriteCSV", "files/duplicate.csv", mockedHeader, mockedData).Return(nil)

		duplicates := []contact.Contact{
			{ContactID: "1", FirstName: "John", LastName: "Doe", Email: "john@example.com", ZipCode: "12345", Address: "123 Main St"},
		}

		err := repo.WriteDuplicateContacts(duplicates)
		assert.Nil(t, err)

		mockCsv.AssertExpectations(t)
	})

	t.Run("when writing duplicate contacts fails, it should return an error", func(t *testing.T) {
		mockCsv := new(mocks.CsvMock)
		repo := contact.NewContactRepository(logger, mockCsv)
		mockedHeader := []string{"ContactID", "FirstName", "LastName", "Email", "ZipCode", "Address"}
		mockedData := [][]string{
			{"1", "John", "Doe", "john@example.com", "12345", "123 Main St"},
		}
		mockCsv.On("WriteCSV", "files/duplicate.csv", mockedHeader, mockedData).Return(errors.New("failed to write CSV"))

		duplicates := []contact.Contact{
			{ContactID: "1", FirstName: "John", LastName: "Doe", Email: "john@example.com", ZipCode: "12345", Address: "123 Main St"},
		}

		err := repo.WriteDuplicateContacts(duplicates)
		assert.NotNil(t, err)
		assert.Equal(t, "failed to write CSV", err.Error())

		mockCsv.AssertExpectations(t)
	})
}
