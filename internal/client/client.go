package client

import (
	"bytes"
	"io"
	"net/http"
)

type HTTPClient struct {
	Client  *http.Client
	Address string
}

func NewHTTPClient(Address string) *HTTPClient {
	return &HTTPClient{
		Address: Address,
		Client:  &http.Client{},
	}
}

func (hc *HTTPClient) CallAPI(APIName string, Buffer *bytes.Buffer, ContentType string) error {

	url := hc.Address + APIName
	//var body []byte
	request, err := http.NewRequest(http.MethodPost, url, Buffer)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", ContentType)
	request.Header.Set("Content-Encoding", "gzip")
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
