package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/app"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/thewolf27/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/thewolf27/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/thewolf27/hw12_13_14_15_calendar/internal/storage/sql"
)

var (
	configFile string
	logFile    string
)

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
	flag.StringVar(&logFile, "log", "./logs/log.txt", "Path to log file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	config := NewConfig()
	logg := logger.New(config.Logger.Level, logFile)

	var storage app.Storage

	switch config.Storage.Type {
	case "memory":
		storage = memorystorage.New()
	case "db":
		storage = sqlstorage.New(config.DB.Dsn)
	default:
		log.Fatalln("Config Storage Type is unknown")
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.Server.Host, config.Server.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

/*
* go run ./cmd/calendar/... --config=/home/thewolf/Golang/Golang-Tutorial/hw12_13_14_15_calendar/configs/config.toml
* go run ./cmd/calendar/...
* export GOPATH=$HOME/go; export PATH=$PATH:$GOPATH/bin;
* goose -dir migrations postgres "user=homestead password=secret dbname=homestead sslmode=disable" up
*
* docker ps -a -f name=dpost
* docker create network postgres
* docker run --network postgres -d --name dpostgres -e POSTGRES_PASSWORD=password -p 5432:5432 postgres
* docker run --network postgres -it --rm -e PGPASSWORD=password postgres psql -h dpostgres -U postgres
*
* docker run --network postgres -it --rm postgres psql -h dpostgres -U homestead;
 */
