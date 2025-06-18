package main

import (
	"log"

	"github.com/ghulevishal/mcp-server/kube"
)

func main() {
	log.Println("Starting MCP server...")
	kube.WatchPods()
}
