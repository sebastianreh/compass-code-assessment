package mocks

import (
	"github.com/stretchr/testify/mock"
)

type CsvMock struct {
	mock.Mock
}

func (m *CsvMock) ReadCSV(filePath string) ([][]string, error) {
	args := m.Called(filePath)
	return args.Get(0).([][]string), args.Error(1)
}

func (m *CsvMock) WriteCSV(filePath string, header []string, data [][]string) error {
	args := m.Called(filePath, header, data)
	return args.Error(0)
}
