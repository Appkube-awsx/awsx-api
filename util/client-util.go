package util

import (
	"fmt"
	"io"
	"net/http"
)

func HandleHttpRequest(httpMethod string, url string, apiKey string, payload io.Reader) ([]byte, int, error) {
	client := &http.Client{}
	httpRequest, err := http.NewRequest(httpMethod, url, payload)
	if err != nil {
		Error("Cannot create a http requests. Requested method: "+httpMethod, err)
		return nil, http.StatusBadRequest, err
	}
	if apiKey != "" {
		httpRequest.Header.Add("api-key", apiKey)
	}
	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Set("Accept", "application/json")
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		Error("Http request failed.", err)
		return nil, http.StatusBadRequest, err
	}
	defer func() {
		_ = httpResponse.Body.Close()
	}()
	data, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		Error("Reading http response failed.", err)
		return nil, http.StatusInternalServerError, err
	}
	if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusCreated {
		fmt.Errorf("Unable to get data from URL: %s due to status code: %d", url, httpResponse.StatusCode)
		return nil, httpResponse.StatusCode, fmt.Errorf("unable to fetch data from url: %s", url)
	}
	return data, httpResponse.StatusCode, nil
}
