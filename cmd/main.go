package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/X0EF/go-product-api/internal/handlers"
)

func main() {
	s := handlers.NewServer()

	go func() {
		s.ErrorLog.Println("Starting server on port " + s.Addr)

		err := s.ListenAndServe()
		if err != nil {
			s.ErrorLog.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
	defer cancel()
}
