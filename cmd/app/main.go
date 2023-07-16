package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/misikdmytro/go-job-runner/internal/consumer"
	"github.com/misikdmytro/go-job-runner/internal/dependency"
)

func main() {
	d, err := dependency.NewDependency()
	if err != nil {
		panic(err)
	}

	interruption, cancel1 := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel1()

	addr := fmt.Sprintf("%s:%s", d.CFG.Server.Host, d.CFG.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: d.E,
	}

	go func() {
		log.Printf("start server on %s", addr)
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			panic(err)
		}
	}()

	go func() {
		log.Printf("start consumer")
		if err := d.JC.Setup(); err != nil {
			panic(err)
		}

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

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Printf("server stopped")

	if err := d.JC.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Printf("consumer stopped")

	<-ctx.Done()
	log.Printf("graceful shutdown")
}
