package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type Producer struct {
	conn *amqp.Connection
}

func NewProducer(conn *amqp.Connection) (Producer, error) {
	producer := Producer{
		conn: conn,
	}
	if err := producer.setup(); err != nil {
		log.Error().Err(err).Msg("failed to setup producer")
		return Producer{}, err
	}

	return producer, nil
}

func (producer *Producer) setup() error {
	channel, err := producer.conn.Channel()
	if err != nil {
		log.Error().Err(err).Msg("failed to create chanel for producer")
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

func (producer *Producer) Push(event, severity string) error {
	channel, err := producer.conn.Channel()
	if err != nil {
		log.Error().Err(err).Msg("failed to create chanel for producer")
		return err
	}
	defer channel.Close()

	log.Debug().Msg("pushing to chanel")

	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		log.Error().Err(err).Msgf("failed to publish event: %s to rabbitMQ", event)
		return err
	}

	return nil
}
