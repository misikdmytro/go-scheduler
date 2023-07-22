package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/misikdmytro/go-job-worker/internal/consumer"
	"github.com/misikdmytro/go-job-worker/internal/dependency"
)

func main() {
	d, err := dependency.NewDependency()
	if err != nil {
		panic(err)
	}

	interruption, cancel1 := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel1()

	go func() {
		log.Printf("start consumer")
		if err := d.JC.Consume(); err != nil {
			if errors.Is(err, consumer.ErrConsumerClosed) {
				return
			}

			panic(err)
		}
	}()

	<-interruption.Done()

	log.Printf("interrupted")

	ctx, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	if err := d.JC.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Printf("consumer stopped")

	<-ctx.Done()
	log.Printf("graceful shutdown")
}
