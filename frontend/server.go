package main

import (
	"log"
	"net/http"
	"os"

	handle "github.com/dfontana/Cryptical/frontend/handlers"
	"github.com/go-macaron/binding"
	macaron "gopkg.in/macaron.v1"
)

func main() {
	// CD into the build dir for serving
	if err := os.Chdir("build"); err != nil {
		log.Fatal("failed to open build dir")
	}

	// Router
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Static("build"))
	m.Use(macaron.Renderer())

	// Builds the REST API
	m.Group("/api", func() {
		m.Get("/", handle.SendPing)
		m.Group("/model", func() {
			m.Post("/macd", binding.Json(handle.MacdModelRequest{}), handle.ModelMACD)
		})
		m.Group("/simulate", func() {
			m.Post("/macd", handle.SimulateMACD)
		})
	})

	// Serve.
	log.Println("Listening on 8080")
	http.ListenAndServe(":8080", m)
}
