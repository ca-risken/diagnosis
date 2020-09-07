package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
	"go.uber.org/zap"
)

type jiraAPI interface {
	listProjects() (*[]jiraProject, error)
	listIssues(string, string, string) (*jiraIssues, error)
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

func (j *jiraClient) listProjects() (*[]jiraProject, error) {
	var projects []jiraProject
	url := j.config.JiraUrl + `rest/api/2/project`
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(j.config.JiraUserId, j.config.JiraUserPassword)
	client := new(http.Client)
	res, err := client.Do(req)

	if err != nil {
		logger.Error("Failed to list projects", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		logger.Error("Returned error code when get list projects", zap.Int("resCode", res.StatusCode))
		return &projects, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("Failed to read list projects response", zap.Error(err))
		return nil, err
	}

	if err := json.Unmarshal(body, &projects); err != nil {
		logger.Error("Failed to parse projects", zap.Error(err))
		return nil, err
	}

	return &projects, nil
}

func (j *jiraClient) listIssues(jiraKey, jiraID, recordID string) (*jiraIssues, error) {
	var issues jiraIssues
	url := j.config.JiraUrl + `rest/api/2/search`
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(j.config.JiraUserId, j.config.JiraUserPassword)
	jql := ""
	if !zero.IsZeroVal(jiraKey) {
		jql += fmt.Sprintf(`project="%s"`, jiraKey)
	} else if !zero.IsZeroVal(jiraID) {
		jql += fmt.Sprintf(`project="%s"`, jiraID)
	} else if !zero.IsZeroVal(recordID) {
		pj, err := j.getProjectByRecordID(recordID)
		if err != nil {
			logger.Error("Failed to get project by recordID", zap.Error(err))
			return nil, err
		}
		logger.Info("Project has found by recordID.", zap.String("project", pj))
		jql += fmt.Sprintf(`project="%s"`, pj)
	} else {
		logger.Error("Need to set search_key")
		return nil, errors.New("Need to set search_key")
	}
	jql += ` AND issuetype = 10023`
	logger.Info("", zap.String("jql", jql))
	params := req.URL.Query()
	params.Add("jql", jql)
	params.Add("maxResults", "1000")
	req.URL.RawQuery = params.Encode()

	client := new(http.Client)
	res, err := client.Do(req)

	if err != nil {
		logger.Error("Failed to list issues", zap.Error(err))
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		logger.Error("Returned error code when get list issues", zap.Int("resCode", res.StatusCode))
		return &issues, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("Failed to read list issues response", zap.Error(err))
		return nil, err
	}

	if err := json.Unmarshal(body, &issues); err != nil {
		logger.Error("Failed to parse issues", zap.Error(err))
		return nil, err
	}

	return &issues, nil
}

func (j *jiraClient) getProjectByRecordID(recordID string) (string, error) {
	var issues jiraIssues
	url := j.config.JiraUrl + `rest/api/2/search`
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(j.config.JiraUserId, j.config.JiraUserPassword)
	jql := fmt.Sprintf(`cf[10035]="%s"`, recordID)
	jql += ` AND issuetype = 10021`
	logger.Info("", zap.String("jql", jql))
	params := req.URL.Query()
	params.Add("jql", jql)
	params.Add("maxResults", "1000")
	req.URL.RawQuery = params.Encode()

	client := new(http.Client)
	res, err := client.Do(req)

	if err != nil {
		logger.Error("Failed to list issues", zap.Error(err))
		return "", err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		logger.Error("Returned error code when get list issues", zap.Int("resCode", res.StatusCode))
		return "", fmt.Errorf("Cannot get project by recordID. resCode: %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("Failed to read list issues response", zap.Error(err))
		return "", err
	}

	if err := json.Unmarshal(body, &issues); err != nil {
		logger.Error("Failed to parse issues", zap.Error(err))
		return "", err
	}
	issueList := issues.Issues
	if zero.IsZeroVal(issueList) {
		logger.Error("Cannot find project by recordID")
		return "", errors.New("project: Cannot find project by recordID")
	}

	return issueList[0].Fields.Project.Key, nil
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
		Project struct {
			ID  string `json:"id"`
			Key string `json:"key"`
		} `json:"project"`
	} `json:"Fields"`
	URL string `json:"self"`
}
