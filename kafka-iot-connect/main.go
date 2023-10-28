package main

import (
	"context"
	"kafka-iot-connect/client/connect/energy"
	"kafka-iot-connect/controller"
	_ "kafka-iot-connect/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
) // gin-swagger middleware
// swagger embed files

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	// go func() {
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			log.Printf("Error occur: %v", err)
	// 		}
	// 	}()
	// 	// initializeKafkaConnect()
	// 	ctx := context.Background()
	// 	initializeMqttConnectClient(ctx)
	// }()

	ctx := context.Background()
	// go initializeMqttConnectClient(ctx)
	go energy.EnergyᚖMillsᚋMqttᚋClient(ctx)
	// go energy.EnergyᚖMillsᚋKafkaᚋClient(ctx)
	// go initializeKafkaConnectClient(ctx)

	r := gin.Default()

	c := controller.NewController()

	v1 := r.Group("/api/v1")
	{
		test := v1.Group("/test")
		{
			test.GET("", c.GetMap)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
