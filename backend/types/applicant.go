package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ApplicantResponse represents the standardized applicant data structure
type Applicant struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName   string            `json:"firstName" bson:"firstName"`
	LastName    string            `json:"lastName" bson:"lastName"`
	Major       string            `json:"major" bson:"major"`
	Year        string            `json:"year" bson:"year"`
	Rating      float64           `json:"rating" bson:"rating"`
	RatingCount int               `json:"ratingCount" bson:"ratingCount"`
	Timestamp   string            `json:"timestamp" bson:"timestamp"`
	Resume      *FileInfo         `json:"resume,omitempty" bson:"resume,omitempty"`
	CoverLetter *FileInfo         `json:"coverLetter,omitempty" bson:"coverLetter,omitempty"`
	Image       *FileInfo         `json:"image,omitempty" bson:"image,omitempty"`
}

// FileInfo represents file metadata and content
type FileInfo struct {
	FileID      string    `json:"fileId" bson:"fileId"`
	FileName    string    `json:"fileName" bson:"fileName"`
	MimeType    string    `json:"mimeType" bson:"mimeType"`
	DriveFileID string  `json:"driveFileId" bson:"driveFileId"`
	UniqueName  string    `json:"uniqueName" bson:"uniqueName"`
	UploadedAt  time.Time `json:"uploadedAt" bson:"uploadedAt"`
	Data        string    `json:"data,omitempty" bson:"-"` 
}

type FormResponses struct {
	Timestamp string     `json:"submission_timestamp"`
	Responses []Response `json:"responses"`
}

type Response struct {
	Question string      `json:"question"`
	Answer   interface{} `json:"answer"`
}

type FileAnswer struct {
	Type        string    `json:"type"`
	DriveFileId []string  `json:"driveFileId"`
	FileId      string    `json:"fileId"`
	FileName    string    `json:"fileName"`
	MimeType    string    `json:"mimeType"`
	UniqueName  string    `json:"uniqueName"`
	UploadedAt  time.Time `json:"uploadedAt"`
} 