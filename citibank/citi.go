package citibank

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/certificate"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/issuer"
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

type bankCiti struct {
	validate  *models.Validator
	log       *log.Log
	transport *http.Transport
}

func New() (bankCiti, error) {
	var err error
	b := bankCiti{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}

	certificates := certificate.TLSCertificate{
		Crt: config.Get().CitibankCertificateSSLName,
		Key: config.Get().CitibankCertificateSSLName,
	}

	onceTransport.Do(func() {
		transportTLS, err = util.BuildTLSTransport(certificates)
	})
	b.transport = transportTLS

	if err != nil || (b.transport == nil && !config.Get().MockMode) {
		return bankCiti{}, fmt.Errorf("fail on load TLSTransport: %v", err)
	}

	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)
	b.validate.Push(citiValidateAgency)
	b.validate.Push(citiValidateAccount)
	b.validate.Push(citiValidateWallet)

	return b, nil
}

//Log retorna a referencia do log
func (b bankCiti) Log() *log.Log {
	return b.log
}

func (b bankCiti) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {

	boleto.Title.BoletoType, boleto.Title.BoletoTypeCode = getBoletoType()

	boleto.Title.OurNumber = calculateOurNumber(boleto)
	r := flow.NewFlow()
	serviceURL := config.Get().URLCiti
	from := getResponseCiti()
	to := getAPIResponseCiti()
	bod := r.From("message://?source=inline", boleto, getRequestCiti(), tmpl.GetFuncMaps())
	bod.To("log://?type=request&url="+serviceURL, b.log)
	var responseCiti string
	var status int
	var err error
	duration := util.Duration(func() {
		responseCiti, status, err = b.sendRequest(bod.GetBody().(string))
	})
	if err != nil {
		return models.BoletoResponse{}, err
	}
	metrics.PushTimingMetric("citibank-register-boleto-online", duration.Seconds())
	bod.To("set://?prop=header", map[string]string{"status": strconv.Itoa(status)})
	bod.To("set://?prop=body", responseCiti)
	bod.To("log://?type=response&url="+serviceURL, b.log)
	ch := bod.Choice()
	ch.When(flow.Header("status").IsEqualTo("200"))
	ch.To("transform://?format=xml", from, to, tmpl.GetFuncMaps())
	ch.Otherwise()
	ch.To("log://?type=response&url="+serviceURL, b.log).To("apierro://")

	switch t := bod.GetBody().(type) {
	case string:
		response := util.ParseJSON(t, new(models.BoletoResponse)).(*models.BoletoResponse)
		issuer := issuer.NewIssuer(response.BarCodeNumber, response.DigitableLine)
		if !((issuer.IsValidBarCode() && issuer.IsValidDigitableLine()) || response.HasErrors()) {
			return models.BoletoResponse{}, models.NewBadGatewayError("BadGateway")
		}
		return *response, nil
	case models.BoletoResponse:
		return t, nil
	}
	return models.BoletoResponse{}, models.NewInternalServerError("MP500", "Internal error")
}

func (b bankCiti) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}
	return b.RegisterBoleto(boleto)
}

func (b bankCiti) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(boleto))
}

func (b bankCiti) sendRequest(body string) (string, int, error) {
	serviceURL := config.Get().URLCiti
	if config.Get().MockMode {
		return util.Post(serviceURL, body, config.Get().TimeoutDefault, map[string]string{"Soapaction": "RegisterBoleto"})
	} else {
		return util.PostTLS(serviceURL, body, config.Get().TimeoutDefault, map[string]string{"Soapaction": "RegisterBoleto"}, b.transport)
	}
}

//GetBankNumber retorna o codigo do banco
func (b bankCiti) GetBankNumber() models.BankNumber {
	return models.Citibank
}

func calculateOurNumber(boleto *models.BoletoRequest) uint {
	ourNumberWithDigit := strconv.Itoa(int(boleto.Title.OurNumber)) + util.OurNumberDv(strconv.Itoa(int(boleto.Title.OurNumber)), util.MOD11)
	value, _ := strconv.Atoi(ourNumberWithDigit)
	return uint(value)
}

func (b bankCiti) GetBankNameIntegration() string {
	return "Citibank"
}

func (b bankCiti) GetErrorsMap() map[string]int {
	return nil
}

func getBoletoType() (bt string, btc string) {
	return "DMI", "03"
}
