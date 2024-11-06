package echohttp

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"todo-list-api/internal/models"

	"github.com/labstack/echo/v4"
)

func (s *Server) CreateNote() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var note models.Note
		validate := validator.New()

		if err := ctx.Bind(&note); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		if err := validate.Struct(&note); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		if err := s.service.CreateNote(&note); err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return ctx.JSON(http.StatusCreated, note)
	}
}

func (s *Server) GetNotes() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var note models.NoteList
		page, _ := strconv.Atoi(ctx.QueryParam("page"))
		limit, _ := strconv.Atoi(ctx.QueryParam("limit"))

		if page&limit == 0 {
			page = 1
			limit = 10
		}

		if err := s.service.GetNotes(&note, page, limit); err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return ctx.JSON(http.StatusOK, note)
	}
}

func (s *Server) GetNotesByID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		//var note *models.Note
		id := ctx.Param("id")

		note, err := s.service.GetNoteByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ctx.JSON(http.StatusNotFound, map[string]string{"message": "note not found"})
			}
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return ctx.JSON(http.StatusOK, note)
	}
}

func (s *Server) UpdateNoteByID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var note *models.Note
		id := ctx.Param("id")

		if err := ctx.Bind(&note); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		err := s.service.UpdateNoteByID(id, note)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ctx.JSON(http.StatusNotFound, map[string]string{"message": "note not found"})
			}

			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return ctx.JSON(http.StatusOK, note)
	}
}

func (s *Server) DeleteNoteByID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id := ctx.Param("id")

		err := s.service.DeleteNoteByID(id)
		if err != nil {
			if err.Error() == "record not found" {
				return ctx.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
			}

			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return ctx.JSON(http.StatusNoContent, map[string]string{"message": "Successfully deleted note"})
	}
}
