// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameEkuiperDataSource = "ekuiperDataSource"

// EkuiperDataSource mapped from table <ekuiperDataSource>
type EkuiperDataSource struct {
	ID             int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ProjectID      string `gorm:"column:projectId;not null" json:"projectId"`
	DataSourceName string `gorm:"column:dataSourceName;not null" json:"dataSourceName"`
	SourceType     string `gorm:"column:sourceType;not null" json:"sourceType"`
	TargetType     string `gorm:"column:targetType;not null" json:"targetType"`
}

// TableName EkuiperDataSource's table name
func (*EkuiperDataSource) TableName() string {
	return TableNameEkuiperDataSource
}
