package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"backend/db"
	"backend/elo"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ApplicantController struct {
	collection *mongo.Collection
}

func NewApplicantController() *ApplicantController {
	return &ApplicantController{
		collection: db.GetCollection("applicants"),
	}
}

// elo and comparison helper functions

func contains(matchesPlayed []primitive.ObjectID, id primitive.ObjectID) bool {
	for _, matchID := range matchesPlayed {
		if matchID == id {
			return true
		}
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func resetMatchHistory(ctx context.Context, collection *mongo.Collection) {
	_, _ = collection.UpdateMany(ctx, bson.M{}, bson.M{"$set": bson.M{"matches_played": []primitive.ObjectID{}}})
	log.Println("Reset all applicants' match history")
}

func (ac *ApplicantController) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := ac.collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch applicants", http.StatusInternalServerError)
		log.Println("MongoDB Find applicants error: ", err)
		return
	}
	defer cursor.Close(ctx)

	var applicants []models.Applicant
	if err = cursor.All(ctx, &applicants); err != nil {
		http.Error(w, "Error decoding applicants", http.StatusInternalServerError)
		log.Println("Cursor decode error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applicants)
}

func (ac *ApplicantController) GetById(w http.ResponseWriter, r *http.Request) {
	// log.Println("Fetching applicant by ID")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get ID from URL path parameter
	applicantIDStr := r.URL.Query().Get("id")
	// log.Println("Applicant ID from path:", applicantIDStr)

	if applicantIDStr == "" {
		http.Error(w, "Applicant ID required", http.StatusBadRequest)
		return
	}

	applicantID, err := primitive.ObjectIDFromHex(applicantIDStr)
	if err != nil {
		// log.Println("Error converting applicant ID to ObjectID:", err)
		http.Error(w, "Invalid Applicant ID", http.StatusBadRequest)
		return
	}

	var applicant models.Applicant
	err = ac.collection.FindOne(ctx, bson.M{"_id": applicantID}).Decode(&applicant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// log.Println("Applicant not found in database")
			http.Error(w, "Applicant not found", http.StatusNotFound)
			return
		}
		// log.Println("Failed to fetch applicant from database:", err)
		http.Error(w, "Failed to fetch applicant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applicant)
	// log.Println("Applicant fetched successfully")
}

func (ac *ApplicantController) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement creating new applicant
}

func (ac *ApplicantController) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement updating applicant
}

func (ac *ApplicantController) GetTwoForComparison(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "elo", Value: -1}})
	cursor, err := ac.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		http.Error(w, "Failed to fetch applicants", http.StatusInternalServerError)
		log.Println("MongoDB Find applicants error:", err)
		return
	}
	defer cursor.Close(ctx)

	var applicants []models.Applicant
	if err = cursor.All(ctx, &applicants); err != nil {
		http.Error(w, "Error decoding applicants", http.StatusInternalServerError)
		log.Println("Cursor decode error:", err)
		return
	}

	if len(applicants) < 2 {
		http.Error(w, "Not enough applicants for comparison", http.StatusInternalServerError)
		return
	}

	var applicant1, applicant2 models.Applicant
	minEloDiff := int(^uint(0) >> 1)
	for i := 0; i < len(applicants) - 1; i++ {
		for j := i + 1; j < len(applicants); j++ {
			if contains(applicants[i].MatchesPlayed, applicants[j].ID) {
				continue
			}

			diff := abs(applicants[i].Elo - applicants[j].Elo)
			if diff < minEloDiff {
				minEloDiff = diff
				applicant1 = applicants[i]
				applicant2 = applicants[j]
			}
		}
	}

	if applicant1.ID.IsZero() || applicant2.ID.IsZero() {
		resetMatchHistory(ctx, ac.collection)
		http.Error(w, "All applicants have already played, match history reset", http.StatusConflict)
		return
	}

	// Append file data to Applicant response
	bucket, err := gridfs.NewBucket(db.Client.Database("akpsi-ucsb"))
	if err != nil {
		http.Error(w, "Error creating GridFS bucket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	applicant1.Image = fetchFile(bucket, applicant1.Image)
	applicant1.CoverLetter = fetchFile(bucket, applicant1.CoverLetter)
	applicant1.Resume = fetchFile(bucket, applicant1.Resume)

	applicant2.Image = fetchFile(bucket, applicant2.Image)
	applicant2.CoverLetter = fetchFile(bucket, applicant2.CoverLetter)
	applicant2.Resume = fetchFile(bucket, applicant2.Resume)


	applicant1.MatchesPlayed = append(applicant1.MatchesPlayed, applicant2.ID)
	applicant2.MatchesPlayed = append(applicant2.MatchesPlayed, applicant1.ID)

	_, _ = ac.collection.UpdateOne(ctx, bson.M{"_id": applicant1.ID}, bson.M{"$set": bson.M{"matches_played" : applicant1.MatchesPlayed}})
	_, _ = ac.collection.UpdateOne(ctx, bson.M{"_id": applicant2.ID}, bson.M{"$set": bson.M{"matches_played" : applicant2.MatchesPlayed}})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]models.Applicant{applicant1, applicant2})
}

