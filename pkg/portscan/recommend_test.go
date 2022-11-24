package portscan

import (
	"reflect"
	"testing"
)

func TestGetRecommend(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  recommend
	}{
		{
			name:  "Exists recommend",
			input: recommendTypeNmap,
			want: recommend{
				Risk: `Port opens to pubilc
			- Determine if target TCP or UDP port is open to the public
			- While some ports are required to be open to the public to function properly, Restrict to trusted IP addresses.`,
				Recommendation: `Restrict target TCP and UDP port to trusted IP addresses.`,
			},
		},
		{
			name:  "Unknown recommend",
			input: "typeUnknown",
			want: recommend{
				Risk:           "",
				Recommendation: "",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getRecommend(c.input)
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: want=%v, got=%v", c.want, got)
			}
		})
	}
}
