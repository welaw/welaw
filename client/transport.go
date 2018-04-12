package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func encodeHTTPRequestWithBearer(_ context.Context, req *http.Request, request interface{}, bearer string) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = int64(len(buf.Bytes()))
	req.Body = ioutil.NopCloser(&buf)
	req.Header.Add("authorization", bearer)
	return nil
}

func encodeHTTPRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = int64(len(buf.Bytes()))
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func copyUrl(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
