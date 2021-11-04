package token

import (
	"encoding/json"

	jwt "github.com/dgrijalva/jwt-go"
)

type JwtGenerator interface {
	Sign(payload interface{}, privateKey []byte) (string, error)
}

func GetJwtGenerator(alg string) JwtGenerator {
	switch alg {
	case "RS256":
		return NewRS256()
	default:
		return NewRS256()
	}
}

func createClaims(payload interface{}) jwt.MapClaims {
	claimsMap := make(jwt.MapClaims)
	switch p := payload.(type) {
	case string:
		err := json.Unmarshal([]byte(p), &claimsMap)
		if err != nil {
			panic(err)
		}
		dump(claimsMap)

	case map[string]interface{}:
		for k, v := range p {
			claimsMap[k] = v
		}
	}
	return claimsMap
}

func dump(m map[string]interface{}) {
	for _, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			dump(mv)
		}
	}
}
