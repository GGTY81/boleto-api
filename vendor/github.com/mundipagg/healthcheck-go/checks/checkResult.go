package checks

import "fmt"

type CheckResult struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

func NewCheckResult(name string, err error) CheckResult {
	if err == nil {
		return newHealthy(name)
	}

	return newUnhealthy(err)
}

func newHealthy(name string) CheckResult {
	return CheckResult{
		Status:      Healthy,
		Description: fmt.Sprintf("%s is healthy", name),
	}
}

func newUnhealthy(err error) CheckResult {
	return CheckResult{
		Status:      Unhealthy,
		Description: fmt.Sprintf("ERROR: %s", err),
	}
}
