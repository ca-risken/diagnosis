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

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-diagnosis/pkg/message"
	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
)

type wpscanConfig struct {
	ResultPath         string `required:"true" split_words:"true"`
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

func (w *wpscanConfig) run(target string, wpscanSettingID uint32) (*wpscanResult, error) {
	now := time.Now().UnixNano()
	filePath := fmt.Sprintf("%s/%v_%v.json", w.ResultPath, wpscanSettingID, now)
	if !zero.IsZeroVal(w.WpscanVulndbApikey) {
		cmd := exec.Command("wpscan", "--clear-cache", "--disable-tls-checks", "--url", target, "-e", "vp,u1-5", "--wp-version-all", "-f", "json", "-o", filePath, "--api-token", w.WpscanVulndbApikey)
		err := execWPScan(cmd)
		if err != nil {
			appLogger.Warn("APIKey doesn't work. Try scanning without apikey.")
			cmd := exec.Command("wpscan", "--clear-cache", "--disable-tls-checks", "--url", target, "-e", "vp,u1-5", "--wp-version-all", "-f", "json", "-o", filePath)
			err = execWPScan(cmd)
			if err != nil {
				appLogger.Error("Scan also failed without apikey.")
				return nil, err
			}
		}

	} else {
		cmd := exec.Command("wpscan", "--clear-cache", "--disable-tls-checks", "--url", target, "-e", "vp,u1-5", "--wp-version-all", "-f", "json", "-o", filePath)

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
		if resp.StatusCode == 200 && strings.Index(resp.Request.URL.String(), goal) > -1 {
			target.IsAccess = true
			retList = append(retList, target)
		} else {
			target.IsAccess = false
			retList = append(retList, target)
		}
	}
	return retList, nil
}

func tmpRun() (*wpscanResult, error) {
	bytes, err := tmpReadFile("./tmp/wpscan.json")
	if err != nil {
		return nil, err
	}
	var wpscanResult wpscanResult
	if err := json.Unmarshal(bytes, &wpscanResult); err != nil {
		appLogger.Errorf("Failed to parse scan results. error: %v", err)
		return nil, err
	}
	return &wpscanResult, nil
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

func tmpReadFile(fileName string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func makeFindings(wpscanResult *wpscanResult, message *message.WpscanQueueMessage) ([]*finding.FindingForUpsert, error) {
	var findings []*finding.FindingForUpsert
	for _, interstingFinding := range wpscanResult.InterestingFindings {
		data, err := json.Marshal(map[string]interestingFindings{"data": interstingFinding})
		if err != nil {
			return nil, err
		}
		desc, score := getInterestingFindingInformation(interstingFinding)
		findings = append(findings, makeFinding(desc, fmt.Sprintf("interesting_findings_%v", interstingFinding.ToS), score, &data, message))
	}
	desc, score := getVersionFindingInformation(wpscanResult.Version)
	if !zero.IsZeroVal(desc) {
		data, err := json.Marshal(map[string]version{"data": wpscanResult.Version})
		if err != nil {
			return nil, err
		}
		findings = append(findings, makeFinding(desc, fmt.Sprintf("version_%v", message.TargetURL), score, &data, message))
	}
	isUserFound := false
	for key, val := range wpscanResult.Users {
		isUserFound = true
		data, err := json.Marshal(map[string]interface{}{"data": val})
		if err != nil {
			return nil, err
		}
		desc := fmt.Sprintf("User %v was found.", key)
		score := float32(3.0)
		findings = append(findings, makeFinding(desc, fmt.Sprintf("username_%v", key), score, &data, message))
	}
	for _, access := range wpscanResult.AccessList {
		desc, score := getAccessFindingInformation(access, isUserFound)
		if !zero.IsZeroVal(desc) {
			data, err := json.Marshal(map[string]interface{}{"data": map[string]string{
				"url": access.Target,
			}})
			if err != nil {
				return nil, err
			}
			findings = append(findings, makeFinding(desc, fmt.Sprintf("Accesible_%v", access.Target), score, &data, message))
		}
	}
	return findings, nil
}

func makeFinding(description, dataSourceID string, score float32, data *[]byte, message *message.WpscanQueueMessage) *finding.FindingForUpsert {
	return &finding.FindingForUpsert{
		Description:      description,
		DataSource:       message.DataSource,
		DataSourceId:     generateDataSourceID(dataSourceID),
		ResourceName:     message.TargetURL,
		ProjectId:        message.ProjectID,
		OriginalScore:    score,
		OriginalMaxScore: MaxScore,
		Data:             string(*data),
	}
}

func getInterestingFindingInformation(ie interestingFindings) (string, float32) {
	switch ie.ToS {
	case "Headers":
		return "Software version found by Headers", 1.0
	default:
		return ie.ToS, 1.0
	}
}

func getVersionFindingInformation(version version) (string, float32) {
	if zero.IsZeroVal(version.Number) {
		return "", 0.0
	}
	if version.Status == "insecure" {
		return fmt.Sprintf("WordPress version %v identified (Insecure)", version.Number), 6.0
	}
	return fmt.Sprintf("WordPress version %v identified", version.Number), 1.0
}

func getAccessFindingInformation(access checkAccess, isUserFound bool) (string, float32) {
	switch access.Type {
	case "Login":
		if !access.IsAccess {
			return "WordPress login page is closed.", 1.0
		}
		if isUserFound {
			return "WordPress login page is open. And username was found.", 9.0
		}
		return "WordPress login page is open.", 8.0

	default:
		return "", 0.0
	}

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
	if strings.HasSuffix(wpURL, "/") {
		wpURL = wpURL[:len(wpURL)-1]
	}
	checkList := []checkAccess{
		checkAccess{Target: wpURL + "/wp-admin/", Goal: "wp-login.php", Method: "GET", Type: "Login"},
		checkAccess{Target: wpURL + "/admin/", Goal: "wp-login.php", Method: "GET", Type: "Login"},
		checkAccess{Target: wpURL + "/wp-login.php", Goal: "", Method: "GET", Type: "Login"},
	}
	return checkList
}
