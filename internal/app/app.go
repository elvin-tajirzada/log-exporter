package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/elvin-tajirzada/log-exporter/internal/config"
	"github.com/elvin-tajirzada/log-exporter/internal/service"
	"github.com/elvin-tajirzada/log-exporter/pkg/containerization"
	"github.com/elvin-tajirzada/log-exporter/pkg/database"
)

type App struct {
	Config             *config.Config
	Loki               *database.Loki
	Docker             *containerization.Docker
	LogExporterService *service.LogExporter
	wg                 sync.WaitGroup
}

func New() (*App, error) {
	conf, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("unable to create config: %v", err)
	}

	loki := database.NewLoki(conf.Loki)

	docker, err := containerization.NewDocker(conf.Docker)
	if err != nil {
		return nil, fmt.Errorf("unable to create docker: %v", err)
	}

	logExporterService := service.NewLogExporter(loki, docker, conf.Docker.ContainerName, conf.Docker.ReconnectTime)

	return &App{
		Config:             conf,
		Loki:               loki,
		Docker:             docker,
		LogExporterService: logExporterService,
	}, nil
}

func (a *App) Start(ctx context.Context) {
	a.wg.Add(1)
	// start log exporter service
	go a.startLogExporterService(ctx)
}

func (a *App) startLogExporterService(ctx context.Context) {
	log.Println("Log exporter started successfully")
	a.LogExporterService.Start(ctx)
	a.wg.Done()
}

func (a *App) Shutdown(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Printf("Received signal: %v\n", sig)

		log.Println("Waiting for log-exporter service shutdown gracefully...")
		cancel()
		a.wg.Wait()

		if err := a.Docker.Client.Close(); err != nil {
			log.Printf("Unable to close docker client: %v\n", err)
		} else {
			log.Println("Docker client closed successfully")
		}

		log.Println("Log exporter exited. Bye...")
		done <- struct{}{}
	}()

	<-done
}
