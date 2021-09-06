package rabbit

func newStubRabbitConfig() *Config {
	return &Config{
		ConnectionString: "amqp://guest:guest@localhost:5672/",
	}
}

func (config *Config) withConnectionString(connectionString string) *Config {
	config.ConnectionString = connectionString

	return config
}
