package main

import (
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/networtraffic/plugin"
	"os"
)

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(plugin.Name, plugin.New),
	)

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
