// kube/monitor.go
package kube

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ghulevishal/mcp-server/llm"
	"github.com/ghulevishal/mcp-server/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func WatchPods() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	for {
		pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Error listing pods: %v", err)
			continue
		}

		for _, pod := range pods.Items {
			for _, cs := range pod.Status.ContainerStatuses {
				if cs.RestartCount > 0 {
					logs := utils.GetPodLogs(clientset, pod.Namespace, pod.Name, cs.Name)
					events := utils.GetPodEvents(clientset, pod.Namespace, pod.Name)
					prompt := fmt.Sprintf("Logs: %s\nEvents: %s", logs, events)
					result, err := llm.AnalyzePodIssue(prompt)
					if err != nil {
						log.Printf("LLM Analysis error: %v", err)
					} else {
						log.Printf("Pod %s analysis: %s", pod.Name, result)
					}
				}
			}
		}

		time.Sleep(30 * time.Second)
	}
}
