package main

import (
	"testing"
)

func TestGetDescription(t *testing.T) {
	cases := []struct {
		name   string
		alert  *zapResultAlert
		target string
		want   string
	}{
		{
			name: "OK",
			alert: &zapResultAlert{
				Alert: "hogehoge",
			},
			target: "http://example.com/",
			want:   "hogehoge found in http://example.com/.",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getDescription(c.alert, c.target)
			if c.want != got {
				t.Fatalf("Unexpected resource name: want=%v, got=%v", c.want, got)
			}
		})
	}
}

func TestScore(t *testing.T) {
	cases := []struct {
		name  string
		alert *zapResultAlert
		want  float32
	}{
		{
			name: "OK (Information)",
			alert: &zapResultAlert{
				RiskDesc: "Information (Medium)",
			},
			want: ScoreInformation,
		},
		{
			name: "OK (Low)",
			alert: &zapResultAlert{
				RiskDesc: "Low (Medium)",
			},
			want: ScoreLow,
		},
		{
			name: "OK (Medium)",
			alert: &zapResultAlert{
				RiskDesc: "Medium (High)",
			},
			want: ScoreMedium,
		},
		{
			name: "OK (High)",
			alert: &zapResultAlert{
				RiskDesc: "High (Medium)",
			},
			want: ScoreHigh,
		},
		{
			name: "OK (Other)",
			alert: &zapResultAlert{
				RiskDesc: "HogeFuga (Medium)",
			},
			want: ScoreOther,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getScore(c.alert)
			if c.want != got {
				t.Fatalf("Unexpected resource name: want=%v, got=%v", c.want, got)
			}
		})
	}
}
