package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ca-risken/datasource-api/proto/diagnosis"
)

func (z *applicationScanClient) handleBasicScan(ctx context.Context, setting *diagnosis.ApplicationScanBasicSetting, applicationScanID, projectID uint32, name string) (*zapResult, error) {
	contextName := fmt.Sprintf("%v_%v_%v", projectID, applicationScanID, time.Now().Unix())
	z.targetURL = setting.Target
	err := z.HandleBasicSetting(ctx, contextName, setting.MaxDepth, setting.MaxChildren)
	if err != nil {
		return nil, err
	}
	time.Sleep(1 * time.Second)
	err = z.HandleSpiderScan(ctx, contextName, setting.MaxChildren)
	if err != nil {
		return nil, err
	}
	err = z.HandleActiveScan(ctx)
	if err != nil {
		return nil, err
	}
	report, err := z.getJsonReport(ctx)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (c *applicationScanClient) HandleBasicSetting(ctx context.Context, name string, maxDepth, maxChildren uint32) error {

	// Create Session and Context
	_, err := c.NewSession(name)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to create session, error: %v", err)
		return err
	}
	retNewContext, err := c.newContext(name)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to create context, error: %v", err)
		return err
	}
	if retNewContext["contextId"] == nil {
		err = errors.New("ContextID is null")
		appLogger.Errorf(ctx, "Failed to get context ID, error: %v", err)
		return err
	}
	c.contextID = retNewContext["contextId"].(string)
	_, err = c.IncludeInContext(name, c.targetURL+".*")
	if err != nil {
		appLogger.Errorf(ctx, "Failed to include context, error: %v", err)
		return err
	}

	// Set Scan Config
	_, err = c.SetOptionMaxChildren(int(maxChildren))
	if err != nil {
		appLogger.Errorf(ctx, "Failed to set max children option, error: %v", err)
		return err
	}
	_, err = c.SetOptionMaxDepth(int(maxDepth))
	if err != nil {
		appLogger.Errorf(ctx, "Failed to set max depth option, error: %v", err)
		return err
	}

	return nil
}

func (c *applicationScanClient) HandleSpiderScan(ctx context.Context, contextName string, maxChildren uint32) error {
	// Exec Spider
	retSpiderScan, err := c.SpiderScan(c.targetURL, fmt.Sprint(maxChildren), "True", contextName, "True")
	if err != nil {
		appLogger.Errorf(ctx, "Failed to execute spider, error: %v", err)
		return err
	}
	spiderScanID := retSpiderScan["scan"]
	if spiderScanID == nil {
		err = errors.New("SpiderScanID is null")
		appLogger.Errorf(ctx, "Failed to get spider scan ID, error: %v", err)
		return err
	}
	for {
		retSpiderStatus, err := c.SpiderStatus(spiderScanID.(string))
		if err != nil {
			appLogger.Errorf(ctx, "Failed to get status spider, error: %v", err)
			return err
		}
		status := retSpiderStatus["status"]
		if status == nil {
			err = errors.New("SpiderScanStatus is null")
			appLogger.Errorf(ctx, "Failed to get spider scan status, error: %v", err)
			return err
		}
		spiderStatus, err := strconv.Atoi(status.(string))
		if err != nil {
			appLogger.Errorf(ctx, "Failed to convert spider status to int. status: %v, error: %v", status, err)
			return fmt.Errorf("invalid spider status. err: %w", err)
		}
		if spiderStatus >= 100 {
			break
		}
	}
	return nil

}

func (c *applicationScanClient) HandleActiveScan(ctx context.Context) error {
	retAscanScan, err := c.AscanScan("", "True", "True", "Default Policy", "", "", c.contextID)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to execute active scan, error: %v", err)
		return err
	}
	ascanScanID := retAscanScan["scan"]
	if ascanScanID == nil {
		err = errors.New("ActiveScanID is null")
		appLogger.Errorf(ctx, "Failed to get active scan ID, error: %v", err)
		return err
	}
	for {
		retAscanStatus, err := c.AscanStatus(ascanScanID.(string))
		if err != nil {
			appLogger.Errorf(ctx, "Failed to get active scan status, error: %v", err)
			return err
		}
		status := retAscanStatus["status"]
		if status == nil {
			err = errors.New("ActiveScanStatus is null")
			appLogger.Errorf(ctx, "Failed to get active scan status, error: %v", err)
			return err
		}
		ascanStatus, err := strconv.Atoi(status.(string))
		if err != nil {
			appLogger.Errorf(ctx, "Failed to convert active scan status to int. status: %v, error: %v", status, err)
			return fmt.Errorf("invalid active scan status. err: %w", err)
		}
		if ascanStatus >= 100 {
			break
		}
	}
	return nil
}

func (c *applicationScanClient) getJsonReport(ctx context.Context) (*zapResult, error) {
	retJsonReport, err := c.Jsonreport()
	if err != nil {
		appLogger.Errorf(ctx, "Failed to get json report, error: %v", err)
		return nil, err
	}
	var jsonReport zapResult
	err = json.Unmarshal(retJsonReport, &jsonReport)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to marshal jsonReport, error: %v", err)
		return nil, err
	}

	return &jsonReport, nil
}
