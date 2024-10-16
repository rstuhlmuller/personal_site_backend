package models

import (
	"time"
)

const (
	CountItemID = "visitor_count"
)

type VisitorItem struct {
	ID        string    `json:"id" dynamodbav:"id"`
	Type      string    `json:"type" dynamodbav:"type"`
	Count     int       `json:"count,omitempty" dynamodbav:"count,omitempty"`
	IP        string    `json:"ip,omitempty" dynamodbav:"ip,omitempty"`
	UserAgent string    `json:"userAgent,omitempty" dynamodbav:"userAgent,omitempty"`
	Referer   string    `json:"referer,omitempty" dynamodbav:"referer,omitempty"`
	Timestamp time.Time `json:"timestamp" dynamodbav:"timestamp"`
}

func NewVisitorLog(ip, userAgent, referer string) *VisitorItem {
	return &VisitorItem{
		Type:      "log",
		IP:        ip,
		UserAgent: userAgent,
		Referer:   referer,
		Timestamp: time.Now().UTC(),
	}
}
