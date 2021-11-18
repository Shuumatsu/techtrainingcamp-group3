package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/database"
	"techtrainingcamp-group3/pkg/logger"
	"techtrainingcamp-group3/rdb"

	pb "techtrainingcamp-group3/proto/pkg/user"
)

type userServer struct {
	rds *redis.Client
}

func (s *userServer) SnatchEnevelope(ctx context.Context, req *pb.SnatchEnevelopeReq) (*pb.SnatchEnevelopeReply, error) {
	return &pb.SnatchEnevelopeReply{}, nil
}

func (s *userServer) OpenEnvelope(ctx context.Context, req *pb.OpenEnvelopeReq) (*pb.OpenEnvelopeReply, error) {
	var envelope *database.Envelope
	// txf := func(tx *redis.Tx) error {
	// 	envelope, err := rdb.GetEnvelope(tx, req.EnvelopeId)
	// 	if err != nil && err != redis.Nil {
	// 		return err
	// 	}

	// 	if err == redis.Nil {

	// 	}

	// 	if !envelope.Opened {
	// 		envelope.Opened = true
	// 	}

	// 	return rdb.SetEnvelope(tx, envelope, 300*time.Second)
	// }

	// err := rdb.Client.Watch(txf, rdb.GetEnvelopeKey(req.EnvelopeId))
	// if err != nil && err != redis.TxFailedErr {
	// 	return nil, err
	// }

	// if err == redis.TxFailedErr {
	// 	return &proto.OpenEnvelopeReply{
	// 		ErrorType: proto.ErrorType_TxFailed,
	// 		Value:     envelope.Value,
	// 	}, nil
	// }

	return &pb.OpenEnvelopeReply{
		ErrorType: pb.ErrorType_NoError,
		Value:     envelope.Value,
	}, nil
}

func (s *userServer) ListEnvelopes(ctx context.Context, req *pb.ListEnvelopesReq) (*pb.ListEnvelopesReply, error) {
	user, err := rdb.GetUser(ctx, req.UserId)
	if err == redis.Nil {
		user, err = database.GetUser(req.UserId)
	}
	if err != nil {
		logger.Sugar.Errorf("ListEnvelopes", "req", req, "error", err)
		return &pb.ListEnvelopesReply{ErrorType: pb.ErrorType_Internal}, nil
	}
	_ = rdb.SetUser(ctx, user, time.Minute)

	list, err := database.ParseEnvelopeList(user.EnvelopeList)
	if err != nil {
		logger.Sugar.Errorf("ListEnvelopes", "req", req, "error", err)
		return &pb.ListEnvelopesReply{ErrorType: pb.ErrorType_Internal}, nil
	}

	envelopes := []*pb.Envelope{}
	for _, id := range list {
		envelope, err := rdb.GetEnvelope(ctx, id)
		if err == redis.Nil {
			envelope, err = database.GetEnvelope(id)
		}
		if err != nil {
			logger.Sugar.Errorf("ListEnvelopes", "req", req, "error", err)
			return &pb.ListEnvelopesReply{ErrorType: pb.ErrorType_Internal}, nil
		}
		_ = rdb.SetEnvelope(ctx, envelope, time.Minute)

		envelopes = append(envelopes, &pb.Envelope{
			EnvelopeId: envelope.EnvelopeId,
			Opened:     envelope.Opened,
			Value:      envelope.Value,
			SnatchTime: envelope.SnatchTime,
		})
	}

	return &pb.ListEnvelopesReply{
		ErrorType: pb.ErrorType_NoError,
		Envelopes: envelopes,
	}, nil
}

func newServer() *userServer {
	return &userServer{}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Env.RpcHost, config.Env.RpcPort))
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterUserServer(grpcServer, newServer())

	logger.Sugar.Info("RPC Server initialized.")

	grpcServer.Serve(lis)
}
