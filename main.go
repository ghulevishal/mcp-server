// main.go
package main

import (
	"log"
	"mcp-server/kube"
)

func main() {
	log.Println("Starting MCP server...")
	kube.StartMonitoring()
}