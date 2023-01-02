// Path: main_test.go
// Test the HTTP server by starting it and querying it
// for custom.bar and custom.foo, then verifying that their
// custom mime-types match the ones defined in mime.go.

package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeFile(t *testing.T) {
	log.Println("Testing file serving ...")

	ts := httptest.NewServer(http.HandlerFunc(serveFile))
	defer ts.Close()

	if err := addMimeTypes(); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		path, want string
	}{
		{"/custom.bar", "application/json"},
		{"/custom.foo", "foo/bar"},
	} {
		log.Println("Testing", tt.path, "...")
		res, err := http.Get(ts.URL + tt.path)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if got := res.Header.Get("Content-Type"); got != tt.want {
			t.Errorf("Content-Type: got %q, want %q", got, tt.want)
		} else {
			log.Printf("Content-Type: got %q, want %q", got, tt.want)
		}
	}
}
