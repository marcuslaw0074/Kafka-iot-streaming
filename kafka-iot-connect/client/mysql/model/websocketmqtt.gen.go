// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameWebsocketMqtt = "websocketMqtt"

// WebsocketMqtt mapped from table <websocketMqtt>
type WebsocketMqtt struct {
	ID            int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	WebsocketName string `gorm:"column:websocketName;not null" json:"websocketName"`
	ProjectID     string `gorm:"column:projectId;not null" json:"projectId"`
	ClientID      int32  `gorm:"column:clientId;not null" json:"clientId"`
	Status        int32  `gorm:"column:status;not null" json:"status"`
	Type          string `gorm:"column:type;not null" json:"type"`
	WebsocketType string `gorm:"column:websocketType;not null" json:"websocketType"`
	Args          string `gorm:"column:args;not null" json:"args"`
	StreamID      int32  `gorm:"column:streamId;not null" json:"streamId"`
	SinkID        int32  `gorm:"column:sinkId;not null" json:"sinkId"`
}

// TableName WebsocketMqtt's table name
func (*WebsocketMqtt) TableName() string {
	return TableNameWebsocketMqtt
}
