package repository

import (
	"todo-list-api/internal/models"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) InterfaceRepository {
	return &PostgresRepository{
		db: db,
	}
}

type InterfaceRepository interface {
	Create(note *models.Note) error
	List(page int, limit int) (*models.NoteList, error)
	GetByID(id string) (*models.Note, error)
	Update(id string, updatedNote *models.Note) error
	Delete(id string) error

	CreateNoteChild(note *models.NoteChild) error
	ListChildren(page int, limit int) (*models.NoteChildList, error)
	GetNoteChildByID(id string) (*models.NoteChild, error)
	UpdateNoteChildByID(id string, updatedNote *models.NoteChild) error
	DeleteNoteChildByID(id string) error
}
