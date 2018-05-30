package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/couchbase/gocb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

// Movie represents a Movie on database
type Movie struct {
	ID      string      `json:"id,omitempty"`
	Name    string      `json:"name,omitempty"`
	Genre   string      `json:"genre,omitempty"`
	Formats MovieFormat `json:"formats,omitempty"`
}

// MovieFormat is list of formats that has a movie
type MovieFormat struct {
	Digital bool `json:"digital,omitempty"`
	Bluray  bool `json:"bluray,omitempty"`
	Dvd     bool `json:"dvd,omitempty"`
}

var bucket *gocb.Bucket
var bucketName string

// ListEndpoint returns all movies on the bucket
func ListEndpoint(w http.ResponseWriter, req *http.Request) {
	var movies []Movie
	query := gocb.NewN1qlQuery("SELECT `" + bucketName + "`.* FROM `" + bucketName + "`")
	query.Consistency(gocb.RequestPlus)
	rows, _ := bucket.ExecuteN1qlQuery(query, nil)
	var row Movie
	for rows.Next(&row) {
		movies = append(movies, row)
		row = Movie{}
	}

	if movies == nil {
		movies = make([]Movie, 0)
	}
	json.NewEncoder(w).Encode(movies)
}

// SearchEndpoint returns all movies which satisfy the search
func SearchEndpoint(w http.ResponseWriter, req *http.Request) {
	var movies []Movie
	params := mux.Vars(req)
	var n1qlParams []interface{}
	n1qlParams = append(n1qlParams, strings.ToLower(params["title"]))

	query := gocb.NewN1qlQuery("SELECT `" + bucketName + "`.* FROM `" + bucketName +
		"` WHERE LOWER(name) LIKE LOWER('%' || $1 || '%')")
	query.Consistency(gocb.RequestPlus)

	rows, err := bucket.ExecuteN1qlQuery(query, n1qlParams)
	if err != nil {
		panic(err)
	}
	var row Movie
	for rows.Next(&row) {
		movies = append(movies, row)
		row = Movie{}
	}

	if movies == nil {
		movies = make([]Movie, 0)
	}
	json.NewEncoder(w).Encode(movies)
}

// CreateEndpoint creates a movie entry on database
func CreateEndpoint(w http.ResponseWriter, req *http.Request) {
	var movie Movie
	_ = json.NewDecoder(req.Body).Decode(&movie)
	bucket.Insert(uuid.Must(uuid.NewV4()).String(), movie, 0)
	json.NewEncoder(w).Encode(movie)
}

func main() {
	fmt.Println("Starting server at http://localhost:3333...")

	cluster, _ := gocb.Connect("couchbase://localhost")
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "admin123",
	})
	bucketName = "restfull-sample"
	var err error
	bucket, err = cluster.OpenBucket(bucketName, "")
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/movies", ListEndpoint).Methods("GET")
	router.HandleFunc("/movies", CreateEndpoint).Methods("POST")
	router.HandleFunc("/search/{title}", SearchEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":3333",
		handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT",
			"HEAD"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
