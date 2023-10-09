package main

import (
	"context"
	"kafka-iot-connect/client/kafka/consumer"
	"kafka-iot-connect/client/kafka/producer"
	"kafka-iot-connect/controller"
	_ "kafka-iot-connect/docs"
	"kafka-iot-connect/model"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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
	kp := producer.InitializeKafkaConnectProducer("localhost:9092", "3.3.2", true, 1)
	kp.Produce(
		&sarama.ProducerMessage{Topic: "kafka-streams-energy-raw-data", Key: model.BMSDataTypeEncoder("user1", "elec"), Value: model.BMSRawDataEncoder("user1", 1, 839.543)},
		&sarama.ProducerMessage{Topic: "kafka-streams-energy-raw-data", Key: model.BMSDataTypeEncoder("user1", "elec"), Value: model.BMSRawDataEncoder("user1", 1, 869.543)},
	)

	kc := consumer.InitializeKafkaConnectConsumer("localhost:9092", "3.3.2", "example", "range", true, true)
	kc.Consume(context.Background(), "kafka-streams-energy-realtime-data", func(cm *sarama.ConsumerMessage) error {
		return nil
	})

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