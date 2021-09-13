package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func (c *zapClient) HandleBasicSetting(name string, maxDepth, maxChildren uint32) error {

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

func (c *zapClient) HandleSpiderScan(contextName string, maxChildren uint32) error {
	// Exec Spider
	retSpiderScan, err := c.SpiderScan(c.targetURL, fmt.Sprint(maxChildren), "True", contextName, "True")
	if err != nil {
		appLogger.Errorf("Failed to execute spider, error: %v", err)
		return err
	}
	spiderScanID := retSpiderScan["scan"]
	for {
		retSpiderStatus, err := c.SpiderStatus(spiderScanID.(string))
		if err != nil {
			appLogger.Errorf("Failed to get status spider, error: %v", err)
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

func (c *zapClient) HandleActiveScan() error {
	// Exec Scan

	retAscanScan, err := c.AscanScan(c.targetURL, "True", "True", "Default Policy", "", "", c.contextID)
	if err != nil {
		appLogger.Errorf("Failed to execute active scan, error: %v", err)
		return err
	}
	AscanScanID := retAscanScan["scan"]
	for {
		retAscanStatus, err := c.AscanStatus(AscanScanID.(string))
		if err != nil {
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

func (c *zapClient) getJsonReport() (*zapResult, error) {
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
