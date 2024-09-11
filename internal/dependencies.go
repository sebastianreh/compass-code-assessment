package internal

import (
	"github.com/sebastianreh/compass-code-assessment/internal/contact"
	"github.com/sebastianreh/compass-code-assessment/pkg"
	"github.com/sirupsen/logrus"
)

type Dependencies struct {
	Logger  *logrus.Logger
	Service contact.Service
}

func Build() Dependencies {
	logger := logrus.New()
	csvConnector := pkg.NewCSVConnector()
	repository := contact.NewContactRepository(logger, csvConnector)
	service := contact.NewContactService(logger, repository)

	return Dependencies{
		Logger:  logger,
		Service: service,
	}
}
