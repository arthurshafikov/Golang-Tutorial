package tests

import (
	"os"

	sender "github.com/arthurshafikov/hw12_13_14_15_calendar/internal/app"
	"github.com/arthurshafikov/hw12_13_14_15_calendar/internal/logger"
	"github.com/stretchr/testify/suite"
)

type SenderSuite struct {
	suite.Suite
	Sender *sender.Sender
}

func (sch *SenderSuite) SetupSuite() {
	logger := logger.New("DEBUG", "./int-log.txt")

	sch.Sender = sender.NewSender(logger, os.Getenv("RABBITMQ_URL"))
}

func (sch *SenderSuite) TearDownSuite() {
}
