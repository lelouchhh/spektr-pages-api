package domain

import "context"

type City struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty" validate:"required"`
}
type CityTariff struct {
	City   int `json:"city_id"`
	Tariff int `json:"tariff_id"`
}

type CityUsecase interface {
	GetCities(ctx context.Context) ([]City, error)
	RemoveCity(ctx context.Context, Id int) error
	AddCity(ctx context.Context, city City) error
	RemoveCityTariff(ctx context.Context, tariff CityTariff) error
}

type CityRepository interface {
	GetCities(ctx context.Context) ([]City, error)
	RemoveCity(ctx context.Context, Id int) error
	AddCity(ctx context.Context, city City) error
	RemoveCityTariff(ctx context.Context, tariff CityTariff) error
}
