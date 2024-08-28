package mvc

import (
	"net/http"
)

type View interface {
	Render(w http.ResponseWriter, data interface{}) error
}

type ViewFactory interface {
	NewView(name string) View
}
