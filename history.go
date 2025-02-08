package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var url = "https://history.muffinlabs.com/date"

type HistoryResponse struct {
	Date string `json:"date"`
	Data struct {
		Events []HistoricalEntry `json:"Events"`
		Births []HistoricalEntry `json:"Births"`
		Deaths []HistoricalEntry `json:"Deaths"`
	} `json:"data"`
}

type HistoricalEntry struct {
	Year string `json:"year"`
	Text string `json:"text"`
}

func getAllEventsFrom(date time.Time) (*HistoryResponse, error) {
	urlWithDate := fmt.Sprintf("%s/%s", url, date.Format("01/02"))
	resp, err := http.Get(urlWithDate)
	if err != nil {
		log.Println("Error while fetching data", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading response:", err)
		return nil, err
	}

	var history *HistoryResponse
	err = json.Unmarshal(body, &history)
	if err != nil {
		log.Println("Error while parsing JSON:", err)
		return nil, err
	}

	return history, nil
}
