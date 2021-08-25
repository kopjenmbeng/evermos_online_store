package test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
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

func GetToken(phone_number string, password string) (GenerallResponse,int, error) {
	uri:=fmt.Sprintf("http://localhost:8080/v1/authentication/get_token?phone_number=%s&password=%s",phone_number,password)
	var data GenerallResponse
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Get(uri)
	if err != nil {
		// logger.Errorf(" error %s",err.Error())
		return data,resp.StatusCode, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	err = json.Unmarshal(body, &data)
	if err != nil {
		// logger.Errorf("failed marhsalling response : %s wtih error %s", string(body),err.Error())
		return data,resp.StatusCode, err
	}
	// _=body
	// logger.Info("succussfully connected to host : %s ", uri)
	return data,resp.StatusCode, nil
}

type AddChartReq struct{
	ProductId string `json:"product_id"`
	Qty int `json:"qty"`
}
func AddChart(product_id string,qty int,token string)int{
	uri:="http://localhost:8080/v1/chart/add"
	var data GenerallResponse
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	
	client := &http.Client{Transport: tr}
	param:=AddChartReq{ProductId: product_id,Qty: qty}
	reqBody,err:=json.Marshal(param)

	key:=fmt.Sprintf("bearer %s",token)
	request, err := http.NewRequest("POST", uri,bytes.NewBuffer(reqBody))
	request.Header.Set("Authorization", key)
    request.Header.Set("Content-Type", "application/json")
	
	resp, err := client.Do(request)
	if err != nil {
		// logger.Errorf(" error %s",err.Error())
		return resp.StatusCode
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	err = json.Unmarshal(body, &data)
	if err != nil {
		// logger.Errorf("failed marhsalling response : %s wtih error %s", string(body),err.Error())
		return resp.StatusCode
	}
	return resp.StatusCode
}

type UpdateChartReq struct{
	ChartId string `json:"chart_id"`
	Qty int `json:"qty"`
}

func UpdateChart(chart_id string,qty int,token string)int{
	uri:="http://localhost:8080/v1/chart/update"
	var data GenerallResponse
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	
	client := &http.Client{Transport: tr}
	param:=UpdateChartReq{ChartId: chart_id,Qty: qty}
	reqBody,err:=json.Marshal(param)

	key:=fmt.Sprintf("bearer %s",token)
	request, err := http.NewRequest("PUT", uri,bytes.NewBuffer(reqBody))
	request.Header.Set("Authorization", key)
    request.Header.Set("Content-Type", "application/json")
	
	resp, err := client.Do(request)
	if err != nil {
		// logger.Errorf(" error %s",err.Error())
		return resp.StatusCode
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	err = json.Unmarshal(body, &data)
	if err != nil {
		// logger.Errorf("failed marhsalling response : %s wtih error %s", string(body),err.Error())
		return resp.StatusCode
	}
	return resp.StatusCode
}
type AddOrderRequest struct{
	ChartIds []string `json:"charts"`
}
func AddOrder(chart_ids []string,token string)int{
	uri:="http://localhost:8080/v1/order/add"
	var data GenerallResponse
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	
	client := &http.Client{Transport: tr}
	param:=AddOrderRequest{ChartIds: chart_ids}
	reqBody,err:=json.Marshal(param)

	key:=fmt.Sprintf("bearer %s",token)
	request, err := http.NewRequest("POST", uri,bytes.NewBuffer(reqBody))
	request.Header.Set("Authorization", key)
    request.Header.Set("Content-Type", "application/json")
	
	resp, err := client.Do(request)
	if err != nil {
		// logger.Errorf(" error %s",err.Error())
		return resp.StatusCode
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	err = json.Unmarshal(body, &data)
	if err != nil {
		// logger.Errorf("failed marhsalling response : %s wtih error %s", string(body),err.Error())
		return resp.StatusCode
	}
	return resp.StatusCode
}
func TestAddToChart(t *testing.T) {

	data,code, err := GetToken("6287777000057", "bambang@12345")
	if code==200{
		code=AddChart("9ea602af-48c1-4570-b81d-48eb3fa740e1",40,data.Content.Token)
		assert.Equal(t,201, code)
	}
	// t.Log(data)
	assert.Equal(t, nil, err)

}

func TestUpdateChart(t *testing.T) {

	data,code, err := GetToken("6287777000057", "bambang@12345")
	if code==200{
		code=UpdateChart("95d69135-b125-40a8-837f-41d7b0615469",3,data.Content.Token)
		assert.Equal(t,200, code)
	}
	// t.Log(data)
	assert.Equal(t, nil, err)

}

func TestAddOrder(t *testing.T) {
	// ids:=
	data,code, err := GetToken("6287777000057", "bambang@12345")
	if code==200{
		code=AddOrder([]string{"5e8a0e30-22a8-4e17-8f01-1c95e24629af","95d69135-b125-40a8-837f-41d7b0615469"},data.Content.Token)
		assert.Equal(t,201, code)
	}
	// t.Log(data)
	assert.Equal(t, nil, err)

}
