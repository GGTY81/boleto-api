package test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

//Parameter Parâmetro de teste com input generico
type Parameter struct {
	Input    interface{}
	Expected interface{}
	Length   int
}

type XmlNode struct {
	XMLName xml.Name
	Content []byte    `xml:",innerxml"`
	Nodes   []XmlNode `xml:",any"`
}

//UInt64TestParameter Parâmetro de teste com input do tipo uint64
type UInt64TestParameter struct {
	Input    uint64
	Expected string
}

// ExpectNoError falha o teste se e != nil
func ExpectNoError(e error, t *testing.T) {
	if e != nil {
		t.Fail()
	}
}

// ExpectError falha o teste se e == nil
func ExpectError(e error, t *testing.T) {
	if e == nil {
		t.Fail()
	}
}

// ExpectTrue falha o teste caso a condição não seja verdadeira
func ExpectTrue(condition bool, t *testing.T) {
	if !condition {
		t.Fail()
	}
}

// ExpectFalse falha o teste caso a condição não seja falsa
func ExpectFalse(condition bool, t *testing.T) {
	if condition {
		t.Fail()
	}
}

// ExpectNil falha o teste caso obj seja diferente de nil
func ExpectNil(obj interface{}, t *testing.T) {
	if obj != nil {
		t.Fail()
	}
}

//AssertProcessBoletoWithSuccess Valida se o boleto foi gerado com sucesso
func AssertProcessBoletoWithSuccess(t *testing.T, response models.BoletoResponse) {
	assert.Empty(t, response.Errors, "Não deve ocorrer erros")
	assert.NotEmpty(t, response.BarCodeNumber, "Deve haver um Barcode")
	assert.NotEmpty(t, response.DigitableLine, "Deve haver uma linha digitável")
}

//AssertProcessBoletoFailed Valida se o houve um erro no processamento do boleto
func AssertProcessBoletoFailed(t *testing.T, response models.BoletoResponse) {
	assert.Greater(t, len(response.Errors), 0, "Devem ocorrer erros ")
	assert.Empty(t, response.BarCodeNumber, "Não deve haver um Barcode")
	assert.Empty(t, response.DigitableLine, "Não deve haver uma linha digitável")
}

//AssertError Valida a existência de erros internos
func AssertError(t *testing.T, err error, errType interface{}) {
	assert.NotNil(t, err, "Deve haver um erro")
	assert.NotEmpty(t, err.Error(), "Deve haver uma mensagem de erro")
	assert.IsType(t, errType, err, "Deve ser um erro do tipo "+fmt.Sprintf("%T", err))
}

//CreateClientIP cria IP no contexto
func CreateClientIP(c *gin.Context) {
	c.Request = new(http.Request)
	c.Request.Header = make(map[string][]string)
	c.Request.Header.Add("X-Forwarded-For", "0.0.0.0")
}

func WalkThroughXml(nodes []XmlNode, f func(XmlNode) bool) {
	for _, n := range nodes {
		if f(n) {
			WalkThroughXml(n.Nodes, f)
		}
	}
}

func GetNode(bodyXml string, tagName string) string {
	nodeContent := ""

	buf := bytes.NewBuffer([]byte(bodyXml))
	dec := xml.NewDecoder(buf)

	var n XmlNode
	err := dec.Decode(&n)
	if err != nil {
		panic(err)
	}

	WalkThroughXml([]XmlNode{n}, func(n XmlNode) bool {
		if n.XMLName.Local == tagName {
			nodeContent += string(n.Content)
		}
		return true
	})
	return nodeContent
}
