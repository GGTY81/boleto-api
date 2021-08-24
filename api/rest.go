package api

import (
	"context"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/queue"
)

//InstallRestAPI "instala" e sobe o servico de rest
func InstallRestAPI() {

	l := log.CreateLog()
	l.Operation = "InstallAPI"

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
			l.ErrorWithBasic("got an error when trying serve.ListAndServe()", "Error", err)
			stdlog.Println("err: ", err)
		}
	}()

	<-interrupt
	l.InfoWithBasic("start shutdown server...", "Information", nil)
	stdlog.Println("start shutdown server...")

	// Server Shutdown
	err := server.Shutdown(context.Background())
	if err != nil {
		l.ErrorWithBasic("shutdown server with error", "Error", err)
	}

	// Close DB Connection
	err = db.CloseConnection()
	if err != nil {
		l.ErrorWithBasic("error closing Mongodb connection", "Error", err)
	} else {
		l.InfoWithBasic("mongodb connection successfully closed", "Information", nil)
	}

	// Close RabbitMQ Connection
	err = queue.CloseConnection()
	if err != nil {
		l.ErrorWithBasic("error closing rabbitmq connection", "Error", err)
	} else {
		l.InfoWithBasic("rabbitmq connection successfully closed", "Information", nil)
	}

	l.InfoWithBasic("shutdown completed", "Information", nil)
	stdlog.Println("shutdown completed")
	// time.Sleep(10 * time.Second)
}
