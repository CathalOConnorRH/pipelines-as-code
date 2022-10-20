package main

import (
	"github.com/openshift-pipelines/pipelines-as-code/pkg/adapter"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/kubeinteraction"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/params"
	evadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/signals"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	PACControllerLogKey = "pipelinesascode"
)

func main() {
	ctx := signals.NewContext()

	run := params.New()
	err := run.Clients.NewClients(ctx, &run.Info)
	if err != nil {
		log.Fatal("failed to init clients : ", err)
	}

	kinteract, err := kubeinteraction.NewKubernetesInteraction(run)
	if err != nil {
		log.Fatal("failed to init kinit client : ", err)
	}

	go func() {
		if err := run.WatchConfigMapChanges(ctx, run, ctrl.Request{}); err != nil {
			log.Fatal(err)
		}
	}()

	run.Info.Pac.LogURL = run.Clients.ConsoleUI.URL()

	evadapter.MainWithContext(ctx, PACControllerLogKey, adapter.NewEnvConfig, adapter.New(run, kinteract))
}
