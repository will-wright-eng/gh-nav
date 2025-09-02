package models

import "time"

// Repository represents a GitHub repository
type Repository struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	LastUpdated time.Time `json:"last_updated"`
}
