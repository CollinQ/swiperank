package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var firstNames = []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Heidi", "Ivan", "Judy"}
var lastNames = []string{"Smith", "Johnson", "Williams", "Jones", "Brown", "Davis", "Miller", "Wilson", "Moore", "Taylor"}

func main() {
	clientOptions := options.Client().ApplyURI("pass")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("akpsi-ucsb")
	projectsCollection := db.Collection("projects")
	applicantsCollection := db.Collection("applicants")

	project := bson.M {
		"name":					"Spring 2025 Recruitment",
		"totalApplicants": 		50,
		"completedComparisons": 0,
		"totalComparisons":		100,
	}

	projectInsertRes, err := projectsCollection.InsertOne(context.TODO(), project)
	if err != nil {
		log.Fatal(err)
	}

	projectID := projectInsertRes.InsertedID

	rand.Seed(time.Now().UnixNano())
	var applicants []interface{}
	for i := 0; i < 16; i++ {
		applicant := bson.M{
			"first_name": 	firstNames[rand.Intn(len(firstNames))],
			"last_name": 	lastNames[rand.Intn(len(lastNames))],
			"project_id":	projectID,
			"wins":			0,
			"losses":		0,
			"elo":			0,
		}
		applicants = append(applicants, applicant)
	}

	// insert applicants into applicants collection mongo
	_, err = applicantsCollection.InsertMany(context.TODO(), applicants)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted test project and applicants successfully")
}