package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
	"golang.org/x/net/context/ctxhttp"
)

type jiraAPI interface {
	getJiraProject(context.Context, string, string, string, string) (string, map[string]string)
	listIssues(context.Context, string) (*jiraIssues, error)
}

type jiraClient struct {
	config jiraConfig
	//	Sess *session.Session
	//	Svc  *jira.Jira
}

type jiraConfig struct {
	JiraUrl          string `split_words:"true" default:"http://localhost"`
	JiraUserId       string `split_words:"true" default:"admin"`
	JiraUserPassword string `split_words:"true" default:"password"`
}

func newJiraClient() *jiraClient {
	var conf jiraConfig
	err := envconfig.Process("diagnosis", &conf)
	if err != nil {
		panic(err)
	}
	return &jiraClient{config: conf}
}

func (j *jiraClient) getJiraProject(ctx context.Context, jiraKey, jiraID, IdentityField, IdentityValue string) (string, map[string]string) {
	errDetail := make(map[string]string)
	if !zero.IsZeroVal(jiraKey) {
		_, err := j.searchProjectByJiraKeyID(ctx, jiraKey)
		if err != nil {
			errDetail["jiraKey"] = err.Error()
		} else {
			return jiraKey, nil
		}
	}
	if !zero.IsZeroVal(jiraID) {
		_, err := j.searchProjectByJiraKeyID(ctx, jiraID)
		if err != nil {
			errDetail["jiraID"] = err.Error()
		} else {
			return jiraID, nil
		}
	}
	if !zero.IsZeroVal(IdentityField) && !zero.IsZeroVal(IdentityValue) {
		pj, err := j.getProjectByIdentityKey(ctx, IdentityField, IdentityValue)
		if err != nil {
			errDetail["IdentityField"] = err.Error()
		} else {
			return pj, nil
		}
	}

	return "", errDetail

}
func (j *jiraClient) searchProjectByJiraKeyID(ctx context.Context, search string) (bool, error) {
	var issues jiraIssues
	url := j.config.JiraUrl + `rest/api/2/search`
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(j.config.JiraUserId, j.config.JiraUserPassword)
	jql := fmt.Sprintf(`project="%s" AND issuetype = 10021`, search)
	params := req.URL.Query()
	params.Add("jql", jql)
	params.Add("maxResults", "1000")
	req.URL.RawQuery = params.Encode()
	client := new(http.Client)
	res, err := ctxhttp.Do(ctx, xray.Client(client), req)
	if err != nil {
		appLogger.Errorf("Failed to list projects. sk the System Administrator, error: %v", err)
		return false, err
	}

	defer res.Body.Close()
	if res.StatusCode == 400 {
		return false, fmt.Errorf(`%v projects found. Please check your value`, issues.Total)
	}
	if res.StatusCode != 200 {
		appLogger.Errorf("Returned error code when get list issues, error: %v", err)
		return false, fmt.Errorf("Cannot get project by jiraID or jiraKey. Ask the System Administrator")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		appLogger.Errorf("Failed to read list projects response, error: %v", err)
		return false, err
	}
	if err := json.Unmarshal(body, &issues); err != nil {
		appLogger.Errorf("Failed to parse issues. Ask the System Administrator, error: %v", err)
		return false, err
	}

	if issues.Total != 1 {
		appLogger.Warnf("Unexpect number issues found. issues Total: %v", issues.Total)
		return false, fmt.Errorf(`%v issues found. Please check your recordID,recordKey`, issues.Total)
	}

	return true, nil
}

func (j *jiraClient) getProjectByIdentityKey(ctx context.Context, identityField, identityValue string) (string, error) {
	var issues jiraIssues
	url := j.config.JiraUrl + `rest/api/2/search`
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(j.config.JiraUserId, j.config.JiraUserPassword)
	jql := fmt.Sprintf(`cf[%s]="%s"`, identityField, identityValue)
	jql += ` AND issuetype = 10021`
	params := req.URL.Query()
	params.Add("jql", jql)
	params.Add("maxResults", "2")
	req.URL.RawQuery = params.Encode()

	client := new(http.Client)
	res, err := ctxhttp.Do(ctx, xray.Client(client), req)

	if err != nil {
		appLogger.Errorf("Failed to list projects, error: %v", err)
		return "", err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		appLogger.Errorf("Returned error code when get list issues, responseCode: %v", res.StatusCode)
		return "", fmt.Errorf("Cannot get project by IdentityKey,Field")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		appLogger.Errorf("Failed to read list issues response, error: %v", err)
		return "", err
	}

	if err := json.Unmarshal(body, &issues); err != nil {
		appLogger.Errorf("Failed to parse issues, error: %v", err)
		return "", err
	}

	if issues.Total != 1 {
		appLogger.Warnf("Unexpect number issues found. issues.Total: %v", issues.Total)
		return "", fmt.Errorf(`%v issues found. Please check your recordID,recordKey`, issues.Total)
	}

	issueList := issues.Issues
	if zero.IsZeroVal(issueList) {
		appLogger.Error("Cannot find project by IdentityKey,Field.")
		return "", errors.New("project: Cannot find project by IdentityKey,Field.IdentityKey,Field")
	}

	return issueList[0].Fields.Project.Key, nil
}

func (j *jiraClient) listIssues(ctx context.Context, project string) (*jiraIssues, error) {
	var issues jiraIssues
	url := j.config.JiraUrl + `rest/api/2/search`
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(j.config.JiraUserId, j.config.JiraUserPassword)
	jql := fmt.Sprintf(`project="%s"`, project)
	jql += ` AND issuetype = 10023`
	params := req.URL.Query()
	params.Add("jql", jql)
	params.Add("maxResults", "1000")
	req.URL.RawQuery = params.Encode()

	client := new(http.Client)
	res, err := ctxhttp.Do(ctx, xray.Client(client), req)
	if err != nil {
		appLogger.Errorf("Failed to list Issues, error: %v", err)
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		appLogger.Errorf("Returned error code when get list Issues, resCode: %v", res.StatusCode)
		return &issues, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		appLogger.Errorf("Failed to read list issues response, error: %v", err)
		return nil, err
	}

	if err := json.Unmarshal(body, &issues); err != nil {
		appLogger.Errorf("Failed to parse issues, error: %v", err)
		return nil, err
	}

	return &issues, nil
}

type jiraProject struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
	URL  string `json:"self"`
}

type jiraIssues struct {
	MaxResults int         `json:"maxResults"`
	Total      int         `json:"total"`
	Issues     []jiraIssue `json:"issues"`
}

type jiraIssue struct {
	Key    string `json:"key"`
	Fields struct {
		Date     string   `json:"created"`
		Priority struct { // <- 構造体の中にネストさせて構造体を定義
			Name string `json:"name"`
		} `json:"priority"`
		Summary string `json:"summary"`
		Target  string `json:"customfield_10042"`
		Status  struct {
			Name string `json:"name"`
		} `json:"status"`
		Project jiraProject `json:"project"`
	} `json:"Fields"`
	URL string `json:"self"`
}
