package models

import (
	"github.com/ooyeku/grav-lsm/internal/model"
)

type TestModel struct {
	model.DefaultModel
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func (t *TestModel) TableName() string {
	return "testmodels"
}
