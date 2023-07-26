package app

import (
	"fmt"
	"github.com/elvin-tacirzade/log-exporter/pkg/config"
	"github.com/elvin-tacirzade/log-exporter/pkg/db"
	"github.com/elvin-tacirzade/log-exporter/pkg/docker"
	"github.com/elvin-tacirzade/log-exporter/pkg/reader"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type (
	App interface {
		Start()
		Shutdown()
	}

	app struct {
		Config *config.Config
		Loki   *db.Loki
		Reader *reader.Reader
		Docker *docker.Docker
	}
)

func New() (App, error) {
	// create config
	conf, confErr := config.New()
	if confErr != nil {
		return nil, fmt.Errorf("failed to create a new config: %v", confErr)
	}

	// create loki
	loki := db.NewLoki(conf)

	// create docker
	dock, dockErr := docker.New(conf)
	if dockErr != nil {
		return nil, fmt.Errorf("failed to create a new docker: %v", dockErr)
	}

	// create reader
	r, rErr := reader.New(loki, dock)
	if rErr != nil {
		return nil, fmt.Errorf("failed to create a new reader: %v", rErr)
	}

	return &app{
		Config: conf,
		Loki:   loki,
		Reader: r,
		Docker: dock,
	}, nil
}

func (a *app) Start() {
	// start reader
	go a.Reader.Handle()

	log.Println("exporter started successfully")
}

func (a *app) Shutdown() {
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
