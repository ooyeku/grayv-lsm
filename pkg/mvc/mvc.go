package mvc

import (
	"github.com/ooyeku/grav-lsm/pkg/config"
	"net/http"
)

type MVC struct {
	Config            *config.Config
	ModelFactory      ModelFactory
	ControllerFactory ControllerFactory
	ViewFactory       ViewFactory
}

func NewMVC(config *config.Config) *MVC {
	return &MVC{
		Config: config,
	}
}

func (m *MVC) SetModelFactory(factory ModelFactory) {
	m.ModelFactory = factory
}

func (m *MVC) SetControllerFactory(factory ControllerFactory) {
	m.ControllerFactory = factory
}

func (m *MVC) SetViewFactory(factory ViewFactory) {
	m.ViewFactory = factory
}

func (m *MVC) HandleRequest(controllerName string, action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controller := m.ControllerFactory.NewController(controllerName)
		controller.Init(m.Config, m.ModelFactory)

		switch action {
		case "create":
			controller.Create(w, r)
		case "read":
			controller.Read(w, r)
		case "update":
			controller.Update(w, r)
		case "delete":
			controller.Delete(w, r)
		case "list":
			controller.List(w, r)
		default:
			http.Error(w, "Invalid action", http.StatusBadRequest)
		}
	}
}
