// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameEkuiperProject = "ekuiperProject"

// EkuiperProject mapped from table <ekuiperProject>
type EkuiperProject struct {
	ID          int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ProjectName string    `gorm:"column:projectName;not null" json:"projectName"`
	ProjectID   string    `gorm:"column:projectId;not null" json:"projectId"`
	CreatedTime time.Time `gorm:"column:createdTime" json:"createdTime"`
	UpdatedTime time.Time `gorm:"column:updatedTime" json:"updatedTime"`
}

// TableName EkuiperProject's table name
func (*EkuiperProject) TableName() string {
	return TableNameEkuiperProject
}
