package routes

import (
	macaron "gopkg.in/macaron.v1"
)

// Build routes onto the given router for the API
func Build(m *macaron.Macaron) func() {
	return func() {
		m.Get("/", pingHandler)
	}
}
