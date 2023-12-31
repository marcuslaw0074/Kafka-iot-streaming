// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameFunctionRule = "functionRule"

// FunctionRule mapped from table <functionRule>
type FunctionRule struct {
	ID           int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ClientID     int32     `gorm:"column:clientId;not null" json:"clientId"`
	DataSourceID int32     `gorm:"column:dataSourceId;not null" json:"dataSourceId"`
	RuleName     string    `gorm:"column:ruleName;not null" json:"ruleName"`
	RuleType     string    `gorm:"column:ruleType;not null" json:"ruleType"`
	RuleSubType  string    `gorm:"column:ruleSubType;not null" json:"ruleSubType"`
	RuleGraphID  int32     `gorm:"column:ruleGraphId;not null" json:"ruleGraphId"`
	Order        int32     `gorm:"column:order;not null" json:"order"`
	ScheduleID   int32     `gorm:"column:scheduleId;not null" json:"scheduleId"`
	StreamID     int32     `gorm:"column:streamId;not null" json:"streamId"`
	SinkID       int32     `gorm:"column:sinkId;not null" json:"sinkId"`
	Activated    int32     `gorm:"column:activated;not null" json:"activated"`
	RawSql       string    `gorm:"column:rawSql;not null" json:"rawSql"`
	CallbackID   int32     `gorm:"column:callbackId;not null" json:"callbackId"`
	Status       int32     `gorm:"column:status;not null" json:"status"`
	StartTime    time.Time `gorm:"column:startTime;not null" json:"startTime"`
	EndTime      time.Time `gorm:"column:endTime;not null" json:"endTime"`
	CreatedTime  time.Time `gorm:"column:createdTime;not null" json:"createdTime"`
	UpdatedTime  time.Time `gorm:"column:updatedTime;not null" json:"updatedTime"`
}

// TableName FunctionRule's table name
func (*FunctionRule) TableName() string {
	return TableNameFunctionRule
}
