package repository

import (
	"errors"
	"todo-list-api/internal/models"

	"gorm.io/gorm"
)

// Create Note
func (repo *PostgresRepository) Create(note *models.Note) error {
	return repo.db.Create(note).Error
}

// Get Note by ID
func (repo *PostgresRepository) GetByID(id string) (*models.Note, error) {
	var note models.Note
	err := repo.db.Preload("NoteChildren").Preload("NoteChildren.NoteChildFiles").
		Preload("NoteFiles").First(&note, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("note not found")
	}
	return &note, err
}

// Update Note
func (repo *PostgresRepository) Update(id string, updatedNote *models.Note) error {
	var note models.Note
	if err := repo.db.First(&note, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("note not found")
		}
		return err
	}
	return repo.db.Model(&note).Updates(updatedNote).Error
}

// Delete Note
func (repo *PostgresRepository) Delete(id string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Delete related NoteChild and NoteFile records
		if err := tx.Where("note_id = ?", id).Delete(&models.NoteChild{}).Error; err != nil {
			return err
		}
		if err := tx.Where("note_id = ?", id).Delete(&models.NoteFile{}).Error; err != nil {
			return err
		}
		// Delete the main Note record
		return tx.Delete(&models.Note{}, "id = ?", id).Error
	})
}

// List Notes with Pagination
func (repo *PostgresRepository) List(page int, limit int) (*models.NoteList, error) {
	var notes []*models.Note
	var total int64

	// Count total records
	repo.db.Model(&models.Note{}).Count(&total)

	// Calculate pagination
	offset := (page - 1) * limit
	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	// Get notes with pagination
	err := repo.db.Preload("NoteChildren").Preload("NoteChildren.NoteChildFiles").
		Preload("NoteFiles").Offset(offset).Limit(limit).Find(&notes).Error
	if err != nil {
		return nil, err
	}

	return &models.NoteList{
		Total:       int32(total),
		Count:       int32(len(notes)),
		TotalPage:   totalPages,
		CurrentPage: int32(page),
		Next:        "", // Add next URL if needed
		Notes:       notes,
	}, nil
}
