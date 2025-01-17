package mock

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const contentApplication = "text/json"

const success = `
   {
	"account_id": "946b50ce-ed5d-45ab-8c86-ce3baf90a73a",
	"amount": 5000,
	"barcode": "19799900100000050000000038911660052705632042",
	"beneficiary": {
		"account_code": "23172018",
		"branch_code": "1",
		"document": "14994237000140",
		"document_type": "cnpj",
		"legal_name": "MUNDIPAGG TECNOLOGIA EM PAGAMENTOS S.A.",
		"trade_name": "MUNDIPAGG TECNOLOGIA EM PAGAMENTOS S.A."
	},
	"created_at": "2021-07-01T14:32:21Z",
	"created_by": "application:3279b005-5e40-41c1-996e-8cec24f8006b",
	"customer": {
		"document": "13621248773",
		"document_type": "cpf",
		"legal_name": "Matheus Palanowski",
		"trade_name": null
	},
	"discounts": [],
	"expiration_date": "2022-05-30",
	"expired_at": null,
	"fee": 0,
	"fee_metadata": {
		"billing_exemption_participant": true,
		"fee": 0,
		"max_free": 5,
		"original_fee": 200,
		"remaining_free": 5
	},
	"fine": null,
	"id": "46e902e7-05e6-4efb-a28f-cf8b16ce9eed",
	"interest": null,
	"invoice_type": "bill_of_exchange",
	"issuance_date": "2021-07-01",
	"limit_date": "2022-05-30",
	"our_number": "38911660052705632042",
	"receiver": null,
	"registered_at": null,
	"settled_at": null,
	"status": "CREATED",
	"writable_line": "19790000053891166005827056320420990010000005000"
}
`
const successWithoutOurNumber = `
   {
	"account_id": "946b50ce-ed5d-45ab-8c86-ce3baf90a73a",
	"amount": 5000,
	"barcode": "19799900100000050000000038911660052705632042",
	"beneficiary": {
		"account_code": "23172018",
		"branch_code": "1",
		"document": "14994237000140",
		"document_type": "cnpj",
		"legal_name": "MUNDIPAGG TECNOLOGIA EM PAGAMENTOS S.A.",
		"trade_name": "MUNDIPAGG TECNOLOGIA EM PAGAMENTOS S.A."
	},
	"created_at": "2021-07-01T14:32:21Z",
	"created_by": "application:3279b005-5e40-41c1-996e-8cec24f8006b",
	"customer": {
		"document": "13621248773",
		"document_type": "cpf",
		"legal_name": "Matheus Palanowski",
		"trade_name": null
	},
	"discounts": [],
	"expiration_date": "2022-05-30",
	"expired_at": null,
	"fee": 0,
	"fee_metadata": {
		"billing_exemption_participant": true,
		"fee": 0,
		"max_free": 5,
		"original_fee": 200,
		"remaining_free": 5
	},
	"fine": null,
	"id": "46e902e7-05e6-4efb-a28f-cf8b16ce9eed",
	"interest": null,
	"invoice_type": "bill_of_exchange",
	"issuance_date": "2021-07-01",
	"limit_date": "2022-05-30",
	"our_number": "",
	"receiver": null,
	"registered_at": null,
	"settled_at": null,
	"status": "CREATED",
	"writable_line": "19790000053891166005827056320420990010000005000"
}
`

const successWithFineOnly = `
{
    "account_id": "7e785e98-c859-46c6-9dc7-37ea522ceadf",
    "amount": 5000,
    "barcode": "19791895800000050000000088115105980288509901",
    "beneficiary": {
        "account_code": "88266374",
        "address": {
            "city": "Rio de Janeiro",
            "country": "Brasil",
            "extra": null,
            "neighborhood": "Vila Isabel",
            "postal_code": "20551080",
            "state": "RJ",
            "street": "Rua Visconde de Abaete",
            "street_number": "97"
        },
        "branch_code": "1",
        "document": "92961229022",
        "document_type": "cpf",
        "legal_name": "Nome do CPF 92961229022",
        "trade_name": null
    },
    "cancelled_at": null,
    "cancelled_reason": null,
    "cancelled_reason_description": null,
    "created_at": "2022-04-12T17:47:33Z",
    "created_by": "application:3279b005-5e40-41c1-996e-8cec24f8006b",
    "customer": {
        "document": "13621248773",
        "document_type": "cpf",
        "legal_name": "Matheus Palanowski",
        "trade_name": null
    },
    "discounts": [],
    "expiration_date": "2022-04-17",
    "expired_at": null,
    "fee": 0,
    "fee_metadata": {
        "billing_exemption_participant": true,
        "fee": 0,
        "max_free": 5,
        "original_fee": 200,
        "remaining_free": 5
    },
    "fine": {
        "date": "2022-04-18",
        "value": "2.0"
    },
    "id": "ea95560a-6455-4cf9-83dc-af11f9ccf734",
    "interest": null,
    "invoice_type": "bill_of_exchange",
    "issuance_date": "2022-04-12",
    "limit_date": "2022-04-22",
    "our_number": "88115105980288509901",
    "receiver": null,
    "registered_at": null,
    "settled_at": null,
    "status": "CREATED",
    "writable_line": "19790000058811510598802885099016189580000005000"
}
`

