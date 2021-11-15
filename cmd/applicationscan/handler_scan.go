package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ca-risken/diagnosis/proto/diagnosis"
)

func (z *applicationScanClient) handleBasicScan(setting *diagnosis.ApplicationScanBasicSetting, applicationScanID, projectID uint32, name string) (*zapResult, error) {
	contextName := fmt.Sprintf("%v_%v_%v", projectID, applicationScanID, time.Now().Unix())
	z.targetURL = setting.Target
	err := z.HandleBasicSetting(contextName, setting.MaxDepth, setting.MaxChildren)
	if err != nil {
		return nil, err
	}
	time.Sleep(1 * time.Second)
	err = z.HandleSpiderScan(contextName, setting.MaxChildren)
	if err != nil {
		return nil, err
	}
	err = z.HandleActiveScan()
	if err != nil {
		return nil, err
	}
	report, err := z.getJsonReport()
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (c *applicationScanClient) HandleBasicSetting(name string, maxDepth, maxChildren uint32) error {

	// Create Session and Context
	_, err := c.NewSession(name)
	if err != nil {
		appLogger.Errorf("Failed to create session, error: %v", err)
		return err
	}
	retNewContext, err := c.newContext(name)
	if err != nil {
		appLogger.Errorf("Failed to create context, error: %v", err)
		return err
	}
	if retNewContext["contextId"] == nil {
		err = errors.New("ContextID is null")
		appLogger.Errorf("Failed to get context ID, error: %v", err)
		return err
	}
	c.contextID = retNewContext["contextId"].(string)
	_, err = c.IncludeInContext(name, c.targetURL+".*")
	if err != nil {
		appLogger.Errorf("Failed to include context, error: %v", err)
		return err
	}

	// Set Scan Config
	_, err = c.SetOptionMaxChildren(int(maxChildren))
	if err != nil {
		appLogger.Errorf("Failed to set max children option, error: %v", err)
		return err
	}
	_, err = c.SetOptionMaxDepth(int(maxDepth))
	if err != nil {
		appLogger.Errorf("Failed to set max depth option, error: %v", err)
		return err
	}

	return nil
}

func (c *applicationScanClient) HandleSpiderScan(contextName string, maxChildren uint32) error {
	// Exec Spider
	retSpiderScan, err := c.SpiderScan(c.targetURL, fmt.Sprint(maxChildren), "True", contextName, "True")
	if err != nil {
		appLogger.Errorf("Failed to execute spider, error: %v", err)
		return err
	}
	spiderScanID := retSpiderScan["scan"]
	if spiderScanID == nil {
		err = errors.New("SpiderScanID is null")
		appLogger.Errorf("Failed to get spider scan ID, error: %v", err)
		return err
	}
	for {
		retSpiderStatus, err := c.SpiderStatus(spiderScanID.(string))
		if err != nil {
			appLogger.Errorf("Failed to get status spider, error: %v", err)
			return err
		}
		if retSpiderStatus["status"] == nil {
			err = errors.New("SpiderScanStatus is null")
			appLogger.Errorf("Failed to get spider scan status, error: %v", err)
			return err
		}
		spiderStatus, err := strconv.Atoi(retSpiderStatus["status"].(string))
		if err != nil {
			appLogger.Errorf("Failed to convert spider status. error: %v", err)
			break
		}
		if spiderStatus >= 100 {
			break
		}
	}
	return nil

}

func (c *applicationScanClient) HandleActiveScan() error {
	retAscanScan, err := c.AscanScan("", "True", "True", "Default Policy", "", "", c.contextID)
	if err != nil {
		appLogger.Errorf("Failed to execute active scan, error: %v", err)
		return err
	}
	ascanScanID := retAscanScan["scan"]
	if ascanScanID == nil {
		err = errors.New("ActiveScanID is null")
		appLogger.Errorf("Failed to get active scan ID, error: %v", err)
		return err
	}
	for {
		retAscanStatus, err := c.AscanStatus(ascanScanID.(string))
		if err != nil {
			appLogger.Errorf("Failed to get active scan status, error: %v", err)
			return err
		}
		if retAscanStatus["status"] == nil {
			err = errors.New("ActiveScanStatus is null")
			appLogger.Errorf("Failed to get active scan status, error: %v", err)
			return err
		}
		ascanStatus, err := strconv.Atoi(retAscanStatus["status"].(string))
		if err != nil {
			appLogger.Errorf("Failed to convert active scan status. error: %v", err)
			break
		}
		if ascanStatus >= 100 {
			break
		}
	}
	return nil
}

func (c *applicationScanClient) getJsonReport() (*zapResult, error) {
	retJsonReport, err := c.Jsonreport()
	if err != nil {
		appLogger.Errorf("Failed to get json report, error: %v", err)
		return nil, err
	}
	var jsonReport zapResult
	err = json.Unmarshal(retJsonReport, &jsonReport)
	if err != nil {
		appLogger.Errorf("Failed to marshal jsonReport, error: %v", err)
		return nil, err
	}

	return &jsonReport, nil
}
