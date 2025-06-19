package main

import (
	"log"

	"github.com/ghulevishal/mcp-server/server"
)

func main() {
	log.Println("🚀 Starting MCP HTTP server...")
	server.StartHTTPServer()
}
