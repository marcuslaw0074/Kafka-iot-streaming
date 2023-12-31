// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameSyncDataSource = "syncDataSource"

// SyncDataSource mapped from table <syncDataSource>
type SyncDataSource struct {
	ID             int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ProjectID      string    `gorm:"column:projectId;not null" json:"projectId"`
	DataSourceName string    `gorm:"column:dataSourceName;not null" json:"dataSourceName"`
	SourceType     string    `gorm:"column:sourceType;not null" json:"sourceType"`
	Status         int32     `gorm:"column:status;not null" json:"status"`
	DataSourceID   int32     `gorm:"column:dataSourceId;not null" json:"dataSourceId"`
	DbType         string    `gorm:"column:dbType;not null" json:"dbType"`
	DbName         string    `gorm:"column:dbName;not null" json:"dbName"`
	DbServer       string    `gorm:"column:dbServer;not null" json:"dbServer"`
	DbTable        string    `gorm:"column:dbTable;not null" json:"dbTable"`
	CreatedTime    time.Time `gorm:"column:createdTime;not null" json:"createdTime"`
	UpdatedTime    time.Time `gorm:"column:updatedTime;not null" json:"updatedTime"`
}

// TableName SyncDataSource's table name
func (*SyncDataSource) TableName() string {
	return TableNameSyncDataSource
}
