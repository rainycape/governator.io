// +build appengine

package main

import (
	"net/http"
	"sync"
	"time"

	_ "gnd.la/blobstore/driver/gcs"  // enable Google Could Storage blobstore driver
	_ "gnd.la/cache/driver/memcache" // enable memcached cache driver
	_ "gnd.la/commands"              // required for make-assets command
	// Uncomment the following line to use Google Cloud SQL
	//_ "gnd.la/orm/driver/mysql"
)

var (
	wg sync.WaitGroup
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
	wg.Done()
}

// Only executed on the development server. Required for
// precompiling assets.
func main() {
	wg.Wait()
}

func init() {
	wg.Add(1)
	go _app_engine_app_init()
}
