package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"golang-demo/handler"
	"golang-demo/user"
	"net/http"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	var logger = log.New()
	logFormatter := new(log.TextFormatter)
	logFormatter.FullTimestamp = true
	logger.SetFormatter(logFormatter)
	logLevel, _ := log.ParseLevel(config.LogLevel)
	logger.SetLevel(logLevel)

	dbDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		config.DbHost, config.DbUser, config.DbPassword, config.DbName)

	db, err := sql.Open("postgres", dbDsn)
	if err != nil {
		logger.Fatalln("failed to connect db", err)
	}
	defer db.Close()
	logger.Infoln("connected to db instance")

	mqDsn := fmt.Sprintf("amqp://%s:%s@%s:5672/", config.MqPassword, config.MqUser, config.MqHost)
	conn, err := amqp.Dial(mqDsn)
	if err != nil {
		logger.Fatalln("failed to connect mq", err)
	}
	defer conn.Close()
	logger.Infoln("connected to mq instance")

	h, err := handler.Health(dbDsn, mqDsn)
	if err != nil {
		logger.Panicln("failed to register status", err)
	}

	userRepository := user.NewRepository(db)
	mQ := user.NewMQ(conn, logger)
	userService := user.NewService(userRepository, mQ, logger)
	userHandler := handler.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Use(handler.LoggerWithLevel(logger, log.InfoLevel))
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

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Fatalln("failed to start", err)
	}
}
