package main

import (
	"fmt"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/networtraffic/plugin"
	"time"
)

const (
	nodeMeasureQueryTemplate = "irate(node_network_receive_bytes_total{device='eth0'}[1m])*1024*1024"
)

func main() {
	prometheusHandle := plugin.NewProme("http://10.106.61.245:9090", "eth0", time.Minute)

	value, err := prometheusHandle.Query(nodeMeasureQueryTemplate)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(value.String())
}
