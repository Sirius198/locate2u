package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func doSync(c *gin.Context) {
	syncCustomers()
}

type UpdateCustentityLocate2UIDRequest struct {
	CustentityClxLocate2UId int `json:"custentity_clx_locate2uid"`
}

func syncCustomers() {

	fmt.Println("Starting sync...")

	myClient := http.Client{Timeout: time.Second * 30}
	pageNumber := 1
	const perPage string = "10"

	for {
		req, _ := http.NewRequest(http.MethodGet, config.ApiConfig.BaseUrl+"/api/customer", nil)
		q := req.URL.Query()
		q.Add("page", strconv.Itoa(pageNumber))
		q.Add("per_page", perPage)
		req.URL.RawQuery = q.Encode()
		req.Header.Add("Authorization", "Bearer "+config.ApiConfig.Token)

		res, err := myClient.Do(req)
		if err != nil {
			log.Println(err)
			pageNumber++
			continue
		}

		data, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Println(err)
			pageNumber++
			continue
		}

		var customers []Customer
		unmarshalErr := json.Unmarshal(data, &customers)
		if unmarshalErr != nil {
			log.Println(unmarshalErr)
			pageNumber++
			continue
		}

		var wg sync.WaitGroup
		for i := 0; i < len(customers); i++ {
			// fmt.Println(customers[i].ID, customers[i].CustentityClxLocate2UId)
			wg.Add(1)

			i := i
			go func() {
				defer wg.Done()
				CreateCustomerInLocate2U(customers[i])
			}()
		}

		wg.Wait()

		// If reaches the end of customer page from API, break the loop
		if len(customers) == 0 {
			break
		}

		pageNumber++
		log.Println("pageNumber: ", pageNumber)
	}

	fmt.Println("Sync finished!")
}

func CreateCustomerInLocate2U(c Customer) {

	// Prepare form data
	ltuData := Locate2UCustomer{
		Name:    c.CompanyName,
		Company: c.CompanyName,
		Address: "",
		Email:   c.Email,
		Phone:   c.Phone,
		Notes:   c.Comments,
	}

	// Save the customer's default shipping address as the main address in Locate2u
	for i := 0; i < len(c.Addresses); i++ {
		if c.Addresses[i].DefaultShipping {
			ltuData.Address = c.Addresses[i].AddrText
			break
		}
	}

	// Convert struct data into json
	customerJSON, _ := json.Marshal(ltuData)

	//
	isNewCustomer := false
	if c.CustentityClxLocate2UId == 0 {
		isNewCustomer = true
		fmt.Println("Yes new customer")
	} else {
		fmt.Println("already exist: ", c.CustentityClxLocate2UId, c.ID)
	}

	var req *http.Request

	// Create a customer if it is new or update Locate2U customer fields
	if isNewCustomer {
		req, _ = http.NewRequest(http.MethodPost, config.Locate2U.BaseUrl+"/api/v1/customers", bytes.NewBuffer(customerJSON))
	} else {
		req, _ = http.NewRequest(http.MethodPut, config.Locate2U.BaseUrl+"/api/v1/customers/"+strconv.Itoa(c.CustentityClxLocate2UId), bytes.NewBuffer(customerJSON))
	}
	req.Header.Add("Authorization", "Bearer "+Locate2UAccessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	myClient := http.Client{Timeout: time.Second * 100}
	res, err := myClient.Do(req)
	if err != nil {
		log.Println("CreateCustomerInLocate2U, Send Request Error: ", err)
		return
	}

	if isNewCustomer {
		data, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Println(err)
			return
		}

		l2uResponse := Locate2UCustomer{}
		unmarshalErr := json.Unmarshal(data, &l2uResponse)
		if unmarshalErr != nil {
			log.Println(unmarshalErr)
			return
		}

		// if new customer, update "custentity_clx_locate2uid" value

		log.Println("Updating locate2id...")

		reqBody := UpdateCustentityLocate2UIDRequest{
			CustentityClxLocate2UId: l2uResponse.CustomerID,
		}
		jsonData, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, config.ApiConfig.BaseUrl+"/netsuite/customer/"+strconv.Itoa(c.ID), bytes.NewBuffer(jsonData))
		req.Header.Add("Authorization", "Bearer "+config.ApiConfig.Token)
		res, err := myClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}

		data1, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Println(err)
		} else {
			fmt.Println("Updated api customer id", res.Status, len(data1), string(data1))
		}
	} else {
		fmt.Println("Updated")
	}
}
