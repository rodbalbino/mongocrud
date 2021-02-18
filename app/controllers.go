package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	Name string `json:"name"`
	City string `json:"city"`
	Age  int    `json:"age"`
}

var userCollection = db().Database("goTest").Collection("users")

func createProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // for add content type

	var person user
	err := json.NewDecoder(r.Body).Decode(&person) // storing in person variable of type user

	if err != nil {
		fmt.Print(err)
	}

	insertResult, err := userCollection.InsertOne(context.TODO(), person)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(insertResult.InsertedID) // return the mongodb ID of generated document
}

func getUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body user

	e := json.NewDecoder(r.Body).Decode(&body)

	if e != nil {
		fmt.Print(e)
	}

	var result primitive.M // an unordered representation of a BSON document which is a Map
	err := userCollection.FindOne(context.TODO(), bson.D{{"name", body.Name}}).Decode(&result)

	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(result) // return a map containing mongodb document
}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		Name string `json:"name"`
		City string `json:"city"`
	}

	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {
		fmt.Print(e)
	}

	filter := bson.D{{"name", body.Name}}
	after := options.After // to returning updated document

	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	update := bson.D{{"$set", bson.D{{"city", body.City}}}}

	updateResult := userCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)

}

func deleteProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get parameters value as string

	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb hex ID

	if err != nil {
		fmt.Printf(err.Error())
	}

	opts := options.Delete().SetCollation(&options.Collation{}) // to especify language-specific rules for string comparison, such as rules for lettercase
	res, err := userCollection.DeleteOne(context.TODO(), bson.D{{"_id", _id}}, opts)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(res.DeletedCount) // return number of documents deleted
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []primitive.M                                   //slice for multiple documents
	cur, err := userCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {

		fmt.Println(err)

	}
	for cur.Next(context.TODO()) { //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	json.NewEncoder(w).Encode(results)
}
