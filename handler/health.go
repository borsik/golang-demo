package handler

import (
	"github.com/hellofresh/health-go/v5"
	healthPostgres "github.com/hellofresh/health-go/v5/checks/postgres"
	healthRabbit "github.com/hellofresh/health-go/v5/checks/rabbitmq"
	"time"
)

func Health(dbDsn string, mqDsn string) (*health.Health, error) {
	h, _ := health.New(health.WithComponent(health.Component{
		Name:    "golang-demo",
		Version: "v1.0",
	}))
	err := h.Register(health.Config{
		Name:      "postgres",
		Timeout:   time.Second * 2,
		SkipOnErr: false,
		Check: healthPostgres.New(healthPostgres.Config{
			DSN: dbDsn,
		}),
	})
	err = h.Register(health.Config{
		Name:      "rabbitmq",
		Timeout:   time.Second * 2,
		SkipOnErr: false,
		Check: healthRabbit.New(healthRabbit.Config{
			DSN: mqDsn,
		}),
	})

	if err != nil {
		return h, err
	}

	return h, nil
}
