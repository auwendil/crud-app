package models

import "time"

type Book struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
