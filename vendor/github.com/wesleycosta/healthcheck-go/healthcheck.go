package healthcheck

import (
	"github.com/wesleycosta/healthcheck-go/checks"
	mongo "github.com/wesleycosta/healthcheck-go/checks/mongo"
	rabbit "github.com/wesleycosta/healthcheck-go/checks/rabbit"
)

type HealthCheck interface {
	Execute() checks.HealthCheckResult
	AddService(config checks.Config)
	addCheck(check checks.Check)
}

type healthCheck struct {
	checks map[string]checks.Check
}

func New() HealthCheck {
	return &healthCheck{
		checks: make(map[string]checks.Check),
	}
}

func (healthCheck *healthCheck) Execute() checks.HealthCheckResult {
	healthCheckResult := checks.NewHealthCheckResult()

	for key := range healthCheck.checks {
		healthCheckResult.AddCheckResult(key, healthCheck.checks[key].Execute())
	}

	return healthCheckResult
}

func (healthCheck *healthCheck) AddService(config checks.Config) {
	switch config.(type) {
	case *rabbit.Config:
		rabbitCheck := config.(*rabbit.Config).CreateCheck()
		healthCheck.addCheck(rabbitCheck)

	case *mongo.Config:
		mongoCheck := config.(*mongo.Config).CreateCheck()
		healthCheck.addCheck(mongoCheck)
	}
}

func (healthCheck *healthCheck) addCheck(check checks.Check) {
	healthCheck.checks[check.GetName()] = check
}
