package main

import (
	"context"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

type kubernetesExample struct {
	clientSet *kubernetes.Clientset
	logger    klog.Logger
}

func newKubernetesExample(config *restclient.Config) (k kubernetesExample, err error) {
	k.logger = klog.LoggerWithName(klog.Background(), "kubernetesExample")
	k.clientSet, err = kubernetes.NewForConfig(config)
	return k, errors.WithStack(err)
}

func (k kubernetesExample) getPod(ctx context.Context) {
	pod, err := k.clientSet.CoreV1().Pods("kube-system").Get(ctx, "kube-apiserver-minikube", metav1.GetOptions{})
	if err != nil {
		k.logger.Error(err, "getPod")
		return
	}

	k.logger.Info("getPod:", "pod", pod)
}

func (k kubernetesExample) patchPod(ctx context.Context) {
	//给pod的/metadata/annotations，增加一个alex:gu的标签
	pod, err := k.clientSet.CoreV1().Pods("kube-system").
		Patch(ctx, "kube-apiserver-minikube", types.JSONPatchType,
			[]byte(`[{"op": "add", "path": "/metadata/annotations/alex", "value": "gu"}]`),
			metav1.PatchOptions{})
	if err != nil {
		k.logger.Error(err, "patch add")
		return
	}
	k.logger.Info("patch add:", "alex", pod.GetObjectMeta().GetAnnotations()["alex"])

	//给pod的/metadata/annotations，移除alex:gu的标签
	pod, err = k.clientSet.CoreV1().Pods("kube-system").
		Patch(ctx, "kube-apiserver-minikube", types.JSONPatchType,
			[]byte(`[{"op": "remove", "path": "/metadata/annotations/alex"}]`),
			metav1.PatchOptions{})
	if err != nil {
		k.logger.Error(err, "patch remove")
		return
	}

	k.logger.Info("patch remove:", "alex", pod.GetObjectMeta().GetAnnotations()["alex"])
}

func (k kubernetesExample) listWatchPod(ctx context.Context) {
	podList, err := k.clientSet.CoreV1().Pods("kube-system").List(ctx, metav1.ListOptions{})
	if err != nil {
		k.logger.Error(err, "list err")
		return
	}
	k.logger.Info("list", "len", len(podList.Items))

	//从list获得的最后记录开始watch
	watcher, err := k.clientSet.CoreV1().Pods("kube-system").Watch(ctx, metav1.ListOptions{ResourceVersion: podList.GetResourceVersion()})
	if err != nil {
		k.logger.Error(err, "watch err")
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.ResultChan():
				if !ok {
					k.logger.Info("watcher", "watcher", "close")
					return
				}
				k.logger.Info("watcher", "event", event)
			}
		}
	}()
}
