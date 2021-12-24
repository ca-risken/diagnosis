package main

const (
	recommendTypeNmap = "Nmap"
)

type recommend struct {
	Risk           string `json:"risk,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
}

func getRecommend(recommendType string) recommend {
	return recommendMap[recommendType]
}

var recommendMap = map[string]recommend{
	recommendTypeNmap: {
		Risk: `Port opens to pubilc
			- Determine if target TCP or UDP port is open to the public
			- While some ports are required to be open to the public to function properly, Restrict to trusted IP addresses.`,
		Recommendation: `Restrict target TCP and UDP port to trusted IP addresses.`,
	},
}
