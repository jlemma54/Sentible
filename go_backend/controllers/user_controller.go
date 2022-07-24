package controllers

import (
    "context"
    "encoding/json"
    "mux-mongo-api/configs"
    "mux-mongo-api/models"
    "mux-mongo-api/responses"
    "net/http"
    "time"
    "github.com/go-playground/validator/v10"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

var commentCollection *mongo.Collection = configs.GetCollection(configs.DB, "comments")
var validate = validator.New()

func CreateComment() http.HandlerFunc {
    return func(rw http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var comment models.Comment
        defer cancel()

        //validate the request body
        if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
            rw.WriteHeader(http.StatusBadRequest)
            response := responses.CommentResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
            json.NewEncoder(rw).Encode(response)
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&comment); validationErr != nil {
            rw.WriteHeader(http.StatusBadRequest)
            response := responses.CommentResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
            json.NewEncoder(rw).Encode(response)
            return
        }

        newUser := models.Comment{
            Id:       primitive.NewObjectID(),
			Author: comment.Author,
			Body: comment.Body,
			Created_Utc: comment.Created_Utc,
			Id1: comment.Id1,
			Is_Submitter: comment.Is_Submitter,
			Link_Id: comment.Link_Id,
			Locked: comment.Locked,
			Parent_Id: comment.Parent_Id,
			Permalink: comment.Permalink,
			Retrieved_On: comment.Retrieved_On,
			Subreddit: comment.Subreddit,
			Subreddit_Id: comment.Subreddit_Id,
			Total_Awards_Received: comment.Total_Awards_Received,
			Author_Flair_Type: comment.Author_Flair_Type,
			Author_Fullname: comment.Author_Fullname,
			Author_Patreon_Flair: comment.Author_Patreon_Flair,
			Author_Premium: comment.Author_Premium,
			Stock: comment.Stock,
			Positive_Sentiment: comment.Positive_Sentiment,
			Neutral_Sentiment: comment.Neutral_Sentiment,
			Negative_Sentiment: comment.Negative_Sentiment,
			Compound_Sentiment: comment.Compound_Sentiment,
			Analyzed: comment.Analyzed,
			

        }
        result, err := commentCollection.InsertOne(ctx, newUser)
        if err != nil {
            rw.WriteHeader(http.StatusInternalServerError)
            response := responses.CommentResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
            json.NewEncoder(rw).Encode(response)
            return
        }

        rw.WriteHeader(http.StatusCreated)
        response := responses.CommentResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}}
        json.NewEncoder(rw).Encode(response)
    }
}

func GetComment() http.HandlerFunc {
    return func(rw http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        params := mux.Vars(r)
        id1ID := params["id1"]
        var user models.Comment
        defer cancel()

        // objId, _ := primitive.ObjectIDFromHex(userId)

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
