package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID           string    `json:"id" dynamodbav:"id"`
	Title        string    `json:"title" dynamodbav:"title"`
	Description  string    `json:"description" dynamodbav:"description"`
	Technologies []string  `json:"technologies" dynamodbav:"technologies"`
	ImageURL     string    `json:"imageUrl" dynamodbav:"imageUrl"`
	GitHubURL    string    `json:"githubUrl" dynamodbav:"githubUrl"`
	LiveURL      string    `json:"liveUrl" dynamodbav:"liveUrl"`
	StartDate    time.Time `json:"startDate" dynamodbav:"startDate"`
	EndDate      time.Time `json:"endDate,omitempty" dynamodbav:"endDate,omitempty"`
	IsOngoing    bool      `json:"isOngoing" dynamodbav:"isOngoing"`
	CreatedAt    time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

func NewProject(title, description string, technologies []string) *Project {
	now := time.Now()
	return &Project{
		ID:           uuid.New().String(),
		Title:        title,
		Description:  description,
		Technologies: technologies,
		StartDate:    now,
		IsOngoing:    true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (p *Project) Update(updates map[string]interface{}) {
	// Update fields based on the provided updates
	for key, value := range updates {
		switch key {
		case "title":
			p.Title = value.(string)
		case "description":
			p.Description = value.(string)
		case "technologies":
			p.Technologies = value.([]string)
		case "imageUrl":
			p.ImageURL = value.(string)
		case "githubUrl":
			p.GitHubURL = value.(string)
		case "liveUrl":
			p.LiveURL = value.(string)
		case "endDate":
			p.EndDate = value.(time.Time)
			p.IsOngoing = false
		case "isOngoing":
			p.IsOngoing = value.(bool)
		}
	}
	p.UpdatedAt = time.Now()
}
