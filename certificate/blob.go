package certificate

import (
	"fmt"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/infrastructure/storage"
	"github.com/mundipagg/boleto-api/log"
)

func InstanceStoreCertificatesFromAzureBlob(certificatesName ...string) error {
	l := log.CreateLog()
	l.Operation = "InstanceStoreCertificatesFromAzureBlob"

	azureBlobInst, err := storage.NewAzureBlob(
		config.Get().AzureStorageAccount,
		config.Get().AzureStorageAccessKey,
		config.Get().AzureStorageContainerName,
		config.Get().DevMode,
	)

	if err != nil {
		return err
	}

	for _, certificateName := range certificatesName {
		l.InfoWithBasic(fmt.Sprintf("Start loading [%s] PK from blob", certificateName), "LoadFromAzureBlob", nil)
		skBytes, err := azureBlobInst.Download(
			config.Get().AzureStorageOpenBankSkPath,
			certificateName,
		)

		if err != nil {
			return err
		}

		SetCertificateOnStore(certificateName, skBytes)
		l.InfoWithBasic(fmt.Sprintf("Success loading [%s] PK from blob", certificateName), "LoadFromAzureBlob", nil)

	}

	return nil
}
