package main

import (
	pb "techtrainingcamp-group3/proto/pkg/user"
)

type UserServer struct{}

func (s *UserServer) SnatchEnevelope(req pb.SnatchEnevelopeReq) (pb.SnatchEnevelopeReply, error) {
	return pb.SnatchEnevelopeReply{}, nil
}

func (s *UserServer) OpenEnvelope(req pb.OpenEnvelopeReq) (pb.OpenEnvelopeReply, error) {
	return pb.OpenEnvelopeReply{}, nil
}

func (s *UserServer) ListEnvelopes(req pb.ListEnvelopesReq) (pb.ListEnvelopesReply, error) {
	return pb.ListEnvelopesReply{}, nil
}
