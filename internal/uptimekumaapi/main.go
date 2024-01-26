package uptimekumaapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type UptimeKumaAPI struct {
	Host     string
	User     string
	Password string
	Token    string
}

func NewUptimeKumaAPI(host string, user string, password string) (*UptimeKumaAPI, error) {
	endpoint := fmt.Sprintf("%s/login/access-token/", host)

	data := url.Values{}
	data.Set("username", user)
	data.Set("password", password)

	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Println("Error while getting access token")
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		log.Println(resp.Body)

		log.Println("Error while getting access token")
		return nil, fmt.Errorf("Error while getting access token")
	}

	var body AccessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&body)

	return &UptimeKumaAPI{
		Host:     host,
		User:     user,
		Password: password,
		Token:    body.AccessToken,
	}, nil
}

func (api *UptimeKumaAPI) getRequest(url string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.Token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error while getting tags")
		return nil, err
	}

	return resp, nil
}

func (api *UptimeKumaAPI) postRequest(url string, body []byte) (*http.Response, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.Token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error while processing request to %s\n", url)
		return nil, err
	}

	return resp, nil
}

func (api *UptimeKumaAPI) deleteRequest(url string, body []byte) (*http.Response, error) {
	req, _ := http.NewRequest("DELETE", url, bytes.NewReader(body))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.Token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error while processing request to %s\n", url)
		return nil, err
	}

	return resp, nil
}

func (api *UptimeKumaAPI) patchRequest(url string, body []byte) (*http.Response, error) {
	req, _ := http.NewRequest("PATCH", url, bytes.NewReader(body))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.Token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error while processing request to %s\n", url)
		return nil, err
	}

	return resp, nil
}
