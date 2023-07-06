package main

import(
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"math/rand"
	"strconv"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"net/url"

)

type Movie struct{
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`

}

type Director struct{
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`

}

var (
	moviesCollection *mongo.Collection
	client *mongo.Client
)

// var movies []Movie


func connectToMongoDB() *mongo.Client {
	// Set the MongoDB connection URI
	// uri := "mongodb+srv://shivali432d:Kittu%402000@cluster0.3dfuuak.mongodb.net/?retryWrites=true&w=majority"
	username := "shivali432d
	password := "Kittu@2000"
	escapedUsername := url.QueryEscape(username)
	escapedPassword := url.QueryEscape(password)
	uri := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.3dfuuak.mongodb.net/?retryWrites=true&w=majority", escapedUsername, escapedPassword)


	// Create a new MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
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

	return client
}



func getMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-type", "application/json")
	// json.NewEncoder(w).Encode(movies)

	cursor, err := moviesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var movies []Movie
	for cursor.Next(context.Background()) {
		var movie Movie
		if err := cursor.Decode(&movie); err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	movieID := params["id"]

	filter := bson.M{"id": movieID}

	_, err := moviesCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	cursor, err := moviesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	var movies []Movie
	for cursor.Next(context.Background()) {
		var movie Movie
		if err := cursor.Decode(&movie); err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	json.NewEncoder(w).Encode(movies)

// 	for index, item := range movies{


// 		if item.ID == params["id"]{
// 			movies = append(movies[:index], movies[index+1:]...)
// 			break
// 		}
// 	}
// 	json.NewEncoder(w).Encode(movies)
	}

func getMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	movieID := params["id"]

	filter:= bson.M{"id": movieID}

	var movie Movie

	err := moviesCollection.FindOne(context.Background(), filter).Decode(&movie)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(movie)

	// for _, item := range movies{

	// 	if item.ID == params["id"]{
	// 		json.NewEncoder(w).Encode(item)
	// 		return
	// 	}
}


func createMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	
	_, err := moviesCollection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}

		// movies = append(movies, movie)
		json.NewEncoder(w).Encode(movie)
}

// func updateMovie(w http.ResponseWriter, r *http.Request){
// 	w.Header().Set("Content-type", "application/json")
// 	params := mux.Vars(r)
// 	for index, item := range movies {
// 		if item.ID == params["id"]{
// 			movies = append(movies[:index], movies[index+1:]...)
// 			var movie Movie
// 			_= json.NewDecoder(r.Body).Decode(&movie)
// 			movie.ID = params["id"]
// 			movies = append(movies, movie)
// 			json.NewEncoder(w).Encode(movie)
// 		}
// 	}
// }


func main(){
	r := mux.NewRouter()

	client := connectToMongoDB()
	defer client.Disconnect(context.Background())


	// movies = append(movies, Movie{ID:"1", Isbn:"438227", Title:"Movie One", Director: &Director{FirstName:"John", LastName: "Doe"}})
	// movies = append(movies, Movie{ID:"2", Isbn:"434556", Title:"Movie Two", Director: &Director{FirstName:"Steve", LastName: "Smith"}})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/createMovies", createMovie).Methods("POST")
	// r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))


}

// mongodb+srv://shivali432d:<password>@cluster0.3dfuuak.mongodb.net/?retryWrites=true&w=majority