package service

import (
	"context"
	"github.com/RVodassa/grpc_geoservice/internal/domain/entity"
	"log"
	"net/url"

	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
)

// GeoServiceProvider Интерфейс для провайдера геосервиса
//
//go:generate mockgen -source=geoservice.go -destination=./mock/geoservice_mock.go -package=mock
type GeoServiceProvider interface {
	Search(ctx context.Context, input string) ([]*entity.Address, error)
	GeoCode(ctx context.Context, lat, lng string) ([]*entity.Address, error)
}

// geoService предоставляет методы для работы с геокодированием и автодополнением адресов
type geoService struct {
	api       *suggest.Api
	apiKey    string
	secretKey string
}

// NewGeoService создает экземпляр геосервиса с поддержкой кеширования и запросов к API.
func NewGeoService(apiKey, secretKey string) GeoServiceProvider {
	endpointUrl, err := url.Parse("https://suggestions.dadata.ru/suggestions/api/4_1/rs/")
	if err != nil {
		log.Println(err)
		return nil
	}

	credentials := client.Credentials{
		ApiKeyValue:    apiKey,
		SecretKeyValue: secretKey,
	}

	api := suggest.Api{
		Client: client.NewClient(endpointUrl, client.WithCredentialProvider(&credentials)),
	}

	return &geoService{
		api:       &api,
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

// Search метод поиска адресов
func (g *geoService) Search(ctx context.Context, input string) ([]*entity.Address, error) {
	var res []*entity.Address

	rawRes, err := g.api.Address(ctx, &suggest.RequestParams{Query: input})
	if err != nil {
		log.Printf("ошибка при обращении к API: %v", err)
		return nil, err
	}

	for _, r := range rawRes {
		if r.Data.City == "" || r.Data.Street == "" {
			continue
		}
		res = append(res, &entity.Address{
			City:   r.Data.City,
			Street: r.Data.Street,
			House:  r.Data.House,
			Lat:    r.Data.GeoLat,
			Lon:    r.Data.GeoLon,
		})
	}

	return res, nil
}

// GeoCode метод геокодирования (поиск адресов по координатам)
func (g *geoService) GeoCode(ctx context.Context, lat, lng string) ([]*entity.Address, error) {
	var res []*entity.Address

	rawRes, err := g.api.GeoLocate(ctx, &suggest.GeolocateParams{Lat: lat, Lon: lng})
	if err != nil {
		log.Printf("ошибка при обращении к API: %v", err)
		return nil, err
	}

	for _, r := range rawRes {
		address := &entity.Address{
			City:   r.Data.City,
			Street: r.Data.Street,
			House:  r.Data.House,
			Lat:    r.Data.GeoLat,
			Lon:    r.Data.GeoLon,
		}
		res = append(res, address)
	}

	return res, nil
}
