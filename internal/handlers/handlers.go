package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/X0EF/go-product-api/internal/database"
)

type ResponseObject struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

type Service struct {
	l  *log.Logger
	db *mongo.Collection
}

type Server struct {
	port int
	db   database.Service
}

func NewServer() *http.Server {
	l := log.New(os.Stdout, "products-api ", log.LstdFlags)
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()
	NewServer := &Server{
		port: port,
		db:   database.New(),
	}

	sm := http.NewServeMux()
	sm.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		health := db.Health()
		fmt.Fprintf(w, "Health Check: %s", health["message"])
	})

	productHandlers := NewProductService(l, db.GetClient("products"))
	sm.Handle("/products/", http.StripPrefix("/products", productHandlers))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      sm, // set the default handler
		ErrorLog:     l,  // set the logger for the server
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
