package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

// 1. start the dynamic informer factory
// go run main.go
// 2. create & update & delete a configmap
// kubectl apply -f cm.yaml
// kubectl delete -f cm.yaml

func main() {
	config := controllerruntime.GetConfigOrDie()
	dynamicclient := dynamic.NewForConfigOrDie(config)

	// factory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicclient, 5*time.Second)
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicclient, 5*time.Second, "default", func(options *metav1.ListOptions) {
		options.LabelSelector = "foo=bar"
	})

	cmInformer := factory.ForResource(corev1.SchemeGroupVersion.WithResource("configmaps"))
	cmInformer.Informer().AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cm := obj.(*unstructured.Unstructured)
			fmt.Printf("Informer event: ConfigMap ADDED %s/%s\n", cm.GetNamespace(), cm.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			old := oldObj.(*unstructured.Unstructured)
			new := newObj.(*unstructured.Unstructured)
			if old.GetResourceVersion() != new.GetResourceVersion() {
				fmt.Printf("Informer event: ConfigMap UPDATED %s/%s\n", new.GetNamespace(), new.GetName())
			}
		},
		DeleteFunc: func(obj interface{}) {
			cm := obj.(*unstructured.Unstructured)
			fmt.Printf("Informer event: ConfigMap DELETED %s/%s\n", cm.GetNamespace(), cm.GetName())
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动	informer 机制
	factory.Start(ctx.Done())

	for gvr, ok := range factory.WaitForCacheSync(ctx.Done()) {
		if !ok {
			panic(fmt.Sprintf("failed to sync cache for %v", gvr))
		}
	}

	// 通过 lister 获取所有的 configmap
	cmobjs, err := cmInformer.Lister().List(labels.Everything())
	// cmobjs, err := cmInformer.Lister().Get("dynamic-informer-cm")
	if err != nil {
		panic(err)
	}

	for _, cmobj := range cmobjs {
		cm := cmobj.(*unstructured.Unstructured)
		fmt.Printf("cmobj: %s/%s\n", cm.GetNamespace(), cm.GetName())
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
