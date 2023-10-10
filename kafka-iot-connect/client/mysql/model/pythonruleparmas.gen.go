// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNamePythonRuleParma = "pythonRuleParmas"

// PythonRuleParma mapped from table <pythonRuleParmas>
type PythonRuleParma struct {
	ID          int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	RuleID      int32     `gorm:"column:ruleId;not null" json:"ruleId"`
	ParamID     int32     `gorm:"column:paramId;not null" json:"paramId"`
	Order       int32     `gorm:"column:order;not null" json:"order"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	ParamType   string    `gorm:"column:paramType;not null" json:"paramType"`
	Label       string    `gorm:"column:label;not null" json:"label"`
	Value       string    `gorm:"column:value;not null" json:"value"`
	ValueType   string    `gorm:"column:valueType;not null" json:"valueType"`
	Type        string    `gorm:"column:type;not null" json:"type"`
	IsMulti     int32     `gorm:"column:isMulti;not null" json:"isMulti"`
	Description string    `gorm:"column:Description;not null" json:"Description"`
	CreatedTime time.Time `gorm:"column:createdTime;not null" json:"createdTime"`
	UpdatedTime time.Time `gorm:"column:updatedTime;not null" json:"updatedTime"`
}

// TableName PythonRuleParma's table name
func (*PythonRuleParma) TableName() string {
	return TableNamePythonRuleParma
}