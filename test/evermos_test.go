package test

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	// "log"

	// "encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	// log "github.com/sirupsen/logrus"
	// "github.com/newrelic/go-agent/internal/logger"
	"github.com/stretchr/testify/assert"
)

const (
	token_uri = "http://localhost:8080/v1/authentication/get_token?phone_number=6287777000057&password=bambang@12345"
)

type GenerallResponse struct {
	RequestId string        `json:"request_id"`
	Content   TokenResponse `json:"content"`
	Status    int           `json:"status"`
}
type TokenResponse struct {
	Token string `json:"token"`
}

func GetToken(phone_number string, password string) (GenerallResponse, error) {
	var data GenerallResponse
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Get(token_uri)
	if err != nil {
		// logger.Errorf(" error %s",err.Error())
		return data, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	err = json.Unmarshal(body, &data)
	if err != nil {
		// logger.Errorf("failed marhsalling response : %s wtih error %s", string(body),err.Error())
		return data, err
	}
	// _=body
	// logger.Info("succussfully connected to host : %s ", uri)
	return data, nil
}
func TestAddToChart(t *testing.T) {

	data, err := GetToken("6287777000057", "bambang@12345")
	t.Log(data)
	assert.Equal(t, nil, err)

}
