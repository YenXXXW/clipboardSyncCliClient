package grpcClient

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/YenXXXW/clipboardSyncCliClient/genproto/clipboardSync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type clipboardGrpcCleint struct {
	client pb.ClipSyncServiceClient
}

func NewGrpcClient(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did nto connect %v", err)
	}
	return conn
}

func (c *clipboardGrpcCleint) SendUpdate(content string, conn *grpc.ClientConn) error {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	reqContent := &pb.ClipboardContent{
		Text: content,
	}

	req := &pb.ClipboardUpdateRequest{
		Content: reqContent,
	}

	if _, err := c.client.SendClipboardUpdate(ctx, req); err != nil {
		return err
	}

	return nil

}

func (c *clipboardGrpcCleint) ReceiveUpdateAndSync(deviceId, roomId string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req := &pb.SubscribeRequest{
		DeviceId: deviceId,
		RoomId:   roomId,
	}

	stream, err := c.client.SubscribeClipboardContentUpdate(ctx, req)
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("End of stream")
			break
		}
		if err != nil {
			log.Fatalf("error receiving from stream: %v", err)
		}
		log.Printf("Received data: %s", resp.GetText())
	}

	return nil

}
