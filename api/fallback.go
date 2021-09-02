package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/infrastructure/storage"
)

const persistenceErrorMessage = "Failure during send boleto to fallback. This boleto can't be recovery until manual insert content into database."

type IFallback interface {
	Save(context *gin.Context, registerId, payload string)
}

type Fallback struct{}

// Save resilience application
func (f *Fallback) Save(context *gin.Context, registerId, payload string) {
	lg := loadBankLog(context)

	client, err := storage.GetClient()

	if err != nil {
		lg.ErrorWithContent("failure to get client", "Error", err, payload)
		return
	}

	elapsedTime, err := client.UploadAsJson(
		context,
		registerId,
		payload)

	if err != nil {
		lg.ErrorWithContent(persistenceErrorMessage, "Error", err, payload)
		return
	}

	props := getLogUploadProperties(elapsedTime, registerId, payload)
	lg.InfoWithParams("loaded the payload into Azure Blob Storage with success", "Information", props)
}

func getLogUploadProperties(totalElapsedTimeInMilliseconds int64, registerId, payload string) map[string]interface{} {
	props := make(map[string]interface{})
	props["TotalElapsedTimeInMilliseconds"] = totalElapsedTimeInMilliseconds
	props["RegisterId"] = registerId
	props["Content"] = payload
	return props
}
