package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/X0EF/go-product-api/database"
	"github.com/X0EF/go-product-api/handlers"
)

func main() {
	l := log.New(os.Stdout, "products-api ", log.LstdFlags)
	db := database.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	// create a new serve mux and register the handlers
	sm := http.NewServeMux()
	sm.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		health := db.Health()
		fmt.Fprintf(w, "Health Check: %s", health["message"])
	})
	productsHandler := handlers.NewProducts(l, db.GetClient("products"))
	sm.Handle("/products", productsHandler)
	s := http.Server{
		Addr:         ":" + port,        // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	go func() {
		l.Println("Starting server on port: " + port)

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
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
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
