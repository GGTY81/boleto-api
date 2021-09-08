package mongo

func newStubMongoConfig() *Config {
	return &Config{
		Url:         "mongodb://localhost:27017",
		User:        "test",
		Password:    "test",
		AuthSource:  "admin",
		Timeout:     3,
		ForceTLS:    false,
		MaxPoolSize: 10,
	}
}

func (config *Config) withUrl(url string) *Config {
	config.Url = url

	return config
}
