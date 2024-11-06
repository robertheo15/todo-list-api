package echohttp_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-list-api/internal/echohttp"
	"todo-list-api/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the InterfaceService for testing purposes.
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateNote(note *models.Note) error {
	//TODO implement me
	return nil
}

func (m *MockService) DeleteNoteByID(id string) error {
	//TODO implement me
	return nil
}

func (m *MockService) CreateNoteChild(note *models.NoteChild) error {
	//TODO implement me
	return nil
}

func (m *MockService) GetNoteChildren(note *models.NoteChildList, page, limit int) error {
	//TODO implement me
	return nil
}

func (m *MockService) GetNoteChildByID(id string) (*models.NoteChild, error) {
	//TODO implement me
	return nil, nil
}

func (m *MockService) UpdateNoteChildByID(id string, updatedNote *models.NoteChild) error {
	//TODO implement me
	return nil
}

func (m *MockService) DeleteNoteChildByID(id string) error {
	//TODO implement me
	return nil
}

func (m *MockService) GetNoteByID(id string) (*models.Note, error) {
	args := m.Called(id)
	note, _ := args.Get(0).(*models.Note)
	return note, args.Error(1)
}

func (m *MockService) UpdateNoteByID(id string, updatedNote *models.Note) error {
	args := m.Called(id, updatedNote)
	return args.Error(0)
}

func (m *MockService) DeleteTodoByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockService) GetNotes(note *models.NoteList, page, limit int) error {
	args := m.Called(note, page, limit)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	// Populate noteList with mock data
	note.Notes = []*models.Note{
		{ID: "9fb800cd-6a3e-4f49-8478-f38ff7e5ec6e", Title: "Note 1", Description: "Description 1"},
		{ID: "2deaff79-b18f-44d4-b5a6-d63317e73056", Title: "Note 2", Description: "Description 2"},
	}
	note.Total = 2
	note.Count = 2
	note.TotalPage = 1

	return nil
}

func (m *MockService) CreateTodo(note *models.Note) error {
	args := m.Called(note)
	return args.Error(0)
}

func TestCreateNoteHandler(t *testing.T) {
	e := echo.New()
	mockService := new(MockService)
	server := echohttp.NewServer(e, mockService)

	t.Run("should return 400 if binding fails", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/notes", bytes.NewReader([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		// Invoke the handler function directly with ctx
		handler := server.CreateNote()
		if assert.NoError(t, handler(ctx)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "error")
		}
	})

	t.Run("should return 500 if service CreateNote fails", func(t *testing.T) {
		note := &models.Note{Title: "Test Note", Description: "Test Content"}
		noteJSON, _ := json.Marshal(note)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/notes", bytes.NewReader(noteJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		mockService.On("CreateNote", note).Return(errors.New("service error")).Once()

		// Invoke the handler function directly with ctx
		handler := server.CreateNote()
		if assert.NoError(t, handler(ctx)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Contains(t, rec.Body.String(), "error")
		}

		mockService.AssertExpectations(t)
	})

	t.Run("should return 201 if note is created successfully", func(t *testing.T) {
		note := &models.Note{Title: "Test Note", Description: "Test Content"}
		noteJSON, _ := json.Marshal(note)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/notes", bytes.NewReader(noteJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		mockService.On("CreateNote", note).Return(nil).Once()

		// Invoke the handler function directly with ctx
		handler := server.CreateNote()
		if assert.NoError(t, handler(ctx)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			var createdNote models.Note
			if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &createdNote)) {
				assert.Equal(t, note.Title, createdNote.Title)
				assert.Equal(t, note.Description, createdNote.Description)
			}
		}

		mockService.AssertExpectations(t)
	})
}

