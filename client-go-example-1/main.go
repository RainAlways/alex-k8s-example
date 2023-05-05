package main

import (
	"context"
	"flag"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var action string
	flag.StringVar(&action, "action", "", "")

	// creates the connection
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/guqingbo/.kube/config")
	if err != nil {
		klog.Fatal(err)
	}
	ke, err := newKubernetesExample(config)
	if err != nil {
		klog.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ke.listWatchPod(ctx)

	de, err := newDynamicExample(config)
	if err != nil {
		klog.Fatal(err)
	}
	de.create(context.Background())
	switch action {
	case "patchPod":
		ke.patchPod(ctx)
	case "dynamicCreate":
		de.create(ctx)
	default:
		ke.getPod(ctx)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	select {
	case <-ch:
		cancel()
	}

	klog.Info("close")
}
