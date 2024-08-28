package mvc

import (
	"github.com/ooyeku/grav-lsm/pkg/config"
	"net/http"
)

// MVC represents the Model-View-Controller pattern in software development. It is a design pattern that separates the application logic into three interconnected components: the model, the view, and the controller. The MVC struct takes in a configuration object, model factory, controller factory, and view factory, which are used to configure and create the necessary components for the application. The MVC struct also provides methods to handle incoming HTTP requests and route them to the appropriate controller method.
type MVC struct {
	Config            *config.Config
	ModelFactory      ModelFactory
	ControllerFactory ControllerFactory
	ViewFactory       ViewFactory
}

// NewMVC creates a new instance of the MVC struct with the given config.
// The config parameter is used to configure the MVC instance.
// The returned MVC instance can be used to set the model, controller, and view factories,
// as well as handle requests by mapping them to the appropriate controller actions.
func NewMVC(config *config.Config) *MVC {
	return &MVC{
		Config: config,
	}
}

// SetModelFactory sets the ModelFactory for the MVC instance. The ModelFactory is responsible for creating new model instances.
// Parameters:
//   - factory: An implementation of the ModelFactory interface.
//
// Example usage:
//
//	m.SetModelFactory(&MyModelFactory{})
//	// Now the MVC instance will use the factory to create models.
func (m *MVC) SetModelFactory(factory ModelFactory) {
	m.ModelFactory = factory
}

// SetControllerFactory sets the controller factory for the MVC instance.
// The factory is used to create instances of the Controller interface.
// Parameters:
//   - factory: An implementation of the ControllerFactory interface.
//
// Example usage:
//
//	m.SetControllerFactory(&MyControllerFactory{})
//	// Now the MVC instance will use the factory to create controllers.
func (m *MVC) SetControllerFactory(factory ControllerFactory) {
	m.ControllerFactory = factory
}

// SetViewFactory sets the ViewFactory for the MVC application. The ViewFactory is responsible for creating new View instances.
func (m *MVC) SetViewFactory(factory ViewFactory) {
	m.ViewFactory = factory
}

// HandleRequest is a method of the MVC struct that returns an http.HandlerFunc.
// It takes in a controllerName string and an action string as parameters,
// and returns a function that handles HTTP requests and delegates them to the appropriate controller method.
// The function initializes the controller using the ControllerFactory, and then performs a switch statement
// on the action parameter to determine which controller method to call.
// If the action is not recognized, it returns a HTTP 400 Bad Request error.
// The method is used to handle incoming HTTP requests and route them to the appropriate controller method.
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
