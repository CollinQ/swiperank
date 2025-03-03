package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"backend/db"
	"backend/elo"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	// TODO: Implement getting single applicant by ID
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

	opts := options.Find().SetSort(bson.D{{Key: "elo", Value: -1}})
	cursor, err := ac.collection.Find(ctx, bson.M{"project_id": projectID}, opts)
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
	for i := range len(applicants) - 1 {
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

	applicant1.MatchesPlayed = append(applicant1.MatchesPlayed, applicant2.ID)
	applicant2.MatchesPlayed = append(applicant2.MatchesPlayed, applicant1.ID)

	_, _ = ac.collection.UpdateOne(ctx, bson.M{"_id": applicant1.ID}, bson.M{"$set": bson.M{"matches_played" : applicant1.MatchesPlayed}})
	_, _ = ac.collection.UpdateOne(ctx, bson.M{"_id": applicant2.ID}, bson.M{"$set": bson.M{"matches_played" : applicant2.MatchesPlayed}})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]models.Applicant{applicant1, applicant2})

}

func (ac *ApplicantController) UpdateElo(w http.ResponseWriter, r *http.Request) {
	var result struct {
		WinnerID primitive.ObjectID `json:"winner_id"`
		LoserID primitive.ObjectID `json:"loser_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var winner, loser models.Applicant
	if err := ac.collection.FindOne(ctx, bson.M{"_id": result.WinnerID}).Decode(&winner); err != nil {
		http.Error(w, "Winner not found", http.StatusNotFound)
		return
	}
	if err := ac.collection.FindOne(ctx, bson.M{"_id": result.LoserID}).Decode(&loser); err != nil {
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
	
	if _, err := ac.collection.UpdateOne(ctx, bson.M{"_id": winner.ID}, updateWinner); err != nil {
		http.Error(w, "Failed to update winner", http.StatusInternalServerError)
		return
	}
	if _, err := ac.collection.UpdateOne(ctx, bson.M{"_id": loser.ID}, updateLoser); err != nil {
		http.Error(w, "Failed to update loser", http.StatusInternalServerError)
		return
	}	

    w.WriteHeader(http.StatusOK)
}

