package rpc

import (
	"fmt"
	"techtrainingcamp-group3/config"

	pb "techtrainingcamp-group3/proto/pkg/user"

	"google.golang.org/grpc"
)

var Client pb.UserClient

func init() {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	addr := fmt.Sprintf("%s:%s", config.Env.RpcHost, config.Env.RpcPort)
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		panic(err)
	}

	Client = pb.NewUserClient(conn)
}
