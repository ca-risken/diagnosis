package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/ca-risken/datasource-api/proto/diagnosis"
)

type applicationScanAPI interface {
	executeZap(context.Context, string) (int, error)
	handleBasicScan(context.Context, *diagnosis.ApplicationScanBasicSetting, uint32, uint32, string) (*zapResult, error)
	terminateZap(int) error
}

type applicationScanClient struct {
	config     *zapConfig
	httpClient *http.Client
	targetURL  string
	contextID  string
}

func newApplicationScanClient(port, path, apiKeyName, apiKeyValue, apiKeyHeader string) (applicationScanAPI, error) {
	conf := &zapConfig{
		ZapPort:         port,
		ZapPath:         path,
		ZapApiKeyName:   apiKeyName,
		ZapApiKeyValue:  apiKeyValue,
		ZapApiKeyHeader: apiKeyHeader,
	}
	conf.ZapProxy = fmt.Sprintf("http://localhost:%v", conf.ZapPort)
	conf.BaseUrlJson = conf.ZapProxy + "/json/"
	conf.BaseUrlOther = conf.ZapProxy + "/other/"

	httpClient := &http.Client{}
	cli := applicationScanClient{
		config:     conf,
		httpClient: httpClient,
	}
	return &cli, nil
}

func (c *applicationScanClient) executeZap(ctx context.Context, apiKeyValue string) (int, error) {
	cmd := exec.Command(c.config.ZapPath, "-daemon", "-port", c.config.ZapPort, "-config", fmt.Sprintf("api.key=%v", apiKeyValue))
	err := cmd.Start()
	if err != nil {
		appLogger.Errorf(ctx, "Failed to execute ZAP. cmd: %v, error: %v", cmd, err)
		return 0, err
	}
	pID := cmd.Process.Pid
	err = c.WaitForStartingZap(ctx)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to execute ZAP. cmd: %v, error: %v", cmd, err)
		return 0, err
	}
	return pID, nil
}

func (c *applicationScanClient) terminateZap(pID int) error {

	process, err := os.FindProcess(pID)
	if err != nil {
		return err
	}
	err = process.Kill()
	if err != nil {
		return err
	}
	return nil
}

func (c *applicationScanClient) Request(path string, queryParams map[string]string) (map[string]interface{}, error) {
	body, err := c.request(c.config.BaseUrlJson+path, queryParams)
	if err != nil {
		return nil, err
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (c *applicationScanClient) RequestOther(path string, queryParams map[string]string) ([]byte, error) {
	return c.request(c.config.BaseUrlOther+path, queryParams)
}

func (c *applicationScanClient) request(path string, queryParams map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if len(queryParams) == 0 {
		queryParams = map[string]string{}
	}
	// Send the API key even if there are no parameters,
	// older ZAP versions might need API key as (query) parameter.
	queryParams[c.config.ZapApiKeyName] = c.config.ZapApiKeyValue

	// add url query parameter
	query := req.URL.Query()
	for k, v := range queryParams {
		if v == "" {
			continue
		}
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	// add HTTP Accept header
	req.Header.Add("Accept", "application/json")
	// add API Key header
	req.Header.Add(c.config.ZapApiKeyHeader, c.config.ZapApiKeyValue)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Errored when sending request to the server: %v", err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
