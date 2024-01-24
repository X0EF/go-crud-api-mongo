package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/X0EF/go-product-api/models"
)

// ProductService is a service for managing products
type ProductService struct {
	db *mongo.Collection
}

// NewProductService creates a new ProductService
func NewProductService(db *mongo.Database) *ProductService {
	return &ProductService{db: db.Collection("products")}
}

// AddProduct adds a new product to the MongoDB collection
func (s *ProductService) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Insert the product into the MongoDB collection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.InsertOne(ctx, product)
	if err != nil {
		log.Printf("Failed to insert product: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetAllProducts retrieves all products from the MongoDB collection
func (s *ProductService) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product

	// Retrieve all products from the MongoDB collection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.db.Find(ctx, bson.D{})
	if err != nil {
		log.Printf("Failed to retrieve products: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Decode the products from the cursor
	if err := cursor.All(ctx, &products); err != nil {
		log.Printf("Failed to decode products: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return the products as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
