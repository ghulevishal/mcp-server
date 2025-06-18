package utils

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPodLogs(clientset *kubernetes.Clientset, namespace, podName, containerName string) string {
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{Container: containerName})
	logs, err := req.Do(context.Background()).Raw()
	if err != nil {
		return fmt.Sprintf("Failed to get logs: %v", err)
	}
	return string(logs)
}

func GetPodEvents(clientset *kubernetes.Clientset, namespace, podName string) string {
	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Sprintf("Failed to get events: %v", err)
	}

	var sb strings.Builder
	for _, event := range events.Items {
		if event.InvolvedObject.Name == podName {
			sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", event.Type, event.Reason, event.Message))
		}
	}

	return sb.String()
}
