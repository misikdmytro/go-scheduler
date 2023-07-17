package broker

import (
	"context"
)

type Broker[T any] interface {
	Publish(context.Context, string, T) error
}
