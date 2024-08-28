package mvc

import (
	"github.com/ooyeku/grav-lsm/pkg/config"
	"net/http"
)

// Controller is an interface that defines the methods for handling CRUD operations and list retrieval
// for a specific resource in an HTTP server.
//
// The Init method initializes the configuration settings and the model factory for the controller.
//
// The Create method creates a new resource based on the provided HTTP request and writes the response to the HTTP writer.
//
// The Read method retrieves an existing resource based on the provided HTTP request and writes the response to the HTTP writer.
//
// The Update method updates an existing resource based on the provided HTTP request and writes the response to the HTTP writer.
//
// The Delete method deletes an existing resource based on the provided HTTP request and writes the response to the HTTP writer.
//
// The List method retrieves a list of resources based on the provided HTTP request and writes the response to the HTTP writer.
type Controller interface {
	Init(config *config.Config, modelFactory ModelFactory)
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

// ControllerFactory is an interface used to create instances of the Controller interface.
// The NewController method takes a name as a parameter and returns an instance of the Controller interface.
// Example usage:
//
//	factory := &SomeControllerFactory{}
//	controller := factory.NewController("someController")
type ControllerFactory interface {
	NewController(name string) Controller
}
