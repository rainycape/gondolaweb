// +build appengine

package main

import (
	"net/http"
	"time"

	_ "gnd.la/admin"                 // required for make-assets command
	_ "gnd.la/cache/driver/memcache" // enable memcached cache driver
	_ "gnd.la/orm/driver/gcs"        // enable Google Could Storage blobstore driver
	// Uncomment the following line to use Google Cloud SQL
	//_ "gnd.la/orm/driver/mysql"
)

func _app_engine_app_init() {
	// Make sure App is initialized before the rest
	// of this function runs.
	for App == nil {
		time.Sleep(5 * time.Millisecond)
	}
	if err := App.Prepare(); err != nil {
		panic(err)
	}
	http.Handle("/", App)
}

func init() {
	go _app_engine_app_init()
}
