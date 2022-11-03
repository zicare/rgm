package db

import (
	"github.com/gin-gonic/gin"
)

//TableMeta exported
type TableMeta struct {
	Fields   []string
	Primary  []string
	Serial   []string
	View     []string
	Writable []string
}

//Table exported
type Table interface {

	// Must return the table name
	Name() string

	// Must fetch foreign table data if available
	//Dig(c *gin.Context)

	// Must set conditions to filter out content
	// not intended for the user making the request
	Scope(c *gin.Context) map[string]string
}

type BaseTable struct{}

// Scope exported
func (BaseTable) Scope(c *gin.Context) map[string]string {

	return make(map[string]string)
}

// Dig exported
//func (BaseTable) Dig(c *gin.Context) {}
