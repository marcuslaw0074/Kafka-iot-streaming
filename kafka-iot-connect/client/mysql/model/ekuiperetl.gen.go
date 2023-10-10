// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameEkuiperEtl = "ekuiperEtl"

// EkuiperEtl mapped from table <ekuiperEtl>
type EkuiperEtl struct {
	ID        int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	EtlName   string `gorm:"column:etlName;not null" json:"etlName"`
	ClientID  int32  `gorm:"column:clientId;not null" json:"clientId"`
	ProjectID string `gorm:"column:projectId;not null" json:"projectId"`
	Type      string `gorm:"column:type;not null" json:"type"`
	EtlType   string `gorm:"column:etlType;not null" json:"etlType"`
	Desc      string `gorm:"column:desc;not null" json:"desc"`
	Status    int32  `gorm:"column:status;not null" json:"status"`
	Args      string `gorm:"column:args;not null" json:"args"`
	StreamID  int32  `gorm:"column:streamId;not null" json:"streamId"`
	SinkID    int32  `gorm:"column:sinkId;not null" json:"sinkId"`
}

// TableName EkuiperEtl's table name
func (*EkuiperEtl) TableName() string {
	return TableNameEkuiperEtl
}