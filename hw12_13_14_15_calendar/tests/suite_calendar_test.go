package tests

import (
	"context"
	"os"

	"github.com/stretchr/testify/suite"
	calendar "github.com/thewolf27/hw12_13_14_15_calendar/internal/app"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/thewolf27/hw12_13_14_15_calendar/internal/storage/sql"
)

type CalendarSuite struct {
	suite.Suite
	cancelContext context.CancelFunc
	App           *calendar.App
}

func (cal *CalendarSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	cal.cancelContext = cancel
	sqlstorage := sqlstorage.New(os.Getenv("DSN"))
	if err := sqlstorage.Connect(ctx); err != nil {
		panic(err)
	}
	if err := sqlstorage.ResetDB(); err != nil {
		panic(err)
	}

	logger := logger.New("DEBUG", "./int-log.txt")

	cal.App = calendar.New(logger, sqlstorage)
}

func (cal *CalendarSuite) TearDownSuite() {
	cal.cancelContext()
}
