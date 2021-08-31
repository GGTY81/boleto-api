package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/storage"
)

func fallback(context *gin.Context, registerId, payload string) {
	lg := loadBankLog(context)

	client, err := storage.NewAzureBlob(
		config.Get().AzureStorageAccount,
		config.Get().AzureStorageAccessKey,
		config.Get().AzureStorageContainerName,
		config.Get().DevMode,
	)

	if err != nil {
		lg.FallbackErrorWithBasic(persistenceErrorMessage, "Error", err, payload)
		return
	}

	filename := registerId + ".json"
	fullpath := config.Get().AzureStorageUploadPath + "/" + config.Get().AzureStorageFallbackFolder + "/" + filename

	elapsedTime, err := client.Upload(
		context,
		fullpath,
		payload)

	if err != nil {
		lg.FallbackErrorWithBasic(persistenceErrorMessage, "Error", err, payload)
		return
	}

	props := getLogUploadProperties(elapsedTime, registerId)
	lg.InfoWithBasic("loaded the payload into Azure Blob Storage with success", "Information", props)
}

func getLogUploadProperties(totalElapsedTimeInMilliseconds int64, registerId string) map[string]interface{} {
	props := make(map[string]interface{})
	props["TotalElapsedTimeInMilliseconds"] = totalElapsedTimeInMilliseconds
	props["RegisterId"] = registerId
	return props
}
