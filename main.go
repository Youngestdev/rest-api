package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)
// Curled from the official mongo docs for Go
type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	//  Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	//Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection successful.")

	collection := client.Database("restapi").Collection("trainers")

	abdul := Trainer{"Abdulazeez Abdulazeez Adeshina", 16, "Lekki, Lagos"}

	// Insert data

	insertData, err := collection.InsertOne(context.TODO(), abdul)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertData.InsertedID)

	//  For multiple.

	trainera := Trainer{"Abdulazeez Abdulazeez Adeshina2", 16, "Lekki, Lagos"}
	trainerb := Trainer{"Abdulazeez Abdulazeez Adeshina3", 16, "Lekki, Lagos"}
	trainers := []interface{}{trainera, trainerb}

	insertManyData, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyData.InsertedIDs)

	//Update document
	filter := bson.D{{"name", "Abdulazeez Abdulazeez Adeshina"}}
	updatedDoc := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, updatedDoc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Matched %v documents and updated %v documents. \n", updateResult.MatchedCount, updateResult.ModifiedCount)

	//Find a document.
	var res Trainer
	err = collection.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found a single document: %v\n", res)

	//  Find multiple - will come in handy soon.

	findoptions := options.Find()
	findoptions.SetLimit(2)

	var reses []*Trainer

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findoptions)
	if err != nil {
		log.Fatal(err)
	}

	//  Let's loop, lol.

	for cur.Next(context.TODO()) {

		//  create a single value into which a single document can be decoded
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		reses = append(reses, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers: %v\n", reses)

	//Delete document

	deleteDocument, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteDocument.DeletedCount)
}
