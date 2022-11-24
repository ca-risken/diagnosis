package applicationscan

import (
	"context"
	"reflect"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
)

func TestGetRecommend(t *testing.T) {
	cases := []struct {
		name  string
		input *zapResultAlert
		want  *recommend
	}{
		{
			name:  "OK",
			input: &zapResultAlert{Alert: "alert", RiskDesc: "desc", RiskCode: "1", Solution: "<p>solution</p>"},
			want: &recommend{
				Risk: `alert
		- Risk: desc <risk_code: 1>`,
				Recommendation: `solution`,
			},
		},
		{
			name:  "OK/not html",
			input: &zapResultAlert{Alert: "alert", RiskDesc: "desc", RiskCode: "1", Solution: "solution"},
			want: &recommend{
				Risk: `alert
		- Risk: desc <risk_code: 1>`,
				Recommendation: `solution`,
			},
		},
		{
			name:  "OK/multi html tags",
			input: &zapResultAlert{Alert: "alert", RiskDesc: "desc", RiskCode: "1", Solution: "<ul><li>1</li><li>2</li></ul>"},
			want: &recommend{
				Risk: `alert
		- Risk: desc <risk_code: 1>`,
				Recommendation: `* 1
* 2`,
			},
		},
		{
			name:  "OK/multi line format1",
			input: &zapResultAlert{Alert: "alert", RiskDesc: "desc", RiskCode: "1", Solution: "<p>1</p><p>2</p>"},
			want: &recommend{
				Risk: `alert
		- Risk: desc <risk_code: 1>`,
				Recommendation: `1
2`,
			},
		},
		{
			name:  "OK/multi line format2",
			input: &zapResultAlert{Alert: "alert", RiskDesc: "desc", RiskCode: "1", Solution: `1<br /><br /><br /><br />2`},
			want: &recommend{
				Risk: `alert
		- Risk: desc <risk_code: 1>`,
				Recommendation: `1
2`,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := SqsHandler{logger: logging.NewLogger()}
			got := h.getRecommend(context.Background(), c.input)
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data: want=%v, got=%v", c.want, got)
			}
		})
	}
}
