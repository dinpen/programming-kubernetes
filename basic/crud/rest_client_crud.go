package crud

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func restclient() *rest.RESTClient {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	// 配置 API 路径和请求资源组/资源版本信息
	config.APIPath = "api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	// 配置数据的编解码器
	config.NegotiatedSerializer = scheme.Codecs
	client, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}
	return client
}

func createPod2() {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "echo",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "echo",
					Image: "busybox",
					Command: []string{
						"/bin/sh",
						"-c",
						"echo hello world",
					},
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}
	_, err := restclient().Post().Namespace("default").Resource("pods").Body(pod).Do(context.TODO()).Get()
	if err != nil {
		panic(err)
	}
}

func listPod2() {
	var pods v1.PodList
	err := restclient().Get().Namespace("default").Resource("pods").Do(context.TODO()).Into(&pods)
	if err != nil {
		panic(err)
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
}

func getPod2() {
	var pod v1.Pod
	err := restclient().Get().Namespace("default").Resource("pods").Name("echo").Do(context.TODO()).Into(&pod)
	if err != nil {
		panic(err)
	}
	fmt.Println(pod.Name)
}

func updatePod2() {
	pod := &v1.PodSpec{}
	err := restclient().Put().Namespace("default").Resource("pods").Name("echo").Body(pod).Do(context.TODO()).Error()
	if err != nil {
		panic(err)
	}

	// pod.Spec.Containers[0].Image = "busybox:1.35.0"
	// _, err = clientset().CoreV1().Pods("default").Update(context.TODO(), pod, metav1.UpdateOptions{})
	// if err != nil {
	// 	panic(err)
	// }
}

func deletePod2() {
	err := restclient().Delete().Namespace("default").Resource("pods").Name("echo").Do(context.TODO()).Error()
	if err != nil {
		panic(err)
	}
}
