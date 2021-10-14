package main

import (
	//	"context"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/diagnosis/pkg/message"
	"github.com/gassara-kys/envconfig"
)

const (
	ZAP_API_KEY_PARAM  = "apikey"
	ZAP_API_KEY_HEADER = "X-ZAP-API-Key"
)

type zapConfig struct {
	ZapPort         string `split_words:"true" default:"8080"`
	ZapProxy        string
	BaseUrlJson     string
	BaseUrlOther    string
	ZapPath         string `split_words:"true" default:"/zap/zap.sh"`
	ZapApiKeyName   string `split_words:"true" default:"apikey"`
	ZapApiKeyValue  string
	ZapApiKeyHeader string `split_words:"true" default:"X-ZAP-API-Key"`
}

type zapResult struct {
	Version   string          `json:"@version"`
	Generated string          `json:"@generated"`
	Site      []zapResultSite `json:"site"`
}

type zapResultSite struct {
	Host   string           `json:"@host"`
	Port   string           `json:"@port"`
	Alerts []zapResultAlert `json:"alerts"`
}

type zapResultAlert struct {
	Alert     string                   `json:"alert"`
	RiskDesc  string                   `json:"riskdesc"`
	Instances []map[string]interface{} `json:"instances"`
	Name      string                   `json:"name"`
	RiskCode  string                   `json:"riskcode"`
	Solution  string                   `json:"solution"`
}

type zapAPI interface {
	getResult(*message.PortscanQueueMessage, bool) ([]*finding.FindingForUpsert, error)
	scan() (zapResult, error)
}

type zapClient struct {
	config     *zapConfig
	httpClient *http.Client
	targetURL  string
	contextID  string
}

func newZapClient(apiKeyValue string) (*zapClient, error) {
	var conf zapConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		return nil, err
	}
	conf.ZapProxy = fmt.Sprintf("http://localhost:%v", conf.ZapPort)
	conf.BaseUrlJson = conf.ZapProxy + "/json/"
	conf.BaseUrlOther = conf.ZapProxy + "/other/"

	conf.ZapApiKeyValue = apiKeyValue

	httpClient := &http.Client{}
	cli := zapClient{
		config:     &conf,
		httpClient: httpClient,
	}
	return &cli, nil
}

func (c *zapClient) executeZap(apiKeyValue string) int {
	cmd := exec.Command(c.config.ZapPath, "-daemon", "-port", c.config.ZapPort, "-config", fmt.Sprintf("api.key=%v", apiKeyValue))
	cmd.Start()
	pID := cmd.Process.Pid
	c.WaitForStartingZap()
	return pID
}

func (c *zapClient) terminateZap(pID int) error {

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

func (c *zapClient) Request(path string, queryParams map[string]string) (map[string]interface{}, error) {
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

func (c *zapClient) RequestOther(path string, queryParams map[string]string) ([]byte, error) {
	return c.request(c.config.BaseUrlOther+path, queryParams)
}

func (c *zapClient) request(path string, queryParams map[string]string) ([]byte, error) {
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

func (c *zapClient) NewSession(name string) (map[string]interface{}, error) {
	m := map[string]string{
		"name":      name,
		"overwrite": "True",
	}
	return c.Request("core/action/newSession/", m)
}

func (c *zapClient) newContext(contextname string) (map[string]interface{}, error) {
	m := map[string]string{
		"contextName": contextname,
	}
	return c.Request("context/action/newContext/", m)
}
func (c *zapClient) IncludeInContext(contextname string, regex string) (map[string]interface{}, error) {
	m := map[string]string{
		"contextName": contextname,
		"regex":       regex,
	}
	return c.Request("context/action/includeInContext/", m)
}
func (c *zapClient) SetOptionMaxChildren(i int) (map[string]interface{}, error) {
	m := map[string]string{
		"Integer": strconv.Itoa(i),
	}
	return c.Request("spider/action/setOptionMaxChildren/", m)
}

func (c *zapClient) SetOptionMaxDepth(i int) (map[string]interface{}, error) {
	m := map[string]string{
		"Integer": strconv.Itoa(i),
	}
	return c.Request("spider/action/setOptionMaxDepth/", m)
}

func (c *zapClient) SpiderStatus(scanid string) (map[string]interface{}, error) {
	m := map[string]string{
		"scanId": scanid,
	}
	return c.Request("spider/view/status/", m)
}

func (c *zapClient) SpiderScan(url string, maxchildren string, recurse string, contextname string, subtreeonly string) (map[string]interface{}, error) {
	m := map[string]string{
		"url":         url,
		"maxChildren": maxchildren,
		"recurse":     recurse,
		"contextName": contextname,
		"subtreeOnly": subtreeonly,
	}
	return c.Request("spider/action/scan/", m)
}

func (c *zapClient) AscanStatus(scanid string) (map[string]interface{}, error) {
	m := map[string]string{
		"scanId": scanid,
	}
	return c.Request("ascan/view/status/", m)
}

func (c *zapClient) AscanProgress(scanid string) (map[string]interface{}, error) {
	m := map[string]string{
		"scanId": scanid,
	}
	return c.Request("ascan/view/scanProgress/", m)
}

func (c *zapClient) AscanScan(url string, recurse string, inscopeonly string, scanpolicyname string, method string, postdata string, contextid string) (map[string]interface{}, error) {
	m := map[string]string{
		"url":            url,
		"recurse":        recurse,
		"inScopeOnly":    inscopeonly,
		"scanPolicyName": scanpolicyname,
		"method":         method,
		"postData":       postdata,
		"contextId":      contextid,
	}
	return c.Request("ascan/action/scan/", m)
}

func (c *zapClient) Jsonreport() ([]byte, error) {
	return c.RequestOther("core/other/jsonreport/", nil)
}

func (c *zapClient) WaitForStartingZap() error {
	for {
		time.Sleep(time.Second * 1)
		_, err := c.RequestOther("", nil)
		if err == nil {
			break
		}
	}
	return nil
}
