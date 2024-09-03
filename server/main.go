package main

import (
	"context"
	pb "github.com/miladjlz/go-grpc/proto"
	"github.com/miladjlz/go-grpc/util"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	pb.FileServiceServer
}

const (
	port      = ":8080"
	chunkSize = 1024
)

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Errorf("Failed to start server %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterFileServiceServer(grpcServer, &Server{})
	log.Printf("Server started at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Errorf("Failed to start: %v", err)
	}
}

func (s *Server) GetFileInfo(_ context.Context, req *pb.FileName) (*pb.FileInfoResponse, error) {

	files, err := util.ReadDirectory()
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.Name() == req.Name {
			fi, err := file.Info()
			if err != nil {
				logrus.Error("Error getting info -> ", err)
				return nil, err
			}
			logrus.WithFields(logrus.Fields{
				"name":     fi.Name(),
				"size":     fi.Size(),
				"mode":     fi.Mode(),
				"mod_time": fi.ModTime(),
				"is_dir":   fi.IsDir(),
			}).Info("File found: -> ")

			return &pb.FileInfoResponse{Name: fi.Name(), Size: fi.Size()}, nil
		}
	}
	return nil, nil
}

func (s *Server) GetFile(req *pb.FileName, stream pb.FileService_GetFileServer) error {

	files, err := util.ReadDirectory()
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.Name() == req.Name {

			chunks, err := util.SplitFileIntoChunks(req.Name, chunkSize)
			if err != nil {
				logrus.Error("Error split file into chunks -> ", err)
				return err
			}
			for _, chunk := range chunks {
				res := &pb.FileResponse{Chunk: chunk}
				if err := stream.Send(res); err != nil {
					logrus.Error("Error streaming -> ", err)
					return err
				}

			}

		}

	}
	return nil
}
