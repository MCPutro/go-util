package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-util/constant"
	"io"
	"net/http"
	"net/url"
	"time"
)

type RestClient interface {
	GetRequest(ctx context.Context, url string, headers map[string]string, responseBody any) error
	PostRequest(ctx context.Context, url string, headers map[string]string, requestBody any, responseBody any) error
	PutRequest(ctx context.Context, url string, headers map[string]string, requestBody any, responseBody any) error
	DeleteRequest(ctx context.Context, url string, headers map[string]string, responseBody any) error
	GenericRequest(ctx context.Context, method string, endpoint string, requestBody any, headers map[string]string,
		queryParams map[string]string, responseBody any) error
}

type restClientImpl struct {
	BaseUrl    string
	HttpClient *http.Client
}

func NewRestClient(baseUrl string, timeout time.Duration) RestClient {
	return &restClientImpl{
		BaseUrl:    baseUrl,
		HttpClient: &http.Client{Timeout: timeout},
	}
}

func (c *restClientImpl) GetRequest(ctx context.Context, url string, headers map[string]string, responseBody any) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+url, nil)
	if err != nil {
		return err
	}
	// Header setup
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Make the request
	resp, err := c.HttpClient.Do(req)
	// Check for errors
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Check for non-200 status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("HTTP Error(%d) : %s", resp.StatusCode, resp.Status))
	}

	// Unmarshal the response body into the provided responseBody variable
	if err = json.Unmarshal(respBodyByte, responseBody); err != nil {
		return err
	}

	return nil
}

func (c *restClientImpl) PostRequest(ctx context.Context, url string, headers map[string]string, requestBody any, responseBody any) error {

	jsonBodyByte, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseUrl+url, bytes.NewBuffer(jsonBodyByte))
	if err != nil {
		return err
	}

	req.Header.Set(constant.HeaderContentType, constant.ContentTypeJSON)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Make the request
	resp, err := c.HttpClient.Do(req)
	// Check for errors
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Check for non-200 status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("HTTP Error(%d) : %s", resp.StatusCode, resp.Status))
	}

	// Unmarshal the response body into the provided responseBody variable
	if err = json.Unmarshal(respBodyByte, responseBody); err != nil {
		return err
	}

	return nil
}

func (c *restClientImpl) PutRequest(ctx context.Context, url string, headers map[string]string, requestBody any, responseBody any) error {
	jsonBodyByte, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.BaseUrl+url, bytes.NewBuffer(jsonBodyByte))
	if err != nil {
		return err
	}

	req.Header.Set(constant.HeaderContentType, constant.ContentTypeJSON)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("HTTP Error(%d) : %s", resp.StatusCode, resp.Status))
	}

	if err := json.Unmarshal(respBody, responseBody); err != nil {
		return err
	}

	return nil
}

func (c *restClientImpl) DeleteRequest(ctx context.Context, url string, headers map[string]string, responseBody any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseUrl+url, nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("HTTP Error(%d) : %s", resp.StatusCode, resp.Status))
	}

	if err := json.Unmarshal(respBody, responseBody); err != nil {
		return err
	}

	return nil
}

func (c *restClientImpl) GenericRequest(ctx context.Context, method string, endpoint string, requestBody any,
	headers map[string]string, queryParams map[string]string, responseBody any) error {

	// Validasi method
	switch method {
	case "GET", "POST", "PUT", "DELETE":
	default:
		return fmt.Errorf("method tidak didukung: %s", method)
	}

	// Tambahkan query parameter ke URL
	urlWithParams := addQueryParams(c.BaseUrl+endpoint, queryParams)

	// Buat request body jika ada
	var reqBody io.Reader
	if requestBody != nil {
		jsonBodyByte, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		reqBody = bytes.NewBuffer(jsonBodyByte)
	}

	// Buat request
	req, err := http.NewRequestWithContext(ctx, method, urlWithParams, reqBody)
	if err != nil {
		return err
	}

	// Set default content-type untuk POST/PUT
	if method == "POST" || method == "PUT" {
		req.Header.Set(constant.HeaderContentType, constant.ContentTypeJSON)
	}

	// Set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Kirim request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Baca body
	respBodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Cek status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("HTTP Error(%d) : %s", resp.StatusCode, resp.Status))
	}

	// Unmarshal ke response
	if responseBody != nil {
		if err = json.Unmarshal(respBodyByte, responseBody); err != nil {
			return err
		}
	}

	return nil
}

func addQueryParams(baseURL string, params map[string]string) string {
	u, _ := url.Parse(baseURL)
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
