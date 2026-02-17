package command

import (
	"github.com/psds-microservice/operator-pool-service/internal/database"
)

func MigrateUp(databaseURL string) error {
	return database.MigrateUp(databaseURL)
}
