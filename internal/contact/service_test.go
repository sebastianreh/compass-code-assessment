package contact_test

import (
	"errors"
	"github.com/sebastianreh/compass-code-assessment/internal/contact"
	"github.com/sebastianreh/mocks"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContactService_Evaluate(t *testing.T) {
	logger := logrus.New()

	t.Run("when evaluating contacts successfully, it should return process output", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryMock)
		mockContacts := []contact.Contact{
			{ContactID: "1", FirstName: "John", LastName: "Doe", Email: "john@example.com", ZipCode: "12345", Address: "123 Main St"},
			{ContactID: "2", FirstName: "Jane", LastName: "Doe", Email: "jane@example.com", ZipCode: "54321", Address: "456 Oak St"},
		}

		mockRepo.On("GetContactData").Return(mockContacts, nil)
		mockRepo.On("WriteContactData", mock.Anything).Return(nil)
		mockRepo.On("WriteDuplicateContacts", mock.Anything).Return(nil)

		service := contact.NewContactService(logger, mockRepo)
		results, err := service.Evaluate()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(results))
		assert.Equal(t, "1", results[0].ContactIDSource)
		assert.Equal(t, "2", results[0].ContactIDMatch)

		mockRepo.AssertExpectations(t)
	})

	t.Run("when repository returns error while getting contacts, it should return an error", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryMock)
		mockRepo.On("GetContactData").Return([]contact.Contact{}, errors.New("error fetching contacts"))

		service := contact.NewContactService(logger, mockRepo)
		_, err := service.Evaluate()

		assert.NotNil(t, err)
		assert.Equal(t, "error fetching contacts", err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("when repository returns error while writing contacts, it should log an error", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryMock)
		mockContacts := []contact.Contact{
			{ContactID: "1", FirstName: "John", LastName: "Doe", Email: "john@example.com", ZipCode: "12345", Address: "123 Main St"},
			{ContactID: "2", FirstName: "Jane", LastName: "Doe", Email: "jane@example.com", ZipCode: "54321", Address: "456 Oak St"},
		}

		mockRepo.On("GetContactData").Return(mockContacts, nil)
		mockRepo.On("WriteContactData", mock.Anything).Return(errors.New("failed to write contact data"))
		mockRepo.On("WriteDuplicateContacts", mock.Anything).Return(nil)

		service := contact.NewContactService(logger, mockRepo)
		_, err := service.Evaluate()

		assert.NotNil(t, err)
		assert.Equal(t, "failed to write contact data", err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("when repository returns error while writing duplicate contacts, it should log an error", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryMock)
		mockContacts := []contact.Contact{
			{ContactID: "1", FirstName: "John", LastName: "Doe", Email: "john@example.com", ZipCode: "12345", Address: "123 Main St"},
			{ContactID: "2", FirstName: "Jane", LastName: "Doe", Email: "jane@example.com", ZipCode: "54321", Address: "456 Oak St"},
			{ContactID: "3", FirstName: "John", LastName: "Doe", Email: "john@example.com", ZipCode: "12345", Address: "123 Main St"},
		}

		mockRepo.On("GetContactData").Return(mockContacts, nil)
		mockRepo.On("WriteDuplicateContacts", mock.Anything).Return(errors.New("failed to write duplicate contacts"))

		service := contact.NewContactService(logger, mockRepo)
		_, err := service.Evaluate()

		assert.NotNil(t, err)
		assert.Equal(t, "failed to write duplicate contacts", err.Error())

		mockRepo.AssertExpectations(t)
	})
}
