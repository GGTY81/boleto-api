package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/boleto"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/queue"
)

var fallback = new(Fallback)

//registerBoleto Realiza o registro online do Boleto
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
			lg.Warn(errMongo.Error(), "Error saving to mongo")

			b := boView.ToMinifyJSON()
			p := queue.NewPublisher(b)

			if !queue.WriteMessage(p) {
				fallback.Save(c, boView.ID.Hex(), b)
			}
		}
	}

	c.JSON(st, resp)
	c.Set("boletoResponse", resp)
}

//getBoleto Recupera um boleto devidamente registrado
func getBoleto(c *gin.Context) {
	var boletoHtml string

	var result = models.NewGetBoletoResult(c)

	if !result.HasValidParameters() {
		setupGetBoletoResultFailResponse(c, result, "Warning", "Not Found")
		return
	}

	var err error
	var boView models.BoletoView

	boView, result.DatabaseElapsedTimeInMilliseconds, err = db.GetBoletoByID(result.Id, result.PrivateKey)

	if err != nil && (err.Error() == db.NotFoundDoc || err.Error() == db.InvalidPK) {
		setupGetBoletoResultFailResponse(c, result, "Warning", "Not Found")
		return
	} else if err != nil {
		setupGetBoletoResultFailResponse(c, result, "Error", err.Error())
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
			setupGetBoletoResultFailResponse(c, result, "Error", err.Error())
			return
		}
	}

	setupGetBoletoSuccessResponse(c, result)
}

func setupGetBoletoResultFailResponse(c *gin.Context, result *models.GetBoletoResult, severity, errorMessage string) {
	result.LogSeverity = severity

	switch severity {
	case "Warning":
		result.SetErrorResponse(c, models.NewErrorResponse("MP404", errorMessage), http.StatusNotFound)
	default:
		result.SetErrorResponse(c, models.NewErrorResponse("MP500", errorMessage), http.StatusInternalServerError)
	}
	c.Set(resultGetBoletoKey, result)
}

func setupGetBoletoSuccessResponse(c *gin.Context, result *models.GetBoletoResult) {
	c.Status(http.StatusOK)
	result.LogSeverity = "Information"
	c.Set(resultGetBoletoKey, result)
}

func getResponseStatusCode(response models.BoletoResponse) int {
	if len(response.Errors) == 0 {
		return http.StatusOK
	}

	if response.StatusCode == 0 {
		return http.StatusBadRequest
	}

	return response.StatusCode
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
