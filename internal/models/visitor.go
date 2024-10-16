package models

import (
	"time"
)

const (
	CountItemID = "visitor_count"
)

type VisitorItem struct {
	ID        string    `json:"id" dynamodbav:"id"`
	Count     int       `json:"count,omitempty" dynamodbav:"count,omitempty"`
	// IP        string    `json:"ip,omitempty" dynamodbav:"ip,omitempty"`
	Referer   string    `json:"referer,omitempty" dynamodbav:"referer,omitempty"`
	Timestamp time.Time `json:"timestamp" dynamodbav:"timestamp"`
}

func NewVisitorLog(referer string) *VisitorItem {
	return &VisitorItem{
		// IP:        ip,
		Referer:   referer,
		Timestamp: time.Now().UTC(),
	}
}
