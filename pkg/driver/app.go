package driver

import (
	"net/http"
	"os"
	"worker-service/internal/app_log"
	"worker-service/internal/controller/app_controller"
	"worker-service/pkg"
	"worker-service/pkg/message_system/rabbitmq"

	"github.com/rs/zerolog/log"
)

func Run() {
	pkg.LoadConfig()
	app_log.InitLogger()
	rabbitMq := &rabbitmq.RabbitMQ{
		Url:      os.Getenv("RABBITMQ_URL"),
		Protocol: os.Getenv("RABBITMQ_PROTOCOL"),
		Username: os.Getenv("RABBITMQ_USERNAME"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
	}
	err := rabbitMq.Connect()
	if err != nil {
		log.Fatal().Str("error", "Error when connect to rabbitmq: "+err.Error()).Msg("")
	}
	log.Info().Msg("Connect to rabbitmq successfully")

	// Register health check endpoint
	http.HandleFunc("/health", app_controller.HealthHandler)

	log.Info().Msg("Server is running on port: " + os.Getenv("SERVER_PORT"))
	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), nil)
}
