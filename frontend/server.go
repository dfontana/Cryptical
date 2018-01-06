package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// CD into the build dir for serving
	if err := os.Chdir("build"); err != nil {
		log.Fatal("failed to open build dir")
	}

	// Handle requests to Static
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve the entry-point
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	log.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}
