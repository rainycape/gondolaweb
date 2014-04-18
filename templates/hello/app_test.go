package main

import (
	"gnd.la/app/tester"
	"testing"
)

func TestApp(t *testing.T) {
	// Set up the tester, passing the *testing.T and our *App, which
	// was initialized, including its handlers, in init().
	tt := tester.New(t, App)
	// Send a GET request to / without any parameters, expect a 200 response
	// code and a body matching "Hello"
	tt.Get("/", nil).Expect(200).Match("Hello")
}
