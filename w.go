package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
//	"strconv"

	"github.com/gorilla/mux"
	"context"

"go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)




// Types
type task struct {
	Name    string
	ID      int
	Content string
}

type allTasks []task

// Persistence
var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some Content",
	},
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wecome the my GO API!")
}

func createTask(w http.ResponseWriter, r *http.Request) {
  	var newTask task
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "Insert a Valid Task Data")
		}
		json.Unmarshal(reqBody, &newTask)
		//tasks = append(tasks, newTask)
		///////////////////

		insertResult, err := collection.InsertOne(context.TODO(), newTask)
		if err != nil {
		log.Fatal(err)
		}
//
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)


	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
	    log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
	    var episode bson.M
	    if err = cursor.Decode(&episode); err != nil {
	        log.Fatal(err)
	    }
	    fmt.Println(episode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episode)
}
}

func getOneTask(w http.ResponseWriter, r *http.Request) {
	/*
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}
	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
	*/
}

func updateTask(w http.ResponseWriter, r *http.Request) {
//	var newTask task

	vars := mux.Vars(r)
	www, _ := primitive.ObjectIDFromHex(vars["id"])

	/*reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please Enter Valid Data")
	}
*/
//
cursor, err := collection.Find(context.TODO(), bson.M{})
if err != nil {
    log.Fatal(err)
}
var episodes []bson.M
if err = cursor.All(context.TODO(), &episodes); err != nil {
    log.Fatal(err)
}
fmt.Println(episodes)
//

	filter := bson.D{{"_id", www}}

		update := bson.M{
			"$set": bson.M{
				"id": 10,
				"name": "wwwwwww",
				"content": "wwwwwwwwwwwwwwwww",
			},
		}

		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	/*
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please Enter Valid Data")
	}
	json.Unmarshal(reqBody, &updatedTask)

	for i, t := range tasks {
		if t.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)

			updatedTask.ID = t.ID
			tasks = append(tasks, updatedTask)

			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(updatedTask)
			fmt.Fprintf(w, "The task with ID %v has been updated successfully", taskID)
		}
	}
*/

}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	www, err := primitive.ObjectIDFromHex(vars["id"])
//	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "NO CORRECT ID")
		return
	}

	filter := bson.D{{"_id", www}}
	collection.DeleteOne(context.TODO(), filter)
}

var collection *mongo.Collection

func main() {
	router := mux.NewRouter().StrictSlash(true)

	host := "localhost"
		port := 27017

		clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port))
		client, err := mongo.Connect(context.TODO(), clientOpts)
		if err != nil {
			log.Fatal(err)
		}

		// Check the connections
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}

		quickstartDatabase := client.Database("tasker")
    //podcastsCollection := quickstartDatabase.Collection("podcasts")
    //episodesCollection := quickstartDatabase.Collection("episodes")
  	collection = quickstartDatabase.Collection("task")

  //ash := task{"Ash", 10, "Pallet Town"}
 	//misty := task{"Misty", 10, "Cerulean City"}
 	//brock := task{"Brock", 15, "Pewter City"}




	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", createTask).Methods("POST")

	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", getOneTask).Methods("GET")

	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	fmt.Println("Congratulations, you're already connected to MongoDB!")

	log.Fatal(http.ListenAndServe(":9000", router))
}
