package util

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mundipagg/boleto-api/certificate"
	"github.com/mundipagg/boleto-api/log"

	s "github.com/fullsailor/pkcs7"
	"github.com/mundipagg/boleto-api/config"
)

var defaultDialer = &net.Dialer{Timeout: 16 * time.Second, KeepAlive: 16 * time.Second}

var (
	client            *http.Client
	onceDefaultClient = &sync.Once{}
	onceTransport     = &sync.Once{}
	icpCert           certificate.ICPCertificate
	transport         *http.Transport
)

// HTTPInterface is an abstraction for HTTP client
type HTTPInterface interface {
	Post(url string, headers map[string]string, body interface{}) (*http.Response, error)
}

// HTTPClient is the struct for making requests
type HTTPClient struct{}

// PostFormEncoded is a function for making requests using Post Http method with content-type application/x-www-form-urlencoded.
//
// It receives an endpoint, params and pointer for log and it creates a new Post request, returning []byte and a error.
func (hc *HTTPClient) PostFormURLEncoded(endpoint string, params map[string]string, log *log.Log) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	uri, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return []byte(""), err
	}

	values := uri.Query()
	for k, v := range params {
		values.Set(k, v)
	}

	req, err := http.NewRequest(http.MethodPost, uri.String(), strings.NewReader(values.Encode())) // URL-encoded payload

	if err != nil {
		return []byte(""), err
	}

	header := map[string]string{
		"content-type":   "application/x-www-form-urlencoded",
		"content-length": strconv.Itoa(len(values.Encode())),
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	log.Request(params, endpoint, header)
	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte(""), fmt.Errorf("stone authentication returns status code %d", resp.StatusCode)
	}

	respByte, err := ioutil.ReadAll(resp.Body)
	log.Response(string(respByte), endpoint, nil)

	return respByte, err
}

