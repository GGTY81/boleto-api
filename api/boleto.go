package api

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/mundipagg/boleto-api/queue"
	"github.com/mundipagg/boleto-api/storage"

	"github.com/gin-gonic/gin"

	"strings"

	"fmt"
	"io/ioutil"

	"github.com/mundipagg/boleto-api/boleto"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
)

//Regista um boleto em um determinado banco
func registerBoleto(c *gin.Context) {

	if _, hasErr := c.Get("error"); hasErr {
		return
	}

	lg := loadBankLog(c)
	bol := getBoletoFromContext(c)
	bank := getBankFromContext(c)

	resp, err := bank.ProcessBoleto(&bol)

	if qualifiedForNewErrorHandling(c, resp) {
		c.Set(responseKey, resp)
		return
	}

	if checkError(c, err, lg) {
		return
	}

	st := getResponseStatusCode(resp)

	if st == http.StatusOK {

		boView := models.NewBoletoView(bol, resp, bank.GetBankNameIntegration())
		resp.ID = boView.ID.Hex()
		resp.Links = boView.Links

		errMongo := db.SaveBoleto(boView)

		if errMongo != nil {
			lg.Warn(errMongo.Error(), fmt.Sprintf("Error saving to mongo - %s", errMongo.Error()))

			b := boView.ToMinifyJSON()
			p := queue.NewPublisher(b)

			if !queue.WriteMessage(p) {
				err := uploadPayloadBlob(
					c,
					boView.ID.Hex(),
					b)

				if err != nil {
					lg.Error(b, persistenceErrorMessage)
				}
			}
		}
	}

	c.JSON(st, resp)
	c.Set("boletoResponse", resp)
}

func uploadPayloadBlob(context *gin.Context, registerId, payload string) (err error) {
	clientBlob, err := getClientBlob()

	if err != nil {
		return
	}

	fileName := registerId + ".json"

	err = clientBlob.Upload(
		context,
		config.Get().AzureStorageUploadPath,
		fileName,
		payload)

	return
}

func getBoleto(c *gin.Context) {
	start := time.Now()
	var boletoHtml string

	c.Status(200)
	log := log.CreateLog()
	log.Operation = "GetBoleto"
	log.IPAddress = c.ClientIP()

	var result = models.NewGetBoletoResult(c)

	defer logResult(result, log, start)

	if !result.HasValidKeys() {
		result.SetErrorResponse(c, models.NewErrorResponse("MP404", "Not Found"), http.StatusNotFound)
		result.LogSeverity = "Warning"
		return
	}

	var err error
	var boView models.BoletoView

	boView, result.DatabaseElapsedTimeInMilliseconds, err = db.GetBoletoByID(result.Id, result.PrivateKey)

	if err != nil && (err.Error() == db.NotFoundDoc || err.Error() == db.InvalidPK) {
		result.SetErrorResponse(c, models.NewErrorResponse("MP404", "Not Found"), http.StatusNotFound)
		result.LogSeverity = "Warning"
		return
	} else if err != nil {
		result.SetErrorResponse(c, models.NewErrorResponse("MP500", err.Error()), http.StatusInternalServerError)
		result.LogSeverity = "Error"
		return
	}

	result.BoletoSource = "mongo"
	boletoHtml = boleto.MinifyHTML(boView)

	if result.Format == "html" {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteString(boletoHtml)
	} else {
		c.Header("Content-Type", "application/pdf")
		if boletoPdf, err := toPdf(boletoHtml); err == nil {
			c.Writer.Write(boletoPdf)
		} else {
			c.Header("Content-Type", "application/json")
			result.SetErrorResponse(c, models.NewErrorResponse("MP500", err.Error()), http.StatusInternalServerError)
			result.LogSeverity = "Error"
			return
		}
	}

	result.LogSeverity = "Information"
}

func getResponseStatusCode(response models.BoletoResponse) int {
	if len(response.Errors) == 0 {
		return http.StatusOK
	}

	if response.StatusCode < 1 {
		return http.StatusBadRequest
	}

	return response.StatusCode
}

func logResult(result *models.GetBoletoResult, log *log.Log, start time.Time) {
	result.TotalElapsedTimeInMilliseconds = time.Since(start).Milliseconds()
	log.GetBoleto(result, result.LogSeverity)
}

func toPdf(page string) ([]byte, error) {
	url := config.Get().PdfAPIURL
	payload := strings.NewReader(page)
	if req, err := http.NewRequest("POST", url, payload); err != nil {
		return nil, err
	} else if res, err := http.DefaultClient.Do(req); err != nil {
		return nil, err
	} else {
		defer res.Body.Close()
		return ioutil.ReadAll(res.Body)
	}
}

func getBoletoByID(c *gin.Context) {
	id := c.Param("id")
	pk := c.Param("pk")
	log := log.CreateLog()
	log.Operation = "GetBoletoV1"

	boleto, _, err := db.GetBoletoByID(id, pk)
	if err != nil {
		checkError(c, models.NewHTTPNotFound("MP404", "Boleto não encontrado"), nil)
		return
	}
	c.JSON(http.StatusOK, boleto)
}

func confirmation(c *gin.Context) {
	if dump, err := httputil.DumpRequest(c.Request, true); err == nil {
		l := log.CreateLog()
		l.BankName = "BradescoShopFacil"
		l.Operation = "BoletoConfirmation"
		l.Request(string(dump), c.Request.URL.String(), nil)
	}
	c.String(200, "OK")
}

func getClientBlob() (*storage.AzureBlob, error) {
	return storage.NewAzureBlob(
		config.Get().AzureStorageAccount,
		config.Get().AzureStorageAccessKey,
		config.Get().AzureStorageContainerName,
		config.Get().DevMode,
	)
}
