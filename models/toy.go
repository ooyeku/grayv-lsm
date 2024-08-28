package models

import (
	"github.com/ooyeku/grav-lsm/internal/model"
)

type Toy struct {
	model.DefaultModel
	Name string `json:"name"`
	Sku  int    `json:" sku"`
}

func (t *Toy) TableName() string {
	return "toys"
}
