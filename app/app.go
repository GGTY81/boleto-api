package app

import (
	"fmt"
	"os"
	"time"

	"github.com/mundipagg/boleto-api/api"
	"github.com/mundipagg/boleto-api/certificate"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/env"
	"github.com/mundipagg/boleto-api/healthcheck"
	"github.com/mundipagg/boleto-api/infrastructure/storage"
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

	if config.Get().MockMode == false && config.Get().EnableFileServerCertificate == false {
		err := certificate.InstanceStoreCertificatesFromAzureVault(config.Get().VaultName, config.Get().CertificateICPName, config.Get().CertificateSSLName)
		if err == nil {
			l.InfoWithBasic("Success in load certificates from azureVault", "LoadFromAzureVault", nil)
		} else {
			l.ErrorWithBasic("Error in load certificates from azureVault", "LoadFromAzureVault", err)
		}
	}

	if config.Get().MockMode == false && config.Get().EnableFileServerCertificate == true {
		err := certificate.InstanceStoreCertificatesFromFileServer(config.Get().CertificateICPName, config.Get().CertificateSSLName)
		if err == nil {
			l.InfoWithBasic("Success in load certificates from fileserver", "LoadFromFileServer", nil)
		} else {
			l.ErrorWithBasic("Error in load certificates from fileServer", "LoadFromFileServer", err)
		}
	}

	sk, err := openBankSkFromBlob()
	if err != nil {
		l.ErrorWithBasic("Error loading open bank secret key from blob", "LoadFromAzureBlob", err)
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}

	l.InfoWithBasic(fmt.Sprintf("Success loading [%s] PK from blob", config.Get().AzureStorageOpenBankSkName), "LoadFromAzureBlob", nil)
	certificate.SetCertificateOnStore(config.Get().AzureStorageOpenBankSkName, sk)
}

func openBankSkFromBlob() ([]byte, error) {
	azureBlobInst, err := storage.NewAzureBlob(
		config.Get().AzureStorageAccount,
		config.Get().AzureStorageAccessKey,
		config.Get().AzureStorageContainerName,
		config.Get().DevMode,
	)
	if err != nil {
		return []byte(""), err
	}

	skBytes, err := azureBlobInst.Download(
		config.Get().AzureStorageOpenBankSkPath,
		config.Get().AzureStorageOpenBankSkName,
	)
	if err != nil {
		return []byte(""), err
	}

	return skBytes, nil
}

func getLoadDependenciesLogProp(start time.Time) map[string]interface{} {
	props := make(map[string]interface{})
	props["totalElapsedTimeInMilliseconds"] = time.Since(start).Milliseconds()
	props["Operation"] = "RunApp"
	return props
}
