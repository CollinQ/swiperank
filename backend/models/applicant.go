package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Applicant struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    FirstName string            `json:"first_name" bson:"first_name"`
    LastName  string            `json:"last_name" bson:"last_name"`
    ProjectID primitive.ObjectID `json:"project_id" bson:"project_id"`
    Wins      int               `json:"wins" bson:"wins"`
    Losses    int               `json:"losses" bson:"losses"`
    Elo       int               `json:"elo" bson:"elo"`
	MatchesPlayed []primitive.ObjectID `json:"matches_played" bson:"matches_played"`
}