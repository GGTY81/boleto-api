package checks

type HealthCheckResult struct {
	Status  string                 `json:"status"`
	Results map[string]CheckResult `json:"results"`
}

func NewHealthCheckResult() HealthCheckResult {
	return HealthCheckResult{
		Results: make(map[string]CheckResult),
		Status:  Healthy,
	}
}

func (healthCheckResult *HealthCheckResult) AddCheckResult(key string, checkResult CheckResult) {
	if checkResult.Status == Unhealthy {
		healthCheckResult.Status = Unhealthy
	}

	healthCheckResult.Results[key] = checkResult
}
