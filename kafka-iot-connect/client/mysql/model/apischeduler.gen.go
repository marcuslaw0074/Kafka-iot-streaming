// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameAPIScheduler = "apiScheduler"

// APIScheduler mapped from table <apiScheduler>
type APIScheduler struct {
	ID           int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name         string    `gorm:"column:name;not null" json:"name"`
	Type         string    `gorm:"column:type;not null" json:"type"`
	Description  string    `gorm:"column:Description;not null" json:"Description"`
	DataSourceID int32     `gorm:"column:dataSourceId;not null" json:"dataSourceId"`
	StreamID     int32     `gorm:"column:streamId;not null" json:"streamId"`
	ProjectID    string    `gorm:"column:projectId;not null" json:"projectId"`
	ScheduleID   int32     `gorm:"column:scheduleId;not null" json:"scheduleId"`
	TemplateID   int32     `gorm:"column:templateId;not null" json:"templateId"`
	Status       int32     `gorm:"column:Status;not null" json:"Status"`
	Activated    int32     `gorm:"column:activated;not null" json:"activated"`
	CreatedTime  time.Time `gorm:"column:createdTime;not null" json:"createdTime"`
	UpdatedTime  time.Time `gorm:"column:updatedTime;not null" json:"updatedTime"`
}

// TableName APIScheduler's table name
func (*APIScheduler) TableName() string {
	return TableNameAPIScheduler
}
