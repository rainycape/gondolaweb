// +build appengine

package main

import (
	"net/http"
	"time"

	_ "gnd.la/admin" // required for make-assets command
	_ "gnd.la/cache/driver/memcache"
)

func _app_engine_app_init() {
	// Make sure App is initialized before the rest
	// of this function runs.
	for App == nil {
		time.Sleep(5 * time.Millisecond)
	}
	http.Handle("/", App)
}

func init() {
	go _app_engine_app_init()
}