const failedFineValue = `
{
    "reason": [
        {
            "error": "Percentage must be equal or lower than 2.0",
            "path": [
                "fine",
                "value"
            ]
        }
    ],
    "type": "srn:error:validation"
}
`

const failedFineDate = `
{
    "reason": [
        {
            "error": "fine date should be greater than expiration date",
            "path": [
                "fine"
            ]
        }
    ],
    "type": "srn:error:validation"
}
`

const successWithInterestOnly = `
{
    "account_id": "7e785e98-c859-46c6-9dc7-37ea522ceadf",
    "amount": 5000,
    "barcode": "19797895800000050000000015924667218367048068",
    "beneficiary": {
        "account_code": "88266374",
        "address": {
            "city": "Rio de Janeiro",
            "country": "Brasil",
            "extra": null,
            "neighborhood": "Vila Isabel",
            "postal_code": "20551080",
            "state": "RJ",
            "street": "Rua Visconde de Abaete",
            "street_number": "97"
        },
        "branch_code": "1",
        "document": "92961229022",
        "document_type": "cpf",
        "legal_name": "Nome do CPF 92961229022",
        "trade_name": null
    },
    "cancelled_at": null,
    "cancelled_reason": null,
    "cancelled_reason_description": null,
    "created_at": "2022-04-12T19:15:20Z",
    "created_by": "application:3279b005-5e40-41c1-996e-8cec24f8006b",
    "customer": {
        "document": "13621248773",
        "document_type": "cpf",
        "legal_name": "Matheus Palanowski",
        "trade_name": null
    },
    "discounts": [],
    "expiration_date": "2022-04-17",
    "expired_at": null,
    "fee": 0,
    "fee_metadata": {
        "billing_exemption_participant": true,
        "fee": 0,
        "max_free": 5,
        "original_fee": 200,
        "remaining_free": 5
    },
    "fine": null,
    "id": "b4541f0a-0e76-4bc4-82bd-2897b1ca1563",
    "interest": {
        "date": "2022-04-18",
        "value": "1.0"
    },
    "invoice_type": "bill_of_exchange",
    "issuance_date": "2022-04-12",
    "limit_date": "2022-04-22",
    "our_number": "15924667218367048068",
    "receiver": null,
    "registered_at": null,
    "settled_at": null,
    "status": "CREATED",
    "writable_line": "19790000051592466721383670480686789580000005000"
}
`

const failedInterestValue = `
{
    "reason": [
        {
            "error": "Percentage must be equal or lower than 1.0",
            "path": [
                "interest",
                "value"
            ]
        }
    ],
    "type": "srn:error:validation"
}
`

const failedInterestDate = `
{
    "reason": [
        {
            "error": "interest date should be greater than expiration date",
            "path": [
                "interest"
            ]
        }
    ],
    "type": "srn:error:validation"
}
`

