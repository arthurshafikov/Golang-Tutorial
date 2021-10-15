package tests

import (
	"github.com/stretchr/testify/suite"
)

type AppSuite struct {
	suite.Suite
	CalendarSuite
	SchedulerSuite
	SenderSuite
}

func (app *AppSuite) SetupSuite() {
	app.CalendarSuite = CalendarSuite{}
	app.CalendarSuite.SetupSuite()
	app.SchedulerSuite = SchedulerSuite{}
	app.SchedulerSuite.SetupSuite()
	app.SenderSuite = SenderSuite{}
	app.SenderSuite.SetupSuite()
}

func (app *AppSuite) TearDownSuite() {
}
