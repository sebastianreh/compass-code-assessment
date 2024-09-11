package mocks

import (
	"github.com/sebastianreh/compass-code-assessment/internal/contact"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) GetContactData() ([]contact.Contact, error) {
	args := m.Called()
	return args.Get(0).([]contact.Contact), args.Error(1)
}

func (m *RepositoryMock) WriteContactData(data []contact.ProcessOutput) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *RepositoryMock) WriteDuplicateContacts(duplicate []contact.Contact) error {
	args := m.Called(duplicate)
	return args.Error(0)
}
