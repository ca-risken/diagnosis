package main

import (
	"context"

	"github.com/CyberAgent/mimosa-common/pkg/portscan"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-diagnosis/pkg/message"
	"github.com/aws/aws-xray-sdk-go/xray"
)

type portscanAPI interface {
	getResult(*message.PortscanQueueMessage, bool) ([]*finding.FindingForUpsert, error)
	scan() ([]*portscan.NmapResult, error)
}

type portscanClient struct {
	target []target
}

func newPortscanClient() (*portscanClient, error) {
	//	var conf portscanConfig
	//	err := envconfig.Process("", &conf)
	//	if err != nil {
	//		return nil, err
	//	}

	p := portscanClient{}
	return &p, nil
}

func (p *portscanClient) getResult(ctx context.Context, message *message.PortscanQueueMessage) ([]*finding.FindingForUpsert, error) {
	putData := []*finding.FindingForUpsert{}

	_, segment := xray.BeginSubsegment(ctx, "scanTargets")
	nmapResults, err := p.scan()
	segment.Close(err)
	if err != nil {
		appLogger.Errorf("Faild to Portscan: err=%+v", err)
		return putData, err
	}
	putData, err = makeFindings(nmapResults, message)
	if err != nil {
		appLogger.Errorf("Faild to make findings: err=%+v", err)
		return putData, err
	}

	return putData, nil
}

func (p *portscanClient) scan() ([]*portscan.NmapResult, error) {
	var nmapResults []*portscan.NmapResult
	for _, target := range p.target {
		results, err := portscan.Scan(target.Target, target.Protocol, target.FromPort, target.ToPort)
		if err != nil {
			appLogger.Warnf("Error occured when scanning. err: %v", err)
			return nmapResults, nil
		}
		for _, result := range results {
			result.ResourceName = target.Target
			nmapResults = append(nmapResults, result)
		}
	}

	return nmapResults, nil
}

func makeTargets(targetIPFQDN string) []target {
	return []target{
		target{
			Target:   targetIPFQDN,
			FromPort: 0,
			ToPort:   0,
			Protocol: "udp",
		},
		target{
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
