package resolver

import (
	"context"
	"fmt"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/thewolf27/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/thewolf27/hw12_13_14_15_calendar/internal/storage/sql"
)

func ResolveStorage(ctx context.Context, config config.Config) (interface{}, error) {
	var storage interface{}

	switch config.Storage.Type {
	case "memory":
		storage = memorystorage.New()
	case "db":
		storage = sqlstorage.New(config.DB.Dsn)
		storage.(*sqlstorage.Storage).Connect(ctx)
	default:
		return nil, fmt.Errorf("config Storage Type is unknown")
	}

	return storage, nil
}
