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

func (s *Server) CreateNoteChild() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var noteChild models.NoteChild
		validate := validator.New()

		if err := ctx.Bind(&noteChild); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		if err := validate.Struct(&noteChild); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		if err := s.service.CreateNoteChild(&noteChild); err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return ctx.JSON(http.StatusCreated, noteChild)
	}
}

func (s *Server) GetNoteChildren() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var noteChildList models.NoteChildList
		page, _ := strconv.Atoi(ctx.QueryParam("page"))
		limit, _ := strconv.Atoi(ctx.QueryParam("limit"))

		if page&limit == 0 {
			page = 1
			limit = 10
		}

		if err := s.service.GetNoteChildren(&noteChildList, page, limit); err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return ctx.JSON(http.StatusOK, noteChildList)
	}
}

func (s *Server) GetNoteChildByID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id := ctx.Param("id")

		noteChild, err := s.service.GetNoteChildByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ctx.JSON(http.StatusNotFound, map[string]string{"message": "note not found"})
			}
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return ctx.JSON(http.StatusOK, noteChild)
	}
}

func (s *Server) UpdateNoteChildByID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var noteChild *models.NoteChild
		id := ctx.Param("id")

		if err := ctx.Bind(&noteChild); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		err := s.service.UpdateNoteChildByID(id, noteChild)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ctx.JSON(http.StatusNotFound, map[string]string{"message": "note not found"})
			}

			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return ctx.JSON(http.StatusOK, noteChild)
	}
}

func (s *Server) DeleteNoteChildByID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id := ctx.Param("id")

		err := s.service.DeleteNoteChildByID(id)
		if err != nil {
			if err.Error() == "record not found" {
				return ctx.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
			}

			return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return ctx.JSON(http.StatusNoContent, map[string]string{"message": "Successfully deleted note"})
	}
}
