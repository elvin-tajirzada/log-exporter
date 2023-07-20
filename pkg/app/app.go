package app

import (
	"fmt"
	"github.com/elvin-tacirzade/log-exporter/pkg/config"
	"github.com/elvin-tacirzade/log-exporter/pkg/db"
	"github.com/elvin-tacirzade/log-exporter/pkg/reader"
	"log"
	"os"
	"os/signal"
	"sync"
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
		Wg     sync.WaitGroup
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

	// create reader
	r, rErr := reader.New(conf.Docker.ContainerLogFilePath, loki)
	if rErr != nil {
		return nil, fmt.Errorf("failed to create a new reader: %v", rErr)
	}

	return &app{
		Config: conf,
		Loki:   loki,
		Reader: r,
	}, nil
}

func (a *app) Start() {
	a.Wg.Add(1)

	// start reader
	go a.startReader()
}

func (a *app) startReader() {
	a.Reader.Handle()
	a.Wg.Done()
}

func (a *app) Shutdown() {
	signals := make(chan os.Signal)
	done := make(chan struct{})
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		readerStopErr := a.Reader.Tail.Stop()
		if readerStopErr != nil {
			log.Printf("failed to stop reader: %v\n", readerStopErr)
		} else {
			log.Println("log reader: stopped the tailing activity")
		}

		a.Wg.Wait()

		log.Printf("received signal: %v", sig)
		done <- struct{}{}
	}()

	<-done
}
