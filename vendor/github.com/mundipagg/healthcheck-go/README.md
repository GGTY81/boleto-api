
[![GoDoc](https://godoc.org/github.com/mundipagg/goseq?status.svg)](https://godoc.org/github.com/mundipagg/goseq)
# Golang health check

Bliblioteca de health check em Golang.

### Lista de serviços
1. MongoDB;
2. RabbitMQ.
2. SQL Server.

## Instalação

### Usando *go get*

    $ go get github.com/mundipagg/healthcheck-go

### Usando govendor
	$ govendor add github.com/mundipagg/healthcheck-go
	$ govendor add github.com/mundipagg/healthcheck-go/checks
	$ govendor add github.com/mundipagg/healthcheck-go/checks/mongo
	$ govendor add github.com/mundipagg/healthcheck-go/checks/rabbit       

## Exemplo

### HealthCheck
```go
package healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/config"

	HealthCheckLib "github.com/mundipagg/healthcheck-go"
	"github.com/mundipagg/healthcheck-go/checks/mongo"
	"github.com/mundipagg/healthcheck-go/checks/rabbit"
)

func createHealthCheck() HealthCheckLib.HealthCheck {
	mongoConfig := &mongo.Config{
		Url:        config.Get().MongoURL,
		User:       config.Get().MongoUser,
		Password:   config.Get().MongoPassword,
		Database:   config.Get().MongoDatabase,
		AuthSource: config.Get().MongoAuthSource,
		Timeout:    3,
		ForceTLS:   config.Get().ForceTLS,
		MaxPoolSize:   100,
	}

	rabbitConfig := &rabbit.Config{
		ConnectionString: config.Get().ConnQueue,
	}

	healthCheck := HealthCheckLib.New()
	healthCheck.AddService(mongoConfig)
	healthCheck.AddService(rabbitConfig)

	return healthCheck
}

func Endpoint(c *gin.Context) {
	healthcheck := createHealthCheck()
	c.JSON(200, healthcheck.Execute())
}
```

### API

```go
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/healthcheck"
)

func Base(router *gin.Engine) {
	router.GET("/healthcheck", healthcheck.Endpoint)
}
```

### Response do endpoint
```json
{
    "status": "Healthy",
    "results": {
        "mongo": {
            "status": "Healthy",
            "description": "mongo is healthy"
        },
        "rabbit": {
            "status": "Healthy",
            "description": "rabbit is healthy"
        }
    }
}
```
