package main

import (
	"context"

	pb "github.com/blidd/fractr-proto/service"
	pbStore "github.com/blidd/fractr-proto/storage"
)

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

	return &pb.GetArtworksByTitleResponse{
		Artworks: resp.Artworks,
	}, nil

	// // query libstore
	// artworks, err := server.ls.Query(
	// 	"artwork",
	// 	[]libstore.Filter{{
	// 		Field: "title",
	// 		Value: req.Title,
	// 	}},
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// resp := &pb.GetArtworksByTitleResponse{
	// 	Artworks: make([]*pb.Artwork, 0),
	// }

	// for _, aw := range artworks {
	// 	artwork := aw.(pb.Artwork)
	// 	resp.Artworks = append(resp.Artworks, &artwork)
	// }

	// var getReq pbStore.GetRequest

	// return resp, nil
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

	return &pb.GetArtworksByArtistResponse{
		Artworks: resp.Artworks,
	}, nil

	// // query libstore
	// artworks, err := server.ls.Query(
	// 	"artwork",
	// 	[]libstore.Filter{{
	// 		Field: "artistName",
	// 		Value: req.ArtistName,
	// 	}},
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// resp := &pb.GetArtworksByArtistResponse{
	// 	Artworks: make([]*pb.Artwork, 0),
	// }

	// for _, aw := range artworks {
	// 	artwork := aw.(pb.Artwork)
	// 	resp.Artworks = append(resp.Artworks, &artwork)
	// }

	// return resp, nil
}
