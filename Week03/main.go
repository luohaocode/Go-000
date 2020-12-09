package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	server *http.Server
}

func NewServer() *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "hello, world")
	})

	return &Server{
		server: &http.Server{
			Addr:    "127.0.0.1:8001",
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	fmt.Println("server start")
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("server shutdown")
	return s.server.Shutdown(ctx)
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	s := NewServer()

	g.Go(func() error {
		go func() {
			<-ctx.Done()
			fmt.Println("server ctx done")
			s.Shutdown(ctx)
		}()
		return s.Start()
	})

	g.Go(func() error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)
		select {
		case <-sig:
			_, cancel := context.WithCancel(ctx)
			cancel()
			return errors.New("interrupt signal")
		case <-ctx.Done():
			fmt.Println("sig ctx done")
			return ctx.Err()
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("error is %v\n", err)
	}
}
