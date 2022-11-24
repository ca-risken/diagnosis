package applicationscan

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/datasource-api/proto/diagnosis"
)

type applicationScanAPI interface {
	executeZap(context.Context, string) (int, error)
	handleBasicScan(*diagnosis.ApplicationScanBasicSetting, uint32, uint32, string) (*zapResult, error)
	terminateZap(int) error
	setApiKey(string)
}

type ApplicationScanClient struct {
	config     *zapConfig
	httpClient *http.Client
	targetURL  string
	contextID  string
	logger     logging.Logger
}

func NewApplicationScanClient(port, path, apiKeyName, apiKeyHeader string, l logging.Logger) (applicationScanAPI, error) {
	conf := &zapConfig{
		ZapPort:         port,
		ZapPath:         path,
		ZapApiKeyName:   apiKeyName,
		ZapApiKeyHeader: apiKeyHeader,
	}
	conf.ZapProxy = fmt.Sprintf("http://localhost:%v", conf.ZapPort)
	conf.BaseUrlJson = conf.ZapProxy + "/json/"
	conf.BaseUrlOther = conf.ZapProxy + "/other/"

	httpClient := &http.Client{}
	return &ApplicationScanClient{
		config:     conf,
		httpClient: httpClient,
		logger:     l,
	}, nil
}

func (a *ApplicationScanClient) executeZap(ctx context.Context, apiKeyValue string) (int, error) {
	cmd := exec.Command(a.config.ZapPath, "-daemon", "-port", a.config.ZapPort, "-config", fmt.Sprintf("api.key=%v", apiKeyValue))
	err := cmd.Start()
	if err != nil {
		a.logger.Errorf(ctx, "Failed to execute ZAP. cmd: %v, error: %v", cmd, err)
		return 0, err
	}
	pID := cmd.Process.Pid
	err = a.WaitForStartingZap(ctx)
	if err != nil {
		a.logger.Errorf(ctx, "Failed to execute ZAP. cmd: %v, error: %v", cmd, err)
		return 0, err
	}
	return pID, nil
}

func (a *ApplicationScanClient) terminateZap(pID int) error {
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

func (a *ApplicationScanClient) Request(path string, queryParams map[string]string) (map[string]interface{}, error) {
	body, err := a.request(a.config.BaseUrlJson+path, queryParams)
	if err != nil {
		return nil, err
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (a *ApplicationScanClient) RequestOther(path string, queryParams map[string]string) ([]byte, error) {
	return a.request(a.config.BaseUrlOther+path, queryParams)
}

func (a *ApplicationScanClient) request(path string, queryParams map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if len(queryParams) == 0 {
		queryParams = map[string]string{}
	}
	// Send the API key even if there are no parameters,
	// older ZAP versions might need API key as (query) parameter.
	queryParams[a.config.ZapApiKeyName] = a.config.ZapApiKeyValue

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
	req.Header.Add(a.config.ZapApiKeyHeader, a.config.ZapApiKeyValue)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Errored when sending request to the server: %v", err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (a *ApplicationScanClient) setApiKey(key string) {
	a.config.ZapApiKeyValue = key
}
