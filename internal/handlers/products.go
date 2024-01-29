package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/X0EF/go-product-api/internal/models"
)

func NewProductService(l *log.Logger, db *mongo.Collection) *Service {
	return &Service{l, db}
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	switch r.Method {
	case http.MethodGet:
		s.getProducts(w, r)
	case http.MethodPost:
		s.addProduct(w, r)
	case http.MethodPut:
		fmt.Fprintln(w, "This is a PUT request.")
	case http.MethodDelete:
		if len(segments) == 2 && segments[1] != "" {
			s.removeOne(w, r, segments[0])
			return
		}
		s.clearCollection(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Unsupported HTTP method.")
	}
}

func (s *Service) getProducts(w http.ResponseWriter, r *http.Request) {
	s.l.Println("Handle GET Products")
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

	if err := cursor.All(ctx, &products); err != nil {
		log.Printf("Failed to decode products: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response := ResponseObject{Data: &products, Message: "Documents retrieved successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service) clearCollection(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := s.db.DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Service) removeOne(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := s.db.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Printf("Failed to remvoe product: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var response ResponseObject
	if result.DeletedCount > 0 {
		response = ResponseObject{Message: "Document deleted successfully"}
	} else {
		response = ResponseObject{Message: "No document found to delete"}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service) addProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	product.Id = primitive.NewObjectID()
	currentTime := time.Now().Unix()
	product.CreatedAt = currentTime
	product.UpdatedAt = currentTime

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.InsertOne(ctx, product)
	if err != nil {
		log.Printf("Failed to insert product: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response := ResponseObject{Data: product, Message: "Document deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
