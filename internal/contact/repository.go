package contact

import (
	"fmt"
	"github.com/sebastianreh/compass-code-assessment/pkg"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type contactRepository struct {
	log *logrus.Logger
	csv pkg.CSVConnector
}

func NewContactRepository(log *logrus.Logger, csv pkg.CSVConnector) Repository {
	return &contactRepository{
		log: log,
		csv: csv,
	}
}

func (c contactRepository) GetContactData() ([]Contact, error) {
	filePath := filepath.Join("files", "input.csv")
	records, err := c.csv.ReadCSV(filePath)
	if err != nil {
		c.log.Errorf("Error reading CSV file: %v", err)
		return nil, err
	}

	contacts, err := c.parseRecords(records)
	if err != nil {
		c.log.Errorf("Error parsing records: %v", err)
		return nil, err
	}

	return contacts, nil
}

func (c contactRepository) parseRecords(records [][]string) ([]Contact, error) {
	if len(records) < 1 {
		return nil, fmt.Errorf("no records found in CSV file")
	}

	var contacts []Contact
	header := records[0]

	for i, record := range records[1:] {
		if len(record) != len(header) {
			c.log.Errorf("record %d has a different number of fields than header", i+1)
			continue
		}

		contact := Contact{
			ContactID: record[0],
			FirstName: record[1],
			LastName:  record[2],
			Email:     record[3],
			ZipCode:   record[4],
			Address:   record[5],
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (c contactRepository) WriteContactData(data []ProcessOutput) error {
	filePath := filepath.Join("files", "output.csv")
	header, csvData := c.convertProcessOutputToCSV(data)

	err := c.csv.WriteCSV(filePath, header, csvData)
	if err != nil {
		c.log.Errorf("Error writing to CSV file: %v", err)
		return err
	}

	return nil
}

func (c contactRepository) convertProcessOutputToCSV(outputs []ProcessOutput) (header []string, data [][]string) {
	header = []string{"ContactIDSource", "ContactIDMatch", "Accuracy"}
	for _, output := range outputs {
		accuracy, err := MapLevelToAccuracy(int(output.AccuracyLevel))
		if err != nil {
			c.log.Errorf("Error converting process output to csv: %v", err)
		}

		record := []string{
			output.ContactIDSource,
			output.ContactIDMatch,
			string(accuracy),
		}

		data = append(data, record)
	}

	return header, data
}

func (c contactRepository) WriteDuplicateContacts(data []Contact) error {
	filePath := filepath.Join("files", "duplicate.csv")
	header, csvData := c.convertProcessDuplicatesToCSV(data)

	err := c.csv.WriteCSV(filePath, header, csvData)
	if err != nil {
		c.log.Errorf("Error writing to CSV file: %v", err)
		return err
	}

	return nil
}

func (c contactRepository) convertProcessDuplicatesToCSV(contacts []Contact) (header []string, data [][]string) {
	header = []string{"ContactID", "FirstName", "LastName", "Email", "ZipCode", "Address"}

	for _, contact := range contacts {
		record := []string{
			contact.ContactID,
			contact.FirstName,
			contact.LastName,
			contact.Email,
			contact.ZipCode,
			contact.Address,
		}
		data = append(data, record)
	}

	return header, data
}
