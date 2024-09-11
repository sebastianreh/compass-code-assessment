package contact

import (
	"github.com/sirupsen/logrus"
)

type Service interface {
	Evaluate() ([]ProcessOutput, error)
}

type Repository interface {
	GetContactData() ([]Contact, error)
	WriteContactData(data []ProcessOutput) error
	WriteDuplicateContacts(duplicate []Contact) error
}

type contactService struct {
	log        *logrus.Logger
	repository Repository
}

func NewContactService(log *logrus.Logger, repository Repository) Service {
	return &contactService{
		log:        log,
		repository: repository,
	}
}

func (c contactService) Evaluate() ([]ProcessOutput, error) {
	contacts, err := c.repository.GetContactData()
	if err != nil {
		c.log.Errorf("error getting contact data: %v", err)
		return nil, err
	}

	contactMap := make(map[string]Contact)
	var results []ProcessOutput

	// Create compared pairs map to validate if already compared
	comparedPairs := make(map[string]bool)

	var duplicates []Contact

	for _, contact := range contacts {
		contactMap[contact.ContactID] = contact
	}

	for i := 0; i < len(contacts); i++ {
		for j := i + 1; j < len(contacts); j++ {
			compareContacts(contacts[i], contacts[j], contactMap, comparedPairs, &results, &duplicates)
		}
	}

	err = c.repository.WriteDuplicateContacts(duplicates)
	if err != nil {
		c.log.Errorf("error writing duplicate contact data: %v", err)
		return nil, err
	}

	err = c.repository.WriteContactData(results)
	if err != nil {
		c.log.Errorf("error writing contact data: %v", err)
		return nil, err
	}

	return results, nil
}

func compareContacts(contact1, contact2 Contact, contactMap map[string]Contact, comparedPairs map[string]bool, results *[]ProcessOutput, duplicates *[]Contact) {
	pairKey := generatePairKey(contact1.ContactID, contact2.ContactID)

	if comparedPairs[pairKey] {
		return
	}

	accuracyLevel, duplicate := calculateAccuracy(contact1, contact2)
	if accuracyLevel > 0 {
		*results = append(*results, ProcessOutput{
			ContactIDSource: contact1.ContactID,
			ContactIDMatch:  contact2.ContactID,
			AccuracyLevel:   accuracyLevel,
		})
		*results = append(*results, ProcessOutput{
			ContactIDSource: contact2.ContactID,
			ContactIDMatch:  contact1.ContactID,
			AccuracyLevel:   accuracyLevel,
		})
	}

	if duplicate {
		*duplicates = append(*duplicates, contact1)
		*duplicates = append(*duplicates, contact2)
	}

	comparedPairs[pairKey] = true
}

// Create unique pair keys
func generatePairKey(id1, id2 string) string {
	if id1 < id2 {
		return id1 + "_" + id2
	}
	return id2 + "_" + id1
}

func calculateAccuracy(c1, c2 Contact) (int, bool) {
	var score float64
	fieldsToCompare := [][2]string{
		{c1.FirstName, c2.FirstName},
		{c1.LastName, c2.LastName},
		{c1.Email, c2.Email},
		{c1.ZipCode, c2.ZipCode},
		{c1.Address, c2.Address},
	}

	for _, fields := range fieldsToCompare {
		compareAndAddScore(fields[0], fields[1], &score)
	}

	if score == 5 {
		return 0, true
	}

	// Since there is no logic defined for the accuracy score, I decided that if there is no match in any of fields, the value of both names matching should be 0,5
	// so it doesn't generate too much very low accuracy data. Otherwise, it should add 1, since there is some more probably of being the same contact.
	var firstLetterAddValue float64
	if score == 0 {
		firstLetterAddValue = 0.5
	} else {
		firstLetterAddValue = 1
	}

	if c1.FirstName != c2.FirstName {
		compareFirstLetter(c1.FirstName, c2.FirstName, &score, firstLetterAddValue)
	}

	if c1.LastName != c2.LastName {
		compareFirstLetter(c1.LastName, c2.LastName, &score, firstLetterAddValue)
	}

	return int(score), false
}

func compareFirstLetter(name1, name2 string, score *float64, addValue float64) {
	if len(name1) > 0 && len(name2) > 0 {
		field1 := name1[:1]
		field2 := name2[:1]
		if field1 == field2 {
			*score += addValue
		}
	}
}

func compareAndAddScore(field1, field2 string, score *float64) {
	if len(field1) > 1 && len(field2) > 1 {
		if field1 == field2 {
			*score++
		}
	}
}
