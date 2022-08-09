package main

import (
	"context"
	"errors"
	"fmt"

	smpb "github.com/blidd/fractr-proto/marketplace_secondary"
	pb "github.com/blidd/fractr-proto/service"
	stpb "github.com/blidd/fractr-proto/storage"
)

func (server *Server) GetBids(
	ctx context.Context,
	req *pb.GetBidsRequest,
) (*pb.GetBidsResponse, error) {
	resp, err := server.ls.Get(
		false,
		stpb.Type_BID,
		"",
		"",
		0,
		req.BidIds,
	)
	if err != nil {
		return nil, err
	}

	return &pb.GetBidsResponse{
		BidStatuses: resp.BidStatuses,
	}, nil
}

func (server *Server) GetAsks(
	ctx context.Context,
	req *pb.GetAsksRequest,
) (*pb.GetAsksResponse, error) {
	resp, err := server.ls.Get(
		false,
		stpb.Type_ASK,
		"",
		"",
		0,
		req.AskIds,
	)
	if err != nil {
		return nil, err
	}

	return &pb.GetAsksResponse{
		AskStatuses: resp.AskStatuses,
	}, nil
}

func (server *Server) Bid(
	ctx context.Context,
	req *pb.BidRequest,
) (*pb.BidResponse, error) {

	switch req.Market {

	case pb.Market_PRIMARY:
		return nil, nil

	case pb.Market_SECONDARY:
		client, err := server.NewSecondaryMarketServiceClient()
		if err != nil {
			return &pb.BidResponse{}, fmt.Errorf("error occurred while dialing secondary market service: %v", err)
		}
		resp, err := client.PlaceBid(ctx, &smpb.PlaceBidRequest{Bid: req.Bid})
		if err != nil {
			return &pb.BidResponse{}, fmt.Errorf("error occurred while placing bid with secondary market service: %v", err)
		}

		return &pb.BidResponse{
			Status: resp.BidStatus,
		}, nil

	default:
		return &pb.BidResponse{}, errors.New("market type doesn't exist")
	}
}

func (server *Server) Ask(
	ctx context.Context,
	req *pb.AskRequest,
) (*pb.AskResponse, error) {

	switch req.Market {

	case pb.Market_PRIMARY:
		return nil, nil

	case pb.Market_SECONDARY:
		client, err := server.NewSecondaryMarketServiceClient()
		if err != nil {
			return &pb.AskResponse{}, fmt.Errorf("error occurred while dialing secondary market service: %v", err)
		}
		resp, err := client.PlaceAsk(ctx, &smpb.PlaceAskRequest{Ask: req.Ask})
		if err != nil {
			return &pb.AskResponse{}, fmt.Errorf("error occurred while placing bid with secondary market service: %v", err)
		}

		return &pb.AskResponse{
			Status: resp.AskStatus,
		}, nil

	default:
		return &pb.AskResponse{}, errors.New("market type doesn't exist")
	}
}
