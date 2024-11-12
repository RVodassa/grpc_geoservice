package grpc

import (
	"context"
	"fmt"
	"github.com/RVodassa/grpc_geoservice/internal/service"
	"github.com/RVodassa/grpc_geoservice/proto/generated"
)

// Server - структура, реализующая gRPC интерфейс
type Server struct {
	generated.UnimplementedGeoServiceServer
	geoService service.GeoServiceProvider
}

// NewServer - создает новый gRPC-сервер с GeoService
func NewServer(geoService service.GeoServiceProvider) *Server {
	return &Server{geoService: geoService}
}

// Search реализует gRPC метод Search
func (s *Server) Search(ctx context.Context, req *generated.SearchRequest) (*generated.SearchResponse, error) {
	addresses, err := s.geoService.Search(ctx, req.GetInput())
	if err != nil {
		return nil, fmt.Errorf("failed to search addresses: %w", err)
	}

	var responseAddresses []*generated.Address
	for _, addr := range addresses {
		responseAddresses = append(responseAddresses, &generated.Address{
			City:   addr.City,
			Street: addr.Street,
			House:  addr.House,
			Lat:    addr.Lat,
			Lon:    addr.Lon,
		})
	}

	return &generated.SearchResponse{Addresses: responseAddresses}, nil
}

// GeoCode реализует gRPC метод GeoCode
func (s *Server) GeoCode(ctx context.Context, req *generated.GeoCodeRequest) (*generated.GeoCodeResponse, error) {
	addresses, err := s.geoService.GeoCode(ctx, req.GetLat(), req.GetLng())
	if err != nil {
		return nil, fmt.Errorf("failed to geocode: %w", err)
	}

	var responseAddresses []*generated.Address
	for _, addr := range addresses {
		responseAddresses = append(responseAddresses, &generated.Address{
			City:   addr.City,
			Street: addr.Street,
			House:  addr.House,
			Lat:    addr.Lat,
			Lon:    addr.Lon,
		})
	}

	return &generated.GeoCodeResponse{Addresses: responseAddresses}, nil
}
