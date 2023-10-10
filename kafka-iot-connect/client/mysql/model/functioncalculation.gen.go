// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameFunctionCalculation = "functionCalculation"

// FunctionCalculation mapped from table <functionCalculation>
type FunctionCalculation struct {
	ID                int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CalculationFileID int32     `gorm:"column:calculationFileId;not null" json:"calculationFileId"`
	CreatedTime       time.Time `gorm:"column:createdTime" json:"createdTime"`
	DataSourceFromID  int32     `gorm:"column:dataSourceFromId;not null" json:"dataSourceFromId"`
	DataSourceToID    int32     `gorm:"column:dataSourceToId;not null" json:"dataSourceToId"`
	EndTime           time.Time `gorm:"column:endTime" json:"endTime"`
	ExeID             int32     `gorm:"column:exeId;not null" json:"exeId"`
	FileID            int32     `gorm:"column:fileId;not null" json:"fileId"`
	Interval          string    `gorm:"column:interval;not null" json:"interval"`
	Name              string    `gorm:"column:name;not null" json:"name"`
	ProjectID         string    `gorm:"column:projectId;not null" json:"projectId"`
	Spec              string    `gorm:"column:spec;not null" json:"spec"`
	StartTime         time.Time `gorm:"column:startTime" json:"startTime"`
	Status            int32     `gorm:"column:status;not null" json:"status"`
	UpdatedTime       time.Time `gorm:"column:updatedTime" json:"updatedTime"`
}

// TableName FunctionCalculation's table name
func (*FunctionCalculation) TableName() string {
	return TableNameFunctionCalculation
}