const successWithFineAndInterest = `
{
    "account_id": "7e785e98-c859-46c6-9dc7-37ea522ceadf",
    "amount": 5000,
    "barcode": "19798895800000050000000037485189119692884798",
    "beneficiary": {
        "account_code": "88266374",
        "address": {
            "city": "Rio de Janeiro",
            "country": "Brasil",
            "extra": null,
            "neighborhood": "Vila Isabel",
            "postal_code": "20551080",
            "state": "RJ",
            "street": "Rua Visconde de Abaete",
            "street_number": "97"
        },
        "branch_code": "1",
        "document": "92961229022",
        "document_type": "cpf",
        "legal_name": "Nome do CPF 92961229022",
        "trade_name": null
    },
    "cancelled_at": null,
    "cancelled_reason": null,
    "cancelled_reason_description": null,
    "created_at": "2022-04-12T19:36:59Z",
    "created_by": "application:3279b005-5e40-41c1-996e-8cec24f8006b",
    "customer": {
        "document": "13621248773",
        "document_type": "cpf",
        "legal_name": "Matheus Palanowski",
        "trade_name": null
    },
    "discounts": [],
    "expiration_date": "2022-04-17",
    "expired_at": null,
    "fee": 0,
    "fee_metadata": {
        "billing_exemption_participant": true,
        "fee": 0,
        "max_free": 5,
        "original_fee": 200,
        "remaining_free": 5
    },
    "fine": {
        "date": "2022-04-18",
        "value": "2.0"
    },
    "id": "c073397a-b87c-4718-830a-eb50eb6d1f92",
    "interest": {
        "date": "2022-04-18",
        "value": "1.0"
    },
    "invoice_type": "bill_of_exchange",
    "issuance_date": "2022-04-12",
    "limit_date": "2022-04-22",
    "our_number": "37485189119692884798",
    "receiver": null,
    "registered_at": null,
    "settled_at": null,
    "status": "CREATED",
    "writable_line": "19790000053748518911496928847985889580000005000"
}
`

const unauthenticated = `{ "type": "srn:error:unauthenticated" }`

const unauthorized = `{ "type": "srn:error:unauthorized" }`

const conflict = `{	"type": "srn:error:conflict"  }`

const unprocessableEntity = `{"reason":"barcode_payment_invoice_bill_of_exchange is not ena bled on this account","type":"srn:error:product_not_enabled"}`

const customer_doc_invalid = `{"reason":[{"error":"is invalid","path":["customer","document"]}],"type":"srn:error:validation"}`

const customer_blank_name = `{"reason":[{"error":"can't be blank","path":["customer","legal_name"]}],"type":"srn:error:validation"}`

const amount_not_allowed = `{"reason":[{"error":"not allowed","path":["amount"]}],"type": "srn:error:validation"}`

const multipleValidationErrorsPath = `{"reason":[{"error":"is invalid","path":["receiver","document"]}],"type": "srn:error:validation"}`

const multipleValidationErrorsReason = `{"reason":[{"error": "is invalid","path":["account_id"]},{"error": "not allowed","path":["amount"]}],"type":"srn:error:validation"}`

