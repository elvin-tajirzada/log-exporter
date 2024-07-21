package main

import (
	"context"
	"log"

	"github.com/elvin-tajirzada/log-exporter/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("Unable to create app: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.Start(ctx)
	a.Shutdown(cancel)
}
