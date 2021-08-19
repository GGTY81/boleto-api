package api

import (
	"context"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/queue"
)

//InstallRestAPI "instala" e sobe o servico de rest
func InstallRestAPI() {

	l := log.CreateLog()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	useNewRelic(router)

	if config.Get().DevMode && !config.Get().MockMode {
		router.Use(gin.Logger())
	}

	Base(router)
	V1(router)
	V2(router)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:    config.Get().APIPort,
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			interrupt <- syscall.SIGTERM
			l.Error(err, "Got an error when trying serve.ListAndServe()")
			stdlog.Println("err: ", err)
		}
	}()

	<-interrupt
	stdlog.Println("shutdown server")

	// Close DB Connection
	err := db.CloseConnection()
	if err != nil {
		l.Error(err, "error closing Mongodb connection")
	} else {
		l.InfoWithParams("mongodb connection successfully closed", "Information", nil)
	}

	// Close RabbitMQ Connection
	err = queue.CloseConnection()
	if err != nil {
		l.Error(err, "error closing rabbitmq connection")
	} else {
		l.InfoWithParams("rabbitmq connection successfully closed", "Information", nil)
	}

	// Server Shutdown
	err = server.Shutdown(context.Background())
	if err != nil {
		l.Error(err, "shutdown server with error")
	}

	l.InfoWithParams("shutdown completed", "Information", nil)
	stdlog.Println("shutdown completed")
	time.Sleep(10 * time.Second)
}
