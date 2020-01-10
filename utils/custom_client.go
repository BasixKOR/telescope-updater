package utils

import (
	"fmt"
	"net/http"
)

type BearerRoundTripper string

func (key BearerRoundTripper) RoundTrip(r *http.Request) (res *http.Response, err error) {
	r.Header.Add("Authorization", fmt.Sprint("Bearer ", key))
	res, err = http.DefaultTransport.RoundTrip(r)
	return
}

func NewBearerClient(key string) *http.Client {
	return &http.Client{
		Transport: BearerRoundTripper(key),
	}
}
