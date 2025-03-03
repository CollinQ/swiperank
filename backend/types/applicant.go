package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ApplicantResponse represents the standardized applicant data structure
type ApplicantResponse struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	FirstName     string            `json:"first_name" bson:"first_name"`
	LastName      string            `json:"last_name" bson:"last_name"`
	Major         string            `json:"major" bson:"major"`
	Year          string            `json:"year" bson:"year"`
	CoverLetter   string            `json:"coverLetter" bson:"coverLetter"`
	Resume        *FileInfo         `json:"resume,omitempty" bson:"resume,omitempty"`
	Image         *FileInfo         `json:"image,omitempty" bson:"image,omitempty"`
	Rating        int               `json:"rating" bson:"rating"`
	RatingCounter int               `json:"rating_counter" bson:"rating_counter"`
	Timestamp     string            `json:"timestamp" bson:"submission_timestamp"`
}

// FileInfo represents file metadata and content
type FileInfo struct {
	FileID     primitive.ObjectID `json:"fileId" bson:"fileId"`
	FileName   string            `json:"fileName" bson:"fileName"`
	MimeType   string            `json:"mimeType" bson:"mimeType"`
	Data       string            `json:"data,omitempty" bson:"data,omitempty"`
} 