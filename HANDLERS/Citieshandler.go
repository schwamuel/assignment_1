package HANDLERS

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CityResponse struct {
	Cities []string `json:"data"`
}

func getCities(countryCode string) ([]string, error) {

	payload := strings.NewReader(fmt.Sprintf(`{"iso2": "%s"}`, countryCode))

	req, err := http.NewRequest("POST", "https://countriesnow.space/api/v0.1/countries/cities", payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("HTTP Status Code:", resp.StatusCode)
	fmt.Println("HTTP Headers:", resp.Header)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Response Body:", string(body))

	var cityResponse struct {
		Cities []string `json:"data"`
	}

	err = json.Unmarshal(body, &cityResponse)
	if err != nil {
		return nil, err
	}

	if len(cityResponse.Cities) == 0 {
		fmt.Println("Warning: No cities found for country:", countryCode)
	}

	return cityResponse.Cities, nil
}
