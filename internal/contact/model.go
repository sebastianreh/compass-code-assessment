package contact

import "errors"

const (
	InvalidLevelError = "invalid level"
)

type Accuracy string

const (
	VeryLow  Accuracy = "Very Low"
	Low      Accuracy = "Low"
	Medium   Accuracy = "Medium"
	High     Accuracy = "High"
	VeryHigh Accuracy = "Very High"
)

type Contact struct {
	ContactID string `json:"contact_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	ZipCode   string `json:"zip_code"`
	Address   string `json:"address"`
}

type ProcessOutput struct {
	ContactIDSource string `json:"contact_id_source"`
	ContactIDMatch  string `json:"contact_id_match"`
	AccuracyLevel   int    `json:"accuracy"`
}

func MapLevelToAccuracy(level int) (Accuracy, error) {
	if level < 0 || level > 5 {
		return "", errors.New(InvalidLevelError)
	}

	var accuracyMap = map[int]Accuracy{
		1: VeryLow,
		2: Low,
		3: Medium,
		4: High,
		5: VeryHigh,
	}

	return accuracyMap[level], nil
}
