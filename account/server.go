package account

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
		s.server.ErrorLog.Printf("%s - Shutdown signal received...\n", hostname)
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		s.server.SetKeepAlivesEnabled(false)
		if err := s.server.Shutdown(ctx); err != nil {
			s.server.ErrorLog.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	s.server.ErrorLog.Printf("%s - Start server on port %v\n", hostname, s.server.Addr)
	atomic.StoreInt32(&healthy, 1)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.server.ErrorLog.Fatalf("Could not listen on %s: %v\n", s.server.Addr, err)
	}

	<-done
	s.server.ErrorLog.Fatalf("%s - Server gracefully stopped.\n", hostname)

	return nil
}
