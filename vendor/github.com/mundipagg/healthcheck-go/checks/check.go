package checks

const (
	Unhealthy string = "Unhealthy"
	Healthy   string = "Healthy"
)

type Check interface {
	Execute() CheckResult
	GetName() string
}
