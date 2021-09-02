package healthcheck

import (
	stdlog "log"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	HealthCheckLib "github.com/wesleycosta/healthcheck-go"

	"github.com/wesleycosta/healthcheck-go/checks/mongo"
	"github.com/wesleycosta/healthcheck-go/checks/rabbit"
)

const (
	Unhealthy string = "Unhealthy"
)

func createHealthCheck() HealthCheckLib.HealthCheck {
	mongoConfig := &mongo.Config{
		Url:        config.Get().MongoURL,
		User:       config.Get().MongoUser,
		Password:   config.Get().MongoPassword,
		Database:   config.Get().MongoDatabase,
		AuthSource: config.Get().MongoAuthSource,
		Timeout:    3,
		ForceTLS:   config.Get().ForceTLS,
	}

	rabbitConfig := &rabbit.Config{
		ConnectionString: config.Get().ConnQueue,
	}

	healthCheck := HealthCheckLib.New()
	healthCheck.AddService(mongoConfig)
	healthCheck.AddService(rabbitConfig)

	return healthCheck
}

func Endpoint(c *gin.Context) {
	healtcheck := createHealthCheck()
	c.JSON(200, healtcheck.Execute())
}

func ExecuteOnStartup() bool {
	logger := log.CreateLog()
	logger.InfoWithBasic("HealthCheckStart", "Starting HealthCheck", nil)

	healtcheck := createHealthCheck()
	result := healtcheck.Execute()

	if result.Status == Unhealthy {
		stdlog.Println("Healthcheck is Unhealthy")
		logger.ErrorBasicWithContent("Healthcheck is Unhealthy", "ExecuteOnStartup", map[string]interface{}{"Error": result, "Operation": "HealthCheckUnhealthy"})
		gracefullyShutdown()

		return false
	}

	logger.InfoWithBasic("Result of execution", "HealthCheck", map[string]interface{}{"Content": result, "Operation": "HealthCheckResult"})
	return true
}

func gracefullyShutdown() {
	logger := log.CreateLog()
	logger.InfoWithBasic("Shutdown", "The application will be terminated", nil)
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}
