package jpmorgan

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/certificate"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/infrastructure/token"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
	"github.com/mundipagg/boleto-api/validations"
)

var (
	onceTransport = &sync.Once{}
	transportTLS  *http.Transport
)

type bankJPMorgan struct {
	validate  *models.Validator
	log       *log.Log
	transport *http.Transport
	jwtSigner token.JwtGenerator
}

func New() (bankJPMorgan, error) {
	var err error
	b := bankJPMorgan{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}

	certificates := certificate.TLSCertificate{
		Crt: config.Get().AzureStorageJPMorganCrtName,
		Key: config.Get().AzureStorageJPMorganPkName,
	}

	onceTransport.Do(func() {
		transportTLS, err = util.BuildTLSTransport(certificates)
	})
	b.transport = transportTLS

	if err != nil || (b.transport == nil && !config.Get().MockMode) {
		return bankJPMorgan{}, fmt.Errorf("fail on load TLSTransport: %v", err)
	}

	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)

	b.jwtSigner = token.GetJwtGenerator("RS256")

	return b, nil
}

func (b bankJPMorgan) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}

	return b.RegisterBoleto(boleto)
}

func (b bankJPMorgan) ValidateBoleto(request *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(request))
}

func (b bankJPMorgan) GetBankNumber() models.BankNumber {
	return models.JPMorgan
}

func (b bankJPMorgan) GetBankNameIntegration() string {
	return "JPMorgan"
}

func (b bankJPMorgan) GetErrorsMap() map[string]int {
	var erros = map[string]int{
		"BOL-5":                          http.StatusBadRequest,
		"BOL-7":                          http.StatusBadRequest,
		"BOL-144":                        http.StatusBadRequest,
		"GCA-001":                        http.StatusBadGateway,
		"GCA-111":                        http.StatusBadRequest,
		"GCA-003":                        http.StatusBadGateway,
		"BOL-2":                          http.StatusBadRequest,
		"BOL-3":                          http.StatusBadRequest,
		"BOL-4":                          http.StatusBadRequest,
		"BOL-1":                          http.StatusBadGateway,
		"Signature Verification Failure": http.StatusInternalServerError,
		"Authentication Failure":         http.StatusInternalServerError,
		"Internal Error - Contact Service Provider": http.StatusBadGateway,
	}
	return erros
}

func (b bankJPMorgan) Log() *log.Log {
	return b.log
}

func (b bankJPMorgan) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	var response string
	var respHeader map[string]interface{}
	var status int
	var err error

	JPMorganURL := config.Get().URLJPMorgan
	boleto.Title.BoletoType, boleto.Title.BoletoTypeCode = getBoletoType(boleto)

	body := flow.NewFlow().From("message://?source=inline", boleto, templateRequest, tmpl.GetFuncMaps()).GetBody().(string)
	head := hearders()

	bodyEncripted, encryptedErr := b.encriptedBody(body)
	if encryptedErr != nil {
		return *encryptedErr, nil
	}

	b.log.Request(body, JPMorganURL, getLogRequestProperties(head, bodyEncripted))

	duration := util.Duration(func() {
		if config.Get().MockMode {
			response, respHeader, status, err = util.PostWithHeader(JPMorganURL, body, config.Get().TimeoutDefault, head)
		} else {
			response, respHeader, status, err = util.PostTLSWithHeader(JPMorganURL, bodyEncripted, config.Get().TimeoutDefault, head, b.transport)
		}

	})
	metrics.PushTimingMetric("jpmorgan-register-boleto-time", duration.Seconds())

	b.log.Response(response, JPMorganURL, nil)

	return mapJPMorganResponse(boleto, getContentType(respHeader), response, status, err), nil
}

func (b bankJPMorgan) encriptedBody(body string) (string, *models.BoletoResponse) {
	errorResponse := models.GetBoletoResponseError("MP500", "Encript error")
	sk, err := certificate.GetCertificateFromStore(config.Get().AzureStorageJPMorganSignCrtName)
	if err != nil {
		return "", &errorResponse
	}
	bodyEncripted, err := b.jwtSigner.Sign(util.SanitizeBody(body), sk.([]byte))
	if err != nil {
		return "", &errorResponse
	}
	return bodyEncripted, nil
}

func getContentType(respHeader map[string]interface{}) string {
	return respHeader["Content-Type"].(string)
}

func getLogRequestProperties(header map[string]string, encryptedBody string) map[string]string {
	m := make(map[string]string)
	m["encryptedBody"] = encryptedBody
	for k, v := range header {
		m[k] = v
	}
	return m
}

func hearders() map[string]string {
	return map[string]string{"Content-Type": "text/xml"}
}

func getBoletoType(boleto *models.BoletoRequest) (bt string, btc string) {
	return "DM", "02"
}

func mapJPMorganResponse(request *models.BoletoRequest, contentType string, response string, status int, httpErr error) models.BoletoResponse {
	f := flow.NewFlow().To("set://?prop=body", response)
	switch status {
	case 200:
		f.To("transform://?format=json", templateResponse, templateAPI, tmpl.GetFuncMaps())
		f.To("unmarshall://?format=json", new(models.BoletoResponse))
	case 0, 504:
		return models.GetBoletoResponseError("MPTimeout", timeoutMessage(httpErr))
	case 401, 500:
		if contentType == "application/xml" {
			f.To("set://?prop=body", parseXmlErrorToJson(response))
			f.To("transform://?format=json", templateErrorXmltoJson, templateAPI, tmpl.GetFuncMaps())
		} else if contentType == "application/json" {
			f.To("transform://?format=json", templateErrorJson, templateAPI, tmpl.GetFuncMaps())
		}

		f.To("unmarshall://?format=json", new(models.BoletoResponse))
	case 400, 403, 404:
		dataError := util.ParseJSON(response, new(arrayDataError)).(*arrayDataError)
		f.To("set://?prop=body", strings.Replace(util.Stringify(dataError.Error[0]), "\\\"", "", -1))
		f.To("transform://?format=json", templateErrorJson, templateAPI, tmpl.GetFuncMaps())
		f.To("unmarshall://?format=json", new(models.BoletoResponse))
	case 409, 422:
		f.To("transform://?format=json", templateErrorBoletoJson, templateAPI, tmpl.GetFuncMaps())
		f.To("unmarshall://?format=json", new(models.BoletoResponse))
	}

	switch t := f.GetBody().(type) {
	case *models.BoletoResponse:
		if hasOurNumberFail(t) {
			return models.GetBoletoResponseError("MPOurNumberFail", "our number was not returned by the bank")
		} else {
			return *t
		}
	}
	return models.GetBoletoResponseError("MP500", "Internal Error")
}

func parseXmlErrorToJson(xmlBody string) string {
	data := &ServiceMessageError{}
	err := xml.Unmarshal([]byte(xmlBody), data)
	if err != nil {
		panic(err)
	}
	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(json)
}

func timeoutMessage(err error) string {
	var msg string
	if err != nil {
		msg = fmt.Sprintf("%v", err)
	} else {
		msg = "GatewayTimeout"
	}
	return msg
}

func hasOurNumberFail(response *models.BoletoResponse) bool {
	return !response.HasErrors() && (response.OurNumber == "" || response.OurNumber == "000000000000")
}