const tkStone = `{
	"access_token": "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI4d3NUd3BhYTRJWUZIYWV5ZFRubnRoRC1UaVlCaU9kanNmOGx6RUlMR1hVIn0.eyJqdGkiOiIyZTlkNGZkMy0zN2M1LTRjOWUtYTJjYy1lMjQ1N2MxZDgyMWQiLCJleHAiOjE2MjQ4OTQ1NDQsIm5iZiI6MCwiaWF0IjoxNjI0ODkzNjQ0LCJpc3MiOiJodHRwczovL3NhbmRib3gtYWNjb3VudHMub3BlbmJhbmsuc3RvbmUuY29tLmJyL2F1dGgvcmVhbG1zL3N0b25lX2JhbmsiLCJzdWIiOiJkNDY0ZDg3MC1mYzc2LTRjZGMtYWM5OC1hNjcyYjYyOTdhOGYiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiIzMjc5YjAwNS01ZTQwLTQxYzEtOTk2ZS04Y2VjMjRmODAwNmIiLCJhdXRoX3RpbWUiOjAsInNlc3Npb25fc3RhdGUiOiJhM2MyYzY3OC0wOGIxLTRmNmQtYmQ2Yi0wNjgzMjQ1M2UzNmMiLCJhY3IiOiIxIiwic2NvcGUiOiJwYXltZW50YWNjb3VudDpwYXltZW50bGlua3M6d3JpdGUgcGF5bWVudGFjY291bnQ6Y29udGFjdDp3cml0ZSBwaXg6cGF5bWVudF9pbnZvaWNlIHBpeDpwYXltZW50IHBpeDplbnRyeV9jbGFpbSBwYXltZW50YWNjb3VudDpyZWFkIHBpeDplbnRyeSBwYXltZW50YWNjb3VudDp0cmFuc2ZlcnM6aW50ZXJuYWwgcGF5bWVudGFjY291bnQ6ZmVlczpyZWFkIHBheW1lbnRhY2NvdW50OnBheW1lbnRzIHN0b25lX3N1YmplY3RfaWQgcGF5bWVudGFjY291bnQ6Y29udGFjdDpyZWFkIHNpZ251cDpwYXltZW50YWNjb3VudCBwYXltZW50YWNjb3VudDpib2xldG9pc3N1YW5jZSBwYXltZW50YWNjb3VudDpwYXltZW50bGlua3M6cmVhZCBwYXltZW50YWNjb3VudDp0cmFuc2ZlcnM6ZXh0ZXJuYWwiLCJjbGllbnRJZCI6IjMyNzliMDA1LTVlNDAtNDFjMS05OTZlLThjZWMyNGY4MDA2YiIsImNsaWVudEhvc3QiOiIxMC4xMC4zLjE3MiIsInN0b25lX3N1YmplY3RfaWQiOiJhcHBsaWNhdGlvbjozMjc5YjAwNS01ZTQwLTQxYzEtOTk2ZS04Y2VjMjRmODAwNmIiLCJjbGllbnRBZGRyZXNzIjoiMTAuMTAuMy4xNzIifQ.JloXzaTUFW0IVDi191U_WujRLhIIPiZUZngDb1nbhHo9mclG176CIgdSsBPmoOZr35ry47JCLgEq5ZAos8Sts72kpi1BivvVq0rJn5_NrmSyb0zqMSK4sNYzbhBafK7U6wamUZCjDeJmQ_wBUDNvPxGC1gToreMFnhrbak0pQr_CWp9Csgkn-9QUvFFpTRkJ3fdca57YnKoGsEWJWMs8Suq6g097244EWHISlUtO1ZGt01mypDeU8g5Z_eYD8qdN_woUeCGL86QhDoH-V8Dl_NIwbsHGTm8iRDTqjRBid2XH6Cj0RAMH10EpTKI8buSBzJ872bKLoCwXQQUnWIYT4Q",
	"expires_in": 900,
	"refresh_expires_in": 2700,
	"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJiNzJiOTVmZC0zOWVjLTRmZjktYTRkNS1lOGY0YTlmNTNmM2EifQ.eyJqdGkiOiJmNDg5YjU1My0wZjUwLTQ4ZjktYmU3OS0zZWQ3ZDAzZDIyZDQiLCJleHAiOjE2MjQ4OTYzNDQsIm5iZiI6MCwiaWF0IjoxNjI0ODkzNjQ0LCJpc3MiOiJodHRwczovL3NhbmRib3gtYWNjb3VudHMub3BlbmJhbmsuc3RvbmUuY29tLmJyL2F1dGgvcmVhbG1zL3N0b25lX2JhbmsiLCJhdWQiOiJodHRwczovL3NhbmRib3gtYWNjb3VudHMub3BlbmJhbmsuc3RvbmUuY29tLmJyL2F1dGgvcmVhbG1zL3N0b25lX2JhbmsiLCJzdWIiOiJkNDY0ZDg3MC1mYzc2LTRjZGMtYWM5OC1hNjcyYjYyOTdhOGYiLCJ0eXAiOiJSZWZyZXNoIiwiYXpwIjoiMzI3OWIwMDUtNWU0MC00MWMxLTk5NmUtOGNlYzI0ZjgwMDZiIiwiYXV0aF90aW1lIjowLCJzZXNzaW9uX3N0YXRlIjoiYTNjMmM2NzgtMDhiMS00ZjZkLWJkNmItMDY4MzI0NTNlMzZjIiwic2NvcGUiOiJwYXltZW50YWNjb3VudDpwYXltZW50bGlua3M6d3JpdGUgcGF5bWVudGFjY291bnQ6Y29udGFjdDp3cml0ZSBwaXg6cGF5bWVudF9pbnZvaWNlIHBpeDpwYXltZW50IHBpeDplbnRyeV9jbGFpbSBwYXltZW50YWNjb3VudDpyZWFkIHBpeDplbnRyeSBwYXltZW50YWNjb3VudDp0cmFuc2ZlcnM6aW50ZXJuYWwgcGF5bWVudGFjY291bnQ6ZmVlczpyZWFkIHBheW1lbnRhY2NvdW50OnBheW1lbnRzIHN0b25lX3N1YmplY3RfaWQgcGF5bWVudGFjY291bnQ6Y29udGFjdDpyZWFkIHNpZ251cDpwYXltZW50YWNjb3VudCBwYXltZW50YWNjb3VudDpib2xldG9pc3N1YW5jZSBwYXltZW50YWNjb3VudDpwYXltZW50bGlua3M6cmVhZCBwYXltZW50YWNjb3VudDp0cmFuc2ZlcnM6ZXh0ZXJuYWwifQ.hXQSiQ-Bbto35TjlOdtnxUAbNysiS3TZIhqOpxV7A2s",
	"token_type": "bearer",
	"not-before-policy": 1620910623,
	"session_state": "a3c2c678-08b1-4f6d-bd6b-06832453e36c",
	"scope": "paymentaccount:paymentlinks:write paymentaccount:contact:write pix:payment_invoice pix:payment pix:entry_claim paymentaccount:read pix:entry paymentaccount:transfers:internal paymentaccount:fees:read paymentaccount:payments stone_subject_id paymentaccount:contact:read signup:paymentaccount paymentaccount:boletoissuance paymentaccount:paymentlinks:read paymentaccount:transfers:external"
}`

