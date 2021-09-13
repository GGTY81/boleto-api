package sqlServer

func newStubSqlServerConfig() *Config {
	return &Config{
		ConnectionString: "server=localhost;port=1434;user id=sa;password=sa;database=master;connection timeout=130",
		Query:            "SELECT TOP 1 TABLE_NAME from INFORMATION_SCHEMA.TABLES",
	}
}

func (config *Config) withQuery(query string) *Config {
	config.Query = query

	return config
}
