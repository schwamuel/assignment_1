package HANDLERS

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Stat struct {
	Info       int
	Population int
	Version    string
	Uptime     int
}

var tida = time.Now()

const version1 = "v1"

func Status(w http.ResponseWriter, r *http.Request) {

	url1 := "http://localhost:8080/countryinfo/v1/info/no"
	url2 := "http://localhost:8080/countryinfo/v1/population/no"
	resp, err := http.Get(url1)

	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("ERROR WITH INFO HANDLER:", resp.StatusCode)
	} else {
		fmt.Println("Request succeeded with status code:", resp.StatusCode)
	}

	resp1, err1 := http.Get(url2)
	if err1 != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp1.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("ERROR WITH POPULATION HANDLER", resp.StatusCode)
	} else {
		fmt.Println("Request succeeded with status code:", resp.StatusCode)
	}
	tidnaa := time.Since(tida)

	statusResult := Stat{
		Info:       resp.StatusCode,
		Population: resp1.StatusCode,
		Version:    version1,
		Uptime:     int(tidnaa.Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(statusResult)

}
