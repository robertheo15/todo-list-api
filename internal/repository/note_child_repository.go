package repository

import (
	"errors"
	"todo-list-api/internal/models"

	"gorm.io/gorm"
)

func (repo *PostgresRepository) CreateNoteChild(noteChild *models.NoteChild) error {
	return repo.db.Create(noteChild).Error
}

func (repo *PostgresRepository) GetNoteChildByID(id string) (*models.NoteChild, error) {
	var noteChild models.NoteChild
	err := repo.db.First(&noteChild, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("note not found")
	}
	return &noteChild, err
}

// Update Note
func (repo *PostgresRepository) UpdateNoteChildByID(id string, updatedNote *models.NoteChild) error {
	var note models.NoteChild
	if err := repo.db.First(&note, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("note not found")
		}
		return err
	}
	return repo.db.Model(&note).Updates(updatedNote).Error
}

// Delete Note
func (repo *PostgresRepository) DeleteNoteChildByID(id string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Delete related NoteChild and NoteFile records
		if err := tx.Where("note_id = ?", id).Delete(&models.NoteChild{}).Error; err != nil {
			return err
		}
		// Delete the main Note record
		return tx.Delete(&models.NoteChild{}, "id = ?", id).Error
	})
}

// List Notes with Pagination
func (repo *PostgresRepository) ListChildren(page int, limit int) (*models.NoteChildList, error) {
	var noteChildren []*models.NoteChild
	var total int64

	// Count total records
	repo.db.Model(&models.Note{}).Count(&total)

	// Calculate pagination
	offset := (page - 1) * limit
	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	// Get notes with pagination
	err := repo.db.Offset(offset).Limit(limit).Find(&noteChildren).Error
	if err != nil {
		return nil, err
	}

	return &models.NoteChildList{
		Total:        int32(total),
		Count:        int32(len(noteChildren)),
		TotalPage:    totalPages,
		CurrentPage:  int32(page),
		Next:         "", // Add next URL if needed
		NoteChildren: noteChildren,
	}, nil
}
