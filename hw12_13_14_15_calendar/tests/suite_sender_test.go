package tests

import (
	"os"

	"github.com/stretchr/testify/suite"
	sender "github.com/thewolf27/hw12_13_14_15_calendar/internal/app"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/logger"
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
