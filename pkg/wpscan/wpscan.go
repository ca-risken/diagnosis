package wpscan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/vikyd/zero"
)

type WpscanConfig struct {
	ResultPath         string
	WpscanVulndbApikey string
	logger             logging.Logger
}

func NewWpscanConfig(
	resultPath string,
	wpscanVulndbApikey string,
	l logging.Logger,

) *WpscanConfig {
	return &WpscanConfig{
		ResultPath:         resultPath,
		WpscanVulndbApikey: wpscanVulndbApikey,
		logger:             l,
	}
}

func (w *WpscanConfig) run(ctx context.Context, target string, wpscanSettingID uint32, options wpscanOptions) (*wpscanResult, error) {
	now := time.Now().UnixNano()
	filePath := fmt.Sprintf("%s/%v_%v.json", w.ResultPath, wpscanSettingID, now)
	args := []string{"--clear-cache", "--disable-tls-checks", "--url", target, "-e", "vp,u1-5", "--wp-version-all", "-f", "json", "-o", filePath}
	isUseAPIKey := false
	if options.Force {
		args = append(args, "--force")
	}
	if options.RandomUserAgent {
		args = append(args, "--random-user-agent")
	}
	if !zero.IsZeroVal(options.WpContentDir) {
		args = append(args, "--wp-content-dir", options.WpContentDir)
	}
	if !zero.IsZeroVal(w.WpscanVulndbApikey) {
		isUseAPIKey = true
		argsWithApiKey := append(args, "--api-token", w.WpscanVulndbApikey)
		cmd := exec.Command("wpscan", argsWithApiKey...)
		err := w.execWPScan(ctx, cmd)
		if err != nil {
			// ReScan for Invalid APIKey or reaching APIKey Limit
			w.logger.Warnf(ctx, "APIKey doesn't work. Try scanning without apikey, err=%v", err)
			cmd := exec.Command("wpscan", args...)
			err = w.execWPScan(ctx, cmd)
			if err != nil {
				w.logger.Errorf(ctx, "Scan also failed without apikey, err=%v", err)
				return nil, err
			}
		}
	} else {
		cmd := exec.Command("wpscan", args...)
		err := w.execWPScan(ctx, cmd)
		if err != nil {
			w.logger.Errorf(ctx, "Scan failed without apikey, err=%v", err)
			return nil, err
		}
	}

	bytes, err := readAndDeleteFile(filePath)
	if err != nil {
		return nil, err
	}
	var wpscanResult wpscanResult
	if err := json.Unmarshal(bytes, &wpscanResult); err != nil {
		w.logger.Errorf(ctx, "Failed to parse scan results. error: %v", err)
		return nil, err
	}

	if isUseAPIKey {
		remain := map[string]interface{}{
			"api_remaining": wpscanResult.VulnAPI.RequestRemaining,
		}
		w.logger.WithItems(ctx, logging.InfoLevel, remain, "Executed WPScan VulnAPI")
	}
	wpscanResult.CheckAccess, err = checkOpen(target)
	if err != nil {
		return nil, err
	}
	return &wpscanResult, nil
}

type wpscanError struct {
	Code   int
	StdOut *bytes.Buffer
	StdErr *bytes.Buffer
}

func (w *wpscanError) Error() string {
	return fmt.Sprintf("Failed to wpscan, code=%d, stdout=%s, stderr=%s", w.Code, w.StdOut.String(), w.StdErr.String())
}

func (w *WpscanConfig) execWPScan(ctx context.Context, cmd *exec.Cmd) error {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run command, err=%w", err)
	}
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 && exitCode != 5 {
		w.logger.Errorf(ctx, "Failed exec WPScan. exitCode: %v", exitCode)
		return &wpscanError{Code: exitCode, StdOut: &stdout, StdErr: &stderr}
	}
	return nil
}

func checkOpen(wpURL string) (*checkAccess, error) {
	checkAccess := getAccessList(wpURL)
	for i, target := range checkAccess.Target {
		goal := target.Goal
		if zero.IsZeroVal(target.Goal) {
			goal = target.URL
		}
		if target.Method != "GET" && target.Method != "POST" {
			return nil, fmt.Errorf("invalid checkAccessTarget method: %v", target.Method)
		}

		req, err := http.NewRequest(target.Method, target.URL, nil)
		if err != nil {
			return nil, err
		}
		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 && strings.Contains(resp.Request.URL.String(), goal) {
			checkAccess.Target[i].IsAccessible = true
			checkAccess.isFoundAccesibleURL = true
		}
	}
	return checkAccess, nil
}

func readAndDeleteFile(fileName string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if err := os.Remove(fileName); err != nil {
		return nil, err
	}
	return bytes, nil
}

type wpscanOptions struct {
	Force           bool   `json:"force"`
	RandomUserAgent bool   `json:"random-user-agent"`
	WpContentDir    string `json:"wp-content-dir"`
}

type wpscanResult struct {
	InterestingFindings []interestingFindings  `json:"interesting_findings"`
	Version             *version               `json:"version"`
	Maintheme           mainTheme              `json:"main_theme"`
	Users               map[string]interface{} `json:"users"`
	CheckAccess         *checkAccess
	Plugins             map[string]plugin `json:"plugins"`
	VulnAPI             vulnAPI           `json:"vuln_api"`
}

type interestingFindings struct {
	URL               string                 `json:"url"`
	ToS               string                 `json:"to_s"`
	Type              string                 `json:"type"`
	InterstingEntries []string               `json:"intersting_entries"`
	References        map[string]interface{} `json:"references"`
}

type version struct {
	Number            string          `json:"number"`
	Status            string          `json:"status"`
	InterstingEntries []string        `json:"intersting_entries"`
	Vulnerabilities   []vulnerability `json:"vulnerabilities"`
}

type mainTheme struct {
	InterstingEntries []string        `json:"intersting_entries"`
	Vulnerabilities   []vulnerability `json:"vulnerabilities"`
}

type vulnerability struct {
	Title      string                 `json:"title"`
	FixedIn    string                 `json:"fixed_in"`
	References map[string]interface{} `json:"references"`
	URL        []string               `json:"url"`
}

type plugin struct {
	Slug              string          `json:"slug"`
	LatestVersion     string          `json:"latest_version"`
	Location          string          `json:"location"`
	InterstingEntries []string        `json:"intersting_entries"`
	Vulnerabilities   []vulnerability `json:"vulnerabilities"`
	Version           version         `json:"version"`
}

type vulnAPI struct {
	Plan                   string `json:"plan"`
	RequestsDoneDuringScan uint32 `json:"requests_done_during_scan"`
	RequestRemaining       uint32 `json:"requests_remaining"`
}

type checkAccess struct {
	Target              []checkAccessTarget
	isFoundAccesibleURL bool
	isUserFound         bool
}

type checkAccessTarget struct {
	URL          string
	IsAccessible bool
	Goal         string `json:"-"`
	Method       string `json:"-"`
}

func getAccessList(wpURL string) *checkAccess {
	wpURL = strings.TrimSuffix(wpURL, "/")
	checkAccess := &checkAccess{
		Target: []checkAccessTarget{{URL: wpURL + "/wp-admin/", Goal: "wp-login.php", Method: "GET"},
			{URL: wpURL + "/admin/", Goal: "wp-login.php", Method: "GET"},
			{URL: wpURL + "/wp-login.php", Goal: "", Method: "GET"},
		},
	}
	return checkAccess
}
