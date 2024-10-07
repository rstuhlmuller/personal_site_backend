package models

import (
	"time"
)

type VisitorCount struct {
	ID    string    `json:"id" dynamodbav:"id"`
	Count int       `json:"count" dynamodbav:"count"`
	Date  time.Time `json:"date" dynamodbav:"date"`
}

func NewVisitorCount() *VisitorCount {
	return &VisitorCount{
		ID:    "visitor_count",
		Count: 0,
		Date:  time.Now().UTC(),
	}
}
