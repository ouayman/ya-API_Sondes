package helperhttp

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type service struct {
	Name   string
	Server *http.Server
}

type Services struct {
	Logger   *log.Logger
	services []service
}

func (obj *Services) AddServer(name string, server *http.Server) bool {
	if server == nil {
		return false
	}
	obj.services = append(obj.services, service{Name: name, Server: server})
	return true
}

func (obj *Services) Add(name string, config ServerConfig, router http.Handler) bool {
	return obj.AddServer(name, CreateService(config, router))
}

func (obj *Services) Run() {
	if obj.Logger == nil {
		null, _ := os.Open(os.DevNull)
		obj.Logger = log.New(null, "", 0)
	}

	// shutdownChan channel signals the application to finish up
	var shutdownChan = make(chan struct{}, len(obj.services)+1)

	// Run services in go routines
	for _, v := range obj.services {
		go obj.startService(v, shutdownChan)
	}

	// launch a worker whose job it is to always watch for gracefulStop signals
	go obj.osSignalShutdown(shutdownChan)

	// wait for shutdownChan channel write to close the app down
	<-shutdownChan

	// Gracefully shutdown all services
	var timeout = 5
	var waitStopGroup sync.WaitGroup
	timedCtx, cancelFn := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	for _, v := range obj.services {
		go obj.stopService(timedCtx, v, &waitStopGroup)
	}
	waitStopGroup.Wait()
	cancelFn()
}

func (obj *Services) osSignalShutdown(shutdownChan chan struct{}) {
	// Channel to call stop services for os signals
	var gracefulStopChan = make(chan os.Signal)
	signal.Notify(gracefulStopChan, syscall.SIGINT, syscall.SIGTERM)
	// Block until we get an OS signal
	sig := <-gracefulStopChan
	obj.Logger.Printf("caught sig: %+v", sig)

	// send message on "finish up" channel to tell the app to gracefully shutdown
	shutdownChan <- struct{}{}
}

func (obj *Services) startService(svc service, shutdownChan chan struct{}) {
	obj.Logger.Printf("%v listening on %v", svc.Name, svc.Server.Addr)

	if svc.Server.TLSConfig != nil {
		if err := svc.Server.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			obj.Logger.Printf("%v service error %v", svc.Name, err.Error())
			shutdownChan <- struct{}{}
		}
	} else if err := svc.Server.ListenAndServe(); err != http.ErrServerClosed {
		obj.Logger.Printf("%v service error %v", svc.Name, err.Error())
		shutdownChan <- struct{}{}
	}
}

func (obj *Services) stopService(ctx context.Context, svc service, waitStopGroup *sync.WaitGroup) {
	waitStopGroup.Add(1)
	defer waitStopGroup.Done()
	if err := svc.Server.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		obj.Logger.Printf("%v service Shutdown: %v", svc.Name, err)
	}
}
