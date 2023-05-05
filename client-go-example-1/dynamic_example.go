package main

import (
	"context"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

type dynamicExample struct {
	dynamicClient *dynamic.DynamicClient
	logger        klog.Logger
}

func newDynamicExample(config *restclient.Config) (d dynamicExample, err error) {
	d.logger = klog.LoggerWithName(klog.Background(), "dynamicExample")
	d.dynamicClient, err = dynamic.NewForConfig(config)
	return d, errors.WithStack(err)
}

func (d dynamicExample) create(ctx context.Context) {
	rs, err := d.dynamicClient.Resource(schema.GroupVersionResource{
		Group: "examples.alex.com", Version: "v1", Resource: "alexdynamics",
	}).
		Namespace(metav1.NamespaceDefault).
		Create(ctx, &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "examples.alex.com/v1",
				"kind":       "AlexDynamic",
				"metadata": map[string]interface{}{
					"name": "alex-demo",
				},
			},
		}, metav1.CreateOptions{})
	if err != nil {
		d.logger.Error(err, "create")
		return
	}

	d.logger.Info("create:", "rs", rs)
}
