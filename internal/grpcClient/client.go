package grpcClient

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcClient(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did nto connect %v", err)
	}
	return conn
}

func main() {

}
