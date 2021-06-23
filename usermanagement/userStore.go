package usermanagement

import (
	"fmt"
	"sync"

	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
)

var userCredentialStorage = sync.Map{}

func addUser(key string, value interface{}) {
	userCredentialStorage.Store(key, value)
}

//GetUser Busca credenciais de um usuário
func GetUser(key string) (interface{}, bool) {
	if value, ok := userCredentialStorage.Load(key); ok {
		return value, true
	}
	return nil, false
}

//LoadUserCredentials Carrega credenciais salvas no banco de dados
func LoadUserCredentials() {
	log := log.CreateLog()

	c, err := db.GetUserCredentials()
	if err != nil {
		log.Error(err.Error(), fmt.Sprintf("Error in get user crendentials - %s", err.Error()))
		return
	}

	for _, u := range c {
		u.UserKey = u.ID.Hex()
		addUser(u.UserKey, u)
	}
}
