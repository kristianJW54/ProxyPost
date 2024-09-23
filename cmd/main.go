package main

import (
	"ProxPost/server"
	"encoding/json"
	"fmt"
	"net/http"
)

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

		// Combine response data and monitoring information
		response := map[string]interface{}{
			"message": "Here is your sample data!",
			"data":    s.sd,
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
	http.Handle("/", server.HandleTemplate(tmpl.Tmpl["home"]))

	sampleData := &SampleData{sd: "This is a sample string"}
	http.Handle("/data", QuickRequest(sampleData))

	// Start the server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
