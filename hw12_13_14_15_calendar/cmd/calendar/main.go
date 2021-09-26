package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/app"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/logger"
	"github.com/thewolf27/hw12_13_14_15_calendar/internal/server/grpc/api"
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := NewConfig()
	logg := logger.New(config.Logger.Level, logFile)

	var storage app.Storage

	switch config.Storage.Type {
	case "memory":
		storage = memorystorage.New()
	case "db":
		storage = sqlstorage.New(config.DB.Dsn)
		storage.(*sqlstorage.Storage).Connect(ctx)
	default:
		log.Println("Config Storage Type is unknown")
		return
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.HTTPServer.Host, config.HTTPServer.Port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	go api.RunGrpcServer(ctx, fmt.Sprintf("%s:%s", config.GrpcServer.Host, config.GrpcServer.Port), calendar)
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
protoc -I ./protobuf --go_out=plugins:grpc protobuf/calendar.proto
protoc -I=grpc/protobuf --go_out=grpc/generated grpc/protobuf/calendar.proto

protoc -I grpc/protobuf --go-grpc_out=grpc/generated grpc/protobuf/calendar.proto

protoc -I grpc/protobuf --go-grpc_out=grpc grpc/protobuf/calendar.proto

protoc -I=grpc/protobuf --go_out=grpc/generated --go-grpc_out=grpc/generated grpc/protobuf/calendar.proto
* docker run --network postgres -it --rm postgres psql -h dpostgres -U homestead;
*/
