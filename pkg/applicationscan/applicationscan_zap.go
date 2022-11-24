package applicationscan

import (
	"context"
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

func (a *ApplicationScanClient) NewSession(name string) (map[string]interface{}, error) {
	m := map[string]string{
		"name":      name,
		"overwrite": "True",
	}
	return a.Request("core/action/newSession/", m)
}

func (a *ApplicationScanClient) newContext(contextname string) (map[string]interface{}, error) {
	m := map[string]string{
		"contextName": contextname,
	}
	return a.Request("context/action/newContext/", m)
}
func (a *ApplicationScanClient) IncludeInContext(contextname string, regex string) (map[string]interface{}, error) {
	m := map[string]string{
		"contextName": contextname,
		"regex":       regex,
	}
	return a.Request("context/action/includeInContext/", m)
}
func (a *ApplicationScanClient) SetOptionMaxChildren(i int) (map[string]interface{}, error) {
	m := map[string]string{
		"Integer": strconv.Itoa(i),
	}
	return a.Request("spider/action/setOptionMaxChildren/", m)
}

func (a *ApplicationScanClient) SetOptionMaxDepth(i int) (map[string]interface{}, error) {
	m := map[string]string{
		"Integer": strconv.Itoa(i),
	}
	return a.Request("spider/action/setOptionMaxDepth/", m)
}

func (a *ApplicationScanClient) SpiderStatus(scanid string) (map[string]interface{}, error) {
	m := map[string]string{
		"scanId": scanid,
	}
	return a.Request("spider/view/status/", m)
}

func (a *ApplicationScanClient) SpiderScan(url string, maxchildren string, recurse string, contextname string, subtreeonly string) (map[string]interface{}, error) {
	m := map[string]string{
		"url":         url,
		"maxChildren": maxchildren,
		"recurse":     recurse,
		"contextName": contextname,
		"subtreeOnly": subtreeonly,
	}
	return a.Request("spider/action/scan/", m)
}

func (a *ApplicationScanClient) AscanStatus(scanid string) (map[string]interface{}, error) {
	m := map[string]string{
		"scanId": scanid,
	}
	return a.Request("ascan/view/status/", m)
}

func (a *ApplicationScanClient) AscanProgress(scanid string) (map[string]interface{}, error) {
	m := map[string]string{
		"scanId": scanid,
	}
	return a.Request("ascan/view/scanProgress/", m)
}

func (a *ApplicationScanClient) AscanScan(url string, recurse string, inscopeonly string, scanpolicyname string, method string, postdata string, contextid string) (map[string]interface{}, error) {
	m := map[string]string{
		"url":            url,
		"recurse":        recurse,
		"inScopeOnly":    inscopeonly,
		"scanPolicyName": scanpolicyname,
		"method":         method,
		"postData":       postdata,
		"contextId":      contextid,
	}
	return a.Request("ascan/action/scan/", m)
}

func (a *ApplicationScanClient) Jsonreport() ([]byte, error) {
	return a.RequestOther("core/other/jsonreport/", nil)
}

func (a *ApplicationScanClient) WaitForStartingZap(ctx context.Context) error {
	WaitDuration := 300.0
	now := time.Now()
	count := 1
	for {
		time.Sleep(time.Second * 20)
		_, err := a.RequestOther("", nil)
		if err == nil {
			return nil
		}
		a.logger.Infof(ctx, "Waiting to start ZAP. wait time: %d [sec], err: %+v", count*20, err)
		if time.Since(now).Seconds() > WaitDuration {
			return fmt.Errorf("ZAP doesn't start. waitDuration: %v", WaitDuration)
		}
		count++
	}
}
