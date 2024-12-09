package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/dubass83/go-micro-broker/cmd/api"
	"github.com/dubass83/go-micro-broker/event"
	"github.com/dubass83/go-micro-broker/util"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot load configuration")
	}
	if conf.Enviroment == "devel" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		// log.Debug().Msgf("config values: %+v", conf)
	}

	// try to connect to rabbitmq
	rcon, err := rabbitConn(5, conf)
	if err != nil {
		log.Fatal().Msg("failed to establish the connection to RabbitMQ")
	}
	defer rcon.Close()

	log.Info().Msg("connected to RabbitMQ and ready to produce events...")

	producer, err := event.NewProducer(rcon)
	if err != nil {
		log.Fatal().Msg("failed to create a producer for RabbitMQ conection")
	}

	s := api.CreateNewServer(conf, producer)
	s.ConfigureCORS()
	s.AddMiddleware()
	s.MountHandlers()
	log.Info().
		Msgf("start listening on the port %s\n", conf.HTTPAddressString)
	err = http.ListenAndServe(conf.HTTPAddressString, s.Router)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("canot start http service")
	}
}

func rabbitConn(backOffLimit int64, conf util.Config) (*amqp.Connection, error) {
	var counts int64
	backOff := 1 * time.Second
	var connection *amqp.Connection

	// wait until rabbit is ready
	for {
		// check if we reach the backOffLimit the stop trying and return an error
		if counts >= backOffLimit {
			msg := fmt.Sprintf("failed connect to rabbit with %d tries", backOffLimit)
			log.Error().Msg(msg)
			return nil, fmt.Errorf(msg)
		}
		// try connect to rabbitmq
		c, err := amqp.Dial(conf.RabbitURL)
		if err != nil {
			log.Error().Err(err).Msg("rabbitmq not yet ready...")
			if counts++; counts < backOffLimit {
				time.Sleep(backOff)
			}
			backOff = time.Duration(math.Pow(2, float64(counts))) * time.Second
			continue
		}
		log.Debug().Msg("connected to RabbitMQ")
		connection = c
		break
	}
	return connection, nil
}
