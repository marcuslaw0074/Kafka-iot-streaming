// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameAnaloguePoint = "analoguePoints"

// AnaloguePoint mapped from table <analoguePoints>
type AnaloguePoint struct {
	ID           int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	AliasID      string `gorm:"column:alias_id;not null" json:"alias_id"`
	BMSID        string `gorm:"column:BMS_id;not null" json:"BMS_id"`
	PointType    string `gorm:"column:pointType;not null" json:"pointType"`
	FileID       int32  `gorm:"column:FileId;not null" json:"FileId"`
	Description  string `gorm:"column:description;not null" json:"description"`
	ProjectID    string `gorm:"column:projectId;not null" json:"projectId"`
	DataSourceID int32  `gorm:"column:dataSourceId;not null" json:"dataSourceId"`
}

// TableName AnaloguePoint's table name
func (*AnaloguePoint) TableName() string {
	return TableNameAnaloguePoint
}
