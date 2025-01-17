package healthcheck

import (
	"errors"
	stdlog "log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/queue"

	HealthCheckLib "github.com/mundipagg/healthcheck-go"
	checks "github.com/mundipagg/healthcheck-go/checks"
)

const (
	Unhealthy string = "Unhealthy"
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
	mongoConfig := db.GetDatabaseConfiguration()
	rabbitConfig := queue.GetQueueConfiguration()

	healthCheck := HealthCheckLib.New()
	healthCheck.AddService(mongoConfig)
	healthCheck.AddService(rabbitConfig)

	return healthCheck
}

func ExecuteOnAPI(c *gin.Context) {
	healtcheck := createHealthCheck()
	result := healtcheck.Execute()
	status := http.StatusOK

	if result.Status == Unhealthy {
		logInstance("ExecuteOnAPI").ErrorBasicWithContent("Healthcheck is Unhealthy", "HealthCheck", result)
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, newHealthCheckResponse(&result))
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
