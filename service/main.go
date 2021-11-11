package main

import (
	"context"
	"fmt"
	"net"
	user "techtrainingcamp-group3/proto/pkg/user"

	"techtrainingcamp-group3/service/config"

	"google.golang.org/grpc"
)

type userServer struct{}

func (s *userServer) SnatchEnevelope(ctx context.Context, req *user.SnatchEnevelopeReq) (*user.SnatchEnevelopeReply, error) {
	return &user.SnatchEnevelopeReply{}, nil
}

func (s *userServer) OpenEnvelope(ctx context.Context, req *user.OpenEnvelopeReq) (*user.OpenEnvelopeReply, error) {
	return &user.OpenEnvelopeReply{}, nil
}

func (s *userServer) ListEnvelopes(ctx context.Context, req *user.ListEnvelopesReq) (*user.ListEnvelopesReply, error) {
	return &user.ListEnvelopesReply{}, nil
}

func newServer() *userServer {
	return &userServer{}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Env.Host, config.Env.Port))
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	user.RegisterUserServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
