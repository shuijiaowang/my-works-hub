package api

import (
	service2 "htmlhub/service"
)

// HandlerGroup 包含所有处理器的结构
type ApiGroup struct {
	ExampleApi
	AdminApi
}

var (
	exampleService = service2.ExampleService{}
)
