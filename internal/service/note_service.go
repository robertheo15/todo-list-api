package service

import (
	"log"
	"todo-list-api/internal/models"
)

func (s *Service) CreateNote(note *models.Note) error {
	return s.repo.Create(note)
}

func (s *Service) GetNotes(notes *models.NoteList, page, limit int) error {
	list, err := s.repo.List(page, limit)
	if err != nil {
		return err
	}

	notes.Notes = list.Notes

	return nil
}

func (s *Service) GetNoteByID(id string) (*models.Note, error) {
	note, err := s.repo.GetByID(id)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	return note, nil
}

func (s *Service) UpdateNoteByID(id string, updatedNote *models.Note) error {
	err := s.repo.Update(id, updatedNote)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteNoteByID(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
