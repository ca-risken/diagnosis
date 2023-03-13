package wpscan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/cenkalti/backoff/v4"
	"github.com/vikyd/zero"
)

type WpscanConfig struct {
	ResultPath string
	logger     logging.Logger
	retryer    backoff.BackOff
}

func NewWpscanConfig(
	resultPath string,
	l logging.Logger,
) *WpscanConfig {
	// ref: https://pkg.go.dev/github.com/cenkalti/backoff/v4#ExponentialBackOff
	retryer := backoff.NewExponentialBackOff()
	retryer.InitialInterval = 1 * time.Second
	retryer.MaxInterval = 10 * time.Second
	retryer.MaxElapsedTime = 60 * time.Second
	return &WpscanConfig{
		ResultPath: resultPath,
		logger:     l,
		retryer:    retryer,
	}
}

func (w *WpscanConfig) run(ctx context.Context, target string, wpscanSettingID uint32, options wpscanOptions) (*wpscanResult, error) {
	now := time.Now().UnixNano()
	filePath := fmt.Sprintf("%s/%v_%v.json", w.ResultPath, wpscanSettingID, now)
	args := []string{"--clear-cache", "--disable-tls-checks", "--url", target, "-e", "u1-5", "--wp-version-all", "-f", "json", "-o", filePath}
	if options.Force {
		args = append(args, "--force")
	}
	if options.RandomUserAgent {
		args = append(args, "--random-user-agent")
	}
	if !zero.IsZeroVal(options.WpContentDir) {
		args = append(args, "--wp-content-dir", options.WpContentDir)
	}
	cmd := exec.Command("wpscan", args...)
	err := w.execWPScan(ctx, cmd)
	if err != nil {
		w.logger.Errorf(ctx, "Scan failed,target: %s, err: %v", target, err)
		return nil, err
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

	wpscanResult.CheckAccess, err = w.CheckOpen(ctx, target)
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

func (w *WpscanConfig) CheckOpen(ctx context.Context, wpURL string) (*checkAccess, error) {
	operation := func() (*checkAccess, error) {
		return w.checkOpen(wpURL)
	}
	return backoff.RetryNotifyWithData(operation, w.retryer, w.newRetryLogger(ctx, "CheckOpen"))
}

func (w *WpscanConfig) checkOpen(wpURL string) (*checkAccess, error) {
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
			return nil, fmt.Errorf("http request error: target=%+v, err=%+v", target, err)
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
	bytes, err := os.ReadFile(fileName)
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

func (w *WpscanConfig) newRetryLogger(ctx context.Context, funcName string) func(error, time.Duration) {
	return func(err error, t time.Duration) {
		w.logger.Warnf(ctx, "[RetryLogger] %s error: duration=%+v, err=%+v", funcName, t, err)
	}
}
