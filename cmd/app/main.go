package main

import (
	"fmt"
	"log"

	"github.com/misikdmytro/go-job-runner/internal/dependency"
)

func main() {
	d, err := dependency.NewDependency()
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf("%s:%s", d.CFG.Server.Host, d.CFG.Server.Port)
	log.Printf("start server on %s", addr)
	if err := d.E.Run(addr); err != nil {
		panic(err)
	}
}
