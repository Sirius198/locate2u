package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionResponse struct {
	CustomerId int `json:"entity"`
}

func trip_to(c *gin.Context) {

	fulfillmentid := c.Param("fulfillmentid")
	fmt.Println("FulfillmentID:", fulfillmentid)

	// Get transaction from transaction id
	trans := GetTransaction(fulfillmentid)
	if trans == nil {
		log.Println("GetTransaction failed:")
		return
	}

	// Get customer data from customer id of transaction
	customer := GetCustomer(trans.CustomerId)
	if customer == nil {
		log.Println("GetCustomer failed")
		return
	}

	// Create a `stop` in Locate2u based on that fulfillment
	stop := createStopInLocate2U(customer)
	if stop == nil {
		log.Println("CreateStopInLocate2U failed")
		return
	}

	// Create a new 'Link' in Locate2u
	link := CreateNewLink(stop)
	if link == nil {
		log.Println("CreateNewLink failed")
		return
	}

	// Add tracking
	success := AddTrackingLink(fulfillmentid, link.Url)
	if success {
		log.Println("Trip success")
	} else {
		log.Println("Trip failed")
	}
}

func GetTransaction(transactionId string) *TransactionResponse {

	url := fmt.Sprintf("%s/api/transaction/tranid/%s", config.ApiConfig.BaseUrl, transactionId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("GetTransaction:", err)
		return nil
	}

	// Add Bearer Token to header
	req.Header.Add("Authorization", "Bearer "+config.ApiConfig.Token)

	myClient := http.Client{Timeout: time.Second * 100}
	res, err := myClient.Do(req)
	if err != nil {
		log.Println("GetTransaction", err)
		return nil
	}

	data, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println("GetTransaction", err)
		return nil
	}

	transData := TransactionResponse{}
	unmarshalErr := json.Unmarshal(data, &transData)
	if err != nil {
		log.Println("GetTransaction", unmarshalErr)
		return nil
	}
	// fmt.Println(transData)

	return &transData
}

func GetCustomer(customerId int) *Customer {

	url := fmt.Sprintf("%s/api/customer/%d", config.ApiConfig.BaseUrl, customerId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("GetCustomer", err)
		return nil
	}

	// Add Bearer Token to header
	req.Header.Add("Authorization", "Bearer "+config.ApiConfig.Token)

	myClient := http.Client{Timeout: time.Second * 100}
	res, err := myClient.Do(req)
	if err != nil {
		log.Println("GetCustomer", err)
		return nil
	}

	data, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println("GetCustomer", err)
		return nil
	}

	cus := Customer{}
	unmarshalErr := json.Unmarshal(data, &cus)
	if unmarshalErr != nil {
		log.Println("GetCustomer", unmarshalErr)
		return nil
	}
	// fmt.Println(cus)

	return &cus
}

func createStopInLocate2U(cm *Customer) *Locate2UStop {

	dt := time.Now()
	runNumber := 1
	if dt.Hour() > 12 {
		runNumber = 2
	}

	if cm.CustentityClxLocate2UId == 0 {
		log.Println("Customer's data is not synchronized with Locate2U", cm.ID)
		return nil
	}

	// It's only for test
	// cm.CustentityClxLocate2UId = 39443
	stopData := Locate2UStop{
		Name:                 cm.CompanyName,
		Address:              "Default Address",
		Notes:                cm.Comments,
		TripDate:             dt.Format("2006-01-02 15:04:05"),
		AssignedTeamMemberID: config.Locate2U.AssignedTeamMemberID,
		CustomerID:           cm.CustentityClxLocate2UId,
		RunNumber:            runNumber,
	}
	// Set the customer's default shipping address as the main address in Locate2uStop
	for i := 0; i < len(cm.Addresses); i++ {
		if cm.Addresses[i].DefaultShipping {
			stopData.Address = cm.Addresses[i].AddrText
			break
		}
	}

	stopJson, _ := json.Marshal(stopData)
	// fmt.Println("stopJson:", string(stopJson))
	// fmt.Println("token:", Locate2UAccessToken)

	req, err := http.NewRequest(http.MethodPost, config.Locate2U.BaseUrl+"/api/v1/stops", bytes.NewBuffer(stopJson))
	if err != nil {
		log.Println("CreateStopInLocate2U", err)
		return nil
	}

	// Add Bearer Token to header
	req.Header.Add("Authorization", "Bearer "+Locate2UAccessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	myClient := http.Client{Timeout: time.Second * 100}
	res, err := myClient.Do(req)
	if err != nil {
		log.Println("CreateStopInLocate2U", err)
		return nil
	}

	data, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println("CreateStopInLocate2U", err)
		return nil
	}

	// fmt.Println("Crated a new stop")
	// fmt.Println(string(data))

	stopRes := Locate2UStop{}
	unmarshalErr := json.Unmarshal(data, &stopRes)
	if unmarshalErr != nil {
		log.Println("CreateStopInLocate2U", unmarshalErr, string(data))
		return nil
	}

	return &stopRes
}

func CreateNewLink(stop *Locate2UStop) *Locate2ULink {

	linkData := Locate2ULink{
		StopID:  stop.StopID,
		Message: config.Locate2U.LinkMessage,
		Type:    "LR",
	}
	linkJson, _ := json.Marshal(linkData)

	req, err := http.NewRequest(http.MethodPost, config.Locate2U.BaseUrl+"/api/v1/links", bytes.NewBuffer(linkJson))
	if err != nil {
		log.Println(err)
		return nil
	}

	// Add Bearer Token to header
	req.Header.Add("Authorization", "Bearer "+Locate2UAccessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// fmt.Println("linkJson:", string(linkJson))
	// fmt.Println("token:", Locate2UAccessToken)
	// fmt.Println("Url: ", req.URL.String())

	myClient := http.Client{Timeout: time.Second * 100}
	res, err := myClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}

	data, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println(err)
		return nil
	}

	// fmt.Println("CreateNewLink():", res.StatusCode)
	// fmt.Println(string(data))

	linkRes := Locate2ULink{}
	unmarshalErr := json.Unmarshal(data, &linkRes)
	if unmarshalErr != nil {
		log.Println(unmarshalErr)
		return nil
	}

	return &linkRes
}

func AddTrackingLink(transId string, trackingLink string) bool {

	url := fmt.Sprintf("%s/ext/add-tracking/%s?delivery=true&complete=true&tracking=https://%s", config.ApiConfig.BaseUrl, transId, trackingLink)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return false
	}
	req.Header.Add("Authorization", "Bearer "+config.ApiConfig.Token)

	myClient := http.Client{Timeout: time.Second * 100}
	res, err := myClient.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}

	data, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println(err)
		return false
	}

	fmt.Println(url)
	fmt.Println("AddTrackingLink(): ", res.StatusCode)
	fmt.Println(string(data))
	return true
}
