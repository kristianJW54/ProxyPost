package main

import (
	"ProxPost/server"
	"encoding/json"
	"fmt"
	"net/http"
)

// PayloadInput represents the user input payload.
type PayloadInput struct {
	URL string `json:"inputText"`
}

func ProxyPost(client *server.ClientAPI) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Ensure the request's content type is application/json
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		var input PayloadInput

		// Try to decode the JSON request body into input struct
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			// If decoding fails, return a 400 error with the error message
			http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}

		// Print the received input for testing purposes
		fmt.Printf("Received input: %s\n", input.URL)

		client.SendAPIRequest(input.URL, "GET", nil)

	})
}

func main() {

	clientContext := server.NewClientContext()                       // Create a new ClientContext
	user := &server.User{Name: "exampleUser", Token: "exampleToken"} // Initialize User

	clientAPI := &server.ClientAPI{
		Client: clientContext, // Assign the ClientContext
		User:   user,          // Assign the User
	}

	// Example request
	//clientAPI.SendAPIRequest("https://example.com", "GET", nil)

	// Set up the handler to listen for POST requests at "/sendInput"
	http.Handle("/sendInput", ProxyPost(clientAPI))

	// Serve the static HTML files from the "public" directory
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Start the server and listen on port 8080
	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
