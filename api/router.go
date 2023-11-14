package api

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hellofresh/health-go/v5"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"golang-demo/api/user"
)

func NewRouter(db *sql.DB, conn *amqp091.Connection, h *health.Health) *chi.Mux {
	userRepository := user.NewRepository(db)
	mQ := user.NewMQ(conn)
	userService := user.NewService(userRepository, mQ)
	userHandler := user.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Use(LoggerWithLevel(log.InfoLevel))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.Get)
		r.Post("/", userHandler.Store)
		r.Route("/{userId}", func(r chi.Router) {
			r.Use(userHandler.UserCtx)
			r.Get("/", userHandler.GetByID)
			r.Put("/", userHandler.Update)
			r.Delete("/", userHandler.Delete)
		})
	})
	r.Get("/status", h.HandlerFunc)
	return r
}
