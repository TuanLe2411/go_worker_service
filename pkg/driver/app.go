package driver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"worker-service/internal/controller/app_controller"
	"worker-service/pkg"
	"worker-service/pkg/message_system/rabbitmq"
)

func Run() {
	pkg.LoadConfig()
	rabbitMq := &rabbitmq.RabbitMQ{
		Url:      os.Getenv("RABBITMQ_URL"),
		Protocol: os.Getenv("RABBITMQ_PROTOCOL"),
		Username: os.Getenv("RABBITMQ_USERNAME"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
	}
	err := rabbitMq.Connect()
	if err != nil {
		log.Println("Error when connect to rabbitmq: " + err.Error())
		panic(err)
	}
	log.Println("Connect to rabbitmq successfully")

	// Register health check endpoint
	http.HandleFunc("/health", app_controller.HealthHandler)

	fmt.Println("Server is running on port: " + os.Getenv("SERVER_PORT"))
	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), nil)
}
