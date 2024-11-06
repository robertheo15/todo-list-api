package service

import (
	"log"
	"todo-list-api/internal/models"
)

func (s *Service) CreateNoteChild(noteChild *models.NoteChild) error {
	return s.repo.CreateNoteChild(noteChild)
}

func (s *Service) GetNoteChildren(noteChild *models.NoteChildList, page, limit int) error {
	list, err := s.repo.ListChildren(page, limit)
	if err != nil {
		return err
	}

	noteChild.NoteChildren = list.NoteChildren

	return nil
}

func (s *Service) GetNoteChildByID(id string) (*models.NoteChild, error) {
	noteChild, err := s.repo.GetNoteChildByID(id)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	return noteChild, nil
}

func (s *Service) UpdateNoteChildByID(id string, updatedNoteChild *models.NoteChild) error {
	err := s.repo.UpdateNoteChildByID(id, updatedNoteChild)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteNoteChildByID(id string) error {
	err := s.repo.DeleteNoteChildByID(id)
	if err != nil {
		return err
	}

	return nil
}
