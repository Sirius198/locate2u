package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Locate2UAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func getAccessToken(Locate2U Locate2UConfig) string {

	reqBody := url.Values{
		"client_id":     {Locate2U.ClientId},
		"client_secret": {Locate2U.ClientSecret},
		"grant_type":    {Locate2U.GrantType},
		"scope":         {Locate2U.Scope},
	}

	res, err := http.PostForm("https://id.locate2u.com/connect/token", reqBody)
	if err != nil {
		log.Println("GetAccessToken; Connect to Locate2U failed: ", err)

		// Retry
		// log.Fatal(err)
		return getAccessToken(Locate2U)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println("GetAccessToken: ", readErr)

		// Retry
		// log.Fatal(err)
		return getAccessToken(Locate2U)
	}

	accessTokenResp := Locate2UAccessTokenResponse{}
	if err := json.Unmarshal(body, &accessTokenResp); err != nil {
		panic(err)
	}

	fmt.Println("Locate2U Access Token Received!")
	return accessTokenResp.AccessToken
}
