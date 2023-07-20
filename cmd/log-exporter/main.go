package main

import (
	"github.com/elvin-tacirzade/log-exporter/pkg/app"
	"log"
)

func main() {
	// create app
	a, err := app.New()
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	// start app
	a.Start()

	// shutdown app
	a.Shutdown()
}
