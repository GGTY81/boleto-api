package santander

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"net/http"

	. "github.com/PMoneda/flow"
	"github.com/mundipagg/boleto-api/certificate"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/metrics"
	"github.com/mundipagg/boleto-api/models"
	"github.com/mundipagg/boleto-api/tmpl"
	"github.com/mundipagg/boleto-api/util"
	"github.com/mundipagg/boleto-api/validations"
)

var (
	onceMap       = &sync.Once{}
	onceTransport = &sync.Once{}
	transportTLS  *http.Transport
	m             map[string]string
)

type bankSantander struct {
	validate  *models.Validator
	log       *log.Log
	transport *http.Transport
}

//New Create a new Santander Integration Instance
func New() (bankSantander, error) {
	var err error
	b := bankSantander{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}

	certificates := certificate.TLSCertificate{
		Crt: config.Get().SantanderCertificateSSLName,
		Key: config.Get().SantanderCertificateSSLName,
	}

	onceTransport.Do(func() {
		transportTLS, err = util.BuildTLSTransport(certificates)
	})
	b.transport = transportTLS

	if err != nil || (b.transport == nil && !config.Get().MockMode) {
		return bankSantander{}, fmt.Errorf("fail on load TLSTransport: %v", err)
	}

	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)
	b.validate.Push(santanderValidateAgreementNumber)
	b.validate.Push(satanderBoletoTypeValidate)

	return b, nil
}

//Log retorna a referencia do log
func (b bankSantander) Log() *log.Log {
	return b.log
}

func (b bankSantander) GetTicket(boleto *models.BoletoRequest) (string, error) {
	boleto.Title.OurNumber = calculateOurNumber(boleto)
	boleto.Title.BoletoType, boleto.Title.BoletoTypeCode = getBoletoType(boleto)
	pipe := NewFlow()
	url := config.Get().URLTicketSantander
	tlsURL := strings.Replace(config.Get().URLTicketSantander, "https", "tls", 1)
	pipe.From("message://?source=inline", boleto, getRequestTicket(), tmpl.GetFuncMaps())
	pipe.To("log://?type=request&url="+url, b.log)
	duration := util.Duration(func() {
		pipe.To(tlsURL, b.transport, map[string]string{"timeout": config.Get().TimeoutToken})
	})
	metrics.PushTimingMetric("santander-get-ticket-boleto-time", duration.Seconds())
	pipe.To("log://?type=response&url="+url, b.log)
	ch := pipe.Choice()
	ch.When(Header("status").IsEqualTo("200"))
	ch.To("transform://?format=xml", getTicketResponse(), `{{.returnCode}}:::{{.ticket}}`, tmpl.GetFuncMaps())
	ch.When(Header("status").IsEqualTo("403"))
	ch.To("set://?prop=body", errors.New("403 Forbidden"))
	ch.Otherwise()
	ch.To("log://?type=request&url="+url, b.log).To("set://?prop=body", errors.New("integration error"))
	switch t := pipe.GetBody().(type) {
	case string:
		items := pipe.GetBody().(string)
		parts := strings.Split(items, ":::")
		returnCode, ticket := parts[0], parts[1]
		return ticket, checkError(returnCode)
	case error:
		return "", t
	}
	return "", nil
}

func (b bankSantander) RegisterBoleto(input *models.BoletoRequest) (models.BoletoResponse, error) {
	serviceURL := config.Get().URLRegisterBoletoSantander
	fromResponse := getResponseSantander()
	toAPI := getAPIResponseSantander()
	inputTemplate := getRequestSantander()
	santanderURL := strings.Replace(serviceURL, "https", "tls", 1)

	exec := NewFlow().From("message://?source=inline", input, inputTemplate, tmpl.GetFuncMaps())
	exec.To("log://?type=request&url="+serviceURL, b.log)
	duration := util.Duration(func() {
		exec.To(santanderURL, b.transport, map[string]string{"method": "POST", "insecureSkipVerify": "true", "timeout": config.Get().TimeoutRegister})
	})
	metrics.PushTimingMetric("santander-register-boleto-time", duration.Seconds())
	exec.To("log://?type=response&url="+serviceURL, b.log)
	ch := exec.Choice()
	ch.When(Header("status").IsEqualTo("200"))
	ch.To("transform://?format=xml", fromResponse, toAPI, tmpl.GetFuncMaps())
	ch.To("unmarshall://?format=json", new(models.BoletoResponse))
	ch.Otherwise()
	ch.To("log://?type=response&url="+serviceURL, b.log).To("apierro://")
	switch t := exec.GetBody().(type) {
	case *models.BoletoResponse:
		return *t, nil
	case error:
		return models.BoletoResponse{}, t
	}
	return models.BoletoResponse{}, models.NewInternalServerError("MP500", "Internal error")
}
func (b bankSantander) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}
	if ticket, err := b.GetTicket(boleto); err != nil {
		return models.BoletoResponse{Errors: errs}, err
	} else {
		boleto.Authentication.AuthorizationToken = ticket
	}
	return b.RegisterBoleto(boleto)
}

func (b bankSantander) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(boleto))
}

//GetBankNumber retorna o codigo do banco
func (b bankSantander) GetBankNumber() models.BankNumber {
	return models.Santander
}

func (b bankSantander) GetErrorsMap() map[string]int {
	return nil
}

func calculateOurNumber(boleto *models.BoletoRequest) uint {
	ourNumberWithDigit := strconv.Itoa(int(boleto.Title.OurNumber)) + util.OurNumberDv(strconv.Itoa(int(boleto.Title.OurNumber)), util.MOD11)
	value, _ := strconv.Atoi(ourNumberWithDigit)
	return uint(value)
}

func (b bankSantander) GetBankNameIntegration() string {
	return "Santander"
}

func santanderBoletoTypes() map[string]string {
	onceMap.Do(func() {
		m = make(map[string]string)

		m["DM"] = "02"  //Duplicata Mercantil
		m["DS"] = "04"  //Duplicata de serviço
		m["NP"] = "12"  //Nota promissória
		m["RC"] = "17"  //Recibo
		m["BDP"] = "32" //Boleto de proposta
		m["CH"] = "97"  //Cheque
		m["OUT"] = "99" //Outros
	})

	return m
}

func getBoletoType(boleto *models.BoletoRequest) (bt string, btc string) {
	if len(boleto.Title.BoletoType) < 1 {
		return "DM", "02"
	}

	btm := santanderBoletoTypes()

	if btm[strings.ToUpper(boleto.Title.BoletoType)] == "" {
		return "DM", "02"
	}

	return boleto.Title.BoletoType, btm[strings.ToUpper(boleto.Title.BoletoType)]
}
