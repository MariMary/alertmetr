package client

import (
	"bytes"
	"io"
	"net/http"
)

type HttpClient struct {
	Client  *http.Client
	Address string
}

func NewHttpClient(Address string) *HttpClient {
	return &HttpClient{
		Address: Address,
		Client:  &http.Client{},
	}
}

func (hc *HttpClient) CallApi(ApiName string) error {

	url := hc.Address + ApiName
	var body []byte
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := hc.Client.Do(request)
	if err != nil {
		return err
	}
	_, err = io.Copy(io.Discard, response.Body)
	response.Body.Close()
	if err != nil {
		return err
	}
	return nil
}
