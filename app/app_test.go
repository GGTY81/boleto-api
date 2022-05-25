//go:build integration || !unit
// +build integration !unit

package app

import (
	"testing"

	"github.com/mundipagg/boleto-api/certificate"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/mock"
	"github.com/stretchr/testify/assert"
)

func Test_installCertificates(t *testing.T) {
	mock.StartMockService("9093")

	tests := []struct {
		name        string
		certificate string
	}{
		{
			name:        "Fetch sk OpenBank from localStorage successfully",
			certificate: config.Get().AzureStorageOpenBankSkName,
		},
		{
			name:        "Fetch JPMorganPk from localStorage successfully",
			certificate: config.Get().AzureStorageJPMorganPkName,
		},
		{
			name:        "Fetch JPMorganCrt from localStorage successfully",
			certificate: config.Get().AzureStorageJPMorganCrtName,
		},
		{
			name:        "Fetch JPMorganSignCrt from localStorage successfully",
			certificate: config.Get().AzureStorageJPMorganSignCrtName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installCertificates()
		})
		sk, err := certificate.GetCertificateFromStore(tt.certificate)
		assert.Nil(t, err)
		assert.NotNil(t, sk)
	}
}
