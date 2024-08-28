package mvc

type Model interface {
	TableName() string
	PrimaryKey() string
	Validate() error
	BeforeSave() error
	AfterSave() error
	BeforeDelete() error
	AfterDelete() error
}

type ModelManager interface {
	Create(model Model) error
	Read(model Model, id interface{}) error
	Update(model Model) error
	Delete(model Model) error
	List(model Model, conditions ...interface{}) ([]Model, error)
}

type ModelFactory interface {
	NewModel(name string) Model
	GetModelManager() ModelManager
}
