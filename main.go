package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

var client *mongo.Client

type Food struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Tribe  string             `json:"tribe,omitempty" bson:"tribe,omitempty"`
}

func AddFood(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var food Food
	_ = json.NewDecoder(request.Body).Decode(&food)
	collection := client.Database("foodrestapi").Collection("foods")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, food)
	json.NewEncoder(response).Encode(result)
}

func GetFood(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var food Food
	collection := client.Database("foodrestapi").Collection("foods")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Food{ID: id}).Decode(&food)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(food)
}

func GetFoods(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var foods []Food
	collection := client.Database("foodrestapi").Collection("foods")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var food Food
		cursor.Decode(&food)
		foods = append(foods, food)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(foods)
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/food", AddFood).Methods("POST")
	router.HandleFunc("/foods", GetFoods).Methods("GET")
	router.HandleFunc("/food/{id}", GetFood).Methods("GET")
	http.ListenAndServe(":12345", router)
}
