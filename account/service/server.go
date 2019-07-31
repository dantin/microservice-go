package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/dantin/microservice-go/account/dbclient"
	"github.com/dantin/microservice-go/pkg/logutil"
	log "github.com/sirupsen/logrus"
)

var (
	healthy int32

	// DBClient is a client of database.
	DBClient dbclient.IBoltClient
)

// Server implements HTTP server.
type Server struct {
	server *http.Server
}

// initBoltClient creates instance and calls the OpenBoltDB and Seed funcs.
func initBoltClient() {
	DBClient = &dbclient.BoltClient{}
	DBClient.OpenBoltDB()
	DBClient.Seed()
}

// NewServer creates a new instance of HTTP server.
func NewServer(addr string) *Server {
	// init logger.
	logutil.InitLogger(&logutil.LogConfig{Level: "debug", File: logutil.FileLogConfig{}})
	initBoltClient()
	handler := newRouter()

	return &Server{
		server: &http.Server{
			Addr:           addr,
			Handler:        handler,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    30 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

// Run starts HTTP server.
func (s *Server) Run() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	done := make(chan bool)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		sig := <-sc
		fmt.Printf("got signal [%d] to exit.\n", sig)
		log.Infof("%s - Shutdown signal received...", hostname)
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		s.server.SetKeepAlivesEnabled(false)
		if err := s.server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v", err)
		}
		close(done)
	}()

	log.Infof("%s - Start server on port %v", hostname, s.server.Addr)
	atomic.StoreInt32(&healthy, 1)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v", s.server.Addr, err)
	}

	<-done
	log.Fatalf("%s - Server gracefully stopped", hostname)

	return nil
}
