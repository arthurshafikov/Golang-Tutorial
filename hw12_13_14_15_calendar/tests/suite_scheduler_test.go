package tests

import (
	"context"
	"os"

	scheduler "github.com/arthurshafikov/hw12_13_14_15_calendar/internal/app"
	"github.com/arthurshafikov/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/arthurshafikov/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/stretchr/testify/suite"
)

type SchedulerSuite struct {
	suite.Suite
	cancelContext context.CancelFunc
	Scheduler     *scheduler.Scheduler
}

func (sch *SchedulerSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	sch.cancelContext = cancel

	sqlstorage := sqlstorage.New(os.Getenv("DSN"))
	sqlstorage.Connect(ctx)
	if err := sqlstorage.ResetDB(); err != nil {
		panic(err)
	}

	logger := logger.New("DEBUG", "./int-log.txt")

	sch.Scheduler = scheduler.NewScheduler(logger, sqlstorage, os.Getenv("RABBITMQ_URL"))
}

func (sch *SchedulerSuite) TearDownSuite() {
	sch.cancelContext()
}
