package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/godown"
)

type recommend struct {
	Risk           string `json:"risk,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
}

var multiLine = regexp.MustCompile("[\n]+")

func getRecommend(alert *zapResultAlert) *recommend {
	var r string
	var buf bytes.Buffer
	if err := godown.Convert(&buf, strings.NewReader(alert.Solution), &godown.Option{}); err != nil {
		appLogger.Warnf("Failed to convert markdown from html, err=%+v, input=%s", err, alert.Solution)
		r = alert.Solution
	} else {
		r = multiLine.ReplaceAllString(strings.TrimSpace(buf.String()), "\n")
	}

	return &recommend{
		Risk: fmt.Sprintf(`%s
		- Risk: %s <risk_code: %s>`, alert.Alert, alert.RiskDesc, alert.RiskCode),
		Recommendation: r,
	}
}
