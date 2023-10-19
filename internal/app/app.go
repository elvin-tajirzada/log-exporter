package app

import (
	"fmt"
	"github.com/elvin-tacirzade/log-exporter/internal/config"
	"github.com/elvin-tacirzade/log-exporter/internal/service"
	"github.com/elvin-tacirzade/log-exporter/pkg/containerization"
	"github.com/elvin-tacirzade/log-exporter/pkg/db"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type (
	App struct {
		Config             *config.Config
		Loki               *db.Loki
		Docker             *containerization.Docker
		LogExporterService *service.LogExporter
	}
)

func New() (*App, error) {
	// create config
	conf, confErr := config.New()
	if confErr != nil {
		return nil, fmt.Errorf("failed to create a new config: %v", confErr)
	}

	// create database (loki)
	loki := db.NewLoki(conf.Loki)

	// create docker
	docker, dockerErr := containerization.NewDocker(conf.Docker)
	if dockerErr != nil {
		return nil, fmt.Errorf("failed to create a new docker: %v", dockerErr)
	}

	// create service
	logExporter := service.NewLogExporter(loki, docker, conf.Docker.ContainerName)

	return &App{
		Config:             conf,
		Loki:               loki,
		LogExporterService: logExporter,
		Docker:             docker,
	}, nil
}

func (a *App) Start() {
	// start service
	go a.LogExporterService.Start()

	log.Println("exporter started successfully")
}

func (a *App) Shutdown() {
	signals := make(chan os.Signal)
	done := make(chan struct{})
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		dockerCliCloseErr := a.Docker.Client.Close()
		if dockerCliCloseErr != nil {
			log.Printf("failed to close docker client: %v\n", dockerCliCloseErr)
		} else {
			log.Println("docker client closed successfully")
		}

		log.Printf("received signal: %v\n", sig)
		log.Println("exporter exited. Bye...")
		done <- struct{}{}
	}()

	<-done
}
