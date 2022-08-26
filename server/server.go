package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	smpb "github.com/blidd/fractr-proto/marketplace_secondary"
	pb "github.com/blidd/fractr-proto/service"
	pbStore "github.com/blidd/fractr-proto/storage"
	"github.com/blidd/fractr-service/libstore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Options struct {
	PrimaryMarketServiceAddr   string
	SecondaryMarketServiceAddr string
}

// func DefaultServerOptions() Options {
// 	return Options{
// 		PrimaryMarketServiceAddr:   "[::1]:8081",
// 		SecondaryMarketServiceAddr: "[::1]:8082",
// 	}
// }

type Server struct {
	pb.UnimplementedServiceServer
	opts Options
	ls   *libstore.Libstore
}

func New() *Server {
	return &Server{
		opts: Options{
			SecondaryMarketServiceAddr: fmt.Sprintf("[::1]:%d", *secondaryMarketPort),
		},
		ls: libstore.NewLibstore(string(fmt.Sprintf("[::1]:%d", *storageServicePort))),
	}
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

func (server *Server) GetArtworkDetails(
	ctx context.Context,
	req *pb.GetArtworkDetailsRequest,
) (*pb.GetArtworkDetailsResponse, error) {

	resp, err := server.ls.Get(
		false,
		pbStore.Type_ARTWORK,
		"",
		"",
		req.ArtworkId,
		[]uint32{},
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Artworks) > 0 {
		return &pb.GetArtworkDetailsResponse{
			ArtworkDetails: resp.Artworks[0],
		}, nil
	} else {
		return &pb.GetArtworkDetailsResponse{
			ArtworkDetails: &pb.Artwork{
				Id: req.GetArtworkId(),
				Artist: &pb.Artist{
					Id:          2000,
					Name:        "vincent van gogh",
					Description: "crazy guy",
					Handle:      "@vangogh",
				},
				Name:           "a starry night",
				OwnerIds:       []uint32{1, 2, 3},
				Description:    "vincent van gogh made this whe he was insane",
				NumShares:      100,
				Market:         pb.Market_SECONDARY,
				ConversionDate: 10,
			},
		}, nil
	}
}

func (server *Server) GetArtworkLiveData(
	req *pb.GetArtworkLiveDataRequest,
	stream pb.Service_GetArtworkLiveDataServer,
) error {
	// simulate random data updates
	for i := 0; i < 5; i++ {
		data := &pb.ArtworkLiveData{
			ArtworkId:          req.GetArtworkId(),
			HiBidPrice:         uint32(10 - i),
			LoAskPrice:         uint32(i),
			LastTxPrice:        5,
			NumSharesAvailable: uint32(100 - i),
		}
		resp := &pb.GetArtworkLiveDataResponse{
			TimeUpdated: uint32(time.Now().Unix()),
			LiveData:    data,
		}
		stream.Send(resp)
		fmt.Printf("Sent: %+v\n", resp)

		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (server *Server) GetArtworksByTitle(
	ctx context.Context,
	req *pb.GetArtworksByTitleRequest,
) (*pb.GetArtworksByTitleResponse, error) {

	resp, err := server.ls.Get(
		false,
		pbStore.Type_ARTWORK,
		"",
		req.Title,
		0,
		[]uint32{},
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Artworks) > 0 {
		return &pb.GetArtworksByTitleResponse{
			Artworks: resp.Artworks,
		}, nil
	} else {
		return &pb.GetArtworksByTitleResponse{}, nil
	}
}

func (server *Server) GetArtworksByArtist(
	ctx context.Context,
	req *pb.GetArtworksByArtistRequest,
) (*pb.GetArtworksByArtistResponse, error) {

	resp, err := server.ls.Get(
		false,
		pbStore.Type_ARTWORK,
		req.ArtistName,
		"",
		0,
		[]uint32{},
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Artworks) > 0 {
		return &pb.GetArtworksByArtistResponse{
			Artworks: resp.Artworks,
		}, nil
	} else {
		return &pb.GetArtworksByArtistResponse{}, nil
	}
}

func (server *Server) GetUserProfile(
	ctx context.Context,
	req *pb.GetUserProfileRequest,
) (*pb.GetUserProfileResponse, error) {

	resp, err := server.ls.Get(
		false,
		pbStore.Type_USER,
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

var (
	port                = flag.Int("port", 9090, "Server port")
	secondaryMarketPort = flag.Int("market2-port", 8082, "Secondary market port")
	storageServicePort  = flag.Int("storage-port", 8083, "Storage port")
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterServiceServer(
		s,
		New(),
	)
	reflection.Register(s)
	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
