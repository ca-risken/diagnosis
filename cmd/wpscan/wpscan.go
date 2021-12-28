package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gassara-kys/envconfig"
	"github.com/vikyd/zero"
)

type wpscanConfig struct {
	ResultPath         string `split_words:"true" required:"true" default:"/results"`
	WpscanVulndbApikey string `split_words:"true"`
}

func newWpscanConfig() wpscanConfig {
	var conf wpscanConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}
	return conf
}

func (w *wpscanConfig) run(target string, wpscanSettingID uint32, options wpscanOptions) (*wpscanResult, error) {
	now := time.Now().UnixNano()
	filePath := fmt.Sprintf("%s/%v_%v.json", w.ResultPath, wpscanSettingID, now)
	args := []string{"--clear-cache", "--disable-tls-checks", "--url", target, "-e", "vp,u1-5", "--wp-version-all", "-f", "json", "-o", filePath}
	if options.Force {
		args = append(args, "--force")
	}
	if options.RandomUserAgent {
		args = append(args, "--random-user-agent")
	}
	if !zero.IsZeroVal(options.WpContentDir) {
		args = append(args, "–wp-content-dir", options.WpContentDir)
	}
	if !zero.IsZeroVal(w.WpscanVulndbApikey) {
		argsWithApiKey := append(args, "--api-token", w.WpscanVulndbApikey)
		cmd := exec.Command("wpscan", argsWithApiKey...)
		err := execWPScan(cmd)
		if err != nil {
			// ReScan for Invalid APIKey or reaching APIKey Limit
			appLogger.Warn("APIKey doesn't work. Try scanning without apikey.")
			cmd := exec.Command("wpscan", args...)
			err = execWPScan(cmd)
			if err != nil {
				appLogger.Error("Scan also failed without apikey.")
				return nil, err
			}
		}
	} else {
		cmd := exec.Command("wpscan", args...)
		err := execWPScan(cmd)
		if err != nil {
			appLogger.Error("Scan failed without apikey.")
			return nil, err
		}
	}

	bytes, err := readAndDeleteFile(filePath)
	if err != nil {
		return nil, err
	}
	var wpscanResult wpscanResult
	if err := json.Unmarshal(bytes, &wpscanResult); err != nil {
		appLogger.Errorf("Failed to parse scan results. error: %v", err)
		return nil, err
	}
	wpscanResult.AccessList, _ = checkOpen(target)
	return &wpscanResult, nil
}

func execWPScan(cmd *exec.Cmd) error {
	_ = cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 && exitCode != 5 {
		appLogger.Errorf("Failed exec WPScan. exitCode: %v", exitCode)
		return fmt.Errorf("Failed exec WPScan. exitCode: %v", exitCode)
	}
	return nil
}

func checkOpen(wpURL string) ([]checkAccess, error) {
	targetList := getAccessList(wpURL)
	var retList []checkAccess
	for _, target := range targetList {
		goal := target.Goal
		if zero.IsZeroVal(target.Goal) {
			goal = target.Target
		}
		var req *http.Request
		switch target.Method {
		case "GET":
			req, _ = http.NewRequest("GET", target.Target, nil)
		case "POST":
			req, _ = http.NewRequest("POST", target.Target, nil)
		default:
			continue
		}

		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 && strings.Contains(resp.Request.URL.String(), goal) {
			target.IsAccess = true
			retList = append(retList, target)
		} else {
			target.IsAccess = false
			retList = append(retList, target)
		}
	}
	return retList, nil
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
	Version             version                `json:"version"`
	Maintheme           mainTheme              `json:"main_theme"`
	Users               map[string]interface{} `json:"users"`
	AccessList          []checkAccess
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
	Target   string
	Goal     string
	Method   string
	Type     string
	IsAccess bool
}

func getAccessList(wpURL string) []checkAccess {
	wpURL = strings.TrimSuffix(wpURL, "/")
	checkList := []checkAccess{
		checkAccess{Target: wpURL + "/wp-admin/", Goal: "wp-login.php", Method: "GET", Type: "Login"},
		checkAccess{Target: wpURL + "/admin/", Goal: "wp-login.php", Method: "GET", Type: "Login"},
		checkAccess{Target: wpURL + "/wp-login.php", Goal: "", Method: "GET", Type: "Login"},
	}
	return checkList
}
