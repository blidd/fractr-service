package libstore

import (
	"context"
	"fmt"

	pb "github.com/blidd/fractr-proto/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Filter struct {
	Field string
	Value string
}

// type LibstoreAPI interface {
// 	Get(docType string, key string) (string, error)
// 	Put(docType string, key, value string) error
// 	Delete(docType string, key string) error
// 	Query(docType string, filters []Filter) ([]interface{}, error)
// }

type Libstore struct {
	StorageServiceAddr string
}

func NewLibstore(addr string) *Libstore {
	return &Libstore{
		StorageServiceAddr: addr,
	}
}

func (ls *Libstore) NewStorageServiceClient() (pb.StorageClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(ls.StorageServiceAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to secondary market service: %v", err)
	}

	return pb.NewStorageClient(conn), nil
}

func (ls *Libstore) Get(primaryMarket bool, docType pb.Type, artistName, artTitle string, id uint32, ids []uint32) (*pb.GetResponse, error) {
	req := &pb.GetRequest{
		Type:          docType,
		Id:            &id,
		Ids:           ids,
		PrimaryMarket: &primaryMarket,
		ArtistName:    &artistName,
		ArtTitle:      &artTitle,
	}
	client, err := ls.NewStorageServiceClient()
	if err != nil {
		return &pb.GetResponse{}, fmt.Errorf("error occurred while dialing storage service: %v", err)
	}

	return client.Get(context.Background(), req)
}

func (ls *Libstore) Put(docType pb.Type, id uint32, ids []uint32, primaryMarket bool, artistName, artTitle string) (*pb.PutResponse, error) {
	req := &pb.PutRequest{
		Type:          docType,
		Id:            &id,
		Ids:           ids,
		PrimaryMarket: &primaryMarket,
		ArtistName:    &artistName,
		ArtTitle:      &artTitle,
	}
	client, err := ls.NewStorageServiceClient()
	if err != nil {
		return &pb.PutResponse{}, fmt.Errorf("error occurred while dialing storage service: %v", err)
	}

	return client.Put(context.Background(), req)

}

func (ls *Libstore) Delete() {

}
