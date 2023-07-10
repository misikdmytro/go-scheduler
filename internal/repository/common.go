package repository

import (
	"fmt"

	"github.com/misikdmytro/go-job-runner/internal/config"
)

func buildDataSourceName(c config.DBConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
	)
}
