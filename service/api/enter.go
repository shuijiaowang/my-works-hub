package api

import (
	service2 "workhub/service"
)

// HandlerGroup 包含所有处理器的结构
type ApiGroup struct {
	ExampleApi
	AdminApi
	ProjectApi
}

var (
	exampleService = service2.ExampleService{}
)
