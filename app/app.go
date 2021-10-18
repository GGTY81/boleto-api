package app

import (
	"os"
	"time"

	"github.com/mundipagg/boleto-api/api"
	"github.com/mundipagg/boleto-api/certificate"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/env"
	"github.com/mundipagg/boleto-api/healthcheck"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/mundipagg/boleto-api/usermanagement"
)

//Params this struct contains all execution parameters to run application
type Params struct {
	DevMode    bool
	MockMode   bool
	DisableLog bool
}

//NewParams returns new Empty pointer to ExecutionParameters
func NewParams() *Params {
	return new(Params)
}

//Run starts boleto api Application
func Run(params *Params) {
	env.Config(params.DevMode, params.MockMode, params.DisableLog)

	if config.Get().MockMode {
		go mock.Run("9091")
		time.Sleep(2 * time.Second)
	}

	log.Install()

	start := time.Now()

	healthcheck.ExecuteOnStartup()

	installCertificates()

	usermanagement.LoadUserCredentials()

	props := getLoadDependenciesLogProp(start)
	go log.CreateLog().InfoWithBasic("Load Dependencies with success", "Information", props)

	api.InstallRestAPI()
}

func installCertificates() {
	l := log.CreateLog()
	l.Operation = "InstallCertificates"

	if config.Get().MockMode {
		certificate.LoadMockCertificates()
		return
	}
	err := certificate.InstanceStoreCertificatesFromAzureVault(config.Get().VaultName, config.Get().CertificateICPName, config.Get().CertificateSSLName)
	if err == nil {
		l.InfoWithBasic("Success in load certificates from azureVault", "LoadFromAzureVault", nil)
	} else {
		l.ErrorWithBasic("Error in load certificates from azureVault", "LoadFromAzureVault", err)
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}

	err = certificate.InstanceStoreCertificatesFromAzureBlob(config.Get().AzureStorageOpenBankSkName, config.Get().AzureStorageJPMorganPkName, config.Get().AzureStorageJPMorganCrtName, config.Get().AzureStorageJPMorganSignCrtName)
	if err != nil {
		l.ErrorWithBasic("Error loading open bank secret key from blob", "LoadFromAzureBlob", err)
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}

}

func getLoadDependenciesLogProp(start time.Time) map[string]interface{} {
	props := make(map[string]interface{})
	props["totalElapsedTimeInMilliseconds"] = time.Since(start).Milliseconds()
	props["Operation"] = "RunApp"
	return props
}
