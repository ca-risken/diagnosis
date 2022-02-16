package main

import (
	//	"context"

	"fmt"
	"strconv"
	"time"
)

type zapConfig struct {
	ZapPort         string
	ZapProxy        string
	BaseUrlJson     string
	BaseUrlOther    string
	ZapPath         string
	ZapApiKeyName   string
	ZapApiKeyValue  string
	ZapApiKeyHeader string
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

func (c *applicationScanClient) NewSession(name string) (map[string]interface{}, error) {
	m := map[string]string{
		"name":      name,
		"overwrite": "True",
	}
	return c.Request("core/action/newSession/", m)
}

func (c *applicationScanClient) newContext(contextname string) (map[string]interface{}, error) {
	m := map[string]string{
		"contextName": contextname,
	}
	return c.Request("context/action/newContext/", m)
}
func (c *applicationScanClient) IncludeInContext(contextname string, regex string) (map[string]interface{}, error) {
	m := map[string]string{
		"contextName": contextname,
		"regex":       regex,
	}
	return c.Request("context/action/includeInContext/", m)
}
func (c *applicationScanClient) SetOptionMaxChildren(i int) (map[string]interface{}, error) {
	m := map[string]string{
		"Integer": strconv.Itoa(i),
	}
	return c.Request("spider/action/setOptionMaxChildren/", m)
}

func (c *applicationScanClient) SetOptionMaxDepth(i int) (map[string]interface{}, error) {
	m := map[string]string{
		"Integer": strconv.Itoa(i),
	}
	return c.Request("spider/action/setOptionMaxDepth/", m)
}

func (c *applicationScanClient) SpiderStatus(scanid string) (map[string]interface{}, error) {
	m := map[string]string{
		"scanId": scanid,
	}
	return c.Request("spider/view/status/", m)
}

func (c *applicationScanClient) SpiderScan(url string, maxchildren string, recurse string, contextname string, subtreeonly string) (map[string]interface{}, error) {
	m := map[string]string{
		"url":         url,
		"maxChildren": maxchildren,
		"recurse":     recurse,
		"contextName": contextname,
		"subtreeOnly": subtreeonly,
	}
	return c.Request("spider/action/scan/", m)
}

func (c *applicationScanClient) AscanStatus(scanid string) (map[string]interface{}, error) {
	m := map[string]string{
		"scanId": scanid,
	}
	return c.Request("ascan/view/status/", m)
}

func (c *applicationScanClient) AscanProgress(scanid string) (map[string]interface{}, error) {
	m := map[string]string{
		"scanId": scanid,
	}
	return c.Request("ascan/view/scanProgress/", m)
}

func (c *applicationScanClient) AscanScan(url string, recurse string, inscopeonly string, scanpolicyname string, method string, postdata string, contextid string) (map[string]interface{}, error) {
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

func (c *applicationScanClient) Jsonreport() ([]byte, error) {
	return c.RequestOther("core/other/jsonreport/", nil)
}

func (c *applicationScanClient) WaitForStartingZap() error {
	WaitDuration := 300.0
	now := time.Now()
	for {
		time.Sleep(time.Second * 1)
		_, err := c.RequestOther("", nil)
		if err == nil {
			return nil
		}
		if time.Since(now).Seconds() > WaitDuration {
			return fmt.Errorf("ZAP doesn't start. waitDuration: %v", WaitDuration)
		}
	}
}
