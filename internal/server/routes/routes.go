package routes

import (
	"context"
	"net/http"

	"github.com/MihailSergeenkov/GophKeeper/internal/models"
	"github.com/MihailSergeenkov/GophKeeper/internal/server/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Storager интерфейс для обработчиков HTTP запросов приложения.
type Handlerer interface {
	Ping() http.HandlerFunc
	RegisterUser() http.HandlerFunc
	CreateUserToken() http.HandlerFunc
	FetchUserData() http.HandlerFunc
	GetPassword() http.HandlerFunc
	AddPassword() http.HandlerFunc
	GetCard() http.HandlerFunc
	AddCard() http.HandlerFunc
	GetText() http.HandlerFunc
	AddText() http.HandlerFunc
	GetFile() http.HandlerFunc
	AddFile() http.HandlerFunc
}

// Storager интерфейс для хранилища данных.
type Storager interface {
	GetUserByID(ctx context.Context, userID int) (models.User, error)
}

var (
	JSONContentType     = "application/json"
	FormDataContentType = "multipart/form-data"
)

// NewRouter функция инициализации роутинга.
func NewRouter(h Handlerer, settings *config.Settings, l *zap.Logger, s Storager) chi.Router {
	r := chi.NewRouter()

	r.Get("/ping", h.Ping())

	r.Route("/api/user", func(r chi.Router) {
		r.Use(withRequestLogging(l))

		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType(FormDataContentType))
			r.Use(authMiddleware(settings, l, s))

			r.Route("/files", func(r chi.Router) {
				r.Get("/{fileID}", h.GetFile())
				r.Post("/", h.AddFile())
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType(JSONContentType))

			r.Post("/register", h.RegisterUser())
			r.Post("/token", h.CreateUserToken())

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware(settings, l, s))

				r.Get("/data", h.FetchUserData())

				r.Route("/passwords", func(r chi.Router) {
					r.Get("/{passwordID}", h.GetPassword())
					r.Post("/", h.AddPassword())
				})

				r.Route("/cards", func(r chi.Router) {
					r.Get("/{cardID}", h.GetCard())
					r.Post("/", h.AddCard())
				})

				r.Route("/texts", func(r chi.Router) {
					r.Get("/{textID}", h.GetText())
					r.Post("/", h.AddText())
				})
			})
		})
	})

	return r
}
