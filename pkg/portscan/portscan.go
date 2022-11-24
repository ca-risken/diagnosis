package portscan

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/portscan"
	"github.com/ca-risken/datasource-api/pkg/message"
)

type portscanAPI interface {
	makeTargets(string)
	getResult(context.Context, *message.PortscanQueueMessage) ([]*portscan.NmapResult, error)
	scan(context.Context) ([]*portscan.NmapResult, error)
}

type portscanClient struct {
	target []target
}

func newPortscanClient() (portscanAPI, error) {
	p := portscanClient{}
	return &p, nil
}

func (p *portscanClient) getResult(ctx context.Context, message *message.PortscanQueueMessage) ([]*portscan.NmapResult, error) {
	nmapResults, err := p.scan(ctx)
	if err != nil {
		return []*portscan.NmapResult{}, fmt.Errorf("portscan error: err=%w", err)
	}

	return nmapResults, nil
}

func (p *portscanClient) scan(ctx context.Context) ([]*portscan.NmapResult, error) {
	var nmapResults []*portscan.NmapResult
	for _, target := range p.target {
		results, err := portscan.Scan(target.Target, target.Protocol, target.FromPort, target.ToPort)
		if err != nil {
			return nil, err
		}
		for _, result := range results {
			result.ResourceName = target.Target
			nmapResults = append(nmapResults, result)
		}
	}

	return nmapResults, nil
}

func (p *portscanClient) makeTargets(targetIPFQDN string) {
	p.target = []target{
		{
			Target:   targetIPFQDN,
			FromPort: 0,
			ToPort:   0,
			Protocol: "udp",
		},
		{
			Target:   targetIPFQDN,
			FromPort: 1,
			ToPort:   65535,
			Protocol: "tcp",
		},
	}
}

type target struct {
	Target   string
	FromPort int
	ToPort   int
	Protocol string
}
