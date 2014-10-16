// +build !appengine

package main

func main() {
	// Start listening on main(), so the app does not start
	// listening when running tests.
	App.MustListenAndServe()
}
