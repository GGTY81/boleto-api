// +build integration !unit

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mundipagg/boleto-api/models"
	"github.com/stretchr/testify/assert"
)

func TestGetBoleto_WhenNotFoundKeys_ShouldReturnNotFound(t *testing.T) {
	expected := models.ErrorResponse{Code: "MP404", Message: "Not Found"}

	router := mockInstallApi()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/boleto?fmt=html&id=1234567890&pk=1234567890", nil)
	router.ServeHTTP(w, req)

	var response models.BoletoResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, 1, len(response.Errors))
	assert.Equal(t, expected.Code, response.Errors[0].Code, "O erro code deverá ser mapeado corretamente")
	assert.Equal(t, expected.Message, response.Errors[0].Message, "O erro message deverá ser mapeado corretamente")

}
