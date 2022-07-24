package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"bl.com/api/authentication"
	"bl.com/api/models"
	"bl.com/api/responses"
	"bl.com/api/sqlmanager"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _client *mongo.Client

//Handle incoming requests
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/signup", signup)
	// router.HandleFunc("/signin", signin)

	router.HandleFunc("/googlesignin", authentication.Googlesignin)
	router.HandleFunc("/testauth", authentication.JWTAuthTester)
	router.HandleFunc("/refresh", authentication.RefreshToken)
	router.HandleFunc("/comment/{id1}", GetComment(_client)).Methods("GET")
	router.Use(authentication.LoggingMiddleware)

	fmt.Println("Mapped Requests")
	fmt.Println("<------------SERVER SUCCESSFULLY STARTED------------>")

	//Listen on port 6000
	log.Fatal(http.ListenAndServe(":6000", router))

}

//Get comment from database with comment id
func GetComment(client *mongo.Client) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		id1ID := params["id1"]
		var user models.Comment
		defer cancel()

		// objId, _ := primitive.ObjectIDFromHex(userId)
		commentCollection := client.Database("commentDatabase").Collection("comments")

		err := commentCollection.FindOne(ctx, bson.M{"id1": id1ID}).Decode(&user)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.CommentResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.CommentResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(rw).Encode(response)
	}
}

//Connect to database
func connectToMongo() {
	authentication.Db = sqlmanager.InitDB()

	clientOptions := options.Client().
		ApplyURI("XXXXXXXXXXXXX")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	_client = client
	if err != nil {
		log.Fatal(err)
	}
	names, err := client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(names); i++ {
		fmt.Println(names[i])
	}
}

func main() {
	fmt.Println("<------------STARTING SERVER------------>")
	connectToMongo()
	fmt.Println("Initialized Database")
	handleRequests()
}