func TestUpdateNoteByIDHandler(t *testing.T) {
	e := echo.New()
	mockService := new(MockService)
	server := echohttp.NewServer(e, mockService)

	t.Run("should return 404 if note not found", func(t *testing.T) {
		note := &models.Note{Title: "Updated Note", Description: "Updated Description"}
		noteJSON, _ := json.Marshal(note)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/notes/:id", bytes.NewReader(noteJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		mockService.On("UpdateNoteByID", "1", note).Return(gorm.ErrRecordNotFound).Once()

		handler := server.UpdateNoteByID()
		if assert.NoError(t, handler(ctx)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.Contains(t, rec.Body.String(), "note not found")
		}
	})

	t.Run("should return 200 if note is updated successfully", func(t *testing.T) {
		note := &models.Note{Title: "Updated Note", Description: "Updated Description"}
		noteJSON, _ := json.Marshal(note)
		req := httptest.NewRequest(http.MethodPut, "/note/:id", bytes.NewReader(noteJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		mockService.On("UpdateNoteByID", "1", note).Return(nil).Once()

		handler := server.UpdateNoteByID()
		if assert.NoError(t, handler(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}

func TestDeleteNoteByIDHandler(t *testing.T) {
	e := echo.New()
	mockService := new(MockService)
	server := echohttp.NewServer(e, mockService)

	t.Run("should return 404 if note not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/notes/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		mockService.On("DeleteNoteByID", "1").Return(errors.New("note not found")).Once()

		handler := server.DeleteNoteByID()
		err := handler(ctx)

		// Ensure error is handled without crashing
		assert.NoError(t, err)

		// Check response status and body
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "note not found")
	})

	t.Run("should return 204 if note is deleted successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/notes/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		mockService.On("DeleteNoteByID", "1").Return(nil).Once()

		handler := server.DeleteNoteByID()
		err := handler(ctx)

		// Ensure error is handled without crashing
		assert.NoError(t, err)

		// Check response status
		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())
	})
}

func TestGetNotesHandler(t *testing.T) {
	e := echo.New()
	mockService := new(MockService)
	server := echohttp.NewServer(e, mockService)

	t.Run("should return 200 and list of notes", func(t *testing.T) {
		// Prepare request and recorder
		req := httptest.NewRequest(http.MethodGet, "/api/v1/notes?page=1&limit=10", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		// Set up mock expectations
		mockService.On("GetNotes", mock.AnythingOfType("*models.NoteList"), 1, 10).Run(func(args mock.Arguments) {
			noteList := args.Get(0).(*models.NoteList)
			noteList.Notes = []*models.Note{
				{ID: "9fb800cd-6a3e-4f49-8478-f38ff7e5ec6e", Title: "Note 1", Description: "Description 1"},
				{ID: "2deaff79-b18f-44d4-b5a6-d63317e73056", Title: "Note 2", Description: "Description 2"},
			}
			noteList.Total = 2
			noteList.Count = 2
			noteList.TotalPage = 1
			noteList.CurrentPage = 1
		}).Return(nil)

		// Invoke the handler function
		handler := server.GetNotes()
		if assert.NoError(t, handler(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var returnedNoteList models.NoteList
			if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &returnedNoteList)) {
				assert.Equal(t, int32(2), returnedNoteList.Total)
				assert.Len(t, returnedNoteList.Notes, 2)
				assert.Equal(t, "Note 1", returnedNoteList.Notes[0].Title)
				assert.Equal(t, "Note 2", returnedNoteList.Notes[1].Title)
			}
		}

		// Assert that the mock expectations were met
		mockService.AssertExpectations(t)
	})
}

func TestGetNoteByIDHandler(t *testing.T) {
	e := echo.New()
	mockService := new(MockService)
	server := echohttp.NewServer(e, mockService)

	t.Run("should return 200 if note is found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/notes/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		expectedNote := &models.Note{ID: "1", Title: "Test Note", Description: "Test Description"}
		mockService.On("GetNotesByID", "1").Return(expectedNote, nil).Once()

		handler := server.GetNotesByID()
		err := handler(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var returnedNote models.Note
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &returnedNote))
		assert.Equal(t, expectedNote.ID, returnedNote.ID)
		assert.Equal(t, expectedNote.Title, returnedNote.Title)
		assert.Equal(t, expectedNote.Description, returnedNote.Description)

		mockService.AssertExpectations(t)
	})

	t.Run("should return 404 if note is not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/notes/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Define a "note not found" error
		mockService.On("GetNotesByID", "1").Return(nil, gorm.ErrRecordNotFound).Once()

		handler := server.GetNotesByID()
		err := handler(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "note not found")

		mockService.AssertExpectations(t)
	})

	t.Run("should return 500 if there is an internal server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/notes/1", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		// Define a generic error to simulate an internal server error
		mockService.On("GetNotesByID", "1").Return(nil, errors.New("internal server error")).Once()

		handler := server.GetNotesByID()
		err := handler(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "internal server error")

		mockService.AssertExpectations(t)
	})
}
