package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Project struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name                 string             `bson:"name" json:"name"`
	TotalApplicants      int                `bson:"totalApplicants" json:"totalApplicants"`
	CompletedComparisons int                `bson:"completedComparisons" json:"completedComparisons"`
	TotalComparisons     int                `bson:"totalComparisons" json:"totalComparisons"`
}
