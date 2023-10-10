// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNamePythonRuleGraph = "pythonRuleGraph"

// PythonRuleGraph mapped from table <pythonRuleGraph>
type PythonRuleGraph struct {
	ID          int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	GraphName   string    `gorm:"column:graphName" json:"graphName"`
	ProjectID   string    `gorm:"column:projectId;not null" json:"projectId"`
	Status      int32     `gorm:"column:status;not null" json:"status"`
	Activated   int32     `gorm:"column:activated;not null" json:"activated"`
	CreatedTime time.Time `gorm:"column:createdTime;not null" json:"createdTime"`
	UpdatedTime time.Time `gorm:"column:updatedTime;not null" json:"updatedTime"`
}

// TableName PythonRuleGraph's table name
func (*PythonRuleGraph) TableName() string {
	return TableNamePythonRuleGraph
}