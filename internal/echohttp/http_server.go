package echohttp

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"todo-list-api/internal/service"

	"github.com/labstack/echo/v4"
)

type Server struct {
	router  *echo.Echo
	service service.InterfaceService
}

// NewServer ...
func NewServer(router *echo.Echo, service service.InterfaceService) *Server {
	return &Server{
		router:  router,
		service: service,
	}
}

func (s *Server) Run() {

	main := s.router.Group("/api/v1")
	main.GET("/notes", s.GetNotes())
	main.POST("/notes", s.CreateNote())
	main.GET("/notes/:id", s.GetNotesByID())
	main.PUT("/notes/:id", s.UpdateNoteByID())
	main.DELETE("/notes/:id", s.DeleteNoteByID())

	main.GET("/notes/children", s.GetNoteChildren())
	main.POST("/notes/children", s.CreateNoteChild())
	main.GET("/notes/children/:id", s.GetNoteChildByID())
	main.PUT("/notes/children/:id", s.UpdateNoteChildByID())
	main.DELETE("/notes/children/:id", s.DeleteNoteChildByID())

	if err := s.router.Start(fmt.Sprintf(":%d", 8080)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("starting http server :%s", err)
	}
}
