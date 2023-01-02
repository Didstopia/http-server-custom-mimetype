// Path: main.go
// A simple static HTTP server that supports custom mime-types.

package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	port       = flag.String("port", "8080", "port to listen on")
	dir        = flag.String("dir", "./static", "directory to serve")
	mimeConfig = flag.String("mime-types", "./mime_types.json", "path to the JSON file containing the custom mime-types")
)

func main() {
	// Parse the command-line flags.
	log.Println("Parsing flags ...")
	flag.Parse()

	// Setting up the custom mime-types.
	log.Println("Setting up custom mime-types ...")
	if err := addMimeTypes(); err != nil {
		log.Fatal(err)
	}

	// Set up the HTTP server handler.
	log.Println("Setting up file serving from", *dir, "on port", *port)
	http.HandleFunc("/", serveFile)

	// Create a new instance of the server.
	s := &http.Server{
		Addr: ":" + *port,
	}

	// Create a channel to listen for exit signals.
	c := make(chan os.Signal, 1)

	// Listen for exit signals and send them to the channel.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine.
	go func() {
		log.Println("Starting server ...")
		if err := s.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Block until a signal is received.
	<-c

	// Shut down the server.
	log.Println("Shutting down ...")
	if err := s.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}

	// Exit.
	log.Println("Terminating ...")
	os.Exit(0)
}
