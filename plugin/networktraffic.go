package plugin

import (
	"context"
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	fruntime "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
	"time"
)

const Name = "network_traffic"

type NetworkTraffic struct {
	prometheus *PrometheusHandle
	handle     framework.Handle
}

type NetworkTrafficArgs struct {
	IP         string `json:"ip"`
	DeviceName string `json:"device_name"`
	TimeRange  int    `json:"time_range"`
}

func New(plArgs runtime.Object, h framework.Handle) (framework.Plugin, error) {
	args := NetworkTrafficArgs{}
	if err := fruntime.DecodeInto(plArgs, args); err != nil {
		return nil, err
	}

	klog.Infof("[NetworkTraffic] args received. Device: %s; TimeRange: %d, Address: %s", args.DeviceName, args.TimeRange, args.IP)

	return &NetworkTraffic{
		handle:     h,
		prometheus: NewProme(args.IP, args.DeviceName, time.Minute*time.Duration(args.TimeRange)),
	}, nil
}

func (n *NetworkTraffic) Name() string {
	return Name
}

func (n *NetworkTraffic) ScoreExtensions() framework.ScoreExtensions {
	return n
}

// NormalizeScore与ScoreExtensions是固定格式
func (n *NetworkTraffic) NormalizeScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	var higherScores int64
	for _, node := range scores {
		if higherScores < node.Score {
			higherScores = node.Score
		}
	}

	// 计算公式为，满分 - (当前带宽 / 最高最高带宽 * 100)
	// 公式的计算结果为，带宽占用越大的机器，分数越低
	for i, node := range scores {
		scores[i].Score = framework.MaxNodeScore - (node.Score * 100 / higherScores)
		klog.Infof("[NetworkTraffic] Nodes final score: %v", scores)
	}

	klog.Infof("[NetworkTraffic] Nodes final score: %v", scores)
	return nil
}

func (n *NetworkTraffic) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeBandwidth, err := n.prometheus.GetGauge(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("error getting node bandwidth measure: %s", err))
	}

	bandWidth := int64(nodeBandwidth.Value)
	klog.Infof("[NetworkTraffic] node '%s' bandwidth: %s", nodeName, bandWidth)
	return bandWidth, nil
}
