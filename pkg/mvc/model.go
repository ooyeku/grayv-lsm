package mvc

// Model is an interface that defines the behavior of a model in a software system.
// It includes methods for retrieving the table name, primary key, validation,
// and various lifecycle hooks for saving and deleting the model.
type Model interface {
	TableName() string
	PrimaryKey() string
	Validate() error
	BeforeSave() error
	AfterSave() error
	BeforeDelete() error
	AfterDelete() error
}

// ModelManager is an interface that defines CRUD operations (Create, Read, Update, Delete)
// and a List operation for managing models.
// The model type passed to these operations must implement the Model interface.
// The Model interface provides methods for defining the table name, primary key,
// validation, and callback hooks for save and delete operations.
type ModelManager interface {
	Create(model Model) error
	Read(model Model, id interface{}) error
	Update(model Model) error
	Delete(model Model) error
	List(model Model, conditions ...interface{}) ([]Model, error)
}

// ModelFactory is an interface that defines methods for creating new models and retrieving the model manager.
//
// The NewModel method creates a new model with the given name.
//
// The GetModelManager method retrieves the model manager.
type ModelFactory interface {
	NewModel(name string) Model
	GetModelManager() ModelManager
}
