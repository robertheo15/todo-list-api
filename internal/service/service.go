package service

import (
	"todo-list-api/internal/models"
	"todo-list-api/internal/repository"
)

type Service struct {
	repo repository.InterfaceRepository
}

// NewService creates a new isntance of product service.
func NewService(repo repository.InterfaceRepository) InterfaceService {
	return &Service{
		repo: repo,
	}
}

type InterfaceService interface {
	CreateNote(note *models.Note) error
	GetNotes(note *models.NoteList, page, limit int) error
	GetNoteByID(id string) (*models.Note, error)
	UpdateNoteByID(id string, updatedNote *models.Note) error
	DeleteNoteByID(id string) error

	CreateNoteChild(note *models.NoteChild) error
	GetNoteChildren(note *models.NoteChildList, page, limit int) error
	GetNoteChildByID(id string) (*models.NoteChild, error)
	UpdateNoteChildByID(id string, updatedNote *models.NoteChild) error
	DeleteNoteChildByID(id string) error
}
