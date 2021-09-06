package rabbit

import (
	"errors"
	"time"

	checks "github.com/mundipagg/healthcheck-go/checks"
	"github.com/streadway/amqp"
)

func new(config *Config) checks.Check {
	return &healthCheck{
		Config: config,
	}
}

type Config struct {
	ConnectionString string
}

type healthCheck struct {
	Config *Config
}

func (config *Config) CreateCheck() checks.Check {
	return new(config)
}

func (service *healthCheck) GetName() string {
	return "rabbit"
}

func (service *healthCheck) Execute() checks.CheckResult {
	err := service.executeCheck()
	return checks.NewCheckResult(service.GetName(), err)
}

func (service *healthCheck) executeCheck() error {
	connection, err := openConnection(service.Config)
	if err != nil {
		return err
	}

	channel, err := openChannel(connection)
	if err != nil {
		return err
	}

	err = closeChannel(channel)

	if err != nil {
		closeConnection(connection)
		return err
	}

	return closeConnection(connection)
}

func openChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	var channel *amqp.Channel
	var err error

	if conn == nil {
		err = errors.New("failed to open a channel. The connection is closed")
		return nil, err
	}

	if channel, err = conn.Channel(); err != nil {
		return nil, err
	}

	return channel, err
}

func closeChannel(channel *amqp.Channel) error {
	if channel != nil {
		return channel.Close()
	}

	return nil
}

func closeConnection(connection *amqp.Connection) error {
	if connection != nil {
		return connection.Close()
	}

	return nil
}

func openConnection(config *Config) (*amqp.Connection, error) {
	var hb int
	var err error

	conn, err := amqp.DialConfig(config.ConnectionString, amqp.Config{
		Heartbeat: time.Duration(hb) * time.Second,
	})

	return conn, err
}
