// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameEkuiperDataSourceInfo = "ekuiperDataSourceInfo"

// EkuiperDataSourceInfo mapped from table <ekuiperDataSourceInfo>
type EkuiperDataSourceInfo struct {
	ID           int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	DataSourceID int32  `gorm:"column:dataSourceId;not null" json:"dataSourceId"`
	DbType       string `gorm:"column:dbType" json:"dbType"`
	DbName       string `gorm:"column:dbName" json:"dbName"`
	DbServer     string `gorm:"column:dbServer" json:"dbServer"`
	DbTable      string `gorm:"column:dbTable" json:"dbTable"`
}

// TableName EkuiperDataSourceInfo's table name
func (*EkuiperDataSourceInfo) TableName() string {
	return TableNameEkuiperDataSourceInfo
}