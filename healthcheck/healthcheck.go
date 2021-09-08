package healthcheck

import (
	"errors"
	stdlog "log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"

	HealthCheckLib "github.com/mundipagg/healthcheck-go"
	checks "github.com/mundipagg/healthcheck-go/checks"
	"github.com/mundipagg/healthcheck-go/checks/mongo"
	"github.com/mundipagg/healthcheck-go/checks/rabbit"
)

const (
	Unhealthy    string = "Unhealthy"
	MongoTimeout int    = 3
)

type HealthCheckResponse struct {
	Status string `json:"status"`
}

func newHealthCheckResponse(healthCheckResult *checks.HealthCheckResult) HealthCheckResponse {
	return HealthCheckResponse{
		Status: healthCheckResult.Status,
	}
}

func createHealthCheck() HealthCheckLib.HealthCheck {
	mongoConfig := &mongo.Config{
		Url:        config.Get().MongoURL,
		User:       config.Get().MongoUser,
		Password:   config.Get().MongoPassword,
		Database:   config.Get().MongoDatabase,
		AuthSource: config.Get().MongoAuthSource,
		Timeout:    MongoTimeout,
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

func ExecuteOnAPI(c *gin.Context) {
	healtcheck := createHealthCheck()
	result := healtcheck.Execute()

	if result.Status == Unhealthy {
		logInstance("ExecuteOnAPI").ErrorBasicWithContent("Healthcheck is Unhealthy", "HealthCheck", result)
	}

	c.JSON(200, newHealthCheckResponse(&result))
}

func ExecuteOnStartup() bool {
	var logger = logInstance("ExecuteOnStartup")
	logger.InfoWithBasic("Starting HealthCheck", "HealthCheck", nil)

	healtcheck := createHealthCheck()
	result := healtcheck.Execute()

	if result.Status == Unhealthy {
		stdlog.Println("Healthcheck is Unhealthy", result)
		logger.ErrorBasicWithContent("Application Unhealthy. The application will be terminated", "HealthCheck", result)
		shutdown()

		return false
	}

	logger.InfoWithBasic("Result of check dependecies execution", "HealthCheck", map[string]interface{}{"Content": result})
	return true
}

func shutdown() {
	time.Sleep(5 * time.Second)
	panic(errors.New("healthcheck is unhealthy"))
}

func logInstance(operation string) *log.Log {
	logger := log.CreateLog()
	logger.Operation = operation
	return logger
}
