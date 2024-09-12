package main

import (
	"ProxPost/server"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//TODO Look into ip hash
//TODO Look into client channels
//TODO latency and throughput
//TODO time-series data response times etc
//TODO ip logging address etc
//TODO look into context for handling requests

type SampleData struct {
	sd string
}

func QuickRequest(s *SampleData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check the HTTP method if necessary
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Collect headers into a map
		headers := make(map[string]string)
		for name, values := range r.Header {
			// Combine multiple values for a single header
			headers[name] = strings.Join(values, ", ")
		}

		// Gather monitoring information
		monitoringInfo := map[string]interface{}{
			"remote_addr": r.RemoteAddr,
			"headers":     headers,
			"method":      r.Method,
			"url":         r.URL.String(),
		}

		// Combine response data and monitoring information
		response := map[string]interface{}{
			"message":         "Here is your sample data!",
			"data":            s.sd,
			"monitoring_info": monitoringInfo,
		}

		// Convert response to JSON
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

/// Capturing Client connects through channels
//TODO create custom dial and listen
//TODO catch the client connection and wrap with monitoring
//TODO serve the connection and handler

// https://pkg.go.dev/net#Dial
// https://pkg.go.dev/net#example-Listener
// https://pkg.go.dev/net#Conn
// https://dev.to/hgsgtk/how-go-handles-network-and-system-calls-when-tcp-server-1nbd

func main() {
	tmpl := &server.Templates{
		Tmpl: make(map[string]string),
	}
	err := tmpl.AddTemplate("home", "./public/home.html")
	if err != nil {
		return
	}
	fmt.Println(tmpl)

	// Serve files from the static directory
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// Assign the handler to a specific route
	http.Handle("/home", server.HandleTemplate(tmpl.Tmpl["home"]))

	sampleData := &SampleData{sd: "This is a sample string"}
	http.Handle("/data", QuickRequest(sampleData))

	// Start the server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
