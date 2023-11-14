package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"golang-demo/api"
	"golang-demo/config"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln("failed to read config", err)
	}

	dbDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName)

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

	mqDsn := fmt.Sprintf("amqp://%s:%s@%s:5672/", cfg.MqPassword, cfg.MqUser, cfg.MqHost)
	conn, err := amqp.Dial(mqDsn)
	if err != nil {
		log.Fatalln("failed to connect mq", err)
	}
	defer conn.Close()
	log.Infoln("connected to mq instance")

	h, err := api.Health(dbDsn, mqDsn)
	if err != nil {
		log.Panicln("failed to register status", err)
	}

	r := api.NewRouter(db, conn, h)
	http.ListenAndServe(":8080", r)
}
