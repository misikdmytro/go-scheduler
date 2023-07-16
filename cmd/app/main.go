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

	"github.com/misikdmytro/go-job-runner/internal/dependency"
)

func main() {
	d, err := dependency.NewDependency()
	if err != nil {
		panic(err)
	}

	interruption, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

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
		defer d.JC.Close()
		log.Printf("start consumer")
		if err := d.JC.Consume(interruption); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			panic(err)
		}
	}()

exit:
	for {
		select {
		case <-interruption.Done():
			log.Printf("interrupted")
			break exit
		case err := <-d.JC.Err():
			log.Printf("consumer error: %s", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}

	log.Printf("server stopped")
	<-ctx.Done()

	log.Printf("server exited")
}
