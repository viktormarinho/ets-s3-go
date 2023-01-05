package models

import "time"

type File struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Type      string    `json:"type"`
	Path      string    `json:"path"`
	Meta      *string   `json:"meta"`
	CreatedAt time.Time `json:"createdAt"`
	UserId    string    `json:"userId"`
}