// DefaultHTTPClient retorna um cliente http configurado para dar um skip na validação do certificado digital
func DefaultHTTPClient() *http.Client {
	onceDefaultClient.Do(func() {
		client = &http.Client{
			Transport: &http.Transport{
				Dial:                defaultDialer.Dial,
				TLSHandshakeTimeout: 16 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	})
	return client
}

//Post faz um requisição POST para uma URL e retorna o response, status e erro
func PostReponseWithHeader(url, body, timeout string, header map[string]string) (string, string, int, error) {
	return doRequest("POST", url, body, timeout, header)
}

//Post faz um requisição POST para uma URL e retorna o response, status e erro
func Post(url, body, timeout string, header map[string]string) (string, int, error) {
	resp, _, st, err := doRequest("POST", url, body, timeout, header)
	return resp, st, err
}

//PostWithHeader faz um requisição POST para uma URL e retorna o response, status e erro
func PostWithHeader(url, body, timeout string, header map[string]string) (string, map[string]interface{}, int, error) {
	resp, respHeader, st, err := doRequestWithHeaderObject("POST", url, body, timeout, header)
	return resp, respHeader, st, err
}

func doRequest(method, url, body, timeout string, header map[string]string) (string, string, int, error) {
	t := GetDurationTimeoutRequest(timeout) * time.Second

	ctx, cls := context.WithTimeout(context.Background(), t)
	defer cls()

	client := DefaultHTTPClient()

	message := strings.NewReader(body)

	req, err := http.NewRequestWithContext(ctx, method, url, message)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	resp, errResp := client.Do(req)
	if errResp != nil {
		return "", "", 0, errResp
	}
	defer resp.Body.Close()
	respHeader := fmt.Sprintf("%v", resp.Header)

	data, errResponse := ioutil.ReadAll(resp.Body)
	if errResponse != nil {
		return "", respHeader, resp.StatusCode, errResponse
	}
	sData := string(data)
	return sData, respHeader, resp.StatusCode, nil
}

func doRequestWithHeaderObject(method, url, body, timeout string, header map[string]string) (string, map[string]interface{}, int, error) {
	t := GetDurationTimeoutRequest(timeout) * time.Second

	ctx, cls := context.WithTimeout(context.Background(), t)
	defer cls()

	client := DefaultHTTPClient()

	message := strings.NewReader(body)

	req, err := http.NewRequestWithContext(ctx, method, url, message)
	if err != nil {
		return "", nil, http.StatusInternalServerError, err
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, errResp := client.Do(req)
	if errResp != nil {
		return "", nil, 0, errResp
	}
	defer resp.Body.Close()

	respHeader := convertHeader(resp.Header)

	data, errResponse := ioutil.ReadAll(resp.Body)
	if errResponse != nil {
		return "", respHeader, resp.StatusCode, errResponse
	}
	sData := string(data)
	return sData, respHeader, resp.StatusCode, nil
}

func convertHeader(respHeader map[string][]string) map[string]interface{} {
	m2 := make(map[string]interface{}, len(respHeader))
	for k, v := range respHeader {
		m2[k] = v[0]
	}
	return m2
}

// BuildTLSTransport creates a TLS Client Transport from crt, ca and key files
func BuildTLSTransport(con certificate.TLSCertificate) (*http.Transport, error) {

	if config.Get().MockMode {
		return nil, nil
	}

	key, err := certificate.GetCertificateFromStore(con.Key)
	if err != nil {
		return nil, err
	}

	crt, err := certificate.GetCertificateFromStore(con.Crt)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(getCertificateByType(crt), getCertificateByType(key))

	if err != nil {
		return nil, err
	}

	transport = &http.Transport{
		Dial:                defaultDialer.Dial,
		TLSHandshakeTimeout: 16 * time.Second,
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		},
	}

	return transport, nil
}

func getCertificateByType(key interface{}) []byte {
	switch v := key.(type) {
	case certificate.SSLCertificate:
		return v.PemData
	default:
		return v.([]byte)
	}
}

//Sign request
func SignRequest(request string) (string, error) {

	if icpCert == (certificate.ICPCertificate{}) {
		icp, err := certificate.GetCertificateFromStore(config.Get().CertificateICPName)
		if err != nil {
			return "", err
		}
		icpCert = icp.(certificate.ICPCertificate)
	}

	signedData, err := s.NewSignedData([]byte(request))
	if err != nil {
		return "", err
	}

	if err := signedData.AddSigner(icpCert.Certificate, icpCert.RsaPrivateKey, s.SignerInfoConfig{}); err != nil {
		return "", err
	}

	detachedSignature, err := signedData.Finish()
	if err != nil {
		return "", err
	}

	signedRequest := base64.StdEncoding.EncodeToString(detachedSignature)

	return signedRequest, nil
}

//Read privatekey and parse to PKCS#1
func parsePrivateKey() (crypto.PrivateKey, error) {

	pkeyBytes, err := ioutil.ReadFile(config.Get().CertICP_PathPkey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pkeyBytes)
	if block == nil {
		return nil, errors.New("Key Not Found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return rsa, nil
	default:
		return nil, fmt.Errorf("SSH: Unsupported key type %q", block.Type)
	}

}

///Read chainCertificates and adapter to x509.Certificate
func parseChainCertificates() (*x509.Certificate, error) {

	chainCertsBytes, err := ioutil.ReadFile(config.Get().CertICP_PathChainCertificates)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(chainCertsBytes)
	if block == nil {
		return nil, errors.New("Key Not Found")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func doRequestTLS(method, url, body, timeout string, header map[string]string, transport *http.Transport) (string, int, error) {
	tlsClient := &http.Client{}
	tlsClient.Transport = transport
	tlsClient.Timeout = GetDurationTimeoutRequest(timeout) * time.Second
	b := strings.NewReader(body)
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return "", 0, err
	}

	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	resp, err := tlsClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	sData := string(data)
	return sData, resp.StatusCode, nil
}

func doRequestTLSWithHeader(method, url, body, timeout string, header map[string]string, transport *http.Transport) (string, map[string]interface{}, int, error) {
	tlsClient := &http.Client{}
	tlsClient.Transport = transport
	tlsClient.Timeout = GetDurationTimeoutRequest(timeout) * time.Second
	b := strings.NewReader(body)
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return "", nil, 0, err
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := tlsClient.Do(req)
	if err != nil {
		return "", nil, 0, err
	}
	respHeader := convertHeader(resp.Header)
	defer resp.Body.Close()
	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, 0, err
	}
	sData := string(data)
	return sData, respHeader, resp.StatusCode, nil
}

func PostTLS(url, body, timeout string, header map[string]string, transport *http.Transport) (string, int, error) {
	return doRequestTLS("POST", url, body, timeout, header, transport)
}

func PostTLSWithHeader(url, body, timeout string, header map[string]string, transport *http.Transport) (string, map[string]interface{}, int, error) {
	return doRequestTLSWithHeader("POST", url, body, timeout, header, transport)
}

//HeaderToMap converte um http Header para um dicionário string -> string
func HeaderToMap(header http.Header) map[string]string {
	headerMap := make(map[string]string)
	for key, value := range header {
		if key == "Authorization" {
			headerMap[key] = "***"
		} else {
			headerMap[key] = value[0]
		}
	}
	return headerMap
}
