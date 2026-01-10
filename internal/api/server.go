package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/leopoldhub/royal-api-personal/internal/api/handlers"
	"github.com/leopoldhub/royal-api-personal/internal/api/middleware"
	"github.com/leopoldhub/royal-api-personal/internal/database/repository"
)

// Server represents the HTTP API server
type Server struct {
	db       *sql.DB
	apiToken string
	router   *http.ServeMux
	logger   *log.Logger
}

// NewServer creates a new API server
func NewServer(db *sql.DB, apiToken string, logger *log.Logger) *Server {
	if logger == nil {
		logger = log.Default()
	}

	s := &Server{
		db:       db,
		apiToken: apiToken,
		router:   http.NewServeMux(),
		logger:   logger,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	battleRepo := repository.NewBattleRepository(s.db)
	metaRepo := repository.NewMetaDeckRepository(s.db)

	healthHandler := handlers.NewHealthHandler(s.db, battleRepo, metaRepo)
	deckHandler := handlers.NewDeckHandler(battleRepo, metaRepo)
	statsHandler := handlers.NewStatsHandler(battleRepo, metaRepo)

	s.router.HandleFunc("GET /health", healthHandler.Handle)

	s.router.HandleFunc("GET /decks/meta", s.protected(deckHandler.GetMetaDecks))
	s.router.HandleFunc("GET /decks/{signature}", s.protected(deckHandler.GetDeckBySignature))
	s.router.HandleFunc("GET /stats/summary", s.protected(statsHandler.GetSummary))
}

func (s *Server) protected(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.Auth(s.apiToken)(http.HandlerFunc(handler)).ServeHTTP(w, r)
	}
}

// Start starts the HTTP server
func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	s.logger.Printf("Starting API server on %s", addr)

	handler := middleware.Logging(s.logger)(middleware.JSON(s.router))

	return http.ListenAndServe(addr, handler)
}
