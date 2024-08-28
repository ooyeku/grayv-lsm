package mvc

import (
	"github.com/ooyeku/grav-lsm/pkg/config"
	"net/http"
)

type Controller interface {
	Init(config *config.Config, modelFactory ModelFactory)
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

type ControllerFactory interface {
	NewController(name string) Controller
}
