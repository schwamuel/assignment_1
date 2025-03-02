package HANDLERS

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type theName struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
}

type final struct {
	Mean int
	Data struct {
		Values []struct {
			Year  int `json:"year"`
			Value int `json:"value"`
		} `json:"populationCounts"`
	} `json:"data"`
}

const LINEBREAK = "\n"

func Population(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Country code is required", http.StatusBadRequest)
		return
	}
	countryCode := parts[4]

	limitParam := r.URL.Query().Get("limit")
	var startYear, endYear int
	var hasLimit bool
	if limitParam != "" {
		// Expect format "startYear-endYear"
		yearRange := strings.Split(limitParam, "-")
		if len(yearRange) == 2 {
			var err error
			startYear, err = strconv.Atoi(yearRange[0])
			if err != nil {
				http.Error(w, "Invalid start year", http.StatusBadRequest)
				return
			}
			endYear, err = strconv.Atoi(yearRange[1])
			if err != nil {
				http.Error(w, "Invalid end year", http.StatusBadRequest)
				return
			}
			hasLimit = true
		} else {
			http.Error(w, "Invalid limit format. Use 'startYear-endYear'", http.StatusBadRequest)
			return
		}
	}

	url1 := fmt.Sprintf("http://129.241.150.113:8080/v3.1/alpha/%s", countryCode)
	resp, err := http.Get(url1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var apiResponse []theName
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	countryName := apiResponse[0].Name.Common

	url2 := "https://countriesnow.space/api/v0.1/countries/population"

	payload := strings.NewReader(fmt.Sprintf(`{"country": "%s"}`, countryName))

	req, err := http.NewRequest("POST", url2, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp1, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp1.Body.Close()

	body, err := io.ReadAll(resp1.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Population Response Body:", string(body))

	var apiResponse2 final
	if err := json.Unmarshal(body, &apiResponse2); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if hasLimit {
		var filteredValues []struct {
			Year  int `json:"year"`
			Value int `json:"value"`
		}

		for _, value := range apiResponse2.Data.Values {
			if value.Year >= startYear && value.Year <= endYear {
				filteredValues = append(filteredValues, value)
			}
		}
		apiResponse2.Data.Values = filteredValues

	}
	sum := 0

	for _, v := range apiResponse2.Data.Values {
		sum += v.Value
	}
	mean := sum / len(apiResponse2.Data.Values)

	apiResponse2.Mean = mean

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiResponse2)
}
