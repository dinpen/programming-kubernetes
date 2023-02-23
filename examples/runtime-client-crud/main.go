package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	// 从本地 kubeconfig 文件中加载配置
	// config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	config := ctrl.GetConfigOrDie()

	cli, err := client.New(config, client.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		panic(err.Error())
	}

	namespace := "default"
	desired := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      "rest-client-crud-config",
		},
		Data: map[string]string{
			"foo": "bar",
		},
	}

	// 创建 ConfigMap
	err = cli.Create(context.Background(), desired)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Created ConfigMap %s/%s\n", namespace, desired.GetName())
	fmt.Printf("desired.ResourceVersion: %v\n", desired.ResourceVersion)

	// 获取 ConfigMap
	got := &corev1.ConfigMap{}
	err = cli.Get(context.Background(), client.ObjectKey{}, got)

}
