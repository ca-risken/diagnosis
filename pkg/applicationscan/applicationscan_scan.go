package applicationscan

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ca-risken/datasource-api/proto/diagnosis"
)

func (a *ApplicationScanClient) handleBasicScan(setting *diagnosis.ApplicationScanBasicSetting, applicationScanID, projectID uint32, name string) (*zapResult, error) {
	contextName := fmt.Sprintf("%v_%v_%v", projectID, applicationScanID, time.Now().Unix())
	a.targetURL = setting.Target
	err := a.HandleBasicSetting(contextName, setting.MaxDepth, setting.MaxChildren)
	if err != nil {
		return nil, err
	}
	time.Sleep(1 * time.Second)
	err = a.HandleSpiderScan(contextName, setting.MaxChildren)
	if err != nil {
		return nil, err
	}
	err = a.HandleActiveScan()
	if err != nil {
		return nil, err
	}
	report, err := a.getJsonReport()
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (a *ApplicationScanClient) HandleBasicSetting(name string, maxDepth, maxChildren uint32) error {
	// Create Session and Context
	_, err := a.NewSession(name)
	if err != nil {
		return fmt.Errorf("create session error: err=%w", err)
	}
	retNewContext, err := a.newContext(name)
	if err != nil {
		return fmt.Errorf("create context error: err=%w", err)
	}
	if retNewContext["contextId"] == nil {
		return errors.New("ContextID is null")
	}
	a.contextID = retNewContext["contextId"].(string)
	_, err = a.IncludeInContext(name, a.targetURL+".*")
	if err != nil {
		return fmt.Errorf("include in context error: err=%w", err)
	}

	// Set Scan Config
	_, err = a.SetOptionMaxChildren(int(maxChildren))
	if err != nil {
		return fmt.Errorf("set max children option error: err=%w", err)
	}
	_, err = a.SetOptionMaxDepth(int(maxDepth))
	if err != nil {
		return fmt.Errorf("set max depth option error: err=%w", err)
	}
	return nil
}

func (a *ApplicationScanClient) HandleSpiderScan(contextName string, maxChildren uint32) error {
	// Exec Spider
	retSpiderScan, err := a.SpiderScan(a.targetURL, fmt.Sprint(maxChildren), "True", contextName, "True")
	if err != nil {
		return fmt.Errorf("execute spider error: err=%w", err)
	}
	spiderScanID := retSpiderScan["scan"]
	if spiderScanID == nil {
		return errors.New("get spider scan ID error(ID=null)")
	}
	for {
		retSpiderStatus, err := a.SpiderStatus(spiderScanID.(string))
		if err != nil {
			return fmt.Errorf("get status spider error: %w", err)
		}
		status := retSpiderStatus["status"]
		if status == nil {
			return errors.New("SpiderScanStatus is null")
		}
		spiderStatus, err := strconv.Atoi(status.(string))
		if err != nil {
			return fmt.Errorf("invalid spider status: status=%v, err=%w", status, err)
		}
		if spiderStatus >= 100 {
			break
		}
	}
	return nil

}

func (a *ApplicationScanClient) HandleActiveScan() error {
	retAscanScan, err := a.AscanScan("", "True", "True", "Default Policy", "", "", a.contextID)
	if err != nil {
		return fmt.Errorf("execute active scan error: err=%w", err)
	}
	ascanScanID := retAscanScan["scan"]
	if ascanScanID == nil {
		return errors.New("get active scan ID error(ActiveScanID is null)")
	}
	for {
		retAscanStatus, err := a.AscanStatus(ascanScanID.(string))
		if err != nil {
			return fmt.Errorf("get active scan status error: err=%w", err)
		}
		status := retAscanStatus["status"]
		if status == nil {
			return errors.New("get active scan status error(ActiveScanStatus is null)")
		}
		ascanStatus, err := strconv.Atoi(status.(string))
		if err != nil {
			return fmt.Errorf("invalid active scan status: status=%v, err=%w", status, err)
		}
		if ascanStatus >= 100 {
			break
		}
	}
	return nil
}

func (a *ApplicationScanClient) getJsonReport() (*zapResult, error) {
	retJsonReport, err := a.Jsonreport()
	if err != nil {
		return nil, fmt.Errorf("failed to get json report: err=%w", err)
	}
	var jsonReport zapResult
	err = json.Unmarshal(retJsonReport, &jsonReport)
	if err != nil {
		return nil, fmt.Errorf("marshal jsonReport error: err=%w", err)
	}
	return &jsonReport, nil
}
