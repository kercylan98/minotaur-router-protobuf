# minotaur-router-protobuf [`protobufrouter`](https://pkg.go.dev/github.com/kercylan98/minotaur-router-protobuf/protobufrouter)

[![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/kercylan98/minotaur-router-protobuf/protobufrouter)

常使用且易于修改的基于 Protobuf 的路由器服务，它可以帮助开发者快速构建路由器服务。

## 使用示例

```go
package main

import (
	"fmt"
	"github.com/kercylan98/minotaur-router-protobuf/protobufrouter"
	"github.com/kercylan98/minotaur/server"
)

func main() {
	routerService := protobufrouter.NewDefault()
	routerService.Route(protobufrouter.MessageID_MI_Heartbeat, onHeartbeat)
	
	srv := server.New(server.NetworkWebsocket)
	server.BindService(srv, routerService)
	
	if err := srv.Run(":8080"); err != nil {
        	panic(err)
    	}
}


func onHeartbeat(service *protobufrouter.Service[protobufrouter.MessageID, *protobufrouter.Message, *server.Conn], entity *server.Conn, reader protobufrouter.Reader) {
	fmt.Println("onHeartbeat")
}

```
