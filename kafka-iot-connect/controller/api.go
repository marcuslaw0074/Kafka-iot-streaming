package controller

import "github.com/gin-gonic/gin"

type Response struct {
	Title      map[string]string      `json:"title" example:"en:Map,ru:Карта,kk:Карталар"`
	CustomType map[string]interface{} `json:"map_data" swaggertype:"object,string" example:"key:value,key2:value2"`
	Object     Data                   `json:"object"`
}

type Data struct {
	Text string `json:"title" example:"Object data"`
}

// GetMap godoc
//
//	@Summary		Get Map Example
//	@Description	get map
//	@ID				get-map
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Router			/test [get]
func (c *Controller) GetMap(ctx *gin.Context) {
	ctx.JSON(200, Response{
		Title: map[string]string{
			"en": "Map",
		},
		CustomType: map[string]interface{}{
			"key": "value",
		},
		Object: Data{
			Text: "object text",
		},
	})
}