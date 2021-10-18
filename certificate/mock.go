package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/mundipagg/boleto-api/config"
)

func GenerateTestPK() []byte {
	// generate key
	privatekey, _ := rsa.GenerateKey(rand.Reader, 2048)

	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privatekey),
	}

	return pem.EncodeToMemory(privateKeyBlock)
}

func LoadMockCertificates() {
	SetCertificateOnStore(config.Get().AzureStorageOpenBankSkName, GenerateTestPK())
	SetCertificateOnStore(config.Get().AzureStorageJPMorganPkName, GenerateTestPK())
	SetCertificateOnStore(config.Get().AzureStorageJPMorganCrtName, GenerateTestPK())
	SetCertificateOnStore(config.Get().AzureStorageJPMorganSignCrtName, GenerateTestPK())
	SetCertificateOnStore(config.Get().CertificateSSLName, GenerateTestPK())
	SetCertificateOnStore(config.Get().CertificateICPName, GenerateTestPK())
}