func (ac *ApplicantController) UpdateElo(w http.ResponseWriter, r *http.Request) {
	var request struct {
		WinnerID string `json:"winnerId"`
		LoserID  string `json:"loserId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Convert string IDs to ObjectIDs
	winnerID, err := primitive.ObjectIDFromHex(request.WinnerID)
	if err != nil {
		http.Error(w, "Invalid Winner ID format", http.StatusBadRequest)
		return
	}

	loserID, err := primitive.ObjectIDFromHex(request.LoserID)
	if err != nil {
		http.Error(w, "Invalid Loser ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var winner, loser models.Applicant
	if err := ac.collection.FindOne(ctx, bson.M{"_id": winnerID}).Decode(&winner); err != nil {
		http.Error(w, "Winner not found", http.StatusNotFound)
		return
	}
	if err := ac.collection.FindOne(ctx, bson.M{"_id": loserID}).Decode(&loser); err != nil {
		http.Error(w, "Loser not found", http.StatusNotFound)
		return
	}

	winnerElo, loserElo := elo.CalculateElo(winner.Elo, loser.Elo, true)

	winner.Elo = winnerElo
	loser.Elo = loserElo

	winner.Wins += 1
	loser.Losses += 1

	updateWinner := bson.M{"$set": bson.M{"elo": winner.Elo}, "$inc": bson.M{"wins": 1}}
	updateLoser := bson.M{"$set": bson.M{"elo": loser.Elo}, "$inc": bson.M{"losses": 1}}
	
	if _, err := ac.collection.UpdateOne(ctx, bson.M{"_id": winnerID}, updateWinner); err != nil {
		http.Error(w, "Failed to update winner", http.StatusInternalServerError)
		return
	}
	if _, err := ac.collection.UpdateOne(ctx, bson.M{"_id": loserID}, updateLoser); err != nil {
		http.Error(w, "Failed to update loser", http.StatusInternalServerError)
		return
	}	

    w.WriteHeader(http.StatusOK)
}

// Helper function for file fetching
func fetchFile(bucket *gridfs.Bucket, fileInfo *models.FileInfo) *models.FileInfo {
	if fileInfo != nil {
		fileID, err := primitive.ObjectIDFromHex(fileInfo.FileID)
		if err != nil {
			return nil
		}

		var buf bytes.Buffer
		_, err = bucket.DownloadToStream(fileID, &buf)
		if err != nil {
			return nil
		}

		fileInfo.Data = base64.StdEncoding.EncodeToString(buf.Bytes())
	}
	return fileInfo
}

func (ac *ApplicantController) GetRankings(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	projectIDStr := r.URL.Query().Get("project_id")
	if projectIDStr == "" {
		http.Error(w, "Project ID required", http.StatusBadRequest)
		return
	}

	projectID, err := primitive.ObjectIDFromHex(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid Project ID", http.StatusBadRequest)
		return
	}

	log.Println("Project ID: ", projectID)
	// Find all applicants for this project, sorted by Elo
	opts := options.Find().SetSort(bson.D{{Key: "elo", Value: -1}})
	cursor, err := ac.collection.Find(ctx, bson.M{"project_id": projectID}, opts)
	if err != nil {
		http.Error(w, "Failed to fetch rankings", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var rankings []models.Applicant
	if err = cursor.All(ctx, &rankings); err != nil {
		http.Error(w, "Error decoding rankings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rankings)
}

