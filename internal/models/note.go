package models

import (
	"time"
)

type Note struct {
	ID           string      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title        string      `gorm:"type:varchar" json:"title" validate:"required,max=100"`
	Description  string      `gorm:"type:text" json:"description" validate:"required,max=1000"`
	Type         string      `gorm:"type:varchar" json:"type"`
	NoteChildren []NoteChild `gorm:"foreignKey:NoteID" json:"note_children"`
	NoteFiles    []NoteFile  `gorm:"foreignKey:NoteID" json:"note_files"`
	CreatedAt    time.Time   `json:"created_at"`
}

type NoteChild struct {
	ID             string           `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	NoteID         string           `gorm:"type:uuid;index" json:"note_id"`
	Title          string           `gorm:"type:varchar" json:"title" validate:"required,max=100"`
	Description    string           `gorm:"type:text" json:"description" validate:"required,max=1000"`
	Type           string           `json:"type"`
	NoteChildFiles []NoteChildFiles `gorm:"foreignKey:NoteChildID" json:"note_child_files"`
	CreatedAt      time.Time        `json:"created_at"`
}

type NoteFile struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	NoteID    string    `gorm:"type:uuid;index" json:"note_id"`
	Name      string    `gorm:"type:varchar" json:"name"`
	Path      string    `gorm:"type:varchar" json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

type NoteChildFiles struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	NoteChildID string    `gorm:"type:uuid;index" json:"note_child_id"`
	Name        string    `gorm:"type:varchar" json:"name"`
	Path        string    `gorm:"type:varchar" json:"path"`
	CreatedAt   time.Time `json:"created_at"`
}

type NoteList struct {
	Total       int32   `json:"total"`
	Count       int32   `json:"count"`
	TotalPage   int32   `json:"total_page"`
	CurrentPage int32   `json:"current_page"`
	Next        string  `json:"next"`
	Notes       []*Note `json:"notes"`
}

type NoteChildList struct {
	Total        int32        `json:"total"`
	Count        int32        `json:"count"`
	TotalPage    int32        `json:"total_page"`
	CurrentPage  int32        `json:"current_page"`
	Next         string       `json:"next"`
	NoteChildren []*NoteChild `json:"notes"`
}
