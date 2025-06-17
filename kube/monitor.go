package kube

import (
	"context"
	"fmt"
	"log"
	"time"

	"mcp-server/llm"
	"mcp-server/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func StartMonitoring() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		pods, _ := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		for _, pod := range pods.Items {
			for _, cs := range pod.Status.ContainerStatuses {
				if cs.RestartCount > 0 {
					logs := getPodLogs(clientset, pod)
					events := getPodEvents(clientset, pod)
					prompt := utils.BuildPrompt(pod, logs, events)
					llm.Analyze(prompt)
				}
			}
		}
	}
}

func getPodLogs(clientset *kubernetes.Clientset, pod corev1.Pod) string {
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{TailLines: int64Ptr(50)})
	logs, _ := req.Do(context.TODO()).Raw()
	return string(logs)
}

func getPodEvents(clientset *kubernetes.Clientset, pod corev1.Pod) string {
	events, _ := clientset.CoreV1().Events(pod.Namespace).List(context.TODO(), metav1.ListOptions{})
	var relevant string
	for _, e := range events.Items {
		if e.InvolvedObject.Name == pod.Name {
			relevant += fmt.Sprintf("%s: %s\n", e.Reason, e.Message)
		}
	}
	return relevant
}

func int64Ptr(i int64) *int64 { return &i }
