package mvc

import (
	"github.com/ooyeku/grav-lsm/pkg/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockController struct {
	name string
}

func (c *mockController) Init(_ *config.Config, _ ModelFactory) {}
func (c *mockController) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create"))
}
func (c *mockController) Read(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("read"))
}
func (c *mockController) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update"))
}
func (c *mockController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete"))
}
func (c *mockController) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("list"))
}

func TestMVC(t *testing.T) {
	cf := &ControllerFactoryMock{
		callback: func(name string) Controller {
			return &mockController{name: name}
		},
	}

	mvc := NewMVC(&config.Config{})
	mvc.SetControllerFactory(cf)
	mvc.SetModelFactory(&ModelFactoryMock{})
	mvc.SetViewFactory(&ViewFactoryMock{})

	tests := []struct {
		name           string
		controllerName string
		action         string
		want           string
		wantCode       int
	}{
		{"Actions", "mock", "create", "create", http.StatusOK},
		{"Actions", "mock", "read", "read", http.StatusOK},
		{"Actions", "mock", "update", "update", http.StatusOK},
		{"Actions", "mock", "delete", "delete", http.StatusOK},
		{"Actions", "mock", "list", "list", http.StatusOK},
		{"Invalid action", "mock", "noop", "Invalid action\n", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()

			handler := mvc.HandleRequest(tt.controllerName, tt.action)
			handler.ServeHTTP(resp, req)

			got := resp.Body.String()

			if !strings.Contains(got, tt.want) {
				t.Errorf("HandleRequest() = %v, want %v", got, tt.want)
			}
			if resp.Code != tt.wantCode {
				t.Errorf("HandleRequest() code = %v, want %v", resp.Code, tt.wantCode)
			}
		})
	}
}

type ModelFactoryMock struct{}

func (m *ModelFactoryMock) NewModel(name string) Model {
	return nil
}

func (m *ModelFactoryMock) GetModelManager() ModelManager {
	return nil
}

type ControllerFactoryMock struct {
	callback func(name string) Controller
}

func (c *ControllerFactoryMock) NewController(name string) Controller {
	return c.callback(name)
}

type ViewFactoryMock struct{}

func (v *ViewFactoryMock) NewView(name string) View {
	return nil
}
