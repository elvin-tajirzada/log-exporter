package main

import (
	"github.com/elvin-tajirzada/log-exporter/internal/app"
	"log"
)

func main() {
	// create app
	a, err := app.New()
	if err != nil {
		log.Fatalf("failed to create a new app: %v", err)
	}

	// start app
	a.Start()

	// shutdown app
	a.Shutdown()
}
