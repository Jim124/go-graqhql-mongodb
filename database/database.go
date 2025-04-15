package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jim124/go-graqhql-mongodb/graph/model"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongodbUrl := os.Getenv("databaseUrl")
	client, err := mongo.Connect(options.Client().
		ApplyURI(mongodbUrl))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("log cannot connect to mongodb")
		log.Fatal(err)
	}
	return &DB{
		client: client,
	}

}
func (db *DB) GetJob(id string) *model.JobListing {
	jobCollec := db.client.Database("GoGraqh").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println(id)
	objectId, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objectId}}
	var jobListing model.JobListing
	err := jobCollec.FindOne(ctx, filter).Decode(&jobListing)
	if err != nil {
		log.Fatal(err)
	}
	return &jobListing
}

func (db *DB) InsertJob(createJobInput model.CreatingJobListingInput) *model.JobListing {
	jobCollec := db.client.Database("GoGraqh").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	newJob := model.JobListing{
		Title:       createJobInput.Title,
		Company:     createJobInput.Company,
		Description: createJobInput.Description,
		URL:         createJobInput.URL,
	}
	result, err := jobCollec.InsertOne(ctx, newJob)
	if err != nil {
		panic(err)
	}
	id := fmt.Sprintf("", result.InsertedID)
	newJob.ID = id
	return &newJob

}

func (db *DB) GetJobs() []*model.JobListing {
	jobCollec := db.client.Database("GoGraqh").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var jobs []*model.JobListing
	cursor, err := jobCollec.Find(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	if err = cursor.All(context.TODO(), &jobs); err != nil {
		panic(err)
	}
	return jobs
}

func (db *DB) UpdateJob(id string, updateJobInput model.UpdateJobListingInput) *model.JobListing {
	jobCollec := db.client.Database("GoGraqh").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	updateJob := bson.M{}
	objId, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objId}}
	if updateJobInput.Title != nil {
		updateJob["title"] = updateJobInput.Title
	}
	if updateJobInput.Description != nil {
		updateJob["description"] = updateJobInput.Description
	}
	if updateJobInput.Company != nil {
		updateJob["company"] = updateJobInput.Company
	}
	if updateJobInput.URL != nil {
		updateJob["url"] = updateJobInput.URL
	}
	update := bson.D{{Key: "$set", Value: updateJob}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var jobListing model.JobListing
	err := jobCollec.FindOneAndUpdate(ctx, filter, update, opts).Decode(&jobListing)
	if err != nil {
		log.Fatal(err)
	}
	return &jobListing
}

func (db *DB) DeleteJob(id string) *model.DeleteJobResponse {
	jobCollect := db.client.Database("GoGraqh").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objectId, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objectId}}
	_, err := jobCollect.DeleteOne(ctx, filter)
	// Prints a message if any errors occur during the operation
	if err != nil {
		log.Fatal(err)
	}
	return &model.DeleteJobResponse{
		DeleteJobID: id,
	}
}
