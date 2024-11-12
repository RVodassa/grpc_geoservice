package main

import (
	grpc2 "github.com/RVodassa/grpc_geoservice/internal/grpc"
	"github.com/RVodassa/grpc_geoservice/internal/service"
	pb "github.com/RVodassa/grpc_geoservice/proto/generated"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/status"
	"log"
	"net"
	"os"
)

func main() {
	// Загружаем переменные из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	ApiKey, SecretKey := os.Getenv("apiKey"), os.Getenv("secretKey")

	geoService := service.NewGeoService(ApiKey, SecretKey)

	lis, err := net.Listen("tcp", ":44044")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGeoServiceServer(grpcServer, grpc2.NewServer(geoService))

	log.Println("gRPC server is running on port 44044")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
