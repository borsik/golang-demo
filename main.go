package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"golang-demo/handler"
	"golang-demo/user"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	config, err := NewConfig()
	if err != nil {
		log.Fatalln("failed to read config", err)
	}

	dbDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		config.DbHost, config.DbUser, config.DbPassword, config.DbName)

	db, err := sql.Open("postgres", dbDsn)
	if err != nil {
		log.Fatalln("failed to connect db", err)
	}
	defer db.Close()
	log.Infoln("connected to db instance")

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalln("failed to migrate", err)
	}
	log.Infoln("migrated ", n)

	mqDsn := fmt.Sprintf("amqp://%s:%s@%s:5672/", config.MqPassword, config.MqUser, config.MqHost)
	conn, err := amqp.Dial(mqDsn)
	if err != nil {
		log.Fatalln("failed to connect mq", err)
	}
	defer conn.Close()
	log.Infoln("connected to mq instance")

	h, err := handler.Health(dbDsn, mqDsn)
	if err != nil {
		log.Panicln("failed to register status", err)
	}

	userRepository := user.NewRepository(db)
	mQ := user.NewMQ(conn)
	userService := user.NewService(userRepository, mQ)
	userHandler := handler.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Use(handler.LoggerWithLevel(log.InfoLevel))
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

	http.ListenAndServe(":8080", r)
}
