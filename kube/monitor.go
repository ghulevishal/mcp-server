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

// WatchPods continuously scans for restarted pods and analyzes them via LLM
func WatchPods() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatalf("❌ Failed to load kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("❌ Failed to create Kubernetes client: %v", err)
	}

	log.Println("✅ MCP server connected to cluster. Monitoring pod restarts...")

	for {
		analyzeRestartedPods(clientset)
		time.Sleep(30 * time.Second)
	}
}

// TriggerAnalysis performs a one-time scan and returns analysis results (useful for HTTP trigger)
func TriggerAnalysis() string {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return fmt.Sprintf("❌ Failed to load kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Sprintf("❌ Failed to create Kubernetes client: %v", err)
	}

	return analyzeRestartedPods(clientset)
}

// analyzeRestartedPods checks for pods with restarts, analyzes them, and logs the results
func analyzeRestartedPods(clientset *kubernetes.Clientset) string {
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("⚠️ Error listing pods: %v", err)
		return fmt.Sprintf("⚠️ Error listing pods: %v", err)
	}

	var fullOutput string

	for _, pod := range pods.Items {
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.RestartCount > 0 {
				log.Printf("🔁 Pod %s in namespace %s restarted %d times", pod.Name, pod.Namespace, cs.RestartCount)

				logs := utils.GetPodLogs(clientset, pod.Namespace, pod.Name, cs.Name)
				events := utils.GetPodEvents(clientset, pod.Namespace, pod.Name)

				prompt := fmt.Sprintf(`Analyze this Kubernetes pod issue:

Pod: %s
Namespace: %s

Logs:
%s

Events:
%s

What is the root cause and how to fix it?`,
					pod.Name, pod.Namespace, logs, events)

				result, err := llm.AnalyzePodIssue(prompt)
				if err != nil {
					log.Printf("❌ LLM analysis error for pod %s: %v", pod.Name, err)
					fullOutput += fmt.Sprintf("❌ Error analyzing pod %s: %v\n", pod.Name, err)
				} else {
					log.Printf("🧠 Analysis result for pod %s:\n%s", pod.Name, result)
					fullOutput += fmt.Sprintf("🧠 Pod %s analysis:\n%s\n\n", pod.Name, result)
				}
			}
		}
	}

	if fullOutput == "" {
		return "✅ No restarted pods found."
	}
	return fullOutput
}
