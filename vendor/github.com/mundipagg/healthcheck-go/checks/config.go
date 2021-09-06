package checks

type Config interface {
	CreateCheck() Check
}
