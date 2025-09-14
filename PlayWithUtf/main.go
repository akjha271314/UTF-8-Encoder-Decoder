package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"playwithutf/cli"
	"playwithutf/controllers"
)

//go:embed templates/index.html
var indexHTML embed.FS

func main() {
	fmt.Println("Welcome to the UTF-8 Encoder/Decoder Service!")

	// Start the web server
	router := http.NewServeMux()
	router.HandleFunc("/", serveHTML)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Println("Error writing response:", err)
			return
		}
	})

	router.HandleFunc("/playwithutf", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		controllers.PlayWithUTF8Handler(w, r)
	})

	serveMode := flag.Bool("serve", false, "Start the application in web server mode")
	flag.Parse()
	if *serveMode {
		port := ":8080"
		fmt.Printf("Starting server on port: %s\n", port)
		if err := http.ListenAndServe(port, router); err != nil {
			log.Fatal("Error starting server: ", err)
		}
	} else {
		// Call the CLI handler
		if err := cli.StartInteractiveSession(); err != nil {
			log.Fatal("Error calling CLI handler:", err)
		}
	}
}

// serveHTML is a handler function to serve the embedded index.html file.
func serveHTML(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	htmlContent, err := indexHTML.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "Could not read embedded index.html", http.StatusInternalServerError)
		fmt.Println("HERE ", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(htmlContent)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
