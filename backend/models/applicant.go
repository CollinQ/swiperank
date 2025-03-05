package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Applicant struct {
	ID            primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	FirstName     string              `json:"firstName" bson:"firstName"`
	LastName      string              `json:"lastName" bson:"lastName"`
	Major         string              `json:"major" bson:"major"`
	Year          string              `json:"year" bson:"year"`
	Timestamp     string              `json:"timestamp" bson:"timestamp"`
	ProjectID     primitive.ObjectID   `json:"project_id" bson:"project_id"`
	Wins          int                 `json:"wins" bson:"wins"`
	Losses        int                 `json:"losses" bson:"losses"`
	Elo           int                 `json:"elo" bson:"elo"`
	MatchesPlayed []primitive.ObjectID `json:"matches_played" bson:"matches_played"`
	Resume        *FileInfo           `json:"resume,omitempty" bson:"resume,omitempty"`
	CoverLetter   *FileInfo           `json:"coverLetter,omitempty" bson:"coverLetter,omitempty"`
	Image         *FileInfo           `json:"image,omitempty" bson:"image,omitempty"`
}

type FileInfo struct {
	FileID      string    `json:"fileId" bson:"fileId"`
	FileName    string    `json:"fileName" bson:"fileName"`
	MimeType    string    `json:"mimeType" bson:"mimeType"`
	DriveFileID string    `json:"driveFileId" bson:"driveFileId"`
	UniqueName  string    `json:"uniqueName" bson:"uniqueName"`
	UploadedAt  time.Time `json:"uploadedAt" bson:"uploadedAt"`
	Data        string    `json:"data,omitempty" bson:"-"`
}

type FormResponses struct {
	FormID    string     `json:"formId"`
	Timestamp string     `json:"submission_timestamp"`
	Responses []Response `json:"responses"`
}

type Response struct {
	Question string      `json:"question"`
	Answer   interface{} `json:"answer"`
}