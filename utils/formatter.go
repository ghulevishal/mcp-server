package utils

import "mcp-server/kube"
import "mcp-server/llm"
import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
)

func BuildPrompt(pod corev1.Pod, logs, events string) string {
	return fmt.Sprintf(`Analyze the issue with this Kubernetes pod:

Name: %s
Namespace: %s

Events:
%s

Logs:
%s

Please suggest root cause and potential fix.`,
		pod.Name, pod.Namespace, events, logs)
}
