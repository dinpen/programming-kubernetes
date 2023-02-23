package crud

import (
	"context"
	"fmt"

	"pk/basic/client"
	"pk/basic/config"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	clientset kubernetes.Interface
)

func init() {
	clientset = client.Clientset(config.KubeConfigFromCtrlRuntime())
}

// func clientset() *kubernetes.Clientset {
// 	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return kubernetes.NewForConfigOrDie(config)
// }

func CreatePod() {
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
	_, err := clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func ListPod() {
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
}

func GetPod() {
	pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), "echo", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(pod.Name)
}

func UpdatePod() {
	pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), "echo", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	pod.Spec.Containers[0].Image = "busybox:1.35.0"
	_, err = clientset.CoreV1().Pods("default").Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}
}

func DeletePod() {
	err := clientset.CoreV1().Pods("default").Delete(context.TODO(), "echo", metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
}
