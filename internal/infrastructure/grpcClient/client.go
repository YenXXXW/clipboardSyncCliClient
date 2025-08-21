package infrastructure

import (
	"context"
	"io"
	"log"

	pb "github.com/YenXXXW/clipboardSyncCliClient/genproto/clipboardSync"
	"github.com/YenXXXW/clipboardSyncCliClient/internal/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type clipboardGrpcClient struct {
	client pb.ClipSyncServiceClient
}

func NewGrpcClient(addr string) *clipboardGrpcClient {
	log.Printf("Connecting to gRPC server at %s", addr)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}

	log.Println("Connected to gRPC server")

	c := pb.NewClipSyncServiceClient(conn)
	clipboardClient := &clipboardGrpcClient{
		client: c,
	}
	return clipboardClient
}

func (c *clipboardGrpcClient) SendUpdate(ctx context.Context, deviceId, content string) error {

	reqContent := &pb.ClipboardContent{
		Text: content,
	}

	req := &pb.ClipboardUpdate{
		Content:  reqContent,
		DeviceId: deviceId,
	}

	if _, err := c.client.SendClipboardUpdate(ctx, req); err != nil {
		return err
	}

	return nil

}

func (c *clipboardGrpcClient) ReceiveUpdateAndSync(ctx context.Context, deviceId, roomId string, updateChan chan *types.ClipboardUpdate) error {

	req := &pb.SubscribeRequest{
		DeviceId: deviceId,
		RoomId:   roomId,
	}

	stream, err := c.client.SubscribeClipboardContentUpdate(ctx, req)
	if err != nil {
		return err
	}
	log.Println("Subscribed to room: ", roomId)

	go func() {
		defer close(updateChan)

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("End of stream")
				break
			}
			if err != nil {
				log.Fatalf("error receiving from stream: %v", err)
				return
			}

			clipboardDataUpdate := &types.ClipboardUpdate{
				DeviceId: resp.GetDeviceId(),
				Content: &types.ClipboardContent{
					Text: resp.GetContent().GetText(),
				},
			}

			select {
			case updateChan <- clipboardDataUpdate:
			case <-ctx.Done():
				log.Println("Context cancelled while sending update")
				return

			}
		}
	}()

	return nil

}

func (c *clipboardGrpcClient) CreateRoom(ctx context.Context, deviceId string) (string, error) {
	req := &pb.CreateRoomRequest{
		DeviceId: deviceId,
	}

	res, err := c.client.CreateRoom(ctx, req)
	if err != nil {
		log.Printf("Error creating room %v", err)
		return "", err
	}

	return res.GetRoomId(), nil

}

func (c *clipboardGrpcClient) LeaveRoom(ctx context.Context, deviceId, roomId string) error {
	req := &pb.LeaveRoomRequest{
		DeviceId: deviceId,
		RoomId:   roomId,
	}

	if _, err := c.client.LeaveRoom(ctx, req); err != nil {
		log.Printf("Error leaving room %v", err)
		return err
	}

	return nil

}
