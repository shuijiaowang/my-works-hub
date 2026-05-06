package api

import (
	"htmlhub/model/request"
	"htmlhub/util/response"

	"github.com/gin-gonic/gin"
)

type ExampleApi struct{}

func (h *ExampleApi) Test(c *gin.Context) {
	var req request.ExampleRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 检查是否为该用户
	// 添加
	str := exampleService.AddExample()
	response.OkWithData(str, c)
}