const specialCaractereError = `{"type":"srn:error:bad_request"}`

func authStone(c *gin.Context) {
	c.Data(200, "text/json", []byte(tkStone))
}

func registerStone(c *gin.Context) {
	d, _ := ioutil.ReadAll(c.Request.Body)
	json := string(d)

	if strings.Contains(json, `amount": 201,`) {
		c.Data(201, contentApplication, []byte(success))
	} else if strings.Contains(json, `amount": 200,`) {
		c.Data(201, contentApplication, []byte(successWithoutOurNumber))
	} else if strings.Contains(json, `amount": 202,`) {
		c.Data(201, contentApplication, []byte(successWithFineOnly))
	} else if strings.Contains(json, `amount": 203,`) {
		c.Data(201, contentApplication, []byte(successWithInterestOnly))
	} else if strings.Contains(json, `amount": 204,`) {
		c.Data(201, contentApplication, []byte(successWithFineAndInterest))
	} else if strings.Contains(json, `amount": 401,`) {
		c.Data(401, contentApplication, []byte(unauthenticated))
	} else if strings.Contains(json, `amount": 403,`) {
		c.Data(403, contentApplication, []byte(unauthorized))
	} else if strings.Contains(json, `amount": 409,`) {
		c.Data(409, contentApplication, []byte(conflict))
	} else if strings.Contains(json, `amount": 422,`) {
		c.Data(422, contentApplication, []byte(unprocessableEntity))
	} else if strings.Contains(json, `amount": 301,`) {
		c.Data(400, contentApplication, []byte(failedFineValue))
	} else if strings.Contains(json, `amount": 302,`) {
		c.Data(400, contentApplication, []byte(failedInterestValue))
	} else if strings.Contains(json, `amount": 303,`) {
		c.Data(400, contentApplication, []byte(failedFineDate))
	} else if strings.Contains(json, `amount": 304,`) {
		c.Data(400, contentApplication, []byte(failedInterestDate))
	} else if strings.Contains(json, `amount": 4001,`) {
		c.Data(400, contentApplication, []byte(customer_doc_invalid))
	} else if strings.Contains(json, `amount": 4002,`) {
		c.Data(400, contentApplication, []byte(customer_blank_name))
	} else if strings.Contains(json, `amount": 4003,`) {
		c.Data(400, contentApplication, []byte(amount_not_allowed))
	} else if strings.Contains(json, `amount": 4004,`) {
		c.Data(400, contentApplication, []byte(multipleValidationErrorsPath))
	} else if strings.Contains(json, `amount": 4005,`) {
		c.Data(400, contentApplication, []byte(multipleValidationErrorsReason))
	} else if strings.Contains(json, `amount": 504,`) {
		time.Sleep(35 * time.Second)
		c.Data(504, contentApplication, []byte("timeout"))
	} else if strings.Contains(json, `amount": 505,`) {
		c.Data(500, contentApplication, []byte(specialCaractereError))
	} else {
		c.Data(401, contentApplication, []byte(unauthenticated))
	}
}
