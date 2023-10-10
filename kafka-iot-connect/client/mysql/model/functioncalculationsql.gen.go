// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameFunctionCalculationSql = "functionCalculationSql"

// FunctionCalculationSql mapped from table <functionCalculationSql>
type FunctionCalculationSql struct {
	ID           int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	DataSourceID int32     `gorm:"column:dataSourceId;not null" json:"dataSourceId"`
	RuleID       int32     `gorm:"column:ruleId;not null" json:"ruleId"`
	RuleSqlID    int32     `gorm:"column:ruleSqlId;not null" json:"ruleSqlId"`
	Sql          string    `gorm:"column:sql;not null" json:"sql"`
	Status       int32     `gorm:"column:status;not null" json:"status"`
	Order        int32     `gorm:"column:order;not null" json:"order"`
	CreatedTime  time.Time `gorm:"column:createdTime;not null" json:"createdTime"`
	UpdatedTime  time.Time `gorm:"column:updatedTime;not null" json:"updatedTime"`
}

// TableName FunctionCalculationSql's table name
func (*FunctionCalculationSql) TableName() string {
	return TableNameFunctionCalculationSql
}