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


type Blockcube struct {
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Employee *Employee `json:"director"`
}

type Employee struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var (
	client           *mongo.Client
	blockcubeCollection *mongo.Collection
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

	// Get a reference to the blockcube collection
	blockcubeCollection = client.Database("employee").Collection("blockcube")
}

func saveBlockcube(blockcube Blockcube) error {
	// Insert the blockcube document into the collection
	_, err := blockcubeCollection.InsertOne(context.Background(), blockcube)
	if err != nil {
		return err
	}
	return nil
}

func createBlockcube(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var blockcube Blockcube
	_ = json.NewDecoder(r.Body).Decode(&blockcube)
	// blockcube.ID = strconv.Itoa(rand.Intn(100000000))

	err := saveBlockcube(blockcube)
	if err != nil {
		log.Fatal(err)
	}

	// Return the created blockcube document as the response
	json.NewEncoder(w).Encode(blockcube)
}

func deleteBlockcube(w http.ResponseWriter ,r *http.Request){
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	blockcubeID := params["id"]

	objectID, err := primitive.ObjectIDFromHex(blockcubeID)
	if err != nil {
		// Invalid ID format
		http.Error(w, "Invalid blockcube ID", http.StatusBadRequest)
		return
	}
	filter := bson.M{"_id": objectID}

	result, err := blockcubeCollection.DeleteOne(context.Background(), filter)
	if err != nil{
		log.Fatal(err)
	}

	if result.DeletedCount == 0 {
		message := fmt.Sprintf("No blockcube found with ID: %s", blockcubeID)
		http.Error(w, message, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode("Blockcube deleted successfully")


}

func getBlockcube(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-type", "application/json")
    params := mux.Vars(r)
    blockcubeID := params["id"]

    objectID, err := primitive.ObjectIDFromHex(blockcubeID)
    if err != nil {
        http.Error(w, "Invalid blockcube ID", http.StatusBadRequest)
        return
    }
    filter := bson.M{"_id": objectID}

    var blockcube Blockcube

    err = blockcubeCollection.FindOne(context.Background(), filter).Decode(&blockcube)
    if err != nil {
        log.Fatal(err)
    }

    json.NewEncoder(w).Encode(blockcube)
}
	
	func updateBlockcube(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	
		params := mux.Vars(r)
		blockcubeID := params["id"]
	
		// Check if the blockcube exists
		objectID, err := primitive.ObjectIDFromHex(blockcubeID)
		if err != nil {
			http.Error(w, "Invalid blockcube ID", http.StatusBadRequest)
			return
		}
	
		filter := bson.M{"_id": objectID}
		var existingBlockcube Blockcube
		err = blockcubeCollection.FindOne(context.Background(), filter).Decode(&existingBlockcube)
		if err != nil {
			http.Error(w, "Blockcube not found", http.StatusNotFound)
			return
		}
	
		// Parse the request body and update the blockcube
		var updatedBlockcube Blockcube
		err = json.NewDecoder(r.Body).Decode(&updatedBlockcube)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
	
		// Set the updated values
		existingBlockcube.Title = updatedBlockcube.Title
		existingBlockcube.Employee = updatedBlockcube.Employee
	
		// Update the blockcube in the database
		updateResult, err := blockcubeCollection.ReplaceOne(context.Background(), filter, existingBlockcube)
		if err != nil {
			log.Fatal(err)
		}
	
		if updateResult.ModifiedCount == 0 {
			http.Error(w, "Failed to update blockcube", http.StatusInternalServerError)
			return
		}
	
		json.NewEncoder(w).Encode(existingBlockcube)
	}
	json.NewDecoder(w).EN
	


func main() {
	r := mux.NewRouter()

	// Connect to the MongoDB database
	connectToMongoDB()
	defer client.Disconnect(context.Background())

	r.HandleFunc("/blockcube", createBlockcube).Methods("POST")
	r.HandleFunc("/blockcube/{id}", deleteBlockcube).Methods("DELETE")
	r.HandleFunc("/blockcube/{id}", getBlockcube).Methods("GET")
	r.HandleFunc("/blockcube/{id}", updateBlockcube).Methods("PUT")

	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
