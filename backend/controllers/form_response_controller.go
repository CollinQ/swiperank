package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"backend/db"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type FormResponseController struct {
	collection *mongo.Collection
}

func NewFormResponseController() *FormResponseController {
	return &FormResponseController{
		collection: db.GetCollection("applicants"),
	}
}

func (fc *FormResponseController) HandleFormResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		FormId string              `json:"formId"`
		models.FormResponses
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Error parsing request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	projectID, err := primitive.ObjectIDFromHex(requestData.FormId)
	if err != nil {
		projectID = primitive.NewObjectID()
	}

	applicant := models.Applicant{
		ID:            primitive.NewObjectID(),
		ProjectID:     projectID,
		Elo:          1000,
		Wins:         0,
		Losses:       0,
		MatchesPlayed: []primitive.ObjectID{},
		Timestamp:    requestData.Timestamp,
	}

	bucket, err := gridfs.NewBucket(db.Client.Database("akpsi-ucsb"))
	if err != nil {
		http.Error(w, "Error creating GridFS bucket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, resp := range requestData.Responses {
		if err := processFormResponse(&applicant, resp, bucket); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := fc.collection.InsertOne(ctx, applicant)
	if err != nil {
		http.Error(w, "Error inserting document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Form response received successfully",
		"id":     result.InsertedID,
	})

	prettyJSON, _ := json.MarshalIndent(applicant, "", "    ")
	fmt.Printf("Application Received:\n%s\n", string(prettyJSON))
}

func processFormResponse(applicant *models.Applicant, resp models.Response, bucket *gridfs.Bucket) error {
	switch resp.Question {
	case "firstName":
		if strVal, ok := resp.Answer.(string); ok {
			applicant.FirstName = strVal
		}
	case "lastName":
		if strVal, ok := resp.Answer.(string); ok {
			applicant.LastName = strVal
		}
	case "major":
		if strVal, ok := resp.Answer.(string); ok {
			applicant.Major = strVal
		}
	case "year":
		if strVal, ok := resp.Answer.(string); ok {
			applicant.Year = strVal
		}
	case "coverLetter", "resume", "image":
		if err := processFileUpload(applicant, resp, bucket); err != nil {
			return err
		}
	}
	return nil
}

func processFileUpload(applicant *models.Applicant, resp models.Response, bucket *gridfs.Bucket) error {
	answer, ok := resp.Answer.(map[string]interface{})
	if !ok || answer["type"] != "file" {
		return fmt.Errorf("invalid file data format")
	}

	fileInfo, err := uploadFile(answer, bucket)
	if err != nil {
		return err
	}

	switch resp.Question {
	case "coverLetter":
		applicant.CoverLetter = fileInfo
	case "resume":
		applicant.Resume = fileInfo
	case "image":
		applicant.Image = fileInfo
	}
	return nil
}

func uploadFile(answer map[string]interface{}, bucket *gridfs.Bucket) (*models.FileInfo, error) {
	dataStr, ok := answer["data"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid file data format")
	}

	fileData, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return nil, fmt.Errorf("error decoding file data: %v", err)
	}

	timestamp := time.Now().Unix()
	uniqueFileName := fmt.Sprintf("%d_%s", timestamp, answer["filename"].(string))

	fileID, err := uploadToGridFS(bucket, uniqueFileName, fileData)
	if err != nil {
		return nil, err
	}

	return &models.FileInfo{
		FileID:      fileID.Hex(),
		FileName:    answer["filename"].(string),
		MimeType:    answer["mimeType"].(string),
		DriveFileID: answer["fileId"].([]interface{})[0].(string),
		UniqueName:  uniqueFileName,
		UploadedAt:  time.Now(),
	}, nil
}

func uploadToGridFS(bucket *gridfs.Bucket, filename string, data []byte) (primitive.ObjectID, error) {
	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("error opening upload stream: %v", err)
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(data)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("error writing to stream: %v", err)
	}

	return uploadStream.FileID.(primitive.ObjectID), nil
}