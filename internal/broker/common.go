package broker

import (
	"fmt"

	"github.com/misikdmytro/go-job-runner/internal/config"
)

func BuildRabbitMQURL(c config.RabbitMQConfig) string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
	)
}
