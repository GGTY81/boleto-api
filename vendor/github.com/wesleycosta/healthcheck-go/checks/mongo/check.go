package mongo

import (
	"context"
	"crypto/tls"
	"errors"
	"sync"
	"time"

	"github.com/wesleycosta/healthcheck-go/checks"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var (
	conn              *mongo.Client
	ConnectionTimeout = 3 * time.Second
	mu                sync.RWMutex
)

type Config struct {
	Url        string
	User       string
	Password   string
	Database   string
	AuthSource string
	Timeout    int
	ForceTLS   bool
}

type healthCheck struct {
	Config *Config
}

func new(config *Config) checks.Check {
	return &healthCheck{
		Config: config,
	}
}

func (service *healthCheck) GetName() string {
	return "mongo"
}

func (config *Config) CreateCheck() checks.Check {
	return new(config)
}

func (service *healthCheck) Execute() checks.CheckResult {
	err := service.executeCheck()
	return checks.NewCheckResult(service.GetName(), err)
}

func (service *healthCheck) executeCheck() error {
	if service.Config.Url == "" {
		return errors.New("URL is empty")
	}

	_, err := service.connectMongo()
	if err != nil {
		return err
	}

	return ping()
}

func (service *healthCheck) connectMongo() (*mongo.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	if conn != nil {
		return conn, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	var err error
	conn, err = mongo.Connect(ctx, getClientOptions(service.Config))
	if err != nil {
		return conn, err
	}

	return conn, nil
}

func ping() error {
	var err error
	if conn == nil {
		err = errors.New("Connection is empty")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	err = conn.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	return nil
}

func getClientOptions(config *Config) *options.ClientOptions {
	mongoURL := config.Url
	co := options.Client()
	co.SetRetryWrites(true)
	co.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	co.SetConnectTimeout(time.Duration(config.Timeout) * time.Second)
	co.SetMaxConnIdleTime(time.Duration(config.Timeout) * time.Second)
	co.SetMaxPoolSize(512)

	if config.ForceTLS {
		co.SetTLSConfig(&tls.Config{})
	}

	return co.ApplyURI(mongoURL).SetAuth(mongoCredential(config))
}

func mongoCredential(config *Config) options.Credential {
	user := config.User
	password := config.Password
	var database string
	if config.AuthSource != "" {
		database = config.AuthSource
	} else {
		database = config.Database
	}

	credential := options.Credential{
		Username:   user,
		Password:   password,
		AuthSource: database,
	}

	if config.ForceTLS {
		credential.AuthMechanism = "SCRAM-SHA-1"
	}

	return credential
}
