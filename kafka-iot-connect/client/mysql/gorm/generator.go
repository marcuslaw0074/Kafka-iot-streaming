package main

import (
	"fmt"
	"kafka-iot-connect/client/mysql/model"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE status = @status{{end}}
	FilterWithStatus(status int) ([]gen.T, error)
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../query",
		ModelPkgPath: "../model",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	gormdb, _ := gorm.Open(mysql.Open("root:marcus@tcp(192.168.8.125:13306)/neuron_engine?parseTime=true"))
	g.UseDB(gormdb) // reuse your gorm db

	// Generate all tables in <database>
	// g.GenerateAllTable()

	// Generate basic type-safe DAO API for struct `model.RabbitMQMessage` following conventions
	// g.ApplyBasic(model.RabbitMQMessage{})
	// g.ApplyBasic(model.EkuiperClient{})
	// g.ApplyBasic()

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.RabbitMQMessage`
	// g.ApplyInterface(func(Querier) {}, model.RabbitMQMessage{})

	// Generate the code
	// g.Execute()

	var cls []model.EkuiperClient
	var sts []string

	gormdb.Model(&model.EkuiperClient{}).Pluck("serverId", &sts)

	fmt.Println(cls)
	fmt.Println(sts)

}
