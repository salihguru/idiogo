package server

import (
	"context"
	"log"
)

// Server is the interface that must be implemented by a server
// It is used to listen for requests
// The Listen method is used to listen for requests
// The Listen method returns an error
// The error is nil if the server is listening successfully
type Listener interface {
	// Listen is used to listen for requests
	// The Listen method returns an error
	// The error is nil if the server is listening successfully
	Listen() error

	// Shutdown is used to shutdown the server
	// The Shutdown method returns an error
	// The error is nil if the server is shutdown successfully
	Shutdown(context.Context) error
}

func Start(name string, srv Listener, cb func()) {
	go func() {
		defer cb()
		if err := srv.Listen(); err != nil {
			log.Fatal(name, " server error: ", err)
		}
	}()
}
