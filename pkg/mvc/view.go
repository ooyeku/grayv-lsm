package mvc

import (
	"net/http"
)

// View represents a generic interface for rendering views. The Render method is responsible for rendering the view
// and writing the result to the http.ResponseWriter. It takes a data parameter of type interface{} which allows the
// view to accept any type of data for rendering. The Render method returns an error if there was an issue during the
// rendering process.
type View interface {
	Render(w http.ResponseWriter, data interface{}) error
}

// ViewFactory is an interface that defines a method for creating new View instances.
// The NewView method takes a name as a parameter and returns a View instance.
// View instances are responsible for rendering data to an HTTP response writer.
// Implementations of the ViewFactory interface should provide a concrete implementation of the NewView method.
type ViewFactory interface {
	NewView(name string) View
}
