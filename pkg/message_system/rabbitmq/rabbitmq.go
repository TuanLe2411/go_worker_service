package rabbitmq

import (
	"os"
	"strconv"
	"sync"
	"worker-service/pkg/handlers"
	"worker-service/pkg/utils"

	"github.com/rs/zerolog/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	Url      string
	Protocol string
	Username string
	Password string
}

func (r *RabbitMQ) Connect() error {
	connString := r.Protocol + "://" + r.Username + ":" + r.Password + "@" + r.Url
	conn, err := amqp.Dial(connString)
	if err != nil {
		return err
	}
	r.conn = conn
	channel, err := r.conn.Channel()
	if err != nil {
		return err
	}
	r.channel = channel
	return r.init()
}

func (r *RabbitMQ) init() error {
	msgs, err := r.channel.Consume(
		os.Getenv("RABBITMQ_QUEUE_NAME"),
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	numberOfWorkers, _ := strconv.Atoi(os.Getenv("APP_NUMBER_OF_WORKERS"))
	workerPool := &utils.WorkerPool{
		NumberOfWorkers:  numberOfWorkers,
		Mux:              &sync.Mutex{},
		NumberOfWorkings: 0,
	}

	numberOfRetry, _ := strconv.Atoi(os.Getenv("APP_NUMBER_OF_RETRY"))
	go func() {
		for msg := range msgs {
			workerPool.Start(handlers.HandleUserAction, numberOfRetry, string(msg.Body))
			err := msg.Ack(false)
			if err != nil {
				log.Error().Str("error", "Error acknowledging message: "+err.Error()).Msg("")
			} else {
				log.Info().Msg("Acknowledged message")
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
