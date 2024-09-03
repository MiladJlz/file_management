package main

import (
	"context"
	"github.com/miladjlz/go-grpc/util"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"time"

	pb "github.com/miladjlz/go-grpc/proto"
	"google.golang.org/grpc"
)

const (
	addr = "localhost:8080"
)

func main() {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logrus.Errorf("Failed to connect %v", err)
	}
	defer conn.Close()
	client := pb.NewFileServiceClient(conn)

	//callGetFileInfo(client, "verylongtext.txt")
	callGetFile(client, "verylongtext.txt")

}
func callGetFileInfo(client pb.FileServiceClient, fileName string) {
	fn := pb.FileName{Name: fileName}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.GetFileInfo(ctx, &fn)
	if err != nil {
		logrus.Errorf("Failed to get info %v", err)
	}
	logrus.WithFields(logrus.Fields{
		"name": res.Name,
		"size": res.Size,
	}).Info("File received -> : ")

}
func callGetFile(client pb.FileServiceClient, fileName string) {
	fn := pb.FileName{Name: fileName}
	var chunks [][]byte
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := client.GetFile(ctx, &fn)
	if err != nil {
		logrus.Errorf("Colud not send names %v", err)
	}
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Errorf("Error while streaming %v", err)
		}
		chunks = append(chunks, chunk.Chunk)

	}
	util.SaveToStorage(err, chunks, fileName)
}
