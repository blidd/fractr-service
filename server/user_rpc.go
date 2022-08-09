package main

import (
	"context"

	pb "github.com/blidd/fractr-proto/service"
	pbstore "github.com/blidd/fractr-proto/storage"
)

func (server *Server) GetUserProfile(
	ctx context.Context,
	req *pb.GetUserProfileRequest,
) (*pb.GetUserProfileResponse, error) {

	resp, err := server.ls.Get(
		false,
		pbstore.Type_USER,
		"",
		"",
		0,
		[]uint32{req.UserId},
	)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserProfileResponse{
		Profile: resp.Profile,
	}, nil
}
