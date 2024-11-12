package main

import (
	grpc2 "geo-grpc-service/internal/grpc"
	"geo-grpc-service/internal/service"
	pb "geo-grpc-service/proto/generated"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/status"
	"log"
	"net"
	"os"
)

func main() {
	ApiKey, SecretKey := os.Getenv("api_key"), os.Getenv("secret_key")

	geoService := service.NewGeoService(ApiKey, SecretKey)
	
	lis, err := net.Listen("tcp", ":44044")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGeoServiceServer(grpcServer, grpc2.NewServer(geoService))

	log.Println("gRPC server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
