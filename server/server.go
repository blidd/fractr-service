package main

import (
	"fmt"

	smpb "github.com/blidd/fractr-proto/marketplace_secondary"
	pb "github.com/blidd/fractr-proto/service"
	"github.com/blidd/fractr-service/libstore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Options struct {
	PrimaryMarketServiceAddr   string
	SecondaryMarketServiceAddr string
}

func DefaultServerOptions() Options {
	return Options{
		PrimaryMarketServiceAddr:   "[::1]:8081",
		SecondaryMarketServiceAddr: "[::1]:8082",
	}
}

type Server struct {
	pb.UnimplementedServiceServer

	opts Options

	ls libstore.Libstore
}

func (server *Server) NewSecondaryMarketServiceClient() (smpb.MarketplaceSecondaryClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(server.opts.SecondaryMarketServiceAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to secondary market service: %v", err)
	}

	return smpb.NewMarketplaceSecondaryClient(conn), nil
}

func main() {

}
