package bank

import (
	"fmt"

	"github.com/mundipagg/boleto-api/bank/services/jpmorgan"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
)

//Bank é a interface que vai oferecer os serviços em comum entre os bancos
type Bank interface {
	ProcessBoleto(*models.BoletoRequest) (models.BoletoResponse, error)
	RegisterBoleto(*models.BoletoRequest) (models.BoletoResponse, error)
	ValidateBoleto(*models.BoletoRequest) models.Errors
	GetBankNumber() models.BankNumber
	GetBankNameIntegration() string
	GetErrorsMap() map[string]int
	Log() *log.Log
}

//Get retorna estrategia de acordo com o banco ou erro caso o banco não exista
func Get(boleto models.BoletoRequest) (Bank, error) {
	switch boleto.BankNumber {
	case models.BancoDoBrasil:
		return getIntegrationBB(boleto)
	case models.Bradesco:
		return getIntegrationBradesco(boleto)
	case models.Caixa:
		return getIntegrationCaixa(boleto)
	case models.Citibank:
		return getIntegrationCitibank(boleto)
	case models.Santander:
		return getIntegrationSantander(boleto)
	case models.Itau:
		return getIntegrationItau(boleto)
	case models.Pefisa:
		return getIntegrationPefisa(boleto)
	case models.Stone:
		return getIntegrationStone(boleto)
	case models.JPMorgan:
		return jpmorgan.New()
	default:
		return nil, models.NewErrorResponse("MPBankNumber", fmt.Sprintf("Banco %d não existe", boleto.BankNumber))
	}
}
