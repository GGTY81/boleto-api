package sqlServer

import (
	"database/sql"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	checks "github.com/mundipagg/healthcheck-go/checks"
)

var (
	connection  *sql.DB
	errDatabase error
	once        sync.Once
)

type Config struct {
	ConnectionString string
	Query            string
}

func new(config *Config) checks.Check {
	return &healthCheck{
		Config: config,
	}
}

type healthCheck struct {
	Config *Config
}

func (config *Config) CreateCheck() checks.Check {
	return new(config)
}

func (service *healthCheck) GetName() string {
	return "sqlServer"
}

func (service *healthCheck) Execute() checks.CheckResult {
	err := service.executeCheck()
	return checks.NewCheckResult(service.GetName(), err)
}

func (service *healthCheck) executeCheck() error {
	db, err := getInstanceDB(service.Config)

	if err != nil {
		return err
	}

	rows, err := db.Query(service.Config.Query)

	if err != nil {
		return err
	}

	rows.Next()
	defer rows.Close()

	return nil
}

func getInstanceDB(config *Config) (*sql.DB, error) {
	once.Do(func() {
		connection, errDatabase = sql.Open("mssql", config.ConnectionString)
	})

	if errDatabase != nil {
		return nil, errDatabase
	}

	if ping := connection.Ping(); ping != nil {
		err := ping
		return nil, err
	}

	return connection, nil
}
