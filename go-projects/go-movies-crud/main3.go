package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	// "math/rand"/
	"net/http"
	// "strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"

)
import "go.mongodb.org/mongo-driver/bson/primitive"


type Movie struct {
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var (
	client           *mongo.Client
	moviesCollection *mongo.Collection
)

func connectToMongoDB() {
	// Set the MongoDB connection URI
	username := "shivali432d"
	password := "Kittu@2000"
	escapedUsername := url.QueryEscape(username)
	escapedPassword := url.QueryEscape(password)
	uri := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.3dfuuak.mongodb.net/?retryWrites=true&w=majority", escapedUsername, escapedPassword)

	// Create a new MongoDB client
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a reference to the movies collection
	moviesCollection = client.Database("entertainment").Collection("movies")
}

func saveMovie(movie Movie) error {
	// Insert the movie document into the collection
	_, err := moviesCollection.InsertOne(context.Background(), movie)
	if err != nil {
		return err
	}
	return nil
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	// movie.ID = strconv.Itoa(rand.Intn(100000000))

	err := saveMovie(movie)
	if err != nil {
		log.Fatal(err)
	}

	// Return the created movie document as the response
	json.NewEncoder(w).Encode(movie)
}

func deleteMovie(w http.ResponseWriter ,r *http.Request){
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	movieID := params["id"]

	objectID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		// Invalid ID format
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}
	filter := bson.M{"_id": objectID}

	result, err := moviesCollection.DeleteOne(context.Background(), filter)
	if err != nil{
		log.Fatal(err)
	}

	if result.DeletedCount == 0 {
		message := fmt.Sprintf("No movie found with ID: %s", movieID)
		http.Error(w, message, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode("Movie deleted successfully")


}

func getMovie(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-type", "application/json")
    params := mux.Vars(r)
    movieID := params["id"]

    objectID, err := primitive.ObjectIDFromHex(movieID)
    if err != nil {
        http.Error(w, "Invalid movie ID", http.StatusBadRequest)
        return
    }
    filter := bson.M{"_id": objectID}

    var movie Movie

    err = moviesCollection.FindOne(context.Background(), filter).Decode(&movie)
    if err != nil {
        log.Fatal(err)
    }

    json.NewEncoder(w).Encode(movie)
}
	
	func updateMovie(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	
		params := mux.Vars(r)
		movieID := params["id"]
	
		// Check if the movie exists
		objectID, err := primitive.ObjectIDFromHex(movieID)
		if err != nil {
			http.Error(w, "Invalid movie ID", http.StatusBadRequest)
			return
		}
	
		filter := bson.M{"_id": objectID}
		var existingMovie Movie
		err = moviesCollection.FindOne(context.Background(), filter).Decode(&existingMovie)
		if err != nil {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
	
		// Parse the request body and update the movie
		var updatedMovie Movie
		err = json.NewDecoder(r.Body).Decode(&updatedMovie)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
	
		// Set the updated values
		existingMovie.Title = updatedMovie.Title
		existingMovie.Director = updatedMovie.Director
	
		// Update the movie in the database
		updateResult, err := moviesCollection.ReplaceOne(context.Background(), filter, existingMovie)
		if err != nil {
			log.Fatal(err)
		}
	
		if updateResult.ModifiedCount == 0 {
			http.Error(w, "Failed to update movie", http.StatusInternalServerError)
			return
		}
	
		json.NewEncoder(w).Encode(existingMovie)
	}
	


func main() {
	r := mux.NewRouter()

	// Connect to the MongoDB database
	connectToMongoDB()
	defer client.Disconnect(context.Background())

	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")

	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
