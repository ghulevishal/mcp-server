package server

import (
	"log"
	"net/http"

	"github.com/ghulevishal/mcp-server/kube"
)

// StartHTTPServer starts a basic HTTP server with a trigger endpoint
func StartHTTPServer() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// ✅ Register the /analyze route
	http.HandleFunc("/analyze", handleAnalyze)

	log.Println("🌐 Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	result := kube.TriggerAnalysis()
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

// package server

// import (
// 	"fmt"
// 	"net/http"
// )

// func StartHTTPServer() {
// 	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintln(w, "✅ MCP trigger route hit successfully!")
// 	})

// 	fmt.Println("🌐 MCP server listening on http://localhost:8080")
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		panic(fmt.Sprintf("❌ Failed to start server: %v", err))
// 	}
// }